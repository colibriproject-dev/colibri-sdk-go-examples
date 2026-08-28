#!/usr/bin/env bash
#
# Zera o custo ocioso sem destruir nada.
#
#   ./pause.sh     # min_instance_count = 0 nos dois serviços
#   ./pause.sh -r  # volta para o valor do terraform.tfvars
#
# Com min_instance_count = 0 o Cloud Run desliga as instâncias quando não há
# tráfego e a cobrança por CPU/memória para. As URLs continuam funcionando: uma
# requisição sobe uma instância em alguns segundos.
#
# O preço disso é a mensageria: os consumidores do SDK fazem pull do Pub/Sub em
# goroutine de fundo, então com o serviço desligado ninguém consome as filas. As
# mensagens não se perdem — ficam retidas na subscription pelo período de
# retenção e são consumidas na próxima vez que uma instância subir — mas o
# ambiente deixa de processar em tempo real. Pause quando não estiver testando.

source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

require_tools terraform gcloud
require_gcloud_auth

RESUME=false
[[ "${1:-}" == "-r" || "${1:-}" == "--resume" ]] && RESUME=true

PROJECT="$(resolve_project)"
REGION="$(resolve_region)"

if [[ "${RESUME}" == true ]]; then
  target="$(tfvar min_instance_count 1)"
  info "Retomando: min_instance_count = ${target}"
else
  target=0
  info "Pausando: min_instance_count = 0"
fi

# Aplicado pelo Terraform e não pelo gcloud para o state não ficar divergente —
# um `terraform apply` seguinte reverteria uma alteração feita por fora.
tf apply -input=false -auto-approve -var "min_instance_count=${target}"

echo
if [[ "${RESUME}" == true ]]; then
  ok "Serviços de volta com instância sempre ligada — o Pub/Sub volta a ser consumido"
else
  ok "Serviços pausados: as URLs seguem respondendo sob demanda e a cobrança contínua para"
  warn "Enquanto pausado, as mensagens ficam acumuladas nas subscriptions."
  warn "Para retomar: ./pause.sh -r"
fi

warn "Lembre-se: o próximo ./deploy.sh usa o valor do terraform.tfvars,"
warn "então ele desfaz o pause. Ajuste a variável lá se quiser o pause permanente."
