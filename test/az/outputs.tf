output "vm_public_ips" {
  value = [for nic in azurerm_network_interface.nic : nic.private_ip_address]
}

output "subnet_ids" {
  value = azurerm_subnet.subnet[*].id
}

output "nsg_name" {
  value = azurerm_network_security_group.nsg.name
}

