# ─────────────────────────────────────────────────────────────────────────────
# Uma service account por módulo, em vez da default do Compute Engine.
#
# A default vem com roles/editor no projeto inteiro. Separar dá a cada serviço
# só o que ele usa: o school publica em três tópicos e consome um, o finantial
# faz o inverso, e só o school toca no bucket.
# ─────────────────────────────────────────────────────────────────────────────

resource "google_service_account" "school" {
  account_id   = "${local.name}-school"
  display_name = "Colibri examples - school-module"

  depends_on = [google_project_service.required]
}

resource "google_service_account" "finantial" {
  account_id   = "${local.name}-finantial"
  display_name = "Colibri examples - finantial-module"

  depends_on = [google_project_service.required]
}

locals {
  service_accounts = {
    school    = google_service_account.school.email
    finantial = google_service_account.finantial.email
  }

  # Papéis de projeto que os dois precisam: exportar traces e métricas pelo
  # sidecar e escrever log estruturado.
  telemetry_roles = [
    "roles/cloudtrace.agent",
    "roles/monitoring.metricWriter",
    "roles/logging.logWriter",
  ]

  telemetry_bindings = {
    for pair in setproduct(keys(local.service_accounts), local.telemetry_roles) :
    "${pair[0]}-${replace(pair[1], "roles/", "")}" => {
      account = local.service_accounts[pair[0]]
      role    = pair[1]
    }
  }
}

resource "google_project_iam_member" "telemetry" {
  for_each = local.telemetry_bindings

  project = var.project_id
  role    = each.value.role
  member  = "serviceAccount:${each.value.account}"
}

# ─────────────────────────────────────────────────────────────────────────────
# Pub/Sub: publisher só nos tópicos que o serviço publica, subscriber só nas
# subscriptions que ele consome.
# ─────────────────────────────────────────────────────────────────────────────

locals {
  publisher_bindings = merge(
    { for t in local.school_publishes : "school-${t}" => {
      topic   = local.topics[t]
      account = google_service_account.school.email
    } },
    { for t in local.finantial_publishes : "finantial-${t}" => {
      topic   = local.topics[t]
      account = google_service_account.finantial.email
    } },
  )

  subscriber_bindings = merge(
    { for s in local.school_subscribes : "school-${s}" => {
      subscription = s
      account      = google_service_account.school.email
    } },
    { for s in local.finantial_subscribes : "finantial-${s}" => {
      subscription = s
      account      = google_service_account.finantial.email
    } },
  )
}

resource "google_pubsub_topic_iam_member" "publishers" {
  for_each = local.publisher_bindings

  topic  = each.value.topic
  role   = "roles/pubsub.publisher"
  member = "serviceAccount:${each.value.account}"

  depends_on = [google_pubsub_topic.topics]
}

resource "google_pubsub_subscription_iam_member" "subscribers" {
  for_each = local.subscriber_bindings

  subscription = google_pubsub_subscription.subscriptions[each.value.subscription].name
  role         = "roles/pubsub.subscriber"
  member       = "serviceAccount:${each.value.account}"
}

# O SDK também chama Publisher()/Subscriber() do cliente na inicialização, o que
# exige poder ler o tópico da subscription que está consumindo.
resource "google_pubsub_topic_iam_member" "subscriber_topic_viewer" {
  for_each = local.subscriber_bindings

  topic  = google_pubsub_subscription.subscriptions[each.value.subscription].topic
  role   = "roles/pubsub.viewer"
  member = "serviceAccount:${each.value.account}"
}

# ─────────────────────────────────────────────────────────────────────────────
# Bucket: só o school-module faz upload de documento.
# ─────────────────────────────────────────────────────────────────────────────

resource "google_storage_bucket_iam_member" "school_documents" {
  bucket = google_storage_bucket.documents.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.school.email}"
}

# ─────────────────────────────────────────────────────────────────────────────
# Secret Manager: cada serviço lê só os secrets que consome.
# ─────────────────────────────────────────────────────────────────────────────

resource "google_secret_manager_secret_iam_member" "db_password" {
  for_each = local.service_accounts

  secret_id = google_secret_manager_secret.db_password.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${each.value}"
}

resource "google_secret_manager_secret_iam_member" "otel_config" {
  for_each = local.service_accounts

  secret_id = google_secret_manager_secret.otel_config.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${each.value}"
}

resource "google_secret_manager_secret_iam_member" "cache_password" {
  count = var.cache_password == "" ? 0 : 1

  secret_id = google_secret_manager_secret.cache_password[0].secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.school.email}"
}
