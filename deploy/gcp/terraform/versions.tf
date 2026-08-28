terraform {
  required_version = ">= 1.6.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 6.0, < 8.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region

  default_labels = {
    project     = var.name_prefix
    environment = var.environment
    managed_by  = "terraform"
    repository  = "colibri-sdk-go-examples"
  }
}
