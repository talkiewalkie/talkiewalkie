/*

IMPORTANT: two things are not handled by terraform at the moment:
- http2 connection needs to be enabled from the console
- domain mapping needs to be enabled from the console

*/

resource "google_cloud_run_service" "grpc_test" {
  name     = "grpc-hello"
  location = local.region

  template {
    spec {
      containers {
        image = "gcr.io/talkiewalkie-305117/helloworld-grpc"
        ports {
          protocol       = "TCP"
          container_port = 50051
        }
        env {
          name  = "POSTGRES_HOST"
          value = "postgres"
        }
        #        env {
        #          name = "POSTGRES_PASSWORD"
        #          value_from {
        #            secret_key_ref {
        #              key  = ""
        #              name = ""
        #            }
        #          }
        #        }
      }
    }
  }


  traffic {
    percent         = 100
    latest_revision = true
  }
}


data "google_iam_policy" "admin" {
  binding {
    role    = "roles/run.invoker"
    members = ["allUsers"]
  }
}

resource "google_cloud_run_service_iam_policy" "policy" {
  location    = google_cloud_run_service.grpc_test.location
  project     = google_cloud_run_service.grpc_test.project
  service     = google_cloud_run_service.grpc_test.name
  policy_data = data.google_iam_policy.admin.policy_data
  depends_on  = [google_cloud_run_service.grpc_test]
}