provider "helm" {

  kubernetes {
    host                   = "https://${google_container_cluster.primary.endpoint}"
    token                  = data.google_client_config.default.access_token
    cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)
  }
}

resource "helm_release" "metricbeat" {
  count            = var.deployBeat ? 1 : 0
  name             = "metricbeat"
  chart            = "metricbeat"
  repository       = "./charts"
  namespace        = var.namespace
  max_history      = 3
  create_namespace = true
  wait             = true
  reset_values     = true
  set {
    name  = "imageTag"
    value = var.imageTag
  }
  depends_on = [kubernetes_secret.elasticsearch-master-credentials]
}

resource "helm_release" "filebeat" {
  count            = var.deployBeat ? 1 : 0
  name             = "filebeat"
  chart            = "filebeat"
  repository       = "./charts"
  namespace        = var.namespace
  max_history      = 3
  create_namespace = true
  wait             = true
  reset_values     = true
  set {
    name  = "imageTag"
    value = var.imageTag
  }
  depends_on = [kubernetes_secret.elasticsearch-master-credentials]
}

resource "helm_release" "elastic-agent" {
  count            = var.deployAgent ? 1 : 0
  name             = "elastic-agent"
  chart            = "elastic-agent"
  repository       = "./charts"
  namespace        = var.namespace
  max_history      = 3
  create_namespace = true
  wait             = true
  reset_values     = true
  set {
    name  = "imageTag"
    value = var.imageTag
  }
  depends_on = [kubernetes_secret.elasticsearch-master-credentials]
}