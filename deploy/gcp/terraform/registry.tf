resource "google_artifact_registry_repository" "modules" {
  location      = var.region
  repository_id = local.name
  description   = "Imagens dos modulos de exemplo do colibri-sdk-go"
  format        = "DOCKER"

  # Mantém apenas as 10 versões mais recentes de cada módulo.
  cleanup_policies {
    id     = "manter-recentes"
    action = "KEEP"

    most_recent_versions {
      keep_count = 10
    }
  }

  depends_on = [google_project_service.required]
}
