# Bucket usado pelo school-module no upload de documentos do aluno.
resource "google_storage_bucket" "documents" {
  name     = "${local.name}-documents-${var.project_id}"
  location = var.region

  # Permite o destroy remover o bucket mesmo com objetos dentro.
  force_destroy = true

  # Acesso só por IAM: nada de ACL por objeto, e nenhum objeto público.
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"

  versioning {
    enabled = true
  }

  lifecycle_rule {
    condition {
      num_newer_versions = 3
    }
    action {
      type = "Delete"
    }
  }

  depends_on = [google_project_service.required]
}
