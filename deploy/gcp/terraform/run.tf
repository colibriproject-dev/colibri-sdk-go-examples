# ─────────────────────────────────────────────────────────────────────────────
# Blocos reaproveitados pelos dois serviços
# ─────────────────────────────────────────────────────────────────────────────

locals {
  # Variáveis comuns aos dois módulos.
  #
  # PORT fica de fora de propósito: o Cloud Run injeta essa variável com o valor
  # de container_port e RECUSA a revisão se ela for declarada à mão. O SDK lê
  # PORT do ambiente, então a injeção do Cloud Run já resolve.
  common_environment = {
    APP_TYPE    = "service"
    ENVIRONMENT = var.environment
    CLOUD       = "gcp"
    LOG_LEVEL   = var.log_level

    # É esta a variável que o cliente Pub/Sub do SDK lê — não CLOUD_PROJECT_ID,
    # que aparece no .env local mas não é consumida em lugar nenhum do código.
    PUBSUB_PROJECT_ID = var.project_id

    # Usadas pelo Google Cloud SDK dentro do container.
    GOOGLE_CLOUD_PROJECT = var.project_id

    SQL_DB_HOST      = var.db_host
    SQL_DB_PORT      = var.db_port
    SQL_DB_USER      = var.db_user
    SQL_DB_SSL_MODE  = var.db_ssl_mode
    SQL_DB_MIGRATION = tostring(var.db_run_migrations)

    OTEL_EXPORTER_OTLP_ENDPOINT = local.otlp_endpoint
    OTEL_EXPORTER_OTLP_PROTOCOL = "http"
  }
}

# ─────────────────────────────────────────────────────────────────────────────
# finantial-module
#
# Criado primeiro porque o school precisa da URL dele em
# FINANCIAL_MODULE_BASE_URL.
# ─────────────────────────────────────────────────────────────────────────────

resource "google_cloud_run_v2_service" "finantial" {
  name     = "finantial-module"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  # Sem isto o terraform destroy falha: o provider protege o serviço por padrão.
  deletion_protection = false

  template {
    service_account                  = google_service_account.finantial.email
    timeout                          = "${var.request_timeout_seconds}s"
    max_instance_request_concurrency = var.max_concurrency
    execution_environment            = "EXECUTION_ENVIRONMENT_GEN2"

    scaling {
      min_instance_count = var.min_instance_count
      max_instance_count = var.max_instance_count
    }

    volumes {
      name = "otel-config"

      secret {
        secret = google_secret_manager_secret.otel_config.secret_id

        items {
          version = "latest"
          path    = "config.yaml"
        }
      }
    }

    containers {
      name       = "app"
      image      = "${local.repository_url}/finantial-module:${var.finantial_image_tag}"
      depends_on = ["collector"]

      ports {
        container_port = local.container_port
      }

      resources {
        limits = {
          cpu    = var.app_cpu
          memory = var.app_memory
        }

        # CPU sempre alocada. Não é preferência: os consumidores do SDK fazem
        # pull do Pub/Sub numa goroutine de fundo, e o throttling padrão do
        # Cloud Run congela o container entre requisições — o loop pararia de
        # buscar mensagens assim que o tráfego HTTP cessasse.
        cpu_idle          = false
        startup_cpu_boost = true
      }

      dynamic "env" {
        for_each = merge(local.common_environment, {
          APP_NAME          = "Colibri SDK Go Test - Finantial Module"
          SQL_DB_NAME       = var.db_name_finantial
          OTEL_SERVICE_NAME = "finantial-module"
        })

        content {
          name  = env.key
          value = env.value
        }
      }

      env {
        name = "SQL_DB_PASSWORD"

        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.db_password.secret_id
            version = "latest"
          }
        }
      }

      # As migrations rodam no start contra o banco remoto; o failure_threshold
      # alto dá até 5 minutos antes do Cloud Run desistir da revisão.
      startup_probe {
        initial_delay_seconds = 10
        period_seconds        = 10
        timeout_seconds       = 5
        failure_threshold     = 30

        http_get {
          path = "/health"
          port = local.container_port
        }
      }

      liveness_probe {
        period_seconds    = 30
        timeout_seconds   = 5
        failure_threshold = 3

        http_get {
          path = "/health"
          port = local.container_port
        }
      }
    }

    containers {
      name  = "collector"
      image = var.otel_image
      args  = ["--config=/etc/otelcol-google/config.yaml"]

      # Só memória: a CPU da instância é definida pelo container da aplicação e
      # compartilhada com o sidecar. Declarar CPU aqui somaria ao total da
      # instância e dobraria o custo.
      resources {
        limits = {
          memory = var.otel_memory
        }
        cpu_idle = false
      }

      volume_mounts {
        name       = "otel-config"
        mount_path = "/etc/otelcol-google"
      }

      startup_probe {
        period_seconds    = 5
        timeout_seconds   = 3
        failure_threshold = 20

        tcp_socket {
          port = 13133
        }
      }
    }
  }

  depends_on = [
    google_project_service.required,
    google_secret_manager_secret_iam_member.db_password,
    google_secret_manager_secret_iam_member.otel_config,
    google_pubsub_subscription_iam_member.subscribers,
    google_pubsub_topic_iam_member.publishers,
  ]
}

