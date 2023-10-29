variable "gcp_project" {
  type = string
}

variable "gcp_region" {
  type = string
}

variable "gcp_zone" {
  type = string
}

variable "gcp_type" {
  type = string
}

variable "project" {
  type = string
}

variable "profile" {
  type = string
}

output "server_ip" {
  value = google_compute_instance.omgd_dev_instance.network_interface[0].access_config[0].nat_ip
}

provider "google" {
  project = var.gcp_project
  region  = var.gcp_region
  zone    = var.gcp_zone
}

terraform {
  backend "gcs" {
  }
}

resource "google_compute_instance" "omgd_dev_instance" {
  name         = "${var.project}-omgd-dev-instance-${var.profile}"
  machine_type = var.gcp_type

  tags = ["omgd", "nakama"]

  boot_disk {
    initialize_params {
      image = "cos-cloud/cos-109-lts"
    }
  }

  network_interface {
    network = "https://www.googleapis.com/compute/v1/projects/${var.gcp_project}/global/networks/${var.project}-omgd-dev-instance-network"

    access_config {
    }
  }
}
