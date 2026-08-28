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
# As tabelas que as migrations criaram continuam lá.

source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

AUTO_APPROVE=false
CHECK_ONLY=false
SWEEP=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --auto-approve) AUTO_APPROVE=true ;;
    --check)        CHECK_ONLY=true ;;
    --sweep)        SWEEP=true ;;
    -h|--help)      sed -n '2,14p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//'; exit 0 ;;
    *)              fail "Opção desconhecida: $1" ;;
  esac
  shift
done

require_tools aws
require_aws_credentials

REGION="$(resolve_region)"
ACCOUNT_ID="$(aws_account_id)"

PROJECT="${PROJECT_NAME:-colibri-examples}"
if [[ -f "${TF_DIR}/terraform.tfvars" ]]; then
  from_tfvars="$(grep -E '^[[:space:]]*project_name[[:space:]]*=' "${TF_DIR}/terraform.tfvars" \
    | head -n1 | sed -E 's/.*=[[:space:]]*"?([^"[:space:]]+)"?.*/\1/' || true)"
  PROJECT="${from_tfvars:-${PROJECT}}"
fi

# Nomes fixos no código Go dos módulos — não levam o prefixo do projeto.
TOPICS=(SCHOOL_ENROLLMENT SCHOOL_COURSE SCHOOL_STUDENT FINANCIAL_INSTALLMENT)
QUEUES=(SCHOOL_ENROLLMENT_FINANCIAL SCHOOL_COURSE_FINANCIAL SCHOOL_STUDENT_FINANCIAL FINANCIAL_INSTALLMENT_SCHOOL)

info "Projeto: ${PROJECT}   Região: ${REGION}"

# ─────────────────────────────────────────────────────────────────────────────
# Auditoria: o que ainda existe e o que disso custa dinheiro
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

  # NAT Gateway e Elastic IP são o que mais dói se ficar esquecido.
  while read -r id; do
    [[ -n "${id}" && "${id}" != "None" ]] && found "NAT Gateway ${id} (~US\$ 32/mês)" sim
  done < <(aws ec2 describe-nat-gateways --region "${REGION}" \
    --filter "Name=tag:Project,Values=${PROJECT}" "Name=state,Values=available,pending" \
    --query 'NatGateways[].NatGatewayId' --output text 2>/dev/null | tr '\t' '\n')

  while read -r ip; do
    [[ -n "${ip}" && "${ip}" != "None" ]] && found "Elastic IP ${ip} (cobra enquanto não estiver associado)" sim
  done < <(aws ec2 describe-addresses --region "${REGION}" \
    --filters "Name=tag:Project,Values=${PROJECT}" \
    --query 'Addresses[].PublicIp' --output text 2>/dev/null | tr '\t' '\n')

  while read -r arn; do
    [[ -n "${arn}" && "${arn}" != "None" ]] && found "Load Balancer $(basename "${arn}") (~US\$ 17/mês)" sim
  done < <(aws elbv2 describe-load-balancers --region "${REGION}" \
    --query "LoadBalancers[?starts_with(LoadBalancerName, '${PROJECT}')].LoadBalancerArn" \
    --output text 2>/dev/null | tr '\t' '\n')

  # Tasks em execução são o custo por segundo do Fargate.
  if aws ecs describe-clusters --region "${REGION}" --clusters "${PROJECT}" \
      --query 'clusters[0].status' --output text 2>/dev/null | grep -q ACTIVE; then
    running="$(aws ecs list-tasks --region "${REGION}" --cluster "${PROJECT}" \
      --query 'length(taskArns)' --output text 2>/dev/null || echo 0)"
    if [[ "${running}" != "0" && "${running}" != "None" ]]; then
      found "Cluster ECS ${PROJECT} com ${running} task(s) rodando (Fargate por segundo)" sim
    else
      found "Cluster ECS ${PROJECT} sem tasks" nao
    fi
  fi

  # Objetos no S3 e imagens no ECR cobram por armazenamento.
  while read -r bucket; do
    [[ -z "${bucket}" ]] && continue
    objects="$(aws s3api list-objects-v2 --bucket "${bucket}" --max-items 1 \
      --query 'length(Contents)' --output text 2>/dev/null || echo 0)"
    if [[ "${objects}" != "0" && "${objects}" != "None" ]]; then
      found "Bucket S3 ${bucket} com objetos" sim
    else
      found "Bucket S3 ${bucket} vazio" nao
    fi
  done < <(aws s3api list-buckets --query "Buckets[?starts_with(Name, '${PROJECT}-documents')].Name" \
    --output text 2>/dev/null | tr '\t' '\n')

  for module in "${MODULES[@]}"; do
    repo="${PROJECT}/${module}"
    count="$(aws ecr list-images --region "${REGION}" --repository-name "${repo}" \
      --query 'length(imageIds)' --output text 2>/dev/null || echo "")"
    [[ -n "${count}" && "${count}" != "0" ]] && found "Repositório ECR ${repo} com ${count} imagem(ns)" sim
    [[ "${count}" == "0" ]] && found "Repositório ECR ${repo} vazio" nao
  done

  # Log groups cobram por armazenamento retido.
  while read -r lg; do
    [[ -n "${lg}" ]] && found "Log group ${lg}" nao
  done < <(aws logs describe-log-groups --region "${REGION}" \
    --log-group-name-prefix "/ecs/${PROJECT}" \
    --query 'logGroups[].logGroupName' --output text 2>/dev/null | tr '\t' '\n')

  # SNS/SQS não cobram parados, mas ocupam os nomes fixos que o código exige.
  for t in "${TOPICS[@]}"; do
    aws sns get-topic-attributes --region "${REGION}" \
      --topic-arn "arn:aws:sns:${REGION}:${ACCOUNT_ID}:${t}" >/dev/null 2>&1 \
      && found "Tópico SNS ${t}" nao
  done
  for q in "${QUEUES[@]}" "${QUEUES[@]/%/_DLQ}"; do
    aws sqs get-queue-url --region "${REGION}" --queue-name "${q}" >/dev/null 2>&1 \
      && found "Fila SQS ${q}" nao
  done

  while read -r vpc; do
    [[ -n "${vpc}" && "${vpc}" != "None" ]] && found "VPC ${vpc}" nao
  done < <(aws ec2 describe-vpcs --region "${REGION}" \
    --filters "Name=tag:Project,Values=${PROJECT}" \
    --query 'Vpcs[].VpcId' --output text 2>/dev/null | tr '\t' '\n')
}

# ─────────────────────────────────────────────────────────────────────────────
# --check: só audita
# ─────────────────────────────────────────────────────────────────────────────

if [[ "${CHECK_ONLY}" == true ]]; then
  info "Auditando recursos do projeto"
  audit
  echo
  if [[ ${LEFTOVERS} -eq 0 ]]; then
    ok "Nada encontrado — nenhuma cobrança pendente deste projeto"
  else
    warn "${LEFTOVERS} recurso(s) encontrado(s), ${BILLABLE} com custo"
    echo "   Rode ./destroy.sh para remover."
  fi
  exit 0
fi

# ─────────────────────────────────────────────────────────────────────────────
# --sweep: remove por nome, para quando o state do Terraform não existe mais
# ─────────────────────────────────────────────────────────────────────────────

sweep() {
  warn "Modo sweep: removendo por nome, sem consultar o state do Terraform."

  if aws ecs describe-clusters --region "${REGION}" --clusters "${PROJECT}" \
      --query 'clusters[0].status' --output text 2>/dev/null | grep -q ACTIVE; then
    for svc in school-module finantial-module; do
      aws ecs update-service --region "${REGION}" --cluster "${PROJECT}" \
        --service "${svc}" --desired-count 0 >/dev/null 2>&1 \
        && info "Serviço ${svc} zerado"
      aws ecs delete-service --region "${REGION}" --cluster "${PROJECT}" \
        --service "${svc}" --force >/dev/null 2>&1 \
        && ok "Serviço ${svc} removido"
    done
  fi

  # Esvaziar antes de apagar: o S3 recusa deletar bucket com objetos, e este
  # tem versionamento ligado, então as versões antigas também precisam sair.
  while read -r bucket; do
    [[ -z "${bucket}" ]] && continue
    info "Esvaziando bucket ${bucket}"
    aws s3 rm "s3://${bucket}" --recursive >/dev/null 2>&1 || true
    aws s3api delete-objects --bucket "${bucket}" --delete "$(aws s3api list-object-versions \
      --bucket "${bucket}" --output json \
      --query '{Objects: Versions[].{Key:Key,VersionId:VersionId}}' 2>/dev/null)" >/dev/null 2>&1 || true
    aws s3api delete-objects --bucket "${bucket}" --delete "$(aws s3api list-object-versions \
      --bucket "${bucket}" --output json \
      --query '{Objects: DeleteMarkers[].{Key:Key,VersionId:VersionId}}' 2>/dev/null)" >/dev/null 2>&1 || true
    aws s3api delete-bucket --bucket "${bucket}" --region "${REGION}" >/dev/null 2>&1 \
      && ok "Bucket ${bucket} removido"
  done < <(aws s3api list-buckets --query "Buckets[?starts_with(Name, '${PROJECT}-documents')].Name" \
    --output text 2>/dev/null | tr '\t' '\n')

  for module in "${MODULES[@]}"; do
    aws ecr delete-repository --region "${REGION}" \
      --repository-name "${PROJECT}/${module}" --force >/dev/null 2>&1 \
      && ok "Repositório ECR ${PROJECT}/${module} removido"
  done

  for t in "${TOPICS[@]}"; do
    aws sns delete-topic --region "${REGION}" \
      --topic-arn "arn:aws:sns:${REGION}:${ACCOUNT_ID}:${t}" >/dev/null 2>&1 \
      && ok "Tópico ${t} removido"
  done

  for q in "${QUEUES[@]}" "${QUEUES[@]/%/_DLQ}"; do
    url="$(aws sqs get-queue-url --region "${REGION}" --queue-name "${q}" \
      --query QueueUrl --output text 2>/dev/null || true)"
    [[ -n "${url}" && "${url}" != "None" ]] \
      && aws sqs delete-queue --region "${REGION}" --queue-url "${url}" >/dev/null 2>&1 \
      && ok "Fila ${q} removida"
  done

  while read -r lg; do
    [[ -z "${lg}" ]] && continue
    aws logs delete-log-group --region "${REGION}" --log-group-name "${lg}" >/dev/null 2>&1 \
      && ok "Log group ${lg} removido"
  done < <(aws logs describe-log-groups --region "${REGION}" \
    --log-group-name-prefix "/ecs/${PROJECT}" \
    --query 'logGroups[].logGroupName' --output text 2>/dev/null | tr '\t' '\n')

  warn "ALB, NAT Gateway, EIPs e a VPC não são removidos no sweep: a ordem de"
  warn "dependência entre eles é traiçoeira e o Terraform faz isso melhor."
  warn "Se sobrarem, remova pelo console — o ./destroy.sh --check lista quais são."
}

if [[ "${SWEEP}" == true ]]; then
  if [[ "${AUTO_APPROVE}" != true ]]; then
    read -r -p "Remover por nome os recursos do projeto '${PROJECT}' em ${REGION}? [digite 'sim'] " answer
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
warn "Isto remove VPC, ALB, cluster ECS, tópicos SNS, filas SQS (e DLQs),"
warn "bucket S3, repositórios ECR, log groups e parâmetros do SSM."
warn "O PostgreSQL e o Redis remotos NÃO são afetados."
echo

if [[ "${AUTO_APPROVE}" != true ]]; then
  read -r -p "Confirma a remoção do ambiente '${PROJECT}' em ${REGION}? [digite 'sim'] " answer
  [[ "${answer}" == "sim" ]] || fail "Cancelado"
fi

# Zerar os serviços antes acelera bastante o destroy: sem isso o Terraform
# espera as tasks drenarem enquanto já tenta remover o target group.
info "Zerando os serviços ECS"
for svc in school-module finantial-module; do
  aws ecs update-service --region "${REGION}" --cluster "${PROJECT}" \
    --service "${svc}" --desired-count 0 >/dev/null 2>&1 \
    && info "  ${svc} → 0 tasks" || true
done

info "terraform destroy"
tf_args=(-input=false -auto-approve)
if ! tf destroy "${tf_args[@]}"; then
  warn "O destroy falhou. Causa mais comum: uma ENI ainda presa a uma task"
  warn "que estava drenando. Aguarde ~2 minutos e rode de novo."
  warn "Se o state estiver perdido, use: ./destroy.sh --sweep"
  exit 1
fi

echo
info "Auditando o que sobrou"
audit

echo
if [[ ${LEFTOVERS} -eq 0 ]]; then
  ok "Ambiente removido — nenhuma cobrança pendente deste projeto"
else
  warn "${LEFTOVERS} recurso(s) ainda presente(s), ${BILLABLE} com custo."
  warn "Recursos criados fora do Terraform não entram no destroy."
  warn "Para remover por nome: ./destroy.sh --sweep"
fi
