variable "vm_names" {
  description = "A list of names for the virtual machines."
  type        = list(string)
  default     = ["vm-1", "vm-2", "vm-3"]
}

variable "location" {
  description = "The Azure region where resources will be created."
  type        = string
  default     = "East US"
}

variable "client_id" {
  description = "The client ID for the Azure Service Principal."
  type        = string
}

variable "client_secret" {
  description = "The client secret for the Azure Service Principal."
  type        = string
}

variable "tenant_id" {
  description = "The tenant ID for the Azure subscription."
  type        = string
}

variable "subscription_id" {
  description = "The subscription ID for the Azure subscription."
  type        = string
}

