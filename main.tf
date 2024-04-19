# main.tf

terraform {
  required_version = ">= 1.3"

  required_providers {
    google = ">= 5.2"
  }
}

provider "google" {
  project = "websites-394411"
}

data "google_project" "project" {
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_v2_service" "run_service" {
  name     = "activities"
  location = "europe-west4"
  ingress  = "INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER"

  template {
    containers {
      image = "europe-west4-docker.pkg.dev/websites-394411/webstekjes/activities:2.0.0-alpha-10"

      env {
        name  = "PROJECTID"
        value = "websites-394411"
        # value = data.google_project.project.number
      }
      env {
        name = "PASSPHRASE"
        value_source {
          secret_key_ref {
            secret  = data.google_secret_manager_secret.passphrase.secret_id
            version = "1"
          }
        }
      }
      env {
        name = "SIGNINGKEY"
        value_source {
          secret_key_ref {
            secret  = data.google_secret_manager_secret.signingkey.secret_id
            version = "1"
          }
        }
      }
    }
  }


  traffic {
    percent = 100
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
  }


}


resource "google_cloud_run_service_iam_policy" "noauth" {
  location = google_cloud_run_v2_service.run_service.location
  project  = google_cloud_run_v2_service.run_service.project
  service  = google_cloud_run_v2_service.run_service.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

data "google_secret_manager_secret" "passphrase" {
  secret_id = "passphrase"
}

data "google_secret_manager_secret" "signingkey" {
  secret_id = "signingkey"
}

resource "google_secret_manager_secret_iam_member" "passphrase-access" {
  secret_id  = data.google_secret_manager_secret.passphrase.id
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"
  depends_on = [data.google_secret_manager_secret.passphrase]
}

resource "google_secret_manager_secret_iam_member" "signingkey-access" {
  secret_id  = data.google_secret_manager_secret.signingkey.id
  role       = "roles/secretmanager.secretAccessor"
  member     = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"
  depends_on = [data.google_secret_manager_secret.signingkey]
}
