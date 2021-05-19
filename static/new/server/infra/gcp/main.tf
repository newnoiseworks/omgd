variable "gcp_project" {
  type = string
}

variable "gcp_region" {
  type = string
}

variable "gcp_zone" {
  type = string
}

output "server_ip" {
  value = google_compute_instance.nakama_instance.network_interface[0].access_config[0].nat_ip
}

provider "google" {
  project = var.gcp_project
  region  = var.gcp_region
  zone    = var.gcp_zone
}

resource "google_compute_network" "vpc_network" {
  name                    = "nakama-instance-network"
  auto_create_subnetworks = "true"
}

resource "google_compute_firewall" "default" {
  name = "nakama-instance-firewall"
  network = google_compute_network.vpc_network.self_link

  allow {
    protocol = "tcp"
    ports = ["22", "80", "443", "7348-7351"]
  }

  allow {
    protocol = "udp"
    ports = ["7348-7351"]
  }

  target_tags = ["nakama"]
  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_instance" "nakama_instance" {
  name         = "nakama-instance"
  machine_type = "f1-micro"

  tags = ["nakama"]

  boot_disk {
    initialize_params {
      image = "cos-cloud/cos-89-lts"
    }
  }

  network_interface {
    network = google_compute_network.vpc_network.self_link
    access_config {
    }
  }
}
