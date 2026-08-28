data "google_project" "current" {
  project_id = var.project_id
}

locals {
  name           = var.name_prefix
  project_number = data.google_project.current.number

  # ── Nomes de tópicos e subscriptions ───────────────────────────────────────
  # ATENÇÃO: estes nomes são literais no código Go dos dois módulos
  # (messaging.NewProducer("SCHOOL_COURSE"), QueueName() == "SCHOOL_COURSE_FINANCIAL").
  # O SDK chama client.Publisher(nome) e client.Subscriber(nome) direto, então
  # NÃO é possível prefixar nem sufixar nada aqui sem alterar o código.
  # Consequência: um projeto GCP comporta só um ambiente deste stack.
  topics = {
    school_enrollment     = "SCHOOL_ENROLLMENT"
    school_course         = "SCHOOL_COURSE"
    school_student        = "SCHOOL_STUDENT"
    financial_installment = "FINANCIAL_INSTALLMENT"
  }

  # Cada subscription consome exatamente um tópico. O valor é a CHAVE do mapa
  # de tópicos acima.
  subscriptions = {
    "SCHOOL_ENROLLMENT_FINANCIAL"  = "school_enrollment"
    "SCHOOL_COURSE_FINANCIAL"      = "school_course"
    "SCHOOL_STUDENT_FINANCIAL"     = "school_student"
    "FINANCIAL_INSTALLMENT_SCHOOL" = "financial_installment"
  }

  # Quem publica em quê e quem consome o quê — usado para dar a cada service
  # account só os papéis de que ela precisa.
  school_publishes    = ["school_enrollment", "school_course", "school_student"]
  finantial_publishes = ["financial_installment"]

  school_subscribes    = ["FINANCIAL_INSTALLMENT_SCHOOL"]
  finantial_subscribes = ["SCHOOL_ENROLLMENT_FINANCIAL", "SCHOOL_COURSE_FINANCIAL", "SCHOOL_STUDENT_FINANCIAL"]

  # O Cloud Run injeta PORT com o valor de container_port e recusa a variável se
  # ela for declarada à mão. Como cada serviço tem URL própria em :443, a porta
  # interna não precisa ser diferente entre eles.
  container_port = 8080

  # Sidecar do collector escutando OTLP/HTTP dentro da mesma instância.
  otlp_endpoint = "localhost:4318"

  repository_url = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.modules.repository_id}"

  # O service agent do Pub/Sub é quem encaminha para o dead letter topic e quem
  # confirma a mensagem original — precisa de papel nos dois lados.
  pubsub_service_agent = "serviceAccount:service-${local.project_number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}
