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

resource "kubernetes_deployment" "nginx" {
  metadata {
    name = "scalable-nginx-example"
    labels = {
      App = "ScalableNginxExample"
    }
  }

  spec {
    replicas = 1
    selector {
      match_labels = {
        App = "ScalableNginxExample"
      }
    }
    template {
      metadata {
        labels = {
          App = "ScalableNginxExample"
        }
      }
      spec {
        container {
          image = "nginx:1.7.8"
          name  = "example"

          port {
            container_port = 80
          }

          resources {
            limits = {
              cpu    = "0.5"
              memory = "512Mi"
            }
            requests = {
              cpu    = "250m"
              memory = "50Mi"
            }
          }
        }
      }
    }
  }
}


resource "kubernetes_service" "nginx" {
  metadata {
    name = "nginx-example"
  }
  spec {
    selector = {
      App = kubernetes_deployment.nginx.spec.0.template.0.metadata[0].labels.App
    }
    port {
      port        = 80
      target_port = 80
    }

    type = "LoadBalancer"
  }
}

output "lb_ip" {
  value = kubernetes_service.nginx.status.0.load_balancer.0.ingress.0.ip
}

