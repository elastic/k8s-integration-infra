data "google_client_config" "default" {}

provider "kubernetes" {
  host                   = "https://${google_container_cluster.primary.endpoint}"
  token                  = data.google_client_config.default.access_token
  cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)
}

provider "kubectl" {
  host                   = "https://${google_container_cluster.primary.endpoint}"
  token                  = data.google_client_config.default.access_token
  cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth.0.cluster_ca_certificate)
}


resource "kubernetes_namespace" "elastic-system" {
  count = var.provision_elasticsearch ? 1 : 0
  metadata {

    labels = {
      name = "elastic-system"
    }

    name = "elastic-system"
  }
  depends_on = [google_container_node_pool.primary_nodes]
}

data "kubectl_file_documents" "crds_docs" {
  content = file("../eck/crds.yaml")
}

resource "kubectl_manifest" "crds" {
  for_each   = var.provision_elasticsearch ? data.kubectl_file_documents.crds_docs.manifests : {}
  yaml_body  = each.value
  depends_on = [google_container_node_pool.primary_nodes]

}

data "kubectl_file_documents" "op_docs" {
  content = file("../eck/operator.yaml")
}

resource "kubectl_manifest" "operator" {
  for_each   = var.provision_elasticsearch ? data.kubectl_file_documents.op_docs.manifests : {}
  yaml_body  = each.value
  depends_on = [kubernetes_namespace.elastic-system]

}

resource "kubectl_manifest" "standard_storage_class" {
  count      = var.provision_elasticsearch ? 1 : 0
  yaml_body  = file("../eck/standard-storage_class.yaml")
  depends_on = [google_container_node_pool.primary_nodes]
}

resource "kubectl_manifest" "zone-storage_class" {
  count      = var.provision_elasticsearch ? 1 : 0
  yaml_body  = file("../eck/zone-storage_class.yaml")
  depends_on = [google_container_node_pool.primary_nodes]
}

resource "kubectl_manifest" "elasticsearch" {
  count      = var.provision_elasticsearch ? 1 : 0
  yaml_body  = templatefile("../eck/es.yml", { imageTag = var.imageTag, region = var.region })
  depends_on = [kubectl_manifest.operator]
}

resource "kubectl_manifest" "elasticsearch_updated" {
  count      = var.provision_elasticsearch ? 1 : 0
  yaml_body  = templatefile("../eck/es2.yml", { imageTag = var.imageTag, region = var.region })
  depends_on = [kubectl_manifest.elasticsearch]
}

resource "kubectl_manifest" "kibana" {
  count      = var.provision_elasticsearch ? 1 : 0
  yaml_body  = templatefile("../eck/kb.yml", { imageTag = var.imageTag })
  depends_on = [kubectl_manifest.elasticsearch_updated]
}

data "external" "get_elasticsearch_creds" {
  program = ["/bin/bash", "../eck/get_es_cluster.sh"]
  query = {
    run_script   = "${var.provision_elasticsearch}"
    cluster_name = "${var.cluster_name}"
    region       = "${var.region}"
    project      = "${var.project_id}"
  }
  depends_on = [kubectl_manifest.kibana]
}

output "elasticsearch_ip" {
  value = var.provision_elasticsearch ? data.external.get_elasticsearch_creds.result.es_ip : var.es_host
  description = "Elasticsearch IP address and port"

}

output "elasticsearch_password" {
  value = var.provision_elasticsearch ? data.external.get_elasticsearch_creds.result.es_password : var.es_password
  description = "Elasticsearch password"

}

output "kibana_ip" {
  value = var.provision_elasticsearch ? data.external.get_elasticsearch_creds.result.kibana_ip : ""
  description = "Kibana IP address and port"

}

resource "kubernetes_secret" "elasticsearch-master-credentials" {
  metadata {
    name      = "elasticsearch-credentials"
    namespace = var.namespace
  }

  data = {
    username = var.es_user
    password = var.provision_elasticsearch ? data.external.get_elasticsearch_creds.result.es_password : var.es_password
  }

}

resource "kubernetes_secret" "elasticsearch-master-url" {
  metadata {
    name      = "elasticsearch-host"
    namespace = var.namespace
  }

  data = {
    host = var.provision_elasticsearch ? data.external.get_elasticsearch_creds.result.es_ip : var.es_host
  }

}