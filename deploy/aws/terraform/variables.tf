# ─────────────────────────────────────────────────────────────────────────────
# Identificação
# ─────────────────────────────────────────────────────────────────────────────

variable "project_name" {
  description = "Prefixo aplicado aos recursos criados (VPC, ECS, ALB, IAM...)."
  type        = string
  default     = "colibri-examples"
}

variable "environment" {
  description = <<-EOT
    Valor da variável ENVIRONMENT dos serviços. Precisa ser "production" ou
    "sandbox": o SDK só usa a credential chain nativa da AWS (a role da task)
    quando config.IsCloudEnvironment() é verdadeiro. Com "development" ou
    "test" o SDK tenta falar com um LocalStack em CLOUD_HOST.
  EOT
  type        = string
  default     = "sandbox"

  validation {
    condition     = contains(["production", "sandbox"], var.environment)
    error_message = "environment precisa ser production ou sandbox para o SDK usar as credenciais da role da task."
  }
}

variable "aws_region" {
  description = "Região onde todo o ambiente é criado."
  type        = string
  default     = "us-east-1"
}

# ─────────────────────────────────────────────────────────────────────────────
# Rede
# ─────────────────────────────────────────────────────────────────────────────

variable "vpc_cidr" {
  description = "Bloco CIDR da VPC criada para o ambiente."
  type        = string
  default     = "10.42.0.0/16"
}

variable "availability_zones_count" {
  description = "Quantidade de AZs usadas (mínimo 2, exigência do ALB)."
  type        = number
  default     = 2

  validation {
    condition     = var.availability_zones_count >= 2
    error_message = "O ALB exige subnets em pelo menos duas AZs."
  }
}

variable "enable_nat_gateway" {
  description = <<-EOT
    true  → as tasks rodam em subnets privadas e saem pela internet por um NAT
            Gateway (postura mais próxima de produção, ~US$ 32/mês por NAT).
    false → as tasks rodam em subnets públicas com IP público e o Security Group
            só aceita entrada vinda do ALB (mais barato, suficiente para o exemplo).
  EOT
  type        = bool
  default     = false
}

variable "allowed_ingress_cidrs" {
  description = "CIDRs autorizados a chamar o ALB. Restrinja ao seu IP se não quiser o ambiente aberto."
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

# ─────────────────────────────────────────────────────────────────────────────
# Banco de dados remoto (NÃO é criado por este Terraform)
# ─────────────────────────────────────────────────────────────────────────────

variable "db_host" {
  description = "Host do PostgreSQL remoto. Precisa ser alcançável a partir da VPC."
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
  description = "Senha do PostgreSQL remoto. Vai para o SSM Parameter Store como SecureString e é injetada na task como secret."
  type        = string
  sensitive   = true
}

variable "db_ssl_mode" {
  description = "sslmode da string de conexão do PostgreSQL (disable, require, verify-ca, verify-full)."
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
  description = "Valor de SQL_DB_MIGRATION. Com true cada serviço aplica suas migrations no start contra o banco remoto."
  type        = bool
  default     = true
}

# ─────────────────────────────────────────────────────────────────────────────
# Cache remoto (NÃO é criado por este Terraform) — usado só pelo school-module
# ─────────────────────────────────────────────────────────────────────────────

variable "cache_uri" {
  description = "Endereço host:porta do Redis remoto usado pelo school-module (ex.: meu-redis.abc.ng.0001.use1.cache.amazonaws.com:6379)."
  type        = string
}

variable "cache_password" {
  description = "Senha do Redis remoto. Deixe vazio se não houver AUTH; quando preenchida vai para o SSM como SecureString."
  type        = string
  sensitive   = true
  default     = ""
}

# ─────────────────────────────────────────────────────────────────────────────
# Aplicações
# ─────────────────────────────────────────────────────────────────────────────

variable "school_image_tag" {
  description = "Tag da imagem do school-module no ECR."
  type        = string
  default     = "latest"
}

variable "finantial_image_tag" {
  description = "Tag da imagem do finantial-module no ECR."
  type        = string
  default     = "latest"
}

variable "school_desired_count" {
  description = "Número de tasks do school-module."
  type        = number
  default     = 1
}

variable "finantial_desired_count" {
  description = "Número de tasks do finantial-module."
  type        = number
  default     = 1
}

variable "task_cpu" {
  description = "CPU da task Fargate em unidades (256 = 0.25 vCPU). Compartilhada entre a aplicação e o sidecar ADOT."
  type        = number
  default     = 512
}

variable "task_memory" {
  description = "Memória da task Fargate em MiB."
  type        = number
  default     = 1024
}

variable "log_retention_days" {
  description = "Retenção dos log groups do CloudWatch."
  type        = number
  default     = 14
}

variable "log_level" {
  description = "Valor de LOG_LEVEL das aplicações."
  type        = string
  default     = "info"
}

# ─────────────────────────────────────────────────────────────────────────────
# Observabilidade
# ─────────────────────────────────────────────────────────────────────────────

variable "adot_image" {
  description = <<-EOT
    Imagem do AWS Distro for OpenTelemetry Collector usada como sidecar.
    Precisa ser >= 0.34.0: é a versão a partir da qual o exporter awsxray
    aceita trace IDs no formato W3C, que é o que o colibri-sdk-go gera.
  EOT
  type        = string
  default     = "public.ecr.aws/aws-observability/aws-otel-collector:v0.43.0"
}

variable "enable_container_insights" {
  description = "Habilita o CloudWatch Container Insights no cluster ECS."
  type        = bool
  default     = true
}

# ─────────────────────────────────────────────────────────────────────────────
# Mensageria
# ─────────────────────────────────────────────────────────────────────────────

variable "sqs_visibility_timeout_seconds" {
  description = "Visibility timeout das filas SQS."
  type        = number
  default     = 30
}

variable "sqs_message_retention_seconds" {
  description = "Retenção das mensagens nas filas SQS."
  type        = number
  default     = 345600
}

variable "sqs_max_receive_count" {
  description = "Número de entregas antes da mensagem ir para a DLQ."
  type        = number
  default     = 5
}
