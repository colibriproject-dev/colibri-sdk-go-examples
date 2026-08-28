data "aws_caller_identity" "current" {}

data "aws_availability_zones" "available" {
  state = "available"
}

locals {
  name       = var.project_name
  account_id = data.aws_caller_identity.current.account_id
  azs        = slice(data.aws_availability_zones.available.names, 0, var.availability_zones_count)

  # As subnets onde as tasks rodam dependem de haver NAT Gateway. Sem NAT elas
  # ficam nas subnets públicas com IP público — o Security Group continua
  # recusando qualquer entrada que não venha do ALB.
  service_subnet_ids = var.enable_nat_gateway ? aws_subnet.private[*].id : aws_subnet.public[*].id
  service_public_ip  = !var.enable_nat_gateway

  # ── Nomes de tópicos e filas ───────────────────────────────────────────────
  # ATENÇÃO: estes nomes são literais no código Go dos dois módulos
  # (messaging.NewProducer("SCHOOL_COURSE"), QueueName() == "SCHOOL_COURSE_FINANCIAL"...).
  # O SDK monta o ARN do tópico como arn:<partition>:sns:<region>:<account>:<nome>
  # e resolve a fila por GetQueueUrl(nome), então NÃO é possível prefixar nem
  # sufixar nada aqui sem alterar o código dos serviços.
  topics = {
    school_enrollment     = "SCHOOL_ENROLLMENT"
    school_course         = "SCHOOL_COURSE"
    school_student        = "SCHOOL_STUDENT"
    financial_installment = "FINANCIAL_INSTALLMENT"
  }

  # Cada fila assina exatamente um tópico. O valor é a CHAVE do mapa de tópicos
  # acima, o que permite referenciar aws_sns_topic.topics[each.value].arn em vez
  # de montar o ARN na mão.
  queues = {
    "SCHOOL_ENROLLMENT_FINANCIAL"  = "school_enrollment"
    "SCHOOL_COURSE_FINANCIAL"      = "school_course"
    "SCHOOL_STUDENT_FINANCIAL"     = "school_student"
    "FINANCIAL_INSTALLMENT_SCHOOL" = "financial_installment"
  }

  # ── Service discovery interno ──────────────────────────────────────────────
  internal_dns_namespace = "${var.project_name}.internal"
  finantial_internal_url = "http://finantial-module.${local.internal_dns_namespace}:${local.finantial_port}"

  school_port    = 8080
  finantial_port = 8081

  # Endpoint OTLP local: o sidecar ADOT escuta em 4318 dentro da própria task.
  otlp_endpoint = "localhost:4318"
}
