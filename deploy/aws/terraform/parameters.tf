# ─────────────────────────────────────────────────────────────────────────────
# Credenciais do banco e do cache remotos.
#
# Nenhum RDS ou ElastiCache é criado aqui: o Terraform só guarda as credenciais
# no Parameter Store como SecureString e as injeta na task definition pelo campo
# "secrets", de modo que a senha nunca aparece em texto plano na task definition
# nem no console do ECS.
# ─────────────────────────────────────────────────────────────────────────────

resource "aws_ssm_parameter" "db_password" {
  name        = "/${local.name}/${var.environment}/db/password"
  description = "Senha do PostgreSQL remoto usado pelos dois modulos"
  type        = "SecureString"
  value       = var.db_password

  tags = { Name = "${local.name}-db-password" }
}

resource "aws_ssm_parameter" "cache_password" {
  count = var.cache_password == "" ? 0 : 1

  name        = "/${local.name}/${var.environment}/cache/password"
  description = "Senha do Redis remoto usado pelo school-module"
  type        = "SecureString"
  value       = var.cache_password

  tags = { Name = "${local.name}-cache-password" }
}
