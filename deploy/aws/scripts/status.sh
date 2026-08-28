#!/usr/bin/env bash
#
# Estado atual do ambiente: tasks rodando, saúde no ALB e um smoke test HTTP.

source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

require_tools terraform aws curl
require_aws_credentials

REGION="$(resolve_region)"
CLUSTER="$(tf output -raw ecs_cluster_name)"
ALB="$(tf output -raw alb_dns_name)"

info "Cluster ${CLUSTER} (${REGION})"

aws ecs describe-services \
  --region "${REGION}" \
  --cluster "${CLUSTER}" \
  --services school-module finantial-module \
  --query 'services[].{Servico:serviceName,Desejado:desiredCount,Rodando:runningCount,Pendente:pendingCount,Status:status}' \
  --output table

echo
info "Saúde dos targets no ALB"
PROJECT="$(tf output -raw ecs_cluster_name)"
for tg in school finantial; do
  arn="$(aws elbv2 describe-target-groups --region "${REGION}" \
    --query "TargetGroups[?TargetGroupName=='${PROJECT}-${tg}'].TargetGroupArn | [0]" --output text)"
  if [[ "${arn}" != "None" && -n "${arn}" ]]; then
    state="$(aws elbv2 describe-target-health --region "${REGION}" --target-group-arn "${arn}" \
      --query 'TargetHealthDescriptions[].TargetHealth.State' --output text)"
    printf '   %-18s %s\n' "${tg}" "${state:-sem targets registrados}"
  fi
done

echo
info "Smoke test"
for entry in "school-module:8080" "finantial-module:8081"; do
  name="${entry%%:*}"; port="${entry##*:}"
  code="$(curl -s -o /dev/null -w '%{http_code}' --max-time 10 "http://${ALB}:${port}/health" || echo "sem resposta")"
  printf '   %-18s GET :%s/health → %s\n' "${name}" "${port}" "${code}"
done
