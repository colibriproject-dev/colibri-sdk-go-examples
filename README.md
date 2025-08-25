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
│   ├── prometheus/       # Configuração do Prometheus
│   └── postgres/         # Scripts de inicialização do PostgreSQL
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
