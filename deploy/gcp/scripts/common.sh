#!/usr/bin/env bash
# Funções e variáveis compartilhadas pelos scripts de deploy da GCP.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
TF_DIR="${DEPLOY_DIR}/terraform"
REPO_ROOT="$(cd "${DEPLOY_DIR}/../.." && pwd)"
DOCKERFILE="${DEPLOY_DIR}/../docker/Dockerfile"

MODULES=(school-module finantial-module)

# Nomes fixos no código Go dos módulos — não levam prefixo de projeto.
TOPICS=(SCHOOL_ENROLLMENT SCHOOL_COURSE SCHOOL_STUDENT FINANCIAL_INSTALLMENT)
SUBSCRIPTIONS=(SCHOOL_ENROLLMENT_FINANCIAL SCHOOL_COURSE_FINANCIAL SCHOOL_STUDENT_FINANCIAL FINANCIAL_INSTALLMENT_SCHOOL)

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
  [[ ${#missing[@]} -eq 0 ]] || fail "Faltando na máquina: ${missing[*]}"
}

require_gcloud_auth() {
  gcloud auth print-access-token >/dev/null 2>&1 \
    || fail "Sem credenciais válidas. Rode: gcloud auth login && gcloud auth application-default login"
}

# Lê uma variável do terraform.tfvars, com fallback para o ambiente.
tfvar() { # tfvar <nome> <default>
  local name="$1" fallback="${2:-}" value=""
  if [[ -f "${TF_DIR}/terraform.tfvars" ]]; then
    value="$(grep -E "^[[:space:]]*${name}[[:space:]]*=" "${TF_DIR}/terraform.tfvars" \
      | head -n1 | sed -E 's/.*=[[:space:]]*"?([^"[:space:]]+)"?.*/\1/' || true)"
  fi
  printf '%s' "${value:-${fallback}}"
}

resolve_project() { tfvar project_id "${GOOGLE_CLOUD_PROJECT:-$(gcloud config get-value project 2>/dev/null || true)}"; }
resolve_region()  { tfvar region "${GOOGLE_CLOUD_REGION:-southamerica-east1}"; }
resolve_prefix()  { tfvar name_prefix "colibri-examples"; }

tf() { terraform -chdir="${TF_DIR}" "$@"; }
