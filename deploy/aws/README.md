# Deploy na AWS — colibri-sdk-go-examples

Sobe os dois microserviços de exemplo (`school-module` e `finantial-module`), a
estrutura de mensageria entre eles e a observabilidade nativa da AWS, deixando
as duas APIs abertas para chamadas externas.

**O banco de dados não é criado aqui.** Não há RDS nem ElastiCache no Terraform:
as credenciais de um PostgreSQL e de um Redis remotos entram por variável e são
repassadas aos serviços — a senha via SSM Parameter Store, o resto como variável
de ambiente comum.

---

## O que é provisionado

| Camada | Recursos |
|---|---|
| Rede | VPC dedicada, subnets públicas em 2 AZs (privadas + NAT opcionais), Security Groups |
| Compute | Cluster ECS, 2 serviços em Fargate, 1 task definition por módulo |
| Entrada | Application Load Balancer público, 2 target groups, 2 listeners |
| Mensageria | 4 tópicos SNS, 4 filas SQS, 4 DLQs, subscriptions e queue policies |
| Armazenamento | Bucket S3 para upload de documentos do aluno |
| Imagens | 2 repositórios ECR com lifecycle policy |
| Observabilidade | Sidecar ADOT por task → X-Ray (traces) e CloudWatch (métricas EMF); log groups |
| Configuração | SSM Parameter Store (SecureString) para as senhas do banco e do cache |

O que **não** é provisionado, por decisão de projeto: PostgreSQL, Redis,
certificado TLS e domínio próprio.

---

## Topologia

```
                    Internet
                       │
              ┌────────┴────────┐
              │       ALB       │   :8080 → school     :8081 → finantial
              └───┬─────────┬───┘
                  │         │
      ┌───────────▼──┐   ┌──▼──────────────┐
      │ school-module│   │ finantial-module│   ECS Fargate
      │  + ADOT      │   │  + ADOT         │
      └───┬───┬──────┘   └────┬────┬───────┘
          │   │  ▲            │    │  ▲
          │   │  └────────────┘    │  │        HTTP interno (Cloud Map)
          │   │                    │  │
          │   └──────── SNS ───────┼──┘        FINANCIAL_INSTALLMENT
          └───────────► SNS ───────┘           SCHOOL_ENROLLMENT / _COURSE / _STUDENT
                         │
                    PostgreSQL + Redis remotos (fora da conta)
```

Fluxo das mensagens:

| Publica | Tópico SNS | Fila SQS | Consome |
|---|---|---|---|
| school | `SCHOOL_ENROLLMENT` | `SCHOOL_ENROLLMENT_FINANCIAL` | finantial |
| school | `SCHOOL_COURSE` | `SCHOOL_COURSE_FINANCIAL` | finantial |
| school | `SCHOOL_STUDENT` | `SCHOOL_STUDENT_FINANCIAL` | finantial |
| finantial | `FINANCIAL_INSTALLMENT` | `FINANCIAL_INSTALLMENT_SCHOOL` | school |

Cada fila tem uma DLQ `<NOME>_DLQ` que recebe a mensagem depois de 5 entregas
sem sucesso.

> Os nomes acima são literais no código Go (`messaging.NewProducer("SCHOOL_COURSE")`,
> `QueueName() == "SCHOOL_COURSE_FINANCIAL"`). O SDK monta o ARN do tópico como
> `arn:<partition>:sns:<region>:<account>:<nome>` e resolve a fila por
> `GetQueueUrl(nome)` — por isso **não dá para prefixar nem sufixar** esses
> nomes sem alterar o código dos serviços. Como consequência, uma conta AWS
> comporta só um ambiente deste stack por região.

---

## Pré-requisitos

- Terraform >= 1.6
- AWS CLI v2 autenticada (`aws sts get-caller-identity` precisa responder)
- Docker (o build das imagens acontece dentro do container, sem Go local)
- `jq` e `curl` para o smoke test
- Um PostgreSQL alcançável a partir da VPC, com os dois databases já criados
  (`school_module` e `finantial_module`) — as tabelas as migrations criam sozinhas
- Um Redis alcançável a partir da VPC (usado só pelo `school-module`)

---

## Uso

```bash
cd deploy/aws/terraform
cp terraform.tfvars.example terraform.tfvars
$EDITOR terraform.tfvars          # preencha db_host, db_user, db_password, cache_uri

cd ../scripts
./deploy.sh
```

