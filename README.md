# Colibri SDK Go - Exemplos

Este repositório contém projetos de exemplo demonstrando o uso do [Colibri SDK Go](https://github.com/colibriproject-dev/colibri-sdk-go), um framework para desenvolvimento de microserviços em Go com observabilidade integrada.

## 📋 Resumo dos Projetos

### School Module (`school-module`)
- **Domínio**: Sistema de gerenciamento escolar
- **Funcionalidades**: Gestão de cursos, estudantes e operações educacionais
- **Porta**: 8080
- **Banco de Dados**: PostgreSQL (`school_module`)

### Financial Module (`finantial-module`)
- **Domínio**: Sistema financeiro
- **Funcionalidades**: Operações financeiras e transações
- **Porta**: 8081
- **Banco de Dados**: PostgreSQL (`finantial_module`)

Ambos os módulos utilizam:
- ✅ OpenTelemetry para observabilidade (traces)
- ✅ PostgreSQL para persistência de dados
- ✅ Redis para cache
- ✅ LocalStack para simulação de serviços AWS
- ✅ Swagger para documentação de API

## 🛠️ Pré-requisitos

- [Docker](https://www.docker.com/) e Docker Compose
- [Make](https://www.gnu.org/software/make/) (opcional, mas recomendado)

## 🚀 Como Iniciar

### 1. Construir as Imagens Docker
```bash
make build
```

### 2. Iniciar os Serviços
```bash
make start
```

### 3. Parar os Serviços
```bash
make stop
```

### 4. Limpar o Ambiente (remove containers e imagens)
```bash
make clean
```

### Comandos Adicionais
```bash
# Ver logs dos serviços
make logs

# Ver estatísticas dos containers
make stats
```

## 🔍 Acessando os Serviços

### APIs dos Módulos
- **School Module**: http://localhost:8080
  - Swagger UI: http://localhost:8080/api-docs
- **Financial Module**: http://localhost:8081
  - Swagger UI: http://localhost:8081/api-docs

### Ferramentas de Observabilidade

#### 🔍 Jaeger (Distributed Tracing)
- **URL**: http://localhost:16686
- **Descrição**: Interface para visualizar traces distribuídos dos microserviços
- **Como usar**:
  1. Acesse a URL do Jaeger
  2. Selecione o serviço (`school-module` ou `finantial-module`)
  3. Clique em "Find Traces" para visualizar os traces
  4. Clique em um trace específico para ver detalhes da requisição

#### 📊 Prometheus (Métricas)
- **URL**: http://localhost:9090
- **Descrição**: Sistema de monitoramento e coleta de métricas
- **Jobs**: `school-module` e `finantial-module` raspam o `/metrics` de cada módulo, com os
  nomes usados pelos dashboards do SDK; `otel-collector` traz apenas as span metrics geradas
  pelo collector (as demais métricas dele duplicariam as dos módulos).
- **Alertas**: http://localhost:9090/alerts, com as regras do colibri-sdk-go.

#### 📈 Grafana (Dashboards)
- **URL**: http://localhost:3001 (acesso anônimo como admin)
- **Dashboards**: pasta **Colibri**, com HTTP, Messaging, Data stores e Go runtime.

Os dashboards e as regras de alerta são os do
[`colibri-sdk-go/observability`](https://github.com/colibriproject-dev/colibri-sdk-go/tree/master/observability),
montados direto do clone do SDK. Por isso o `colibri-sdk-go` precisa estar ao lado deste
repositório (`../colibri-sdk-go`), que é o mesmo caminho do `replace` nos `go.mod`. Para usar
outro caminho:

```bash
COLIBRI_SDK_PATH=/caminho/para/colibri-sdk-go make start
```

Mudanças nos JSON dos dashboards aparecem no Grafana em até 10 segundos. Mudanças nas
regras exigem recarregar o Prometheus: `curl -X POST localhost:9090/-/reload`.

#### 🚦 Gerando tráfego
Com a stack no ar, gere tráfego para popular os dashboards:

```bash
make traffic                 # roda até Ctrl+C
make traffic DURATION=300    # roda por 5 minutos
```

O script [`scripts/traffic.py`](scripts/traffic.py) usa só a biblioteca padrão do Python 3.
Ele cria e consulta cursos, alunos e matrículas, publicando eventos que o finantial-module
consome. Também envia documentos ao storage e provoca erros 404 e 400 de propósito. No
fim, imprime um resumo por rota e status. Veja `python3 scripts/traffic.py --help`.

#### ☁️ LocalStack Web UI
- **URL**: http://localhost:3000
- **Descrição**: Interface web para gerenciar serviços AWS simulados

### Bancos de Dados
- **PostgreSQL**: `localhost:5432`
  - Usuário: `postgres`
  - Senha: `postgres`
  - Databases: `school_module`, `finantial_module`

- **Redis**: `localhost:6379`

## 🔧 Configuração de Desenvolvimento

### Variáveis de Ambiente
Os módulos utilizam variáveis de ambiente para configuração. Principais configurações:

```bash
# OpenTelemetry
OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4318
OTEL_EXPORTER_OTLP_PROTOCOL=http
OTEL_SERVICE_NAME=<nome-do-servico>

# Banco de Dados
SQL_DB_HOST=postgres
SQL_DB_USER=postgres
SQL_DB_PASSWORD=postgres

# Cache
CACHE_URI=redis:6379

# Cloud Services (LocalStack)
CLOUD_HOST=http://localstack:4566
```

## 📝 Estrutura do Projeto

```
.
├── school-module/          # Módulo de sistema escolar
│   ├── src/               # Código fonte
│   ├── migrations/        # Migrações de banco
│   └── tests/            # Testes automatizados
├── finantial-module/      # Módulo financeiro
│   ├── src/              # Código fonte
│   └── migrations/       # Migrações de banco
├── dev/                  # Configurações de desenvolvimento
│   ├── otel-collector/   # Configuração do OpenTelemetry Collector
│   ├── grafana/          # Provisionamento do Grafana (datasource e dashboards)
│   ├── prometheus/       # Configuração do Prometheus
│   └── postgres/         # Scripts de inicialização do PostgreSQL
├── scripts/
│   └── traffic.py        # Gerador de tráfego para os dashboards
├── docker-compose.yaml   # Orquestração dos serviços
├── Dockerfile           # Imagem Docker multi-stage
└── Makefile            # Comandos de automação
```

## 🎯 Casos de Uso para Observabilidade

### Testando Traces com Jaeger
1. Faça uma requisição para qualquer endpoint dos módulos
2. Acesse o Jaeger em http://localhost:16686
3. Selecione o serviço desejado
4. Visualize o trace da requisição com todos os spans
5. Analise a desempenho e possíveis gargalos
