#!/usr/bin/env bash
#
# Estado atual do ambiente: revisões do Cloud Run, mensagens pendentes nas
# subscriptions e um smoke test HTTP.

source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

require_tools terraform gcloud curl
require_gcloud_auth

PROJECT="$(resolve_project)"
REGION="$(resolve_region)"

info "Serviços Cloud Run em ${PROJECT} / ${REGION}"
gcloud run services list \
  --project "${PROJECT}" --region "${REGION}" \
  --filter "metadata.name=(school-module OR finantial-module)" \
  --format "table(metadata.name:label=SERVICO,
                  status.url:label=URL,
                  spec.template.metadata.annotations['autoscaling.knative.dev/minScale']:label=MIN,
                  status.conditions[0].status:label=PRONTO)" 2>/dev/null \
  || warn "Não foi possível listar os serviços"

echo
info "Mensagens pendentes nas subscriptions"
for sub in "${SUBSCRIPTIONS[@]}" "${SUBSCRIPTIONS[@]/%/_DLQ}"; do
  if gcloud pubsub subscriptions describe "${sub}" --project "${PROJECT}" >/dev/null 2>&1; then
    printf '   %-34s existe\n' "${sub}"
  fi
done
echo "   (a contagem de não confirmadas está em Monitoring:"
echo "    métrica pubsub.googleapis.com/subscription/num_undelivered_messages)"

echo
info "Smoke test"
for svc in school-module finantial-module; do
  url="$(gcloud run services describe "${svc}" --project "${PROJECT}" --region "${REGION}" \
    --format 'value(status.url)' 2>/dev/null || true)"
  if [[ -z "${url}" ]]; then
    printf '   %-18s serviço não encontrado\n' "${svc}"
    continue
  fi
  code="$(curl -s -o /dev/null -w '%{http_code}' --max-time 15 "${url}/health" || echo "sem resposta")"
  printf '   %-18s GET /health → %s\n' "${svc}" "${code}"
done
