# Create a list of resource names
variable "instance_names" {
  type    = list(string)
  default = ["instance-1", "instance-2", "instance-3"]
}

# Create a GCP network
resource "google_compute_network" "vpc_network" {
  name                    = "my-vpc-network"
  auto_create_subnetworks = false
}

# Create multiple subnets using a for_each loop
resource "google_compute_subnetwork" "subnet" {
  count = length(var.instance_names)

  name          = "subnet-${count.index + 1}"
  ip_cidr_range = "10.0.${count.index}.0/24"
  region        = "us-central1"
  network       = google_compute_network.vpc_network.id
}

# Create multiple instances using a count loop
resource "google_compute_instance" "vm_instance" {
  count = length(var.instance_names)

  name         = var.instance_names[count.index]
  machine_type = "e2-medium"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-12"
    }
  }

  network_interface {
    network    = google_compute_network.vpc_network.id
    subnetwork = google_compute_subnetwork.subnet[count.index].id
  }

  tags = ["web", "dev"]

  metadata = {
    ssh-keys = "user:ssh-rsa XXXXXXX"
  }
}

# Create multiple firewall rules using for_each
resource "google_compute_firewall" "firewall" {
  for_each = {
    "allow-http" : "80"
    "allow-https" : "443"
  }

  name    = each.key
  network = google_compute_network.vpc_network.name

  allow {
    protocol = "tcp"
    ports    = [each.value]
  }

  source_ranges = ["0.0.0.0/0"]
}
