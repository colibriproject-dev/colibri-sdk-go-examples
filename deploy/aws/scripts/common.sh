#!/usr/bin/env bash
# Funções e variáveis compartilhadas pelos scripts de deploy.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
TF_DIR="${DEPLOY_DIR}/terraform"
REPO_ROOT="$(cd "${DEPLOY_DIR}/../.." && pwd)"
DOCKERFILE="${DEPLOY_DIR}/../docker/Dockerfile"

MODULES=(school-module finantial-module)

_c_reset=$'\033[0m'; _c_blue=$'\033[34m'; _c_green=$'\033[32m'
_c_yellow=$'\033[33m'; _c_red=$'\033[31m'

info()  { printf '%s==>%s %s\n' "${_c_blue}"   "${_c_reset}" "$*"; }
ok()    { printf '%s ✓ %s %s\n' "${_c_green}"  "${_c_reset}" "$*"; }
warn()  { printf '%s ! %s %s\n' "${_c_yellow}" "${_c_reset}" "$*" >&2; }
fail()  { printf '%s ✗ %s %s\n' "${_c_red}"    "${_c_reset}" "$*" >&2; exit 1; }

require_tools() {
  local missing=()
  for tool in "$@"; do
    command -v "${tool}" >/dev/null 2>&1 || missing+=("${tool}")
  done
  if [[ ${#missing[@]} -gt 0 ]]; then
    fail "Faltando na máquina: ${missing[*]}"
  fi
}

# Confere que há credenciais válidas antes de qualquer coisa demorada.
require_aws_credentials() {
  aws sts get-caller-identity >/dev/null 2>&1 \
    || fail "Credenciais da AWS não configuradas ou expiradas (aws sts get-caller-identity falhou)."
}

aws_account_id() { aws sts get-caller-identity --query Account --output text; }

# Região: primeiro o que estiver no terraform.tfvars, depois o ambiente, por fim
# o padrão us-east-1 — a mesma ordem do default da variável no Terraform.
resolve_region() {
  local from_tfvars=""
  if [[ -f "${TF_DIR}/terraform.tfvars" ]]; then
    from_tfvars="$(grep -E '^[[:space:]]*aws_region[[:space:]]*=' "${TF_DIR}/terraform.tfvars" \
      | head -n1 | sed -E 's/.*=[[:space:]]*"?([^"[:space:]]+)"?.*/\1/' || true)"
  fi
  printf '%s' "${from_tfvars:-${AWS_REGION:-${AWS_DEFAULT_REGION:-us-east-1}}}"
}

tf() { terraform -chdir="${TF_DIR}" "$@"; }
