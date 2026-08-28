# ─────────────────────────────────────────────────────────────────────────────
# Cluster
# ─────────────────────────────────────────────────────────────────────────────

resource "aws_ecs_cluster" "main" {
  name = local.name

  setting {
    name  = "containerInsights"
    value = var.enable_container_insights ? "enabled" : "disabled"
  }
}

resource "aws_ecs_cluster_capacity_providers" "main" {
  cluster_name       = aws_ecs_cluster.main.name
  capacity_providers = ["FARGATE", "FARGATE_SPOT"]

  default_capacity_provider_strategy {
    capacity_provider = "FARGATE"
    weight            = 1
  }
}

# ─────────────────────────────────────────────────────────────────────────────
# Service discovery interno
#
# O school-module chama o finantial-module por HTTP (FINANCIAL_MODULE_BASE_URL).
# Resolver esse endereço pelo Cloud Map em vez de mandar a chamada de volta pelo
# ALB público mantém o tráfego dentro da VPC e evita o hairpin pelo load balancer.
# ─────────────────────────────────────────────────────────────────────────────

resource "aws_service_discovery_private_dns_namespace" "main" {
  name        = local.internal_dns_namespace
  description = "DNS interno dos modulos de exemplo do colibri-sdk-go"
  vpc         = aws_vpc.main.id
}

resource "aws_service_discovery_service" "school" {
  name = "school-module"

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.main.id

    dns_records {
      ttl  = 10
      type = "A"
    }

    routing_policy = "MULTIVALUE"
  }

  health_check_custom_config {
    failure_threshold = 1
  }
}

resource "aws_service_discovery_service" "finantial" {
  name = "finantial-module"

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.main.id

    dns_records {
      ttl  = 10
      type = "A"
    }

    routing_policy = "MULTIVALUE"
  }

  health_check_custom_config {
    failure_threshold = 1
  }
}

# ─────────────────────────────────────────────────────────────────────────────
# Blocos reaproveitados pelas duas task definitions
# ─────────────────────────────────────────────────────────────────────────────

locals {
  # Variáveis comuns aos dois módulos.
  #
  # CLOUD_HOST fica de fora de propósito: fora de um ambiente local o SDK usa a
  # credential chain e o endpoint padrão da AWS, e apontar CLOUD_HOST aqui faria
  # o cliente tentar falar com um emulador que não existe.
  common_environment = [
    # AWS_REGION/AWS_DEFAULT_REGION são obrigatórias e não decorativas: em
    # ambiente de nuvem o SDK chama awsconfig.LoadDefaultConfig sem passar
    # região (CLOUD_REGION só é aplicada fora da nuvem), e o ARN do tópico SNS
    # é montado com cloud.GetAwsConfig().Region. O ECS não injeta essa variável
    # sozinho — sem ela o publish sai com a região vazia no ARN e falha.
    { name = "AWS_REGION", value = var.aws_region },
    { name = "AWS_DEFAULT_REGION", value = var.aws_region },

    { name = "APP_TYPE", value = "service" },
    { name = "ENVIRONMENT", value = var.environment },
    { name = "CLOUD", value = "aws" },
    { name = "CLOUD_REGION", value = var.aws_region },
    { name = "CLOUD_AWS_ROLE_ARN", value = aws_iam_role.task.arn },
    { name = "COLIBRI_MESSAGING", value = "CLOUD_DEFAULT" },
    { name = "LOG_LEVEL", value = var.log_level },
    { name = "SQL_DB_HOST", value = var.db_host },
    { name = "SQL_DB_PORT", value = var.db_port },
    { name = "SQL_DB_USER", value = var.db_user },
    { name = "SQL_DB_SSL_MODE", value = var.db_ssl_mode },
    { name = "SQL_DB_MIGRATION", value = tostring(var.db_run_migrations) },
    { name = "OTEL_EXPORTER_OTLP_ENDPOINT", value = local.otlp_endpoint },
    { name = "OTEL_EXPORTER_OTLP_PROTOCOL", value = "http" },
  ]

  db_secret = [
    { name = "SQL_DB_PASSWORD", valueFrom = aws_ssm_parameter.db_password.arn },
  ]

  cache_secret = var.cache_password == "" ? [] : [
    { name = "CACHE_PASSWORD", valueFrom = aws_ssm_parameter.cache_password[0].arn },
  ]
}

# O sidecar é idêntico nos dois módulos; só o log stream prefix muda.
locals {
  adot_container = {
    name      = "adot-collector"
    image     = var.adot_image
    essential = true

    environment = [
      { name = "AWS_REGION", value = var.aws_region },
      { name = "AWS_DEFAULT_REGION", value = var.aws_region },
    ]

    secrets = [
      { name = "AOT_CONFIG_CONTENT", valueFrom = aws_ssm_parameter.adot_config.arn },
    ]

    portMappings = [
      { containerPort = 4317, protocol = "tcp" },
      { containerPort = 4318, protocol = "tcp" },
    ]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = aws_cloudwatch_log_group.adot.name
        "awslogs-region"        = var.aws_region
        "awslogs-stream-prefix" = "adot"
      }
    }
  }
}

