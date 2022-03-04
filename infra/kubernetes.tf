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

data "google_client_config" "default" {}

provider "kubernetes" {
  host                   = "https://${google_container_cluster.primary.endpoint}"
  token                  = data.google_client_config.default.access_token
  cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)
}


resource "kubernetes_secret" "elasticsearch-master-credentials" {
  metadata {
    name = "elasticsearch-credentials"
  }

  data = {
    username = var.es_user
    password = var.es_password
  }

}

resource "kubernetes_secret" "elasticsearch-master-url" {
  metadata {
    name = "elasticsearch-host"
  }

  data = {
    host = var.es_host
  }

}