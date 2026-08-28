#!/usr/bin/env bash
#
# Remove TODA a infraestrutura criada pelo deploy.sh e confere que não sobrou
# nada cobrando.
#
#   ./destroy.sh                  # destrói e audita (pede confirmação)
#   ./destroy.sh --auto-approve   # sem confirmação
#   ./destroy.sh --check          # só audita, não remove nada
#   ./destroy.sh --sweep          # remove por nome, sem depender do state
#
# O PostgreSQL e o Redis remotos NÃO são tocados: nunca foram criados aqui.
#
# Se o que você quer é só parar de pagar sem perder o ambiente, use ./pause.sh.

source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

AUTO_APPROVE=false
CHECK_ONLY=false
SWEEP=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --auto-approve) AUTO_APPROVE=true ;;
    --check)        CHECK_ONLY=true ;;
    --sweep)        SWEEP=true ;;
    -h|--help)      sed -n '2,16p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//'; exit 0 ;;
    *)              fail "Opção desconhecida: $1" ;;
  esac
  shift
done

require_tools gcloud
require_gcloud_auth

PROJECT="$(resolve_project)"
REGION="$(resolve_region)"
PREFIX="$(resolve_prefix)"

[[ -n "${PROJECT}" ]] || fail "project_id não definido."

info "Projeto: ${PROJECT}   Região: ${REGION}   Prefixo: ${PREFIX}"

# ─────────────────────────────────────────────────────────────────────────────
# Auditoria
# ─────────────────────────────────────────────────────────────────────────────

LEFTOVERS=0
BILLABLE=0

found() { # found <descrição> <custo: sim|nao>
  LEFTOVERS=$((LEFTOVERS + 1))
  if [[ "$2" == "sim" ]]; then
    BILLABLE=$((BILLABLE + 1))
    printf '   %s$%s %s\n' "${_c_red}" "${_c_reset}" "$1"
  else
    printf '   %s·%s %s\n' "${_c_yellow}" "${_c_reset}" "$1"
  fi
}

audit() {
  LEFTOVERS=0
  BILLABLE=0

  # Serviço Cloud Run com instância mínima acima de zero é o que cobra 24/7.
  for svc in school-module finantial-module; do
    min="$(gcloud run services describe "${svc}" --project "${PROJECT}" --region "${REGION}" \
      --format "value(spec.template.metadata.annotations['autoscaling.knative.dev/minScale'])" 2>/dev/null || true)"
    if gcloud run services describe "${svc}" --project "${PROJECT}" --region "${REGION}" \
        >/dev/null 2>&1; then
      if [[ -n "${min}" && "${min}" != "0" ]]; then
        found "Cloud Run ${svc} com ${min} instância(s) sempre ligada(s) — cobrança contínua" sim
      else
        found "Cloud Run ${svc} (escala a zero, cobra só por requisição)" nao
      fi
    fi
  done

  # Imagens no Artifact Registry cobram por armazenamento.
  if gcloud artifacts repositories describe "${PREFIX}" \
      --project "${PROJECT}" --location "${REGION}" >/dev/null 2>&1; then
    count="$(gcloud artifacts docker images list "${REGION}-docker.pkg.dev/${PROJECT}/${PREFIX}" \
      --project "${PROJECT}" --format 'value(IMAGE)' 2>/dev/null | wc -l | tr -d ' ')"
    if [[ "${count}" != "0" ]]; then
      found "Artifact Registry ${PREFIX} com ${count} imagem(ns)" sim
    else
      found "Artifact Registry ${PREFIX} vazio" nao
    fi
  fi

  # Objetos no bucket cobram por armazenamento.
  bucket="${PREFIX}-documents-${PROJECT}"
  if gcloud storage buckets describe "gs://${bucket}" --project "${PROJECT}" >/dev/null 2>&1; then
    objects="$(gcloud storage ls "gs://${bucket}/**" --project "${PROJECT}" 2>/dev/null | head -1 || true)"
    if [[ -n "${objects}" ]]; then
      found "Bucket gs://${bucket} com objetos" sim
    else
      found "Bucket gs://${bucket} vazio" nao
    fi
  fi

  # Pub/Sub não cobra parado, mas ocupa os nomes fixos que o código exige.
  for t in "${TOPICS[@]}" "${SUBSCRIPTIONS[@]/%/_DLQ}"; do
    gcloud pubsub topics describe "${t}" --project "${PROJECT}" >/dev/null 2>&1 \
      && found "Tópico ${t}" nao
  done
  for s in "${SUBSCRIPTIONS[@]}" "${SUBSCRIPTIONS[@]/%/_DLQ}"; do
    gcloud pubsub subscriptions describe "${s}" --project "${PROJECT}" >/dev/null 2>&1 \
      && found "Subscription ${s}" nao
  done

  # Secret Manager cobra por versão ativa (centavos, mas conta).
  while read -r secret; do
    [[ -n "${secret}" ]] && found "Secret ${secret}" sim
  done < <(gcloud secrets list --project "${PROJECT}" \
    --filter "name~${PREFIX}-" --format 'value(name)' 2>/dev/null)

  for sa in "${PREFIX}-school" "${PREFIX}-finantial"; do
    gcloud iam service-accounts describe "${sa}@${PROJECT}.iam.gserviceaccount.com" \
      --project "${PROJECT}" >/dev/null 2>&1 && found "Service account ${sa}" nao
  done
}

