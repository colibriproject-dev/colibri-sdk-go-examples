# ─────────────────────────────────────────────────────────────────────────────
# Estrutura de mensagens entre os dois módulos.
#
#   school-module    ──SCHOOL_ENROLLMENT──────▶ SCHOOL_ENROLLMENT_FINANCIAL   ──▶ finantial-module
#   school-module    ──SCHOOL_COURSE──────────▶ SCHOOL_COURSE_FINANCIAL       ──▶ finantial-module
#   school-module    ──SCHOOL_STUDENT─────────▶ SCHOOL_STUDENT_FINANCIAL      ──▶ finantial-module
#   finantial-module ──FINANCIAL_INSTALLMENT──▶ FINANCIAL_INSTALLMENT_SCHOOL  ──▶ school-module
#
# raw_message_delivery fica FALSE de propósito: o consumidor do SDK
# (pkg/messaging/aws.go) desserializa primeiro o envelope de notificação do SNS
# e só depois o corpo da mensagem. Com raw delivery ligado o payload chega sem
# esse envelope e toda mensagem falha ao ser lida.
# ─────────────────────────────────────────────────────────────────────────────

resource "aws_sns_topic" "topics" {
  for_each = local.topics

  name = each.value

  tags = { Name = each.value }
}

resource "aws_sqs_queue" "dlq" {
  for_each = local.queues

  name                      = "${each.key}_DLQ"
  message_retention_seconds = 1209600 # 14 dias, o máximo do SQS

  tags = { Name = "${each.key}_DLQ" }
}

resource "aws_sqs_queue" "queues" {
  for_each = local.queues

  name                       = each.key
  visibility_timeout_seconds = var.sqs_visibility_timeout_seconds
  message_retention_seconds  = var.sqs_message_retention_seconds

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.dlq[each.key].arn
    maxReceiveCount     = var.sqs_max_receive_count
  })

  tags = { Name = each.key }
}

# Deixa a DLQ aceitar redrive só da fila de origem correspondente.
resource "aws_sqs_queue_redrive_allow_policy" "dlq" {
  for_each = local.queues

  queue_url = aws_sqs_queue.dlq[each.key].id

  redrive_allow_policy = jsonencode({
    redrivePermission = "byQueue"
    sourceQueueArns   = [aws_sqs_queue.queues[each.key].arn]
  })
}

resource "aws_sns_topic_subscription" "queues" {
  for_each = local.queues

  topic_arn = aws_sns_topic.topics[each.value].arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.queues[each.key].arn

  raw_message_delivery = false
}

# Autoriza cada tópico a publicar na sua fila.
data "aws_iam_policy_document" "queue_policy" {
  for_each = local.queues

  statement {
    sid    = "AllowSNSPublish"
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["sns.amazonaws.com"]
    }

    actions   = ["sqs:SendMessage"]
    resources = [aws_sqs_queue.queues[each.key].arn]

    condition {
      test     = "ArnEquals"
      variable = "aws:SourceArn"
      values   = [aws_sns_topic.topics[each.value].arn]
    }
  }
}

resource "aws_sqs_queue_policy" "queues" {
  for_each = local.queues

  queue_url = aws_sqs_queue.queues[each.key].id
  policy    = data.aws_iam_policy_document.queue_policy[each.key].json
}
