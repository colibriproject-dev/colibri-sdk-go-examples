# ─────────────────────────────────────────────────────────────────────────────
# Identificação
# ─────────────────────────────────────────────────────────────────────────────

variable "project_id" {
  description = "ID do projeto GCP onde o ambiente é criado."
  type        = string
}

variable "region" {
  description = "Região do Cloud Run, do Artifact Registry e do bucket."
  type        = string
  default     = "southamerica-east1"
}

variable "name_prefix" {
  description = "Prefixo aplicado aos recursos que aceitam nome livre (bucket, repositório, service accounts)."
  type        = string
  default     = "colibri-examples"
}

variable "environment" {
  description = <<-EOT
    Valor da variável ENVIRONMENT dos serviços. Precisa ser "production" ou
    "sandbox".

    No caminho GCP o SDK usa Application Default Credentials independentemente
    desse valor — quem depende dele é o AWS. Mas ENVIRONMENT também controla o
    formato do log (config.IsLocalEnvironment() liga o output legível para
    humanos), e "development" deixaria o log sem estrutura no Cloud Logging.
  EOT
  type        = string
  default     = "sandbox"

  validation {
    condition     = contains(["production", "sandbox"], var.environment)
    error_message = "environment precisa ser production ou sandbox."
  }
}

variable "enable_apis" {
  description = "Habilita as APIs do projeto pelo Terraform. Deixe false se sua organização gerencia isso à parte."
  type        = bool
  default     = true
}

# ─────────────────────────────────────────────────────────────────────────────
# Banco de dados remoto (NÃO é criado por este Terraform)
# ─────────────────────────────────────────────────────────────────────────────

variable "db_host" {
  description = "Host do PostgreSQL remoto. Precisa aceitar conexão vinda da internet — o Cloud Run sai por IPs do pool do Google, que mudam."
  type        = string
}

variable "db_port" {
  description = "Porta do PostgreSQL remoto."
  type        = string
  default     = "5432"
}

variable "db_user" {
  description = "Usuário do PostgreSQL remoto (compartilhado pelos dois módulos)."
  type        = string
}

variable "db_password" {
  description = "Senha do PostgreSQL remoto. Vai para o Secret Manager e é injetada no container como secret."
  type        = string
  sensitive   = true
}

variable "db_ssl_mode" {
  description = "sslmode da string de conexão. Como a saída é pela internet aberta, require é o mínimo aceitável."
  type        = string
  default     = "require"
}

variable "db_name_school" {
  description = "Nome do database usado pelo school-module."
  type        = string
  default     = "school_module"
}

variable "db_name_finantial" {
  description = "Nome do database usado pelo finantial-module."
  type        = string
  default     = "finantial_module"
}

variable "db_run_migrations" {
  description = "Valor de SQL_DB_MIGRATION. Com true cada serviço aplica suas migrations no start."
  type        = bool
  default     = true
}

# ─────────────────────────────────────────────────────────────────────────────
# Cache remoto (NÃO é criado por este Terraform) — usado só pelo school-module
# ─────────────────────────────────────────────────────────────────────────────

variable "cache_uri" {
  description = "Endereço host:porta do Redis remoto usado pelo school-module."
  type        = string
}

variable "cache_password" {
  description = "Senha do Redis remoto. Vazia se não houver AUTH; quando preenchida vira um secret."
  type        = string
  sensitive   = true
  default     = ""
}

# ─────────────────────────────────────────────────────────────────────────────
# Cloud Run
# ─────────────────────────────────────────────────────────────────────────────

variable "school_image_tag" {
  description = "Tag da imagem do school-module no Artifact Registry."
  type        = string
  default     = "latest"
}

variable "finantial_image_tag" {
  description = "Tag da imagem do finantial-module no Artifact Registry."
  type        = string
  default     = "latest"
}

variable "min_instance_count" {
  description = <<-EOT
    Instâncias mínimas por serviço.

    Precisa ser >= 1 para o ambiente funcionar de verdade: os consumidores do
    SDK fazem pull do Pub/Sub numa goroutine de fundo, e sem instância viva não
    há quem consuma as mensagens. É também o que domina o custo — veja o README.

    Com 0 as APIs continuam respondendo (o Cloud Run sobe sob demanda), mas as
    mensagens só são consumidas enquanto houver requisição em andamento.
  EOT
  type        = number
  default     = 1
}

variable "max_instance_count" {
  description = "Instâncias máximas por serviço."
  type        = number
  default     = 3
}

variable "app_cpu" {
  description = <<-EOT
    vCPU do container da aplicação.

    Não use valor fracionário: o Cloud Run só aceita CPU abaixo de 1 com
    billing por requisição, que é incompatível com o CPU sempre alocado que os
    consumidores do Pub/Sub exigem.
  EOT
  type        = string
  default     = "1"
}

variable "app_memory" {
  description = "Memória do container da aplicação."
  type        = string
  default     = "512Mi"
}

variable "otel_memory" {
  description = "Memória do sidecar do collector. Ele não declara CPU própria: divide a da instância com a aplicação."
  type        = string
  default     = "256Mi"
}

variable "request_timeout_seconds" {
  description = "Timeout de requisição do Cloud Run."
  type        = number
  default     = 300
}

variable "max_concurrency" {
  description = "Requisições simultâneas por instância."
  type        = number
  default     = 80
}

variable "allow_public_access" {
  description = "Concede roles/run.invoker a allUsers, deixando as URLs abertas sem autenticação."
  type        = bool
  default     = true
}

# ─────────────────────────────────────────────────────────────────────────────
# Observabilidade
# ─────────────────────────────────────────────────────────────────────────────

variable "otel_image" {
  description = <<-EOT
    Imagem do collector usada como sidecar. O padrão é o build do Google para
    Cloud Run. A config usa o exporter googlecloud, que também existe no
    otel/opentelemetry-collector-contrib caso você prefira o upstream.
  EOT
  type        = string
  default     = "us-docker.pkg.dev/cloud-ops-agents-artifacts/google-cloud-opentelemetry-collector/otelcol-google:0.156.0"
}

variable "log_level" {
  description = "Valor de LOG_LEVEL das aplicações."
  type        = string
  default     = "info"
}

# ─────────────────────────────────────────────────────────────────────────────
# Mensageria
# ─────────────────────────────────────────────────────────────────────────────

variable "ack_deadline_seconds" {
  description = "Prazo de ack das subscriptions do Pub/Sub."
  type        = number
  default     = 30
}

variable "message_retention_duration" {
  description = "Retenção das mensagens não confirmadas."
  type        = string
  default     = "604800s"
}

variable "max_delivery_attempts" {
  description = "Entregas antes da mensagem ir para o tópico de dead letter."
  type        = number
  default     = 5
}
