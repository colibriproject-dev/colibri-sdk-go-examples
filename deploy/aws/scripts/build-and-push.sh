#!/usr/bin/env bash
#
# Compila as imagens dos dois módulos e publica no ECR.
#
#   ./build-and-push.sh              # tag derivada do commit atual + latest
#   ./build-and-push.sh v1.2.3       # tag explícita + latest
#
# A tag usada fica impressa no fim: passe-a em school_image_tag /
# finantial_image_tag se não quiser depender de "latest".

source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

require_tools aws docker
require_aws_credentials

REGION="$(resolve_region)"
ACCOUNT_ID="$(aws_account_id)"
REGISTRY="${ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com"

# O prefixo dos repositórios acompanha project_name no Terraform.
PROJECT_NAME="${PROJECT_NAME:-colibri-examples}"
if [[ -f "${TF_DIR}/terraform.tfvars" ]]; then
  from_tfvars="$(grep -E '^[[:space:]]*project_name[[:space:]]*=' "${TF_DIR}/terraform.tfvars" \
    | head -n1 | sed -E 's/.*=[[:space:]]*"?([^"[:space:]]+)"?.*/\1/' || true)"
  PROJECT_NAME="${from_tfvars:-${PROJECT_NAME}}"
fi

if [[ $# -ge 1 ]]; then
  TAG="$1"
elif git -C "${REPO_ROOT}" rev-parse --short HEAD >/dev/null 2>&1; then
  TAG="$(git -C "${REPO_ROOT}" rev-parse --short HEAD)"
else
  TAG="$(date +%Y%m%d%H%M%S)"
fi

info "Registry ......: ${REGISTRY}"
info "Projeto .......: ${PROJECT_NAME}"
info "Tag ...........: ${TAG}"

info "Autenticando no ECR"
aws ecr get-login-password --region "${REGION}" \
  | docker login --username AWS --password-stdin "${REGISTRY}" >/dev/null
ok "Autenticado"

for module in "${MODULES[@]}"; do
  repository="${PROJECT_NAME}/${module}"
  image="${REGISTRY}/${repository}"

  # O repositório é criado pelo Terraform; se ainda não existir, o deploy.sh
  # roda a fase de bootstrap antes de chamar este script.
  aws ecr describe-repositories --region "${REGION}" --repository-names "${repository}" >/dev/null 2>&1 \
    || fail "Repositório ECR ${repository} não existe. Rode ./deploy.sh, que cria os repositórios antes do build."

  info "Compilando ${module}"
  docker build \
    --platform linux/amd64 \
    -f "${DOCKERFILE}" \
    --build-arg "APP_SRC=${module}" \
    -t "${image}:${TAG}" \
    -t "${image}:latest" \
    "${REPO_ROOT}"

  info "Publicando ${module}"
  docker push "${image}:${TAG}"
  docker push "${image}:latest"
  ok "${image}:${TAG}"
done

echo
ok "Imagens publicadas com a tag ${TAG}"
echo "   Para fixar esta versão no deploy:"
echo "     terraform apply -var school_image_tag=${TAG} -var finantial_image_tag=${TAG}"
