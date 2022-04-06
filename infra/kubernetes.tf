variable "es_password" {
  default     = "changeme"
  type        = string
  description = "The password of elasticsearch."
}

variable "es_user" {
  default     = "elastic"
  type        = string
  description = "The username of elasticsearch."
}

variable "es_host" {
  default     = "elasticsearch:9200"
  type        = string
  description = "The url:port of elasticsearch."
}

variable "namespace" {
  default     = "default"
  type        = string
  description = "The namespace to deploy beats and agent"
}

data "google_client_config" "default" {}

provider "kubernetes" {
  host                   = "https://${google_container_cluster.primary.endpoint}"
  token                  = data.google_client_config.default.access_token
  cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)
}


resource "kubernetes_secret" "elasticsearch-master-credentials" {
  metadata {
    name = "elasticsearch-credentials"
    namespace = var.namespace
  }

  data = {
    username = var.es_user
    password = var.es_password
  }

}

resource "kubernetes_secret" "elasticsearch-master-url" {
  metadata {
    name = "elasticsearch-host"
    namespace = var.namespace
  }

  data = {
    host = var.es_host
  }

}