# ─────────────────────────────────────────────────────────────────────────────
# school-module
# ─────────────────────────────────────────────────────────────────────────────

resource "google_cloud_run_v2_service" "school" {
  name     = "school-module"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  deletion_protection = false

  template {
    service_account                  = google_service_account.school.email
    timeout                          = "${var.request_timeout_seconds}s"
    max_instance_request_concurrency = var.max_concurrency
    execution_environment            = "EXECUTION_ENVIRONMENT_GEN2"

    scaling {
      min_instance_count = var.min_instance_count
      max_instance_count = var.max_instance_count
    }

    volumes {
      name = "otel-config"

      secret {
        secret = google_secret_manager_secret.otel_config.secret_id

        items {
          version = "latest"
          path    = "config.yaml"
        }
      }
    }

    containers {
      name       = "app"
      image      = "${local.repository_url}/school-module:${var.school_image_tag}"
      depends_on = ["collector"]

      ports {
        container_port = local.container_port
      }

      resources {
        limits = {
          cpu    = var.app_cpu
          memory = var.app_memory
        }

        cpu_idle          = false
        startup_cpu_boost = true
      }

      dynamic "env" {
        for_each = merge(local.common_environment, {
          APP_NAME          = "Colibri SDK Go Test - School Module"
          SQL_DB_NAME       = var.db_name_school
          CACHE_URI         = var.cache_uri
          STORAGE_BUCKET    = google_storage_bucket.documents.name
          OTEL_SERVICE_NAME = "school-module"

          # O school chama o finantial pela URL pública do Cloud Run. Com os
          # dois serviços abertos não há token a injetar; se o finantial passar
          # a exigir IAM, o restclient do SDK precisaria enviar um identity
          # token, o que hoje ele não faz.
          FINANCIAL_MODULE_BASE_URL = google_cloud_run_v2_service.finantial.uri
        })

        content {
          name  = env.key
          value = env.value
        }
      }

      env {
        name = "SQL_DB_PASSWORD"

        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.db_password.secret_id
            version = "latest"
          }
        }
      }

      dynamic "env" {
        for_each = var.cache_password == "" ? [] : [1]

        content {
          name = "CACHE_PASSWORD"

          value_source {
            secret_key_ref {
              secret  = google_secret_manager_secret.cache_password[0].secret_id
              version = "latest"
            }
          }
        }
      }

      startup_probe {
        initial_delay_seconds = 10
        period_seconds        = 10
        timeout_seconds       = 5
        failure_threshold     = 30

        http_get {
          path = "/health"
          port = local.container_port
        }
      }

      liveness_probe {
        period_seconds    = 30
        timeout_seconds   = 5
        failure_threshold = 3

        http_get {
          path = "/health"
          port = local.container_port
        }
      }
    }

    containers {
      name  = "collector"
      image = var.otel_image
      args  = ["--config=/etc/otelcol-google/config.yaml"]

      resources {
        limits = {
          memory = var.otel_memory
        }
        cpu_idle = false
      }

      volume_mounts {
        name       = "otel-config"
        mount_path = "/etc/otelcol-google"
      }

      startup_probe {
        period_seconds    = 5
        timeout_seconds   = 3
        failure_threshold = 20

        tcp_socket {
          port = 13133
        }
      }
    }
  }

  depends_on = [
    google_project_service.required,
    google_secret_manager_secret_iam_member.db_password,
    google_secret_manager_secret_iam_member.otel_config,
    google_storage_bucket_iam_member.school_documents,
    google_pubsub_subscription_iam_member.subscribers,
    google_pubsub_topic_iam_member.publishers,
  ]
}

# ─────────────────────────────────────────────────────────────────────────────
# Acesso público
#
# Equivalente ao ALB aberto do deploy AWS. Com allow_public_access = false as
# URLs passam a exigir um identity token do Google em toda chamada — e o
# school-module deixa de conseguir falar com o finantial, porque o restclient do
# SDK não injeta esse token.
# ─────────────────────────────────────────────────────────────────────────────

resource "google_cloud_run_v2_service_iam_member" "school_public" {
  count = var.allow_public_access ? 1 : 0

  location = google_cloud_run_v2_service.school.location
  name     = google_cloud_run_v2_service.school.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

resource "google_cloud_run_v2_service_iam_member" "finantial_public" {
  count = var.allow_public_access ? 1 : 0

  location = google_cloud_run_v2_service.finantial.location
  name     = google_cloud_run_v2_service.finantial.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