O `deploy.sh` roda em três fases, nessa ordem porque o ECS não sobe uma task
cuja imagem ainda não existe no registry:

1. **bootstrap** — cria os repositórios ECR
2. **build** — compila e publica as imagens dos dois módulos
3. **apply** — cria o resto da infraestrutura

Opções:

```bash
./deploy.sh --plan            # só mostra o plano, não aplica
./deploy.sh --skip-build      # reaproveita as imagens já publicadas
./deploy.sh --auto-approve    # não pergunta antes de aplicar
./deploy.sh --tag v1.2.3      # fixa uma tag em vez de "latest"
```

Depois do deploy:

```bash
./status.sh        # tasks rodando, saúde no ALB e health check HTTP
./smoke-test.sh    # cria curso + aluno + matrícula e confere se a mensagem chegou
```

Os serviços levam de 2 a 4 minutos para ficarem saudáveis no ALB: a task precisa
subir, aplicar as migrations e passar por dois health checks seguidos.

---

## Limpando tudo

Um comando derruba o ambiente inteiro:

```bash
cd deploy/aws/scripts && ./destroy.sh
```

Ele zera os serviços ECS antes do `terraform destroy` — sem isso o Terraform
fica esperando as tasks drenarem enquanto já tenta remover o target group — e no
fim **audita a conta** procurando o que sobrou, marcando com `$` o que continua
cobrando.

Variações:

```bash
./destroy.sh --check          # só audita, não remove nada
./destroy.sh --auto-approve   # sem confirmação
./destroy.sh --sweep          # remove por nome, sem depender do state
```

`--check` é o comando para rodar quando bater a dúvida de "será que deixei algo
ligado?". Ele lista NAT Gateways, Elastic IPs, ALBs, tasks Fargate rodando,
objetos no S3, imagens no ECR, log groups, tópicos, filas e a VPC — e diz quais
desses geram cobrança.

`--sweep` existe para o caso de o `terraform.tfstate` ter se perdido: ele apaga
por nome os recursos que dá para apagar isoladamente (serviços ECS, bucket com
todas as versões, repositórios ECR, tópicos, filas, log groups). ALB, NAT
Gateway, EIPs e VPC ficam de fora de propósito — a ordem de dependência entre
eles é traiçoeira e o Terraform resolve isso melhor; `--check` diz quais
sobraram para você remover pelo console.

O que **nunca** é tocado: o PostgreSQL e o Redis remotos. As tabelas criadas
pelas migrations continuam lá.

> Os itens que mais doem se ficarem esquecidos são o NAT Gateway (~US$ 32/mês
> por AZ, mesmo sem tráfego), o ALB (~US$ 17/mês) e as tasks Fargate, que cobram
> por segundo. A auditoria checa os três primeiro.

---

## Endpoints

O `terraform output` imprime as URLs. As rotas de negócio ficam todas sob
`/public/` — é o prefixo que os controllers declaram com `restserver.PublicApi`.
Só `/health` e `/api-docs` respondem na raiz.

```
http://<alb>:8080/health                     school-module
http://<alb>:8080/api-docs                   contrato OpenAPI
http://<alb>:8080/public/v1/courses
http://<alb>:8080/public/v1/students
http://<alb>:8080/public/v1/enrollments

http://<alb>:8081/health                     finantial-module
http://<alb>:8081/api-docs
http://<alb>:8081/public/accounts
http://<alb>:8081/public/invoices
```

A separação é por porta e não por path porque o ALB não reescreve o caminho da
requisição: um roteamento por prefixo entregaria `/school/public/v1/courses` ao
serviço, que devolveria 404.

O `school-module` chama o `finantial-module` pelo DNS interno do Cloud Map
(`http://finantial-module.colibri-examples.internal:8081`), sem sair pelo ALB.

---

## Observabilidade

O `colibri-sdk-go` já exporta traces e métricas em OTLP/HTTP. Cada task tem um
sidecar do AWS Distro for OpenTelemetry escutando em `localhost:4318`, que
traduz e publica:

| Sinal | Destino | Onde ver |
|---|---|---|
| Traces | AWS X-Ray (exporter `awsxray`) | CloudWatch → X-Ray traces / mapa de serviços |
| Métricas | CloudWatch EMF (exporter `awsemf`), namespace `Colibri/<projeto>` | CloudWatch → Metrics |
| Logs | CloudWatch Logs (log driver `awslogs`) | `/ecs/<projeto>/school-module` e `/finantial-module` |
| Infra ECS | Container Insights | CloudWatch → Container Insights |

