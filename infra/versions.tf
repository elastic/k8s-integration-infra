terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">=3.52.0"
    }

    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.0.1"
    }

    kubectl = {
      source  = "gavinbunney/kubectl"
      version = ">= 1.7.0"
    }

    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.4.1"
    }

    external = {
      source  = "hashicorp/external"
      version = "~> 2.2.2"
    }
  }

  required_version = ">= 0.14"
}

