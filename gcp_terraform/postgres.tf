/*

IMPORTANT: two things are not handled by terraform at the moment:
- http2 connection needs to be enabled from the console
- domain mapping needs to be enabled from the console

*/

resource "google_secret_manager_secret" "postgres_password" {
  secret_id = "postgres_password"

  replication {
    automatic = true
  }
}

resource "random_password" "database_password" {
  length  = 32
  special = false
}
#
#resource "google_secret_manager_secret_version" "initial_postgres" {
#  secret      = google_secret_manager_secret.postgres_password.id
#  secret_data = random_password.database_password.result
#}
#
#data "google_secret_manager_secret" "postgres_password_data" {
#  secret_id = google_secret_manager_secret.postgres_password.secret_id
#}
#
#resource "google_cloud_run_service" "postgres" {
#  provider = google-beta
#  name     = "postgres"
#  location = local.region
#
#  metadata {
#    annotations = {
#      "run.googleapis.com/launch-stage" = "BETA"
#    }
#  }
#
#  template {
#    spec {
#      containers {
#        image = "postgres:13"
#        ports {
#          protocol       = "TCP"
#          container_port = 5432
#        }
#        env {
#          name  = "POSTGRES_USER"
#          value = "postgres"
#        }
#        env {
#          name = "POSTGRES_PASSWORD"
#          value_from {
#            secret_key_ref {
#              key  = google_secret_manager_secret.postgres_password.secret_id
#              name = "latest"
#            }
#          }
#        }
#        #        env {
#        #          name  = "POSTGRES_PASSWORD"
#        #          value = random_password.database_password.result
#        #        }
#        env {
#          name  = "POSTGRES_DB"
#          value = "talkiewalkie"
#        }
#      }
#    }
#  }
#
#  traffic {
#    percent         = 100
#    latest_revision = true
#  }
#}
