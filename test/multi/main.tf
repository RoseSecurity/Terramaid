data "cloudsmith_organization" "my_organization" {
  slug = var.cloudsmith_org
}

resource "cloudsmith_repository" "my_repository" {
  description = "A certifiably-awesome private package repository"
  name        = "My Repository"
  namespace   = data.cloudsmith_organization.my_organization.slug_perm
  slug        = "test-repo"
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_instance" "web" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t3.micro"

  tags = {
    Name = "HelloWorld"
  }
}

resource "google_healthcare_dataset" "dataset" {
  location = "us-central1"
  name     = "my-dataset"
}

resource "google_healthcare_consent_store" "my-consent" {
  dataset = google_healthcare_dataset.dataset.id
  name    = var.consent_store_name
}
