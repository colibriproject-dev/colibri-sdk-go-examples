# Deploy na GCP — colibri-sdk-go-examples

Mesma ideia do [deploy na AWS](../aws/README.md): sobe os dois microserviços de
exemplo, a estrutura de mensageria entre eles e a observabilidade nativa da
nuvem, deixando as duas APIs abertas para chamadas externas.

**O banco de dados não é criado aqui.** Não há Cloud SQL nem Memorystore: as
credenciais de um PostgreSQL e de um Redis remotos entram por variável e são
repassadas aos serviços — as senhas via Secret Manager, o resto como variável de
ambiente comum.

> **Leia a seção de custo antes de aplicar.** Este ambiente custa
> significativamente mais que o da AWS, por um motivo estrutural que não dá para
> contornar sem mexer no código dos serviços.

---

## O que é provisionado

| Camada | Recursos |
|---|---|
| Compute | 2 serviços Cloud Run, cada um com container da aplicação + sidecar do collector |
| Entrada | URL HTTPS pública por serviço, gerada pelo Cloud Run (sem load balancer) |
| Mensageria | 4 tópicos Pub/Sub, 4 subscriptions, 4 tópicos de dead letter + suas subscriptions |
| Armazenamento | Bucket do Cloud Storage para upload de documentos |
| Imagens | Repositório Docker no Artifact Registry com cleanup policy |
| Observabilidade | Sidecar OTel por serviço → Cloud Trace e Cloud Monitoring; logs no Cloud Logging |
| Configuração | Secret Manager para as senhas e para a config do collector |
| Identidade | Uma service account por módulo, com papéis mínimos |

O que **não** é provisionado: PostgreSQL, Redis, VPC, Cloud NAT, domínio e
certificado próprios.

---

## Topologia

```
                    Internet
                       │
        ┌──────────────┴──────────────┐
        │                             │
   URL do Cloud Run             URL do Cloud Run
        │                             │
 ┌──────▼────────┐            ┌───────▼─────────┐
 │ school-module │───HTTP────▶│ finantial-module│   Cloud Run
 │  + collector  │            │  + collector    │   CPU sempre alocada
 └──┬─────────┬──┘            └────┬─────────┬──┘
    │         ▲                    │         ▲
    │         └────────Pub/Sub─────┼─────────┘
    └──────────────────Pub/Sub─────┘
                       │
          PostgreSQL + Redis remotos (fora do projeto)
```

Fluxo das mensagens:

| Publica | Tópico | Subscription | Consome |
|---|---|---|---|
| school | `SCHOOL_ENROLLMENT` | `SCHOOL_ENROLLMENT_FINANCIAL` | finantial |
| school | `SCHOOL_COURSE` | `SCHOOL_COURSE_FINANCIAL` | finantial |
| school | `SCHOOL_STUDENT` | `SCHOOL_STUDENT_FINANCIAL` | finantial |
| finantial | `FINANCIAL_INSTALLMENT` | `FINANCIAL_INSTALLMENT_SCHOOL` | school |

Cada subscription tem um tópico de dead letter `<NOME>_DLQ` — com uma
subscription própria, porque no Pub/Sub o descarte vai para outro tópico e sem
subscription anexada a mensagem se perderia ali.

> Os nomes acima são literais no código Go (`messaging.NewProducer("SCHOOL_COURSE")`,
> `QueueName() == "SCHOOL_COURSE_FINANCIAL"`). O SDK chama `client.Publisher(nome)`
> e `client.Subscriber(nome)` direto, então **não dá para prefixar nem sufixar**
> esses nomes sem alterar o código. Um projeto GCP comporta um ambiente deste
> stack.

---

## Custo — leia antes de aplicar

**Cerca de US$ 90/mês com o ambiente parado**, contra ~US$ 55 na AWS. A razão é
específica e vale entender:

Os consumidores do colibri-sdk-go fazem **pull** do Pub/Sub numa goroutine de
fundo. O Cloud Run, por padrão, congela a CPU do container entre requisições —
esse loop pararia de buscar mensagens assim que o tráfego HTTP cessasse. Para o
ambiente funcionar é preciso CPU sempre alocada (`cpu_idle = false`) e pelo
menos uma instância viva (`min_instance_count = 1`).