Nada muda no código dos serviços — só a variável `OTEL_EXPORTER_OTLP_ENDPOINT`.

A imagem do collector precisa ser **>= 0.34.0**: é a versão a partir da qual o
exporter `awsxray` aceita trace IDs no formato W3C. O SDK usa o gerador de IDs
padrão do OpenTelemetry, então com um collector anterior a essa versão o X-Ray
rejeitaria os spans com `invalid xray traceid`.

Para ver um trace atravessando os dois serviços, chame
`GET /public/v1/courses/{id}` no school-module: esse endpoint faz uma chamada
HTTP para o finantial-module, e o span dos dois aparece no mesmo trace.

---

## Variáveis principais

### Banco remoto (obrigatórias)

| Variável | Descrição |
|---|---|
| `db_host` | Host do PostgreSQL. Precisa aceitar conexão vinda da VPC |
| `db_user` | Usuário compartilhado pelos dois módulos |
| `db_password` | Senha. Vai para o SSM como SecureString e é injetada como secret na task |
| `db_port` | Padrão `5432` |
| `db_ssl_mode` | Padrão `require` |
| `db_name_school` | Padrão `school_module` |
| `db_name_finantial` | Padrão `finantial_module` |
| `db_run_migrations` | Padrão `true` — cada serviço aplica suas migrations no start |

A senha nunca aparece em texto plano na task definition: o ECS resolve o
parâmetro na hora de subir o container.

### Cache remoto (obrigatória)

| Variável | Descrição |
|---|---|
| `cache_uri` | `host:porta` do Redis. Usado só pelo `school-module` |
| `cache_password` | Vazio se não houver AUTH; quando preenchida vira SecureString |

### Rede e custo

| Variável | Padrão | Efeito |
|---|---|---|
| `enable_nat_gateway` | `false` | `true` move as tasks para subnets privadas atrás de NAT (~US$ 32/mês por AZ). Com `false` elas ficam em subnet pública com IP público, e o Security Group só aceita entrada vinda do ALB |
| `allowed_ingress_cidrs` | `["0.0.0.0/0"]` | Restrinja ao seu IP se não quiser o ambiente aberto |
| `availability_zones_count` | `2` | Mínimo exigido pelo ALB |

### Aplicação

| Variável | Padrão |
|---|---|
| `environment` | `sandbox` |
| `task_cpu` / `task_memory` | `512` / `1024` |
| `school_desired_count` / `finantial_desired_count` | `1` |
| `adot_image` | `public.ecr.aws/aws-observability/aws-otel-collector:v0.43.0` |

> `environment` precisa ser `production` ou `sandbox`. O SDK só usa a credential
> chain nativa da AWS — ou seja, a role da task — quando
> `config.IsCloudEnvironment()` é verdadeiro. Com `development` ou `test` ele
> tenta falar com um LocalStack no endereço de `CLOUD_HOST`.

---

## Detalhes que o ambiente depende

Três variáveis parecem decorativas e não são:

- **`AWS_REGION` / `AWS_DEFAULT_REGION`** — em ambiente de nuvem o SDK chama
  `LoadDefaultConfig` sem passar região (`CLOUD_REGION` só vale fora da nuvem) e
  monta o ARN do tópico SNS com `GetAwsConfig().Region`. O ECS não injeta essa
  variável sozinho; sem ela o publish sai com a região vazia no ARN.
- **`CLOUD_AWS_ROLE_ARN`** — o SDK faz parse dela para descobrir a partition e o
  account id com que monta o ARN do tópico. Aqui ela recebe o ARN da role da
  task. Vazia, o publish sai com o account id em branco.
- **`CLOUD_HOST`** — fica deliberadamente **de fora**. Fora de um ambiente local
  o SDK usa o endpoint padrão da AWS; apontar `CLOUD_HOST` faria o cliente
  tentar falar com um emulador que não existe.

E uma configuração da mensageria: as subscriptions SNS→SQS ficam com
`raw_message_delivery = false`. O consumidor do SDK desserializa primeiro o
envelope de notificação do SNS e só depois o corpo da mensagem — com raw
delivery ligado, toda mensagem falharia na leitura.

