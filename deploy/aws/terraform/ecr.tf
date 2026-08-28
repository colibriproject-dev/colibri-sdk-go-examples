locals {
  ecr_repositories = ["school-module", "finantial-module"]
}

resource "aws_ecr_repository" "modules" {
  for_each = toset(local.ecr_repositories)

  name                 = "${local.name}/${each.value}"
  image_tag_mutability = "MUTABLE"
  force_delete         = true

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = { Name = each.value }
}

# Mantém apenas as 10 imagens mais recentes de cada módulo.
resource "aws_ecr_lifecycle_policy" "modules" {
  for_each = aws_ecr_repository.modules

  repository = each.value.name

  policy = jsonencode({
    rules = [{
      rulePriority = 1
      description  = "Mantem as 10 imagens mais recentes"
      selection = {
        tagStatus   = "any"
        countType   = "imageCountMoreThan"
        countNumber = 10
      }
      action = { type = "expire" }
    }]
  })
}
