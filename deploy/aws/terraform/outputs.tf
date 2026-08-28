output "school_module_url" {
  description = "URL publica do school-module."
  value       = "http://${aws_lb.main.dns_name}:${local.school_port}"
}

output "finantial_module_url" {
  description = "URL publica do finantial-module."
  value       = "http://${aws_lb.main.dns_name}:${local.finantial_port}"
}

output "school_module_swagger" {
  description = "Contrato OpenAPI do school-module."
  value       = "http://${aws_lb.main.dns_name}:${local.school_port}/api-docs"
}

output "finantial_module_swagger" {
  description = "Contrato OpenAPI do finantial-module."
  value       = "http://${aws_lb.main.dns_name}:${local.finantial_port}/api-docs"
}

# Todas as rotas de negócio dos dois módulos são registradas com o prefixo
# restserver.PublicApi, ou seja, ficam sob /public/. Só /health e /api-docs
# respondem na raiz.
output "api_paths" {
  description = "Caminhos das rotas expostas por cada modulo."
  value = {
    school_courses     = "${aws_lb.main.dns_name}:${local.school_port}/public/v1/courses"
    school_students    = "${aws_lb.main.dns_name}:${local.school_port}/public/v1/students"
    school_enrollments = "${aws_lb.main.dns_name}:${local.school_port}/public/v1/enrollments"
    finantial_accounts = "${aws_lb.main.dns_name}:${local.finantial_port}/public/accounts"
    finantial_invoices = "${aws_lb.main.dns_name}:${local.finantial_port}/public/invoices"
  }
}

output "alb_dns_name" {
  description = "DNS do Application Load Balancer."
  value       = aws_lb.main.dns_name
}

output "ecr_repository_urls" {
  description = "Repositorios ECR onde as imagens dos modulos devem ser publicadas."
  value       = { for k, v in aws_ecr_repository.modules : k => v.repository_url }
}

output "ecs_cluster_name" {
  description = "Nome do cluster ECS."
  value       = aws_ecs_cluster.main.name
}

output "sns_topic_arns" {
  description = "ARNs dos topicos SNS da estrutura de mensageria."
  value       = { for k, v in aws_sns_topic.topics : k => v.arn }
}

output "sqs_queue_urls" {
  description = "URLs das filas SQS consumidas pelos modulos."
  value       = { for k, v in aws_sqs_queue.queues : k => v.url }
}

output "sqs_dlq_urls" {
  description = "URLs das dead letter queues."
  value       = { for k, v in aws_sqs_queue.dlq : k => v.url }
}

output "documents_bucket" {
  description = "Bucket S3 usado no upload de documentos do school-module."
  value       = aws_s3_bucket.documents.bucket
}

output "task_role_arn" {
  description = "ARN da role da task, tambem usado como CLOUD_AWS_ROLE_ARN pelos servicos."
  value       = aws_iam_role.task.arn
}

output "xray_console_url" {
  description = "Mapa de servicos do X-Ray com os traces dos dois modulos."
  value       = "https://${var.aws_region}.console.aws.amazon.com/cloudwatch/home?region=${var.aws_region}#xray:service-map"
}

output "cloudwatch_logs" {
  description = "Log groups de cada componente."
  value = {
    school    = aws_cloudwatch_log_group.school.name
    finantial = aws_cloudwatch_log_group.finantial.name
    adot      = aws_cloudwatch_log_group.adot.name
    metrics   = aws_cloudwatch_log_group.metrics.name
  }
}

output "smoke_test_commands" {
  description = "Chamadas rapidas para validar o ambiente."
  value = {
    school_health    = "curl -s http://${aws_lb.main.dns_name}:${local.school_port}/health"
    finantial_health = "curl -s http://${aws_lb.main.dns_name}:${local.finantial_port}/health"
    list_courses     = "curl -s http://${aws_lb.main.dns_name}:${local.school_port}/public/v1/courses"
    list_accounts    = "curl -s http://${aws_lb.main.dns_name}:${local.finantial_port}/public/accounts"
  }
}
