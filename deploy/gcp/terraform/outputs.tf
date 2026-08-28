output "school_module_url" {
  description = "URL publica do school-module."
  value       = google_cloud_run_v2_service.school.uri
}

output "finantial_module_url" {
  description = "URL publica do finantial-module."
  value       = google_cloud_run_v2_service.finantial.uri
}

# Todas as rotas de negócio dos dois módulos são registradas com o prefixo
# restserver.PublicApi, ou seja, ficam sob /public/. Só /health e /api-docs
# respondem na raiz.
output "api_paths" {
  description = "Caminhos das rotas expostas por cada modulo."
  value = {
    school_courses     = "${google_cloud_run_v2_service.school.uri}/public/v1/courses"
    school_students    = "${google_cloud_run_v2_service.school.uri}/public/v1/students"
    school_enrollments = "${google_cloud_run_v2_service.school.uri}/public/v1/enrollments"
    school_swagger     = "${google_cloud_run_v2_service.school.uri}/api-docs"
    finantial_accounts = "${google_cloud_run_v2_service.finantial.uri}/public/accounts"
    finantial_invoices = "${google_cloud_run_v2_service.finantial.uri}/public/invoices"
    finantial_swagger  = "${google_cloud_run_v2_service.finantial.uri}/api-docs"
  }
}

output "artifact_registry_url" {
  description = "Repositorio onde as imagens dos modulos devem ser publicadas."
  value       = local.repository_url
}

output "pubsub_topics" {
  description = "Topicos do Pub/Sub da estrutura de mensageria."
  value       = { for k, v in google_pubsub_topic.topics : k => v.name }
}

output "pubsub_subscriptions" {
  description = "Subscriptions consumidas pelos modulos."
  value       = { for k, v in google_pubsub_subscription.subscriptions : k => v.name }
}

output "pubsub_dead_letter" {
  description = "Topicos e subscriptions de dead letter."
  value       = { for k, v in google_pubsub_topic.dead_letter : k => v.name }
}

output "documents_bucket" {
  description = "Bucket usado no upload de documentos do school-module."
  value       = google_storage_bucket.documents.name
}

output "service_accounts" {
  description = "Service accounts de cada modulo."
  value = {
    school    = google_service_account.school.email
    finantial = google_service_account.finantial.email
  }
}

output "trace_console_url" {
  description = "Explorador de traces do Cloud Trace."
  value       = "https://console.cloud.google.com/traces/explorer?project=${var.project_id}"
}

output "logs_console_url" {
  description = "Logs dos dois servicos no Cloud Logging."
  value       = "https://console.cloud.google.com/run?project=${var.project_id}"
}

output "always_on_instances" {
  description = "Instancias mantidas ligadas por servico. Acima de zero o custo e continuo — veja o README."
  value       = var.min_instance_count
}

output "smoke_test_commands" {
  description = "Chamadas rapidas para validar o ambiente."
  value = {
    school_health    = "curl -s ${google_cloud_run_v2_service.school.uri}/health"
    finantial_health = "curl -s ${google_cloud_run_v2_service.finantial.uri}/health"
    list_courses     = "curl -s ${google_cloud_run_v2_service.school.uri}/public/v1/courses"
    list_accounts    = "curl -s ${google_cloud_run_v2_service.finantial.uri}/public/accounts"
  }
}
