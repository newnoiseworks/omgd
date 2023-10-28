variable "gcp_project" {
  type = string
}

variable "gcp_region" {
  type = string
}

variable "gcp_zone" {
  type = string
}

variable "project" {
  type = string
}

variable "profile" {
  type = string
}

provider "google" {
  project = var.gcp_project
  region  = var.gcp_region
  zone    = var.gcp_zone
}

terraform {
  backend "local" {}
}

output "bucket_name" {
  value = google_storage_bucket.default.name
}

resource "google_compute_network" "vpc_network" {
  name                    = "${var.project}-nakama-instance-network"
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

resource "random_id" "bucket_postfix" {
    byte_length = 8
}

resource "google_storage_bucket" "default" {
  name          = "${var.project}-bucket-tfstate-${random_id.bucket_postfix.hex}"
	project 			= var.gcp_project
  force_destroy = true
  location      = "US"
  storage_class = "STANDARD"
  versioning {
    enabled = true
  }
}
