# ─────────────────────────────────────────────────────────────────────────────
# Entrada pública.
#
# Um único ALB com dois listeners, um por porta. A separação é por porta e não
# por path porque os serviços registram as rotas na raiz (/v1/courses,
# /accounts, /health, /api-docs) e o ALB não reescreve o path — um roteamento
# por prefixo entregaria /school/v1/courses ao serviço, que devolveria 404.
#
#   http://<alb>:8080/...  → school-module
#   http://<alb>:8081/...  → finantial-module
# ─────────────────────────────────────────────────────────────────────────────

resource "aws_lb" "main" {
  name               = substr("${local.name}-alb", 0, 32)
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = aws_subnet.public[*].id

  idle_timeout               = 60
  drop_invalid_header_fields = true

  tags = { Name = "${local.name}-alb" }
}

resource "aws_lb_target_group" "school" {
  name        = substr("${local.name}-school", 0, 32)
  port        = local.school_port
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id
  target_type = "ip"

  deregistration_delay = 30

  health_check {
    enabled             = true
    path                = "/health"
    matcher             = "200"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 3
  }

  tags = { Name = "${local.name}-school" }
}

resource "aws_lb_target_group" "finantial" {
  name        = substr("${local.name}-finantial", 0, 32)
  port        = local.finantial_port
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id
  target_type = "ip"

  deregistration_delay = 30

  health_check {
    enabled             = true
    path                = "/health"
    matcher             = "200"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 3
  }

  tags = { Name = "${local.name}-finantial" }
}

resource "aws_lb_listener" "school" {
  load_balancer_arn = aws_lb.main.arn
  port              = local.school_port
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.school.arn
  }
}

resource "aws_lb_listener" "finantial" {
  load_balancer_arn = aws_lb.main.arn
  port              = local.finantial_port
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.finantial.arn
  }
}
