# ─────────────────────────────────────────────────────────────────────────────
# Observabilidade nativa da GCP.
#
# O colibri-sdk-go já exporta traces e métricas em OTLP/HTTP. Cada serviço tem
# um sidecar do collector na mesma instância, que traduz e publica:
#   traces   → Cloud Trace      (exporter googlecloud)
#   métricas → Cloud Monitoring (exporter googlecloud)
#   logs     → Cloud Logging    (stdout do container, coletado pelo Cloud Run)
#
# Nada muda no código dos serviços: eles só apontam OTEL_EXPORTER_OTLP_ENDPOINT
# para localhost:4318, que é o sidecar dentro da mesma instância.
#
# Diferente do X-Ray, o Cloud Trace aceita trace ID no formato W3C sem
# conversão, então não há restrição de versão do collector como havia com o
# ADOT no deploy AWS.
# ─────────────────────────────────────────────────────────────────────────────

locals {
  otel_config = yamlencode({
    receivers = {
      otlp = {
        protocols = {
          http = { endpoint = "0.0.0.0:4318" }
          grpc = { endpoint = "0.0.0.0:4317" }
        }
      }
    }

    processors = {
      # Preenche os atributos de recurso com os metadados do Cloud Run
      # (serviço, revisão, região), que é o que faz o serviço aparecer
      # identificado no Cloud Trace.
      resourcedetection = {
        detectors = ["env", "gcp"]
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
      googlecloud = {
        project = var.project_id
        log     = { default_log_name = "colibri-otel" }
      }
    }

    extensions = {
      # O Cloud Run exige que o sidecar responda a um health check para
      # considerar a instância pronta.
      health_check = { endpoint = "0.0.0.0:13133" }
    }

    service = {
      extensions = ["health_check"]
      pipelines = {
        traces = {
          receivers  = ["otlp"]
          processors = ["memory_limiter", "resourcedetection", "batch"]
          exporters  = ["googlecloud"]
        }
        metrics = {
          receivers  = ["otlp"]
          processors = ["memory_limiter", "resourcedetection", "batch"]
          exporters  = ["googlecloud"]
        }
      }
      telemetry = {
        logs = { level = "info" }
      }
    }
  })
}
