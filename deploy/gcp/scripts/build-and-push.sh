#!/usr/bin/env bash
#
# Compila as imagens dos dois módulos e publica no Artifact Registry.
#
#   ./build-and-push.sh              # tag derivada do commit atual + latest
#   ./build-and-push.sh v1.2.3       # tag explícita + latest

source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

require_tools gcloud docker
require_gcloud_auth

PROJECT="$(resolve_project)"
REGION="$(resolve_region)"
PREFIX="$(resolve_prefix)"

[[ -n "${PROJECT}" ]] || fail "project_id não definido. Preencha no terraform.tfvars ou rode gcloud config set project."

REGISTRY="${REGION}-docker.pkg.dev/${PROJECT}/${PREFIX}"

if [[ $# -ge 1 ]]; then
  TAG="$1"
elif git -C "${REPO_ROOT}" rev-parse --short HEAD >/dev/null 2>&1; then
  TAG="$(git -C "${REPO_ROOT}" rev-parse --short HEAD)"
else
  TAG="$(date +%Y%m%d%H%M%S)"
fi

info "Registry ..: ${REGISTRY}"
info "Tag .......: ${TAG}"

info "Autenticando o docker no Artifact Registry"
gcloud auth configure-docker "${REGION}-docker.pkg.dev" --quiet >/dev/null
ok "Autenticado"

gcloud artifacts repositories describe "${PREFIX}" \
  --project "${PROJECT}" --location "${REGION}" >/dev/null 2>&1 \
  || fail "Repositório ${PREFIX} não existe em ${REGION}. Rode ./deploy.sh, que o cria antes do build."

for module in "${MODULES[@]}"; do
  image="${REGISTRY}/${module}"

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