# ─────────────────────────────────────────────────────────────────────────────
# school-module
# ─────────────────────────────────────────────────────────────────────────────

resource "aws_ecs_task_definition" "school" {
  family                   = "${local.name}-school-module"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.task_cpu
  memory                   = var.task_memory
  execution_role_arn       = aws_iam_role.task_execution.arn
  task_role_arn            = aws_iam_role.task.arn

  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture        = "X86_64"
  }

  container_definitions = jsonencode([
    {
      name      = "school-module"
      image     = "${aws_ecr_repository.modules["school-module"].repository_url}:${var.school_image_tag}"
      essential = true

      portMappings = [
        { containerPort = local.school_port, protocol = "tcp" },
      ]

      environment = concat(local.common_environment, [
        { name = "APP_NAME", value = "Colibri SDK Go Test - School Module" },
        { name = "PORT", value = tostring(local.school_port) },
        { name = "SQL_DB_NAME", value = var.db_name_school },
        { name = "CACHE_URI", value = var.cache_uri },
        { name = "STORAGE_BUCKET", value = aws_s3_bucket.documents.bucket },
        { name = "FINANCIAL_MODULE_BASE_URL", value = local.finantial_internal_url },
        { name = "OTEL_SERVICE_NAME", value = "school-module" },
      ])

      secrets = concat(local.db_secret, local.cache_secret)

      # A aplicação abre o exporter OTLP logo no start, então o collector precisa
      # estar de pé antes.
      dependsOn = [
        { containerName = "adot-collector", condition = "START" },
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.school.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "school"
        }
      }
    },
    local.adot_container,
  ])
}

resource "aws_ecs_service" "school" {
  name            = "school-module"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.school.arn
  desired_count   = var.school_desired_count
  launch_type     = "FARGATE"

  enable_execute_command = true

  network_configuration {
    subnets          = local.service_subnet_ids
    security_groups  = [aws_security_group.service.id]
    assign_public_ip = local.service_public_ip
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.school.arn
    container_name   = "school-module"
    container_port   = local.school_port
  }

  service_registries {
    registry_arn = aws_service_discovery_service.school.arn
  }

  # Dá tempo das migrations rodarem contra o banco remoto antes do primeiro
  # health check contar como falha.
  health_check_grace_period_seconds = 120

  deployment_circuit_breaker {
    enable   = true
    rollback = true
  }

  depends_on = [
    aws_lb_listener.school,
    aws_iam_role_policy.task,
    aws_iam_role_policy.task_execution_secrets,
  ]
}

# ─────────────────────────────────────────────────────────────────────────────
# finantial-module
# ─────────────────────────────────────────────────────────────────────────────

resource "aws_ecs_task_definition" "finantial" {
  family                   = "${local.name}-finantial-module"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.task_cpu
  memory                   = var.task_memory
  execution_role_arn       = aws_iam_role.task_execution.arn
  task_role_arn            = aws_iam_role.task.arn

  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture        = "X86_64"
  }

  container_definitions = jsonencode([
    {
      name      = "finantial-module"
      image     = "${aws_ecr_repository.modules["finantial-module"].repository_url}:${var.finantial_image_tag}"
      essential = true

      portMappings = [
        { containerPort = local.finantial_port, protocol = "tcp" },
      ]

      environment = concat(local.common_environment, [
        { name = "APP_NAME", value = "Colibri SDK Go Test - Finantial Module" },
        { name = "PORT", value = tostring(local.finantial_port) },
        { name = "SQL_DB_NAME", value = var.db_name_finantial },
        { name = "OTEL_SERVICE_NAME", value = "finantial-module" },
      ])

      secrets = local.db_secret

      dependsOn = [
        { containerName = "adot-collector", condition = "START" },
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.finantial.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "finantial"
        }
      }
    },
    local.adot_container,
  ])
}

resource "aws_ecs_service" "finantial" {
  name            = "finantial-module"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.finantial.arn
  desired_count   = var.finantial_desired_count
  launch_type     = "FARGATE"

  enable_execute_command = true

  network_configuration {
    subnets          = local.service_subnet_ids
    security_groups  = [aws_security_group.service.id]
    assign_public_ip = local.service_public_ip
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.finantial.arn
    container_name   = "finantial-module"
    container_port   = local.finantial_port
  }

  service_registries {
    registry_arn = aws_service_discovery_service.finantial.arn
  }

  health_check_grace_period_seconds = 120

  deployment_circuit_breaker {
    enable   = true
    rollback = true
  }

  depends_on = [
    aws_lb_listener.finantial,
    aws_iam_role_policy.task,
    aws_iam_role_policy.task_execution_secrets,
  ]
}