if [[ "${CHECK_ONLY}" == true ]]; then
  info "Auditando recursos do projeto"
  audit
  echo
  if [[ ${LEFTOVERS} -eq 0 ]]; then
    ok "Nada encontrado — nenhuma cobrança pendente deste ambiente"
  else
    warn "${LEFTOVERS} recurso(s) encontrado(s), ${BILLABLE} com custo"
    echo "   Para remover:   ./destroy.sh"
    echo "   Para só pausar: ./pause.sh"
  fi
  exit 0
fi

# ─────────────────────────────────────────────────────────────────────────────
# --sweep: remove por nome, para quando o state do Terraform não existe mais
# ─────────────────────────────────────────────────────────────────────────────

sweep() {
  warn "Modo sweep: removendo por nome, sem consultar o state do Terraform."

  for svc in school-module finantial-module; do
    gcloud run services delete "${svc}" --project "${PROJECT}" --region "${REGION}" --quiet \
      >/dev/null 2>&1 && ok "Cloud Run ${svc} removido"
  done

  # As subscriptions saem antes dos tópicos: o Pub/Sub recusa apagar um tópico
  # que ainda tem subscription anexada.
  for s in "${SUBSCRIPTIONS[@]}" "${SUBSCRIPTIONS[@]/%/_DLQ}"; do
    gcloud pubsub subscriptions delete "${s}" --project "${PROJECT}" --quiet \
      >/dev/null 2>&1 && ok "Subscription ${s} removida"
  done
  for t in "${TOPICS[@]}" "${SUBSCRIPTIONS[@]/%/_DLQ}"; do
    gcloud pubsub topics delete "${t}" --project "${PROJECT}" --quiet \
      >/dev/null 2>&1 && ok "Tópico ${t} removido"
  done

  bucket="${PREFIX}-documents-${PROJECT}"
  if gcloud storage buckets describe "gs://${bucket}" --project "${PROJECT}" >/dev/null 2>&1; then
    info "Esvaziando gs://${bucket} (inclui versões antigas)"
    gcloud storage rm --recursive --all-versions "gs://${bucket}/**" \
      --project "${PROJECT}" --quiet >/dev/null 2>&1 || true
    gcloud storage buckets delete "gs://${bucket}" --project "${PROJECT}" --quiet \
      >/dev/null 2>&1 && ok "Bucket removido"
  fi

  gcloud artifacts repositories delete "${PREFIX}" \
    --project "${PROJECT}" --location "${REGION}" --quiet >/dev/null 2>&1 \
    && ok "Artifact Registry ${PREFIX} removido"

  while read -r secret; do
    [[ -z "${secret}" ]] && continue
    gcloud secrets delete "${secret}" --project "${PROJECT}" --quiet >/dev/null 2>&1 \
      && ok "Secret ${secret} removido"
  done < <(gcloud secrets list --project "${PROJECT}" \
    --filter "name~${PREFIX}-" --format 'value(name)' 2>/dev/null)

  for sa in "${PREFIX}-school" "${PREFIX}-finantial"; do
    gcloud iam service-accounts delete "${sa}@${PROJECT}.iam.gserviceaccount.com" \
      --project "${PROJECT}" --quiet >/dev/null 2>&1 && ok "Service account ${sa} removida"
  done

  warn "As APIs habilitadas não são desligadas: elas não custam nada ligadas e"
  warn "desligá-las quebraria outros recursos do projeto que dependam delas."
}

if [[ "${SWEEP}" == true ]]; then
  if [[ "${AUTO_APPROVE}" != true ]]; then
    read -r -p "Remover por nome os recursos de '${PREFIX}' no projeto ${PROJECT}? [digite 'sim'] " answer
    [[ "${answer}" == "sim" ]] || fail "Cancelado"
  fi
  sweep
  echo
  info "Auditoria pós-sweep"
  audit
  [[ ${LEFTOVERS} -eq 0 ]] && ok "Nada sobrou" || warn "${LEFTOVERS} recurso(s) ainda presente(s)"
  exit 0
fi

# ─────────────────────────────────────────────────────────────────────────────
# Caminho normal: terraform destroy
# ─────────────────────────────────────────────────────────────────────────────

require_tools terraform

[[ -d "${TF_DIR}/.terraform" ]] || tf init -input=false >/dev/null

echo
warn "Isto remove os serviços Cloud Run, tópicos e subscriptions do Pub/Sub"
warn "(inclusive as de dead letter), bucket, Artifact Registry, secrets e as"
warn "service accounts. O PostgreSQL e o Redis remotos NÃO são afetados."
warn "Se você só quer parar de pagar sem perder o ambiente, use ./pause.sh."
echo

if [[ "${AUTO_APPROVE}" != true ]]; then
  read -r -p "Confirma a remoção do ambiente '${PREFIX}' no projeto ${PROJECT}? [digite 'sim'] " answer
  [[ "${answer}" == "sim" ]] || fail "Cancelado"
fi

info "terraform destroy"
if ! tf destroy -input=false -auto-approve; then
  warn "O destroy falhou. Causas comuns:"
  warn "  · deletion_protection ligada em algum serviço criado fora deste Terraform"
  warn "  · tópico do Pub/Sub com subscription que o state não conhece"
  warn "Se o state estiver perdido, use: ./destroy.sh --sweep"
  exit 1
fi

echo
info "Auditando o que sobrou"
audit

echo
if [[ ${LEFTOVERS} -eq 0 ]]; then
  ok "Ambiente removido — nenhuma cobrança pendente deste ambiente"
else
  warn "${LEFTOVERS} recurso(s) ainda presente(s), ${BILLABLE} com custo."
  warn "Para remover por nome: ./destroy.sh --sweep"
fi
