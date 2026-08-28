# Bucket usado pelo school-module no upload de documentos do aluno
# (STORAGE_BUCKET / storage.Initialize()).
resource "aws_s3_bucket" "documents" {
  bucket        = "${local.name}-documents-${local.account_id}"
  force_destroy = true

  tags = { Name = "${local.name}-documents" }
}

# O SDK acessa o bucket com as credenciais da role da task; nada precisa ficar
# público. O upload usa path style, que funciona igual com o bloqueio ligado.
resource "aws_s3_bucket_public_access_block" "documents" {
  bucket = aws_s3_bucket.documents.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_server_side_encryption_configuration" "documents" {
  bucket = aws_s3_bucket.documents.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_versioning" "documents" {
  bucket = aws_s3_bucket.documents.id

  versioning_configuration {
    status = "Enabled"
  }
}
