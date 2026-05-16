output "public_ip" {
  description = "Public IP address of the VM"
  value       = google_compute_address.static_ip.address
}

output "frontend_url" {
  description = "Frontend URL"
  value       = "http://${google_compute_address.static_ip.address}"
}

output "grafana_url" {
  description = "Grafana dashboard URL"
  value       = "http://${google_compute_address.static_ip.address}:3000"
}

output "prometheus_url" {
  description = "Prometheus UI URL"
  value       = "http://${google_compute_address.static_ip.address}:9090"
}

output "ssh_command" {
  description = "SSH command to connect to the VM"
  value       = "ssh ${var.ssh_user}@${google_compute_address.static_ip.address}"
}
