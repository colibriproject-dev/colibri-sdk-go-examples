# APIs necessárias. Habilitar pelo Terraform evita o erro de primeira execução
# em projeto novo, em que o apply falha na metade porque a API do Cloud Run
# ainda não estava ligada.
locals {
  required_apis = [
    "run.googleapis.com",
    "artifactregistry.googleapis.com",
    "pubsub.googleapis.com",
    "storage.googleapis.com",
    "secretmanager.googleapis.com",
    "cloudtrace.googleapis.com",
    "monitoring.googleapis.com",
    "logging.googleapis.com",
    "iam.googleapis.com",
    "cloudresourcemanager.googleapis.com",
  ]
}

resource "google_project_service" "required" {
  for_each = var.enable_apis ? toset(local.required_apis) : toset([])

  project = var.project_id
  service = each.value

  # Desligar a API no destroy quebraria outros recursos do projeto que dependam
  # dela — e ela não custa nada ligada.
  disable_on_destroy = false
}
