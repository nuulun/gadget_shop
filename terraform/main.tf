terraform {
  required_version = ">= 1.5.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}


resource "google_compute_address" "static_ip" {
  name   = "${var.app_name}-ip"
  region = var.region
}


resource "google_compute_firewall" "allow_app" {
  name    = "${var.app_name}-allow-app"
  network = "default"

  allow {
    protocol = "tcp"
    ports    = ["80", "3000", "22", "9090", "30080", "30090", "30300", "30800"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = [var.app_name]
}


resource "google_compute_instance" "vm" {
  name         = "${var.app_name}-vm"
  machine_type = var.machine_type
  zone         = var.zone

  tags = [var.app_name]

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2204-lts"
      size  = 30
      type  = "pd-ssd"
    }
  }

  network_interface {
    network = "default"
    access_config {
      nat_ip = google_compute_address.static_ip.address
    }
  }

  metadata = {
    ssh-keys = "${var.ssh_user}:${var.ssh_public_key}"
  }

  metadata_startup_script = <<-EOF
    #!/bin/bash
    set -e

    curl -fsSL https://get.docker.com | sh
    usermod -aG docker ${var.ssh_user}

    mkdir -p /usr/local/lib/docker/cli-plugins
    curl -SL https://github.com/docker/compose/releases/latest/download/docker-compose-linux-x86_64 \
      -o /usr/local/lib/docker/cli-plugins/docker-compose
    chmod +x /usr/local/lib/docker/cli-plugins/docker-compose

    systemctl enable docker
    systemctl start docker

    git clone https://github.com/nuulun/gadget_shop.git /home/${var.ssh_user}/gadget_shop
    chown -R ${var.ssh_user}:${var.ssh_user} /home/${var.ssh_user}/gadget_shop

    curl -sfL https://get.k3s.io | sh -
    systemctl enable k3s
    systemctl start k3s
    chmod 644 /etc/rancher/k3s/k3s.yaml
  EOF
}
