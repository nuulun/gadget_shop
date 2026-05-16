variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
  default     = "us-central1"
}

variable "zone" {
  description = "GCP zone"
  type        = string
  default     = "us-central1-a"
}

variable "app_name" {
  description = "Application name prefix for all resources"
  type        = string
  default     = "microservices"
}

variable "machine_type" {
  description = "GCP Compute Engine machine type"
  type        = string
  default     = "e2-medium"
}

variable "ssh_user" {
  description = "SSH username for the VM"
  type        = string
  default     = "ubuntu"
}

variable "ssh_public_key" {
  description = "SSH public key content (paste the content of your .pub file)"
  type        = string
}
