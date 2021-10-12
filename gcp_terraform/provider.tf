terraform {
  required_providers {
    google      = {
      source  = "hashicorp/google"
      version = "~> 3.87.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 3.87.0"
    }
  }
}

locals {
  project_id  = "talkiewalkie-305117"
  region      = "europe-west1"
  domain_name = "talkiewalkie.app"
}

provider "google" {
  project = local.project_id
  region  = local.region
  zone    = "eu-west1-b"
}

provider "google-beta" {
  project = local.project_id
  region  = local.region
  zone    = "eu-west1-b"
}
