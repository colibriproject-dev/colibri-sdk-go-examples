#!/usr/bin/env bash
#
# Deploy completo do ambiente de exemplo do colibri-sdk-go na GCP.
#
# Três fases, porque o Cloud Run não sobe uma revisão cuja imagem ainda não
# existe no registry, e o registry só existe depois do Terraform:
#
#   1. bootstrap  habilita as APIs e cria o repositório do Artifact Registry
#   2. build      compila e publica as imagens dos dois módulos
#   3. apply      cria o resto (Pub/Sub, bucket, secrets, IAM, Cloud Run)
#
# Uso:
#   ./deploy.sh                 # ciclo completo
#   ./deploy.sh --plan          # só mostra o plano da fase 3
#   ./deploy.sh --skip-build    # reaproveita as imagens já publicadas
#   ./deploy.sh --auto-approve  # não pergunta antes de aplicar
#   ./deploy.sh --tag v1.2.3    # fixa uma tag em vez de "latest"

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

require_tools terraform gcloud docker
require_gcloud_auth

[[ -f "${TF_DIR}/terraform.tfvars" ]] \
  || fail "Falta ${TF_DIR}/terraform.tfvars. Copie terraform.tfvars.example e preencha project_id e as credenciais do banco remoto."

PROJECT="$(resolve_project)"
REGION="$(resolve_region)"
info "Projeto: ${PROJECT}   Região: ${REGION}"

# ── Fase 0: init ─────────────────────────────────────────────────────────────
info "terraform init"
tf init -input=false
ok "Providers prontos"

if [[ "${PLAN_ONLY}" == true ]]; then
  info "terraform plan"
  tf plan -input=false
  exit 0
fi

# ── Fase 1: APIs e registry ──────────────────────────────────────────────────
# Aplicado isolado com -target para que o build tenha onde publicar. É o único
# lugar em que -target se justifica: o resto do grafo vai inteiro na fase 3.
info "Fase 1/3 — habilitando APIs e criando o Artifact Registry"
tf apply -input=false -auto-approve \
  -target=google_project_service.required \
  -target=google_artifact_registry_repository.modules
ok "Registry pronto"

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
  ok "Imagens no Artifact Registry"
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

always_on="$(tf output -raw always_on_instances 2>/dev/null || echo 1)"
if [[ "${always_on}" != "0" ]]; then
  echo
  warn "min_instance_count = ${always_on}: as instâncias ficam ligadas 24/7 com CPU"
  warn "sempre alocada, o que é necessário para os consumidores do Pub/Sub."
  warn "Isso custa cerca de US\$ 90/mês com os dois serviços."
  warn "Quando não estiver usando: ./pause.sh (zera o custo, preserva a infra)."
fi

echo
info "Valide o ambiente com: ./smoke-test.sh"