E aí entra a restrição que fecha a conta: o Cloud Run **só aceita CPU
fracionária com billing por requisição**, que é justamente o modo incompatível
com CPU sempre alocada. Ou seja, o mínimo é **1 vCPU inteira por serviço**,
ligada 24 horas por dia.

| Item | Estimativa mensal |
|---|---|
| 2 serviços × 1 vCPU sempre alocada | ~US$ 94 |
| 2 serviços × 768 MiB de memória | ~US$ 4 |
| Free tier (240 mil vCPU-s + 450 mil GiB-s) | −US$ 5 |
| Pub/Sub, Storage, Artifact Registry, Trace no volume de exemplo | poucos dólares |
| **Total** | **~US$ 93** |

O sidecar do collector não declara CPU própria de propósito: ele divide a vCPU
da instância com a aplicação. Declarar CPU nele somaria ao total da instância e
**dobraria** essa conta.

### Como não pagar por um ambiente ocioso

```bash
./scripts/pause.sh       # min_instance_count = 0 — a cobrança contínua para
./scripts/pause.sh -r    # volta ao normal
```

Com as instâncias em zero as URLs continuam respondendo (o Cloud Run sobe uma
instância em segundos), e você paga só pelas requisições. O que para é a
mensageria: sem instância viva ninguém faz pull das subscriptions. As mensagens
não se perdem — ficam retidas pelo período de retenção e são consumidas quando
uma instância subir — mas o processamento deixa de ser em tempo real.

Para zerar de vez, `./scripts/destroy.sh`.

### Se o custo for proibitivo

A saída de verdade seria trocar pull por **push**: uma subscription push do
Pub/Sub entrega a mensagem como um POST HTTP, o Cloud Run escala do zero para
atender e volta a dormir — custo praticamente zero em repouso, que é o modelo
que a GCP foi desenhada para servir. O SDK hoje só implementa pull
(`pkg/messaging/gcp.go` usa `sub.Receive`), então isso exigiria um endpoint novo
nos módulos. Não está feito aqui.

---

## Pré-requisitos

- Terraform >= 1.6
- `gcloud` autenticado (`gcloud auth login` **e** `gcloud auth application-default login`)
- Docker
- `jq` e `curl` para o smoke test
- Um projeto GCP com billing ativo
- Um PostgreSQL acessível pela internet, com os dois databases já criados
  (`school_module` e `finantial_module`) — as tabelas as migrations criam sozinhas
- Um Redis acessível pela internet (usado só pelo `school-module`)

> Como escolhemos o egress padrão do Cloud Run, a saída usa IPs do pool do
> Google, que mudam. O banco precisa aceitar conexão de qualquer origem — use
> TLS (`db_ssl_mode = "require"`) e senha forte. Se o seu banco exigir allowlist
> de IP, o caminho é Direct VPC egress + Cloud NAT com IP reservado, que não
> está neste Terraform.

---

## Uso

```bash
cd deploy/gcp/terraform
cp terraform.tfvars.example terraform.tfvars
$EDITOR terraform.tfvars          # project_id, db_host, db_user, db_password, cache_uri

cd ../scripts
./deploy.sh
```

Três fases, na mesma lógica do deploy AWS — o Cloud Run não sobe uma revisão
cuja imagem ainda não existe:

1. **bootstrap** — habilita as APIs e cria o repositório do Artifact Registry
2. **build** — compila e publica as imagens
3. **apply** — cria o resto

```bash
./deploy.sh --plan            # só mostra o plano
./deploy.sh --skip-build      # reaproveita as imagens já publicadas
./deploy.sh --auto-approve    # não pergunta antes de aplicar
./deploy.sh --tag v1.2.3      # fixa uma tag em vez de "latest"
```

Depois:

```bash
./status.sh        # revisões, subscriptions e health check
./smoke-test.sh    # cria curso + aluno + matrícula e confere se a mensagem chegou
./pause.sh         # zera o custo ocioso sem destruir nada
./destroy.sh       # remove tudo e audita o que sobrou
```

---

## Endpoints

O `terraform output` imprime as URLs — o Cloud Run gera uma por serviço, em
HTTPS, sem load balancer no meio. As rotas de negócio ficam todas sob `/public/`
(prefixo `restserver.PublicApi`); só `/health` e `/api-docs` respondem na raiz.

