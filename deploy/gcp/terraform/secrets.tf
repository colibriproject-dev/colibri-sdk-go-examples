# ─────────────────────────────────────────────────────────────────────────────
# Credenciais do banco e do cache remotos.
#
# Nenhum Cloud SQL ou Memorystore é criado aqui: o Terraform só guarda as
# credenciais no Secret Manager e as injeta no container por referência, de modo
# que a senha nunca aparece em texto plano na revisão do Cloud Run.
# ─────────────────────────────────────────────────────────────────────────────

resource "google_secret_manager_secret" "db_password" {
  secret_id = "${local.name}-db-password"

  replication {
    auto {}
  }

  depends_on = [google_project_service.required]
}

resource "google_secret_manager_secret_version" "db_password" {
  secret      = google_secret_manager_secret.db_password.id
  secret_data = var.db_password
}

resource "google_secret_manager_secret" "cache_password" {
  count = var.cache_password == "" ? 0 : 1

  secret_id = "${local.name}-cache-password"

  replication {
    auto {}
  }

  depends_on = [google_project_service.required]
}

resource "google_secret_manager_secret_version" "cache_password" {
  count = var.cache_password == "" ? 0 : 1

  secret      = google_secret_manager_secret.cache_password[0].id
  secret_data = var.cache_password
}

# A configuração do collector também vive no Secret Manager: é assim que o
# sidecar do Google espera receber o arquivo, montado como volume.
resource "google_secret_manager_secret" "otel_config" {
  secret_id = "${local.name}-otel-config"

  replication {
    auto {}
  }

  depends_on = [google_project_service.required]
}

resource "google_secret_manager_secret_version" "otel_config" {
  secret      = google_secret_manager_secret.otel_config.id
  secret_data = local.otel_config
}
