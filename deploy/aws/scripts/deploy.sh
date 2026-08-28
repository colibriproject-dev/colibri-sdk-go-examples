#!/usr/bin/env bash
#
# Deploy completo do ambiente de exemplo do colibri-sdk-go na AWS.
#
# O deploy tem três fases porque o ECS precisa da imagem existindo no ECR antes
# de conseguir subir a task, e o ECR só existe depois do Terraform:
#
#   1. bootstrap  cria os repositórios ECR
#   2. build      compila e publica as imagens dos dois módulos
#   3. apply      cria o resto (rede, SNS/SQS, S3, ALB, ECS, observabilidade)
#
# Uso:
#   ./deploy.sh                 # ciclo completo
#   ./deploy.sh --plan          # só mostra o plano da fase 3, não aplica nada
#   ./deploy.sh --skip-build    # reaproveita as imagens já publicadas
#   ./deploy.sh --auto-approve  # não pergunta antes de aplicar

source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

PLAN_ONLY=false
SKIP_BUILD=false
AUTO_APPROVE=false
IMAGE_TAG=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --plan)         PLAN_ONLY=true ;;
    --skip-build)   SKIP_BUILD=true ;;
    --auto-approve) AUTO_APPROVE=true ;;
    --tag)          IMAGE_TAG="${2:-}"; shift ;;
    -h|--help)      sed -n '2,20p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//'; exit 0 ;;
    *)              fail "Opção desconhecida: $1" ;;
  esac
  shift
done

require_tools terraform aws docker
require_aws_credentials

[[ -f "${TF_DIR}/terraform.tfvars" ]] \
  || fail "Falta ${TF_DIR}/terraform.tfvars. Copie terraform.tfvars.example e preencha as credenciais do banco remoto."

# ── Fase 0: init ─────────────────────────────────────────────────────────────
info "terraform init"
tf init -input=false
ok "Providers prontos"

if [[ "${PLAN_ONLY}" == true ]]; then
  info "terraform plan"
  tf plan -input=false
  exit 0
fi

# ── Fase 1: repositórios ECR ─────────────────────────────────────────────────
# Aplicado isolado com -target para que o build tenha onde publicar. É o único
# lugar em que -target se justifica aqui: o resto do grafo é aplicado inteiro
# na fase 3.
info "Fase 1/3 — criando os repositórios ECR"
tf apply -input=false -auto-approve \
  -target=aws_ecr_repository.modules \
  -target=aws_ecr_lifecycle_policy.modules
ok "Repositórios ECR prontos"

# ── Fase 2: imagens ──────────────────────────────────────────────────────────
if [[ "${SKIP_BUILD}" == true ]]; then
  warn "Fase 2/3 — build pulado (--skip-build); usando as imagens já publicadas"
else
  info "Fase 2/3 — compilando e publicando as imagens"
  if [[ -n "${IMAGE_TAG}" ]]; then
    "${SCRIPT_DIR}/build-and-push.sh" "${IMAGE_TAG}"
  else
    "${SCRIPT_DIR}/build-and-push.sh"
  fi
  ok "Imagens no ECR"
fi

# ── Fase 3: infraestrutura completa ──────────────────────────────────────────
info "Fase 3/3 — aplicando a infraestrutura"

tf_args=(-input=false)
[[ "${AUTO_APPROVE}" == true ]] && tf_args+=(-auto-approve)
[[ -n "${IMAGE_TAG}" ]] && tf_args+=(-var "school_image_tag=${IMAGE_TAG}" -var "finantial_image_tag=${IMAGE_TAG}")

tf apply "${tf_args[@]}"

echo
ok "Ambiente no ar"
echo
tf output

echo
info "Os serviços levam de 2 a 4 minutos para ficarem saudáveis no ALB"
info "Acompanhe com: ./scripts/status.sh"
