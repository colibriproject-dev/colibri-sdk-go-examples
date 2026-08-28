#!/usr/bin/env bash
#
# Exercita o caminho completo do ambiente: cria um curso e um aluno no
# school-module, faz a matrícula (que publica em SCHOOL_ENROLLMENT) e confere se
# o finantial-module consumiu a mensagem e gerou a conta.
#
# Valida de uma vez a API pública, a estrutura de mensageria SNS/SQS e a conexão
# com o banco remoto. A última etapa força ainda a chamada HTTP entre os dois
# serviços, que é o que produz o trace distribuído no X-Ray.
#
# Todas as rotas de negócio ficam sob /public/ — é o prefixo que os controllers
# declaram com restserver.PublicApi. Só /health e /api-docs respondem na raiz.

source "$(dirname "${BASH_SOURCE[0]}")/common.sh"

require_tools terraform curl jq
require_aws_credentials

SCHOOL="$(tf output -raw school_module_url)"
FINANTIAL="$(tf output -raw finantial_module_url)"

info "school-module ....: ${SCHOOL}"
info "finantial-module .: ${FINANTIAL}"

info "1/6 — health check dos dois serviços"
curl -fsS --max-time 10 "${SCHOOL}/health"    | jq -c . || fail "school-module não respondeu"
curl -fsS --max-time 10 "${FINANTIAL}/health" | jq -c . || fail "finantial-module não respondeu"
ok "Ambos no ar"

info "2/6 — criando um curso"
course="$(curl -fsS -X POST "${SCHOOL}/public/v1/courses" \
  -H 'Content-Type: application/json' \
  -d '{"name":"Engenharia de Software","value":1500.00}')" \
  || fail "Falha ao criar o curso"
echo "   ${course}"
course_id="$(echo "${course}" | jq -r '.id')"
[[ "${course_id}" != "null" && -n "${course_id}" ]] || fail "Resposta sem id de curso"
ok "Curso ${course_id}"

info "3/6 — criando um aluno"
student="$(curl -fsS -X POST "${SCHOOL}/public/v1/students" \
  -H 'Content-Type: application/json' \
  -d '{"name":"Aluno Smoke Test","email":"smoke.test@exemplo.com","birthday":"2000-01-01"}')" \
  || fail "Falha ao criar o aluno"
echo "   ${student}"
student_id="$(echo "${student}" | jq -r '.id')"
[[ "${student_id}" != "null" && -n "${student_id}" ]] || fail "Resposta sem id de aluno"
ok "Aluno ${student_id}"

info "4/6 — matriculando (publica em SCHOOL_ENROLLMENT)"
curl -fsS -X POST "${SCHOOL}/public/v1/enrollments" \
  -H 'Content-Type: application/json' \
  -d "{\"studentId\":\"${student_id}\",\"courseId\":\"${course_id}\",\"installments\":12}" \
  || fail "Falha ao criar a matrícula"
echo
ok "Matrícula publicada no tópico"

info "5/6 — aguardando o finantial-module consumir SCHOOL_ENROLLMENT_FINANCIAL"
found=false
for attempt in $(seq 1 12); do
  sleep 5
  accounts="$(curl -fsS --max-time 10 "${FINANTIAL}/public/accounts" || echo '[]')"
  if echo "${accounts}" | jq -e --arg sid "${student_id}" \
      'if type == "array" then any(.[]; .studentId == $sid) else false end' >/dev/null 2>&1; then
    ok "Conta criada a partir da mensagem (tentativa ${attempt})"
    echo "${accounts}" | jq -c --arg sid "${student_id}" '[.[] | select(.studentId == $sid)]'
    found=true
    break
  fi
  printf '   tentativa %2d/12...\n' "${attempt}"
done

if [[ "${found}" != true ]]; then
  warn "A conta não apareceu em 60s. O que checar:"
  warn "  logs do consumidor : aws logs tail /ecs/\$(tf output -raw ecs_cluster_name)/finantial-module --follow"
  warn "  mensagens travadas : fila SCHOOL_ENROLLMENT_FINANCIAL_DLQ"
  exit 1
fi

# GET de um curso por id faz o school-module chamar GET /public/accounts no
# finantial-module. É a chamada que gera o trace com spans dos dois serviços.
info "6/6 — chamada HTTP entre os serviços (gera o trace distribuído)"
curl -fsS --max-time 15 "${SCHOOL}/public/v1/courses/${course_id}" | jq -c . \
  || fail "Falha na chamada school-module → finantial-module"
ok "Trace distribuído gerado"

echo
ok "Ambiente validado de ponta a ponta"
echo "   Mapa de serviços no X-Ray: $(tf output -raw xray_console_url)"
