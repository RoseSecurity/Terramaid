# Create a resource group
resource "azurerm_resource_group" "rg" {
  name     = "my-resource-group"
  location = var.location
}

# Create a virtual network
resource "azurerm_virtual_network" "vnet" {
  name                = "my-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = var.location
  resource_group_name = azurerm_resource_group.rg.name
}

# Create multiple subnets using count
resource "azurerm_subnet" "subnet" {
  count                = length(var.vm_names)
  name                 = "subnet-${count.index + 1}"
  resource_group_name  = azurerm_resource_group.rg.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = ["10.0.${count.index}.0/24"]
}

# Create multiple virtual machines using count
resource "azurerm_virtual_machine" "vm" {
  count               = length(var.vm_names)
  name                = var.vm_names[count.index]
  resource_group_name = azurerm_resource_group.rg.name
  vm_size             = "Standard_DS1_v2"
  location            = var.location

  network_interface_ids = [
    azurerm_network_interface.nic[count.index].id,
  ]

  storage_os_disk {
    name          = "${var.vm_names[count.index]}-osdisk"
    caching       = "ReadWrite"
    create_option = "FromImage"
    disk_size_gb  = 30
  }

  storage_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }

  os_profile {
    computer_name  = var.vm_names[count.index]
    admin_username = "adminuser"
    admin_password = "P@ssw0rd1234!"
  }

  os_profile_windows_config {}
}

# Create multiple network interfaces using count
resource "azurerm_network_interface" "nic" {
  count               = length(var.vm_names)
  name                = "nic-${count.index + 1}"
  location            = var.location
  resource_group_name = azurerm_resource_group.rg.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.subnet[count.index].id
    private_ip_address_allocation = "Dynamic"
  }
}

# Create multiple security rules using for_each
resource "azurerm_network_security_group" "nsg" {
  name                = "my-nsg"
  location            = var.location
  resource_group_name = azurerm_resource_group.rg.name

  security_rule {
    for_each = {
      "allow-http" : 80
      "allow-https" : 443
    }

    name                       = each.key
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_ranges    = [each.value]
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

