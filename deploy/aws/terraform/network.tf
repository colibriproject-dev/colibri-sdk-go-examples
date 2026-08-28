# ─────────────────────────────────────────────────────────────────────────────
# VPC e subnets
# ─────────────────────────────────────────────────────────────────────────────

resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_support   = true
  enable_dns_hostnames = true # exigido pelo Cloud Map para o DNS privado resolver

  tags = { Name = "${local.name}-vpc" }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = { Name = "${local.name}-igw" }
}

resource "aws_subnet" "public" {
  count = var.availability_zones_count

  vpc_id                  = aws_vpc.main.id
  cidr_block              = cidrsubnet(var.vpc_cidr, 8, count.index)
  availability_zone       = local.azs[count.index]
  map_public_ip_on_launch = true

  tags = { Name = "${local.name}-public-${local.azs[count.index]}" }
}

resource "aws_subnet" "private" {
  count = var.enable_nat_gateway ? var.availability_zones_count : 0

  vpc_id            = aws_vpc.main.id
  cidr_block        = cidrsubnet(var.vpc_cidr, 8, count.index + 100)
  availability_zone = local.azs[count.index]

  tags = { Name = "${local.name}-private-${local.azs[count.index]}" }
}

# ─────────────────────────────────────────────────────────────────────────────
# Roteamento
# ─────────────────────────────────────────────────────────────────────────────

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = { Name = "${local.name}-rt-public" }
}

resource "aws_route_table_association" "public" {
  count = var.availability_zones_count

  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

# Um NAT por AZ evita que a queda de uma zona derrube a saída da outra.
resource "aws_eip" "nat" {
  count = var.enable_nat_gateway ? var.availability_zones_count : 0

  domain = "vpc"

  tags = { Name = "${local.name}-nat-eip-${local.azs[count.index]}" }
}

resource "aws_nat_gateway" "main" {
  count = var.enable_nat_gateway ? var.availability_zones_count : 0

  allocation_id = aws_eip.nat[count.index].id
  subnet_id     = aws_subnet.public[count.index].id

  depends_on = [aws_internet_gateway.main]

  tags = { Name = "${local.name}-nat-${local.azs[count.index]}" }
}

resource "aws_route_table" "private" {
  count = var.enable_nat_gateway ? var.availability_zones_count : 0

  vpc_id = aws_vpc.main.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.main[count.index].id
  }

  tags = { Name = "${local.name}-rt-private-${local.azs[count.index]}" }
}

resource "aws_route_table_association" "private" {
  count = var.enable_nat_gateway ? var.availability_zones_count : 0

  subnet_id      = aws_subnet.private[count.index].id
  route_table_id = aws_route_table.private[count.index].id
}

# ─────────────────────────────────────────────────────────────────────────────
# Security Groups
# ─────────────────────────────────────────────────────────────────────────────

resource "aws_security_group" "alb" {
  name        = "${local.name}-alb"
  description = "Entrada publica para os dois modulos"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "school-module"
    from_port   = local.school_port
    to_port     = local.school_port
    protocol    = "tcp"
    cidr_blocks = var.allowed_ingress_cidrs
  }

  ingress {
    description = "finantial-module"
    from_port   = local.finantial_port
    to_port     = local.finantial_port
    protocol    = "tcp"
    cidr_blocks = var.allowed_ingress_cidrs
  }

  egress {
    description = "Encaminhamento para os targets"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = { Name = "${local.name}-alb" }
}

# Todas as regras deste SG ficam em recursos aws_security_group_rule separados.
# Misturar blocos inline com regras separadas no mesmo Security Group faz o
# provider tentar revogar uma a cada plan, gerando diff perpétuo.
resource "aws_security_group" "service" {
  name        = "${local.name}-service"
  description = "Tasks ECS dos modulos"
  vpc_id      = aws_vpc.main.id

  tags = { Name = "${local.name}-service" }
}

resource "aws_security_group_rule" "service_egress" {
  type              = "egress"
  description       = "Saida para SNS, SQS, S3, ECR, CloudWatch, X-Ray e o banco remoto"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.service.id
}

# Só o ALB entra nas portas das aplicações.
resource "aws_security_group_rule" "service_from_alb_school" {
  type                     = "ingress"
  description              = "ALB -> school-module"
  from_port                = local.school_port
  to_port                  = local.school_port
  protocol                 = "tcp"
  security_group_id        = aws_security_group.service.id
  source_security_group_id = aws_security_group.alb.id
}

resource "aws_security_group_rule" "service_from_alb_finantial" {
  type                     = "ingress"
  description              = "ALB -> finantial-module"
  from_port                = local.finantial_port
  to_port                  = local.finantial_port
  protocol                 = "tcp"
  security_group_id        = aws_security_group.service.id
  source_security_group_id = aws_security_group.alb.id
}

# O school-module chama o finantial-module direto pelo DNS interno do Cloud Map,
# sem passar pelo ALB — por isso as tasks precisam se enxergar entre si.
resource "aws_security_group_rule" "service_internal" {
  type                     = "ingress"
  description              = "school-module -> finantial-module (service discovery)"
  from_port                = local.finantial_port
  to_port                  = local.finantial_port
  protocol                 = "tcp"
  security_group_id        = aws_security_group.service.id
  source_security_group_id = aws_security_group.service.id
}
