# ─────────────────────────────────────────────────────────────────────────────
# Observabilidade nativa da AWS.
#
# O colibri-sdk-go já exporta traces e métricas em OTLP/HTTP. Aqui esse OTLP é
# recebido por um sidecar do AWS Distro for OpenTelemetry em cada task, que
# traduz e publica:
#   traces   → AWS X-Ray        (exporter awsxray)
#   métricas → CloudWatch (EMF) (exporter awsemf)
#   logs     → CloudWatch Logs  (log driver awslogs, sem passar pelo collector)
#
# Nada muda no código dos serviços: eles só apontam OTEL_EXPORTER_OTLP_ENDPOINT
# para localhost:4318, que é o sidecar dentro da mesma task.
# ─────────────────────────────────────────────────────────────────────────────

resource "aws_cloudwatch_log_group" "school" {
  name              = "/ecs/${local.name}/school-module"
  retention_in_days = var.log_retention_days
}

resource "aws_cloudwatch_log_group" "finantial" {
  name              = "/ecs/${local.name}/finantial-module"
  retention_in_days = var.log_retention_days
}

resource "aws_cloudwatch_log_group" "adot" {
  name              = "/ecs/${local.name}/adot-collector"
  retention_in_days = var.log_retention_days
}

# Log group de destino das métricas em formato EMF.
resource "aws_cloudwatch_log_group" "metrics" {
  name              = "/aws/ecs/${local.name}/metrics"
  retention_in_days = var.log_retention_days
}

locals {
  adot_config = yamlencode({
    receivers = {
      otlp = {
        protocols = {
          http = { endpoint = "0.0.0.0:4318" }
          grpc = { endpoint = "0.0.0.0:4317" }
        }
      }
    }

    processors = {
      # Enriquece os spans com os metadados da task (cluster, task arn, AZ),
      # que é o que faz o serviço aparecer identificado no mapa do X-Ray.
      resourcedetection = {
        detectors = ["env", "ecs"]
        timeout   = "5s"
        override  = false
      }
      batch = {
        timeout             = "5s"
        send_batch_size     = 512
        send_batch_max_size = 1024
      }
      memory_limiter = {
        check_interval         = "1s"
        limit_percentage       = 70
        spike_limit_percentage = 20
      }
    }

    exporters = {
      awsxray = {
        region               = var.aws_region
        index_all_attributes = true
      }
      awsemf = {
        region         = var.aws_region
        namespace      = "Colibri/${local.name}"
        log_group_name = aws_cloudwatch_log_group.metrics.name
        resource_to_telemetry_conversion = {
          enabled = true
        }
      }
    }

    service = {
      pipelines = {
        traces = {
          receivers  = ["otlp"]
          processors = ["memory_limiter", "resourcedetection", "batch"]
          exporters  = ["awsxray"]
        }
        metrics = {
          receivers  = ["otlp"]
          processors = ["memory_limiter", "resourcedetection", "batch"]
          exporters  = ["awsemf"]
        }
      }
      telemetry = {
        logs = { level = "info" }
      }
    }
  })
}

# O sidecar lê a configuração da variável AOT_CONFIG_CONTENT, que o ECS resolve
# a partir deste parâmetro. Evita ter que construir uma imagem própria do
# collector só para embutir um arquivo de config.
resource "aws_ssm_parameter" "adot_config" {
  name        = "/${local.name}/${var.environment}/adot/config"
  description = "Configuracao do ADOT Collector (X-Ray + CloudWatch EMF)"
  type        = "String"
  value       = local.adot_config
}
