# ─────────────────────────────────────────────────────────────────────────────
# Estrutura de mensagens entre os dois módulos.
#
#   school-module    ──SCHOOL_ENROLLMENT──────▶ SCHOOL_ENROLLMENT_FINANCIAL   ──▶ finantial-module
#   school-module    ──SCHOOL_COURSE──────────▶ SCHOOL_COURSE_FINANCIAL       ──▶ finantial-module
#   school-module    ──SCHOOL_STUDENT─────────▶ SCHOOL_STUDENT_FINANCIAL      ──▶ finantial-module
#   finantial-module ──FINANCIAL_INSTALLMENT──▶ FINANCIAL_INSTALLMENT_SCHOOL  ──▶ school-module
#
# Diferença importante em relação ao deploy AWS: no Pub/Sub o corpo da mensagem
# chega exatamente como o produtor publicou. Não existe o envelope de
# notificação que o SNS coloca por cima do payload, então também não existe o
# equivalente ao raw_message_delivery para configurar errado.
# ─────────────────────────────────────────────────────────────────────────────

resource "google_pubsub_topic" "topics" {
  for_each = local.topics

  name                       = each.value
  message_retention_duration = var.message_retention_duration

  depends_on = [google_project_service.required]
}

# Um tópico de dead letter por subscription. O Pub/Sub não tem DLQ nativa como o
# SQS: o descarte é para outro tópico, que por sua vez precisa de uma
# subscription própria para a mensagem ficar retida e legível.
resource "google_pubsub_topic" "dead_letter" {
  for_each = local.subscriptions

  name = "${each.key}_DLQ"

  depends_on = [google_project_service.required]
}

resource "google_pubsub_subscription" "dead_letter" {
  for_each = local.subscriptions

  name  = "${each.key}_DLQ"
  topic = google_pubsub_topic.dead_letter[each.key].id

  message_retention_duration = "604800s"
  retain_acked_messages      = false

  # Sem isso a subscription é apagada sozinha após 31 dias sem atividade — e uma
  # DLQ saudável fica exatamente assim, sem atividade nenhuma.
  expiration_policy {
    ttl = ""
  }
}

resource "google_pubsub_subscription" "subscriptions" {
  for_each = local.subscriptions

  name  = each.key
  topic = google_pubsub_topic.topics[each.value].id

  ack_deadline_seconds       = var.ack_deadline_seconds
  message_retention_duration = var.message_retention_duration

  expiration_policy {
    ttl = ""
  }

  dead_letter_policy {
    dead_letter_topic     = google_pubsub_topic.dead_letter[each.key].id
    max_delivery_attempts = var.max_delivery_attempts
  }

  retry_policy {
    minimum_backoff = "10s"
    maximum_backoff = "600s"
  }

  # Só o publisher precisa existir antes: sem ele a subscription é criada com
  # uma dead_letter_policy que o Pub/Sub não consegue executar. O subscriber é
  # aplicado depois, sobre a subscription já criada — por isso não entra aqui,
  # sob pena de fechar um ciclo no grafo.
  depends_on = [google_pubsub_topic_iam_member.dead_letter_publisher]
}

# ─────────────────────────────────────────────────────────────────────────────
# Permissões do service agent do Pub/Sub.
#
# Quem move a mensagem para o dead letter topic não é a aplicação: é o próprio
# Pub/Sub, com a identidade service-<numero>@gcp-sa-pubsub. Ele precisa publicar
# no tópico de destino e confirmar a mensagem na subscription de origem. Sem os
# dois papéis a dead_letter_policy é aceita no apply e simplesmente não funciona
# em runtime — a mensagem fica redoando para sempre.
# ─────────────────────────────────────────────────────────────────────────────

resource "google_pubsub_topic_iam_member" "dead_letter_publisher" {
  for_each = local.subscriptions

  topic  = google_pubsub_topic.dead_letter[each.key].name
  role   = "roles/pubsub.publisher"
  member = local.pubsub_service_agent
}

resource "google_pubsub_subscription_iam_member" "dead_letter_subscriber" {
  for_each = local.subscriptions

  subscription = google_pubsub_subscription.subscriptions[each.key].name
  role         = "roles/pubsub.subscriber"
  member       = local.pubsub_service_agent
}