```
https://school-module-<hash>-<região>.run.app/health
https://school-module-<hash>-<região>.run.app/api-docs
https://school-module-<hash>-<região>.run.app/public/v1/courses
https://school-module-<hash>-<região>.run.app/public/v1/students
https://school-module-<hash>-<região>.run.app/public/v1/enrollments

https://finantial-module-<hash>-<região>.run.app/public/accounts
https://finantial-module-<hash>-<região>.run.app/public/invoices
```

O `school-module` chama o `finantial-module` pela URL pública dele — não há rede
interna aqui, e como os dois estão abertos, não há token a injetar.

---

## Observabilidade

| Sinal | Destino | Onde ver |
|---|---|---|
| Traces | Cloud Trace (exporter `googlecloud`) | Trace explorer |
| Métricas | Cloud Monitoring (exporter `googlecloud`) | Metrics explorer |
| Logs | Cloud Logging (stdout do container) | Logs de cada serviço no console do Cloud Run |

Nada muda no código dos serviços — só a variável `OTEL_EXPORTER_OTLP_ENDPOINT`,
que aponta para `localhost:4318`, o sidecar dentro da mesma instância.

Diferente do X-Ray, o **Cloud Trace aceita trace ID no formato W3C** sem
conversão. O deploy AWS precisava de uma versão mínima do ADOT por causa disso;
aqui não há restrição equivalente.

Para ver um trace atravessando os dois serviços, chame
`GET /public/v1/courses/{id}` no school-module: esse endpoint faz uma chamada
HTTP para o finantial-module.

---

## Variáveis principais

### Projeto

| Variável | Descrição |
|---|---|
| `project_id` | **Obrigatória.** Projeto GCP onde tudo é criado |
| `region` | Padrão `southamerica-east1` |
| `name_prefix` | Padrão `colibri-examples` — prefixa bucket, registry e service accounts |
| `enable_apis` | Padrão `true`. Deixe `false` se sua organização gerencia APIs à parte |

### Banco remoto (obrigatórias)

| Variável | Descrição |
|---|---|
| `db_host`, `db_user`, `db_password` | Conexão com o PostgreSQL. A senha vai para o Secret Manager |
| `db_port` | Padrão `5432` |
| `db_ssl_mode` | Padrão `require` — a saída é pela internet aberta |
| `db_name_school` / `db_name_finantial` | Padrões `school_module` / `finantial_module` |
| `db_run_migrations` | Padrão `true` |

### Cache remoto

| Variável | Descrição |
|---|---|
| `cache_uri` | **Obrigatória.** `host:porta` do Redis. Usado só pelo `school-module` |
| `cache_password` | Vazio se não houver AUTH |

### Cloud Run

| Variável | Padrão | Observação |
|---|---|---|
| `min_instance_count` | `1` | Abaixo de 1 a mensageria para de funcionar em tempo real |
| `max_instance_count` | `3` | |
| `app_cpu` | `"1"` | Não use fracionário — incompatível com CPU sempre alocada |
| `app_memory` | `"512Mi"` | |
| `otel_memory` | `"256Mi"` | O sidecar não declara CPU: divide a da instância |
| `allow_public_access` | `true` | `false` exige identity token em toda chamada |

---

## Detalhes que o ambiente depende

Quatro pontos que parecem detalhe e não são:

- **`PUBSUB_PROJECT_ID`** — é esta a variável que o cliente Pub/Sub do SDK lê
  (`pkg/messaging/gcp.go`), e não `CLOUD_PROJECT_ID`, que aparece no `.env` local
  mas não é consumida em lugar nenhum do código. Sem ela o cliente sobe sem
  projeto e nenhuma publicação funciona.
- **`PORT` não é declarada** — o Cloud Run injeta essa variável com o valor de
  `container_port` e **recusa a revisão** se ela vier declarada à mão. O SDK lê
  `PORT` do ambiente, então a injeção do Cloud Run já resolve. Por isso os dois
  serviços escutam na mesma porta interna: cada um tem URL própria em `:443`.
- **`cpu_idle = false`** — sem isso o loop de pull do Pub/Sub congela entre
  requisições. Ver a seção de custo.
- **IAM do service agent do Pub/Sub** — quem move a mensagem para o dead letter
  topic não é a aplicação, é o próprio Pub/Sub, com a identidade
  `service-<número>@gcp-sa-pubsub`. Ele precisa de `publisher` no tópico de
  destino e `subscriber` na subscription de origem. Sem os dois papéis a
  `dead_letter_policy` é aceita no apply e simplesmente não funciona em runtime:
  a mensagem fica redoando para sempre.

