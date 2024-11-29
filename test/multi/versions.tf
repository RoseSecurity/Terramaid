terraform {
  required_version = ">= 1.5.5"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.0"
    }
    cloudsmith = {
      source  = "cloudsmith-io/cloudsmith"
      version = ">= 0.0.40"
    }
    gcp = {
      source  = "hashicorp/google"
      version = ">= 6.0"
    }
  }
}

