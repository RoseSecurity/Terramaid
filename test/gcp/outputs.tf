output "instance_ips" {
  value = google_compute_instance.vm_instance[*].network_interface[0].access_config[0].nat_ip
}

output "subnet_ids" {
  value = google_compute_subnetwork.subnet[*].id
}

output "firewall_rule_names" {
  value = google_compute_firewall.firewall[*].name
}
