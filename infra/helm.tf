provider "helm" {

  kubernetes {
    host                   = "https://${google_container_cluster.primary.endpoint}"
    token                  = data.google_client_config.default.access_token
    cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)
  }
}

resource "helm_release" "metricbeat" {
  name        = "metricbeat"
  chart       = "metricbeat"
  repository  = "./charts"
  namespace   = "default"
  max_history = 3
  create_namespace = true
  wait             = true
  reset_values     = true
}

resource "helm_release" "filebeat" {
  name        = "filebeat"
  chart       = "filebeat"
  repository  = "./charts"
  namespace   = "default"
  max_history = 3
  create_namespace = true
  wait             = true
  reset_values     = true
}