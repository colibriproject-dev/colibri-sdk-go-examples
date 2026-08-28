# ─────────────────────────────────────────────────────────────────────────────
# Role de execução: usada pelo agente do ECS para puxar a imagem do ECR,
# escrever nos log groups e ler os SecureStrings do Parameter Store.
# ─────────────────────────────────────────────────────────────────────────────

data "aws_iam_policy_document" "ecs_assume_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "task_execution" {
  name               = "${local.name}-task-execution"
  assume_role_policy = data.aws_iam_policy_document.ecs_assume_role.json
}

resource "aws_iam_role_policy_attachment" "task_execution_managed" {
  role       = aws_iam_role.task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

data "aws_iam_policy_document" "task_execution_secrets" {
  statement {
    sid    = "ReadSecureParameters"
    effect = "Allow"

    actions = [
      "ssm:GetParameters",
      "ssm:GetParameter",
    ]

    resources = concat(
      [
        aws_ssm_parameter.db_password.arn,
        aws_ssm_parameter.adot_config.arn,
      ],
      aws_ssm_parameter.cache_password[*].arn,
    )
  }

  statement {
    sid       = "DecryptSecureParameters"
    effect    = "Allow"
    actions   = ["kms:Decrypt"]
    resources = ["*"]

    condition {
      test     = "StringEquals"
      variable = "kms:ViaService"
      values   = ["ssm.${var.aws_region}.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy" "task_execution_secrets" {
  name   = "${local.name}-task-execution-secrets"
  role   = aws_iam_role.task_execution.id
  policy = data.aws_iam_policy_document.task_execution_secrets.json
}

# ─────────────────────────────────────────────────────────────────────────────
# Role da task: é com ela que o colibri-sdk-go fala com SNS, SQS, S3 e com o
# sidecar ADOT publica em X-Ray e CloudWatch.
#
# O ARN desta role também vai para a variável CLOUD_AWS_ROLE_ARN. O SDK faz o
# parse dela para descobrir a partition e o account id com que monta o ARN do
# tópico SNS na hora de publicar (pkg/base/cloud/aws.go → getAwsARN).
# Sem essa variável o publish monta um ARN com account id vazio e falha.
# ─────────────────────────────────────────────────────────────────────────────

resource "aws_iam_role" "task" {
  name               = "${local.name}-task"
  assume_role_policy = data.aws_iam_policy_document.ecs_assume_role.json
}

data "aws_iam_policy_document" "task" {
  statement {
    sid       = "PublishToTopics"
    effect    = "Allow"
    actions   = ["sns:Publish", "sns:GetTopicAttributes"]
    resources = [for t in aws_sns_topic.topics : t.arn]
  }

  # ListTopics é chamado no start do módulo de messaging como healthcheck da
  # conexão (newAwsMessaging), e essa API não aceita recurso específico.
  statement {
    sid       = "ListTopics"
    effect    = "Allow"
    actions   = ["sns:ListTopics"]
    resources = ["*"]
  }

  statement {
    sid    = "ConsumeQueues"
    effect = "Allow"

    actions = [
      "sqs:ReceiveMessage",
      "sqs:DeleteMessage",
      "sqs:ChangeMessageVisibility",
      "sqs:GetQueueAttributes",
      "sqs:GetQueueUrl",
    ]

    resources = concat(
      [for q in aws_sqs_queue.queues : q.arn],
      [for q in aws_sqs_queue.dlq : q.arn],
    )
  }

  statement {
    sid       = "BucketLevel"
    effect    = "Allow"
    actions   = ["s3:ListBucket", "s3:GetBucketLocation"]
    resources = [aws_s3_bucket.documents.arn]
  }

  statement {
    sid       = "ObjectLevel"
    effect    = "Allow"
    actions   = ["s3:GetObject", "s3:PutObject", "s3:DeleteObject", "s3:AbortMultipartUpload"]
    resources = ["${aws_s3_bucket.documents.arn}/*"]
  }

  # Usado pelo sidecar ADOT: traces para o X-Ray, métricas para o CloudWatch.
  statement {
    sid    = "Telemetry"
    effect = "Allow"

    actions = [
      "xray:PutTraceSegments",
      "xray:PutTelemetryRecords",
      "xray:GetSamplingRules",
      "xray:GetSamplingTargets",
      "xray:GetSamplingStatisticSummaries",
      "cloudwatch:PutMetricData",
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
      "logs:DescribeLogStreams",
      "logs:DescribeLogGroups",
      "ssm:GetParameters",
    ]

    resources = ["*"]
  }
}

# Necessário para `aws ecs execute-command` abrir shell dentro da task, que é a
# forma mais direta de investigar um problema de conexão com o banco remoto.
data "aws_iam_policy_document" "task_exec_command" {
  statement {
    sid    = "SSMSession"
    effect = "Allow"

    actions = [
      "ssmmessages:CreateControlChannel",
      "ssmmessages:CreateDataChannel",
      "ssmmessages:OpenControlChannel",
      "ssmmessages:OpenDataChannel",
    ]

    resources = ["*"]
  }
}

resource "aws_iam_role_policy" "task_exec_command" {
  name   = "${local.name}-task-exec-command"
  role   = aws_iam_role.task.id
  policy = data.aws_iam_policy_document.task_exec_command.json
}

resource "aws_iam_role_policy" "task" {
  name   = "${local.name}-task"
  role   = aws_iam_role.task.id
  policy = data.aws_iam_policy_document.task.json
}