E uma diferença estrutural em relação à AWS: no Pub/Sub o corpo da mensagem
chega exatamente como o produtor publicou. Não existe o envelope de notificação
que o SNS coloca por cima do payload — logo, não existe o equivalente ao
`raw_message_delivery` para configurar errado.

---

## Limitação conhecida: o bucket do upload de documentos

Mesma do deploy AWS. O Terraform cria o bucket, dá permissão à service account
do school e passa o nome em `STORAGE_BUCKET`, mas o código não lê essa variável:
`school-module/src/infra/storages/document_storage.go` traz o nome fixo.

```go
func NewDocumentS3Storage() *DocumentS3Storage {
	return &DocumentS3Storage{
		bucket: "meu-bucket",
	}
}
```

Tudo funciona menos `POST /public/v1/students/{id}/upload-document`. O ajuste é
de uma linha, no repositório da aplicação:

```go
bucket: os.Getenv("STORAGE_BUCKET"),
```

O ambiente local continua funcionando porque o `.env` já define
`STORAGE_BUCKET=meu-bucket`.

---

## Problemas comuns

**A revisão não fica pronta.** Quase sempre é o banco. Veja os logs:

```bash
gcloud run services logs read school-module --region southamerica-east1 --limit 100
```

Erro de conexão indica que o `db_host` não aceita a origem — lembre que o IP de
saída do Cloud Run muda.

**`The following reserved env names were provided: PORT`.** Alguém declarou
`PORT` nas variáveis. O Cloud Run injeta essa variável sozinho.

**Mensagem publicada e não consumida.** Primeiro cheque
`min_instance_count`: com 0 não há quem faça pull. Depois a subscription
`<NOME>_DLQ`.

**Nenhum trace no Cloud Trace.** Veja os logs do container `collector` do
serviço — ele loga o motivo da recusa do exporter.

**O destroy falha.** `deletion_protection` já vem `false` nos dois serviços; se
o erro for de tópico com subscription anexada, use `./destroy.sh --sweep`, que
remove as subscriptions antes dos tópicos.

---

## Estrutura

```
deploy/
├── docker/
│   └── Dockerfile            build self-contained, compartilhado com o deploy da AWS
├── aws/                      deploy equivalente na AWS
└── gcp/
    ├── README.md
    ├── scripts/
    │   ├── common.sh         funções compartilhadas
    │   ├── deploy.sh         ciclo completo em 3 fases
    │   ├── build-and-push.sh imagens → Artifact Registry
    │   ├── status.sh         revisões, subscriptions e health check
    │   ├── smoke-test.sh     valida API + mensageria + banco
    │   ├── pause.sh          zera o custo ocioso sem destruir nada
    │   └── destroy.sh        remove o ambiente e audita o que sobrou
    └── terraform/
        ├── versions.tf       provider e versões
        ├── variables.tf      entradas
        ├── locals.tf         nomes de tópicos/subscriptions e derivações
        ├── apis.tf           APIs do projeto
        ├── messaging.tf      tópicos, subscriptions, dead letter e IAM do service agent
        ├── storage.tf        bucket
        ├── registry.tf       Artifact Registry
        ├── secrets.tf        Secret Manager
        ├── iam.tf            service accounts e papéis mínimos
        ├── observability.tf  config do collector
        ├── run.tf            serviços Cloud Run e acesso público
        └── outputs.tf        URLs, nomes e comandos de teste
```

---

## Diferenças em relação ao deploy AWS

| | AWS | GCP |
|---|---|---|
| Compute | ECS Fargate | Cloud Run |
| Entrada | ALB público, 2 listeners por porta | URL HTTPS por serviço, sem LB |
| Mensageria | SNS → SQS (com envelope de notificação) | Pub/Sub (payload direto) |
| Dead letter | Fila SQS | Tópico + subscription |
| Rede interna | Cloud Map + VPC | URL pública (sem VPC) |
| Traces | X-Ray via ADOT (exige collector ≥ 0.34) | Cloud Trace (aceita W3C nativamente) |
| Segredos | SSM Parameter Store | Secret Manager |
| Custo parado | ~US$ 55/mês | ~US$ 93/mês |
| Escala a zero | não (Fargate contínuo) | sim, mas para a mensageria — ver `pause.sh` |