---

## Limitação conhecida: o bucket do upload de documentos

O Terraform cria o bucket S3, dá permissão de leitura e escrita à role da task e
passa o nome em `STORAGE_BUCKET`. Só que hoje **o código não lê essa variável**:
`school-module/src/infra/storages/document_storage.go` traz o nome fixo no
construtor, herdado do ambiente local com LocalStack.

```go
func NewDocumentS3Storage() *DocumentS3Storage {
	return &DocumentS3Storage{
		bucket: "meu-bucket",
	}
}
```

Consequência prática: tudo funciona — API, mensageria, banco, observabilidade —
menos `POST /public/v1/students/{id}/upload-document`, que tenta escrever num
bucket chamado `meu-bucket` que não existe na conta. Renomear o bucket criado
aqui para `meu-bucket` não resolve: nomes de bucket são globais na AWS e esse já
está tomado.

O ajuste é de uma linha, no repositório da aplicação e não neste Terraform:

```go
import "os"

func NewDocumentS3Storage() *DocumentS3Storage {
	return &DocumentS3Storage{
		bucket: os.Getenv("STORAGE_BUCKET"),
	}
}
```

Com essa mudança o upload passa a usar o bucket provisionado, e o ambiente local
continua funcionando porque o `.env` já define `STORAGE_BUCKET=meu-bucket`.

---

## Custo aproximado (us-east-1, ambiente parado)

| Item | Estimativa mensal |
|---|---|
| 2 tasks Fargate 0.5 vCPU / 1 GB, 24/7 | ~US$ 36 |
| Application Load Balancer | ~US$ 17 |
| NAT Gateway (só com `enable_nat_gateway = true`) | ~US$ 32 por AZ |
| SNS, SQS, S3, ECR, CloudWatch, X-Ray no volume de exemplo | poucos dólares |

Sem NAT, algo em torno de **US$ 55/mês**. `./destroy.sh` zera tudo.

---

## Problemas comuns

**A task sobe e morre em seguida.** Quase sempre é o banco. Veja os logs:

```bash
aws logs tail /ecs/colibri-examples/school-module --follow
```

Se aparecer erro de conexão, confira se o `db_host` aceita a origem: com
`enable_nat_gateway = false` a saída é o IP público da task, que muda a cada
deploy; com `true` é o IP elástico do NAT, que é fixo e pode entrar numa
allowlist.

**Mensagem publicada mas não consumida.** Olhe a DLQ correspondente
(`<FILA>_DLQ`) e os logs do consumidor. Um erro de decode indica
`raw_message_delivery` ligado por engano.

**Nenhum trace no X-Ray.** Confira o log group `/ecs/<projeto>/adot-collector`.
`invalid xray traceid` significa collector abaixo da versão 0.34.0.

**Abrir um shell dentro da task** (`enable_execute_command` já vem ligado):

```bash
aws ecs execute-command --cluster colibri-examples \
  --task <task-id> --container school-module \
  --interactive --command "/bin/sh"
```

---

## Estrutura

```
deploy/
├── docker/
│   └── Dockerfile            build self-contained, compartilhado com o deploy da GCP
└── aws/
    ├── README.md
    ├── scripts/
    │   ├── common.sh         funções compartilhadas
    │   ├── deploy.sh         ciclo completo em 3 fases
    │   ├── build-and-push.sh imagens → ECR
    │   ├── status.sh         estado das tasks e saúde no ALB
    │   ├── smoke-test.sh     valida API + mensageria + banco
    │   └── destroy.sh        remove o ambiente e audita o que sobrou
    └── terraform/
        ├── versions.tf           provider e versões
        ├── variables.tf          entradas
        ├── locals.tf             nomes de tópicos/filas e derivações
        ├── network.tf            VPC, subnets, NAT, Security Groups
        ├── messaging.tf          SNS, SQS, DLQs, subscriptions
        ├── storage.tf            bucket S3
        ├── parameters.tf         SecureStrings do banco e do cache
        ├── iam.tf                roles de execução e da task
        ├── ecr.tf                repositórios de imagem
        ├── observability.tf      log groups e config do ADOT
        ├── alb.tf                load balancer, target groups, listeners
        ├── ecs.tf                cluster, task definitions, serviços, service discovery
        └── outputs.tf            URLs, ARNs e comandos de teste
```
