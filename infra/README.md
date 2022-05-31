<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 0.14 |
| <a name="requirement_external"></a> [external](#requirement\_external) | ~> 2.2.2 |
| <a name="requirement_google"></a> [google](#requirement\_google) | 3.52.0 |
| <a name="requirement_helm"></a> [helm](#requirement\_helm) | ~> 2.4.1 |
| <a name="requirement_kubectl"></a> [kubectl](#requirement\_kubectl) | >= 1.7.0 |
| <a name="requirement_kubernetes"></a> [kubernetes](#requirement\_kubernetes) | >= 2.0.1 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_external"></a> [external](#provider\_external) | 2.2.2 |
| <a name="provider_google"></a> [google](#provider\_google) | 3.52.0 |
| <a name="provider_helm"></a> [helm](#provider\_helm) | 2.4.1 |
| <a name="provider_kubectl"></a> [kubectl](#provider\_kubectl) | 1.14.0 |
| <a name="provider_kubernetes"></a> [kubernetes](#provider\_kubernetes) | 2.9.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [google_compute_network.vpc](https://registry.terraform.io/providers/hashicorp/google/3.52.0/docs/resources/compute_network) | resource |
| [google_compute_subnetwork.subnet](https://registry.terraform.io/providers/hashicorp/google/3.52.0/docs/resources/compute_subnetwork) | resource |
| [google_container_cluster.primary](https://registry.terraform.io/providers/hashicorp/google/3.52.0/docs/resources/container_cluster) | resource |
| [google_container_node_pool.primary_nodes](https://registry.terraform.io/providers/hashicorp/google/3.52.0/docs/resources/container_node_pool) | resource |
| [helm_release.elastic-agent](https://registry.terraform.io/providers/hashicorp/helm/latest/docs/resources/release) | resource |
| [helm_release.filebeat](https://registry.terraform.io/providers/hashicorp/helm/latest/docs/resources/release) | resource |
| [helm_release.metricbeat](https://registry.terraform.io/providers/hashicorp/helm/latest/docs/resources/release) | resource |
| [kubectl_manifest.crds](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/resources/manifest) | resource |
| [kubectl_manifest.elasticsearch](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/resources/manifest) | resource |
| [kubectl_manifest.elasticsearch_updated](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/resources/manifest) | resource |
| [kubectl_manifest.kibana](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/resources/manifest) | resource |
| [kubectl_manifest.operator](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/resources/manifest) | resource |
| [kubectl_manifest.standard_storage_class](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/resources/manifest) | resource |
| [kubectl_manifest.zone-storage_class](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/resources/manifest) | resource |
| [kubernetes_namespace.elastic-system](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/namespace) | resource |
| [kubernetes_secret.elasticsearch-master-credentials](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/secret) | resource |
| [kubernetes_secret.elasticsearch-master-url](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/secret) | resource |
| [external_external.get_elasticsearch_creds](https://registry.terraform.io/providers/hashicorp/external/latest/docs/data-sources/external) | data source |
| [google_client_config.default](https://registry.terraform.io/providers/hashicorp/google/3.52.0/docs/data-sources/client_config) | data source |
| [kubectl_file_documents.crds_docs](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/data-sources/file_documents) | data source |
| [kubectl_file_documents.op_docs](https://registry.terraform.io/providers/gavinbunney/kubectl/latest/docs/data-sources/file_documents) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_cluster_name"></a> [cluster\_name](#input\_cluster\_name) | cluster name | `any` | n/a | yes |
| <a name="input_deployAgent"></a> [deployAgent](#input\_deployAgent) | deploy elastic-agent | `bool` | `false` | no |
| <a name="input_deployBeat"></a> [deployBeat](#input\_deployBeat) | deploy metricbeat and filebeat | `bool` | `true` | no |
| <a name="input_es_host"></a> [es\_host](#input\_es\_host) | The url:port of elasticsearch. | `string` | `"elasticsearch:9200"` | no |
| <a name="input_es_password"></a> [es\_password](#input\_es\_password) | The password of elasticsearch. | `string` | `"changeme"` | no |
| <a name="input_es_user"></a> [es\_user](#input\_es\_user) | The username of elasticsearch. | `string` | `"elastic"` | no |
| <a name="input_gke_max_num_nodes"></a> [gke\_max\_num\_nodes](#input\_gke\_max\_num\_nodes) | max number of gke nodes | `number` | `10` | no |
| <a name="input_gke_num_nodes"></a> [gke\_num\_nodes](#input\_gke\_num\_nodes) | number of gke nodes | `number` | `2` | no |
| <a name="input_gke_password"></a> [gke\_password](#input\_gke\_password) | gke password | `string` | `""` | no |
| <a name="input_gke_username"></a> [gke\_username](#input\_gke\_username) | gke username | `string` | `""` | no |
| <a name="input_imageTag"></a> [imageTag](#input\_imageTag) | The beats and agent image version | `string` | `"8.2.0-SNAPSHOT"` | no |
| <a name="input_machine_type"></a> [machine\_type](#input\_machine\_type) | machine type for gke cluster | `string` | `"e2-standard-4"` | no |
| <a name="input_namespace"></a> [namespace](#input\_namespace) | The namespace to deploy beats and agent | `string` | `"default"` | no |
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | project id | `any` | n/a | yes |
| <a name="input_provision_elasticsearch"></a> [provision\_elasticsearch](#input\_provision\_elasticsearch) | provision an elasticsearch cluster on gke with eck | `bool` | `false` | no |
| <a name="input_region"></a> [region](#input\_region) | region | `any` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_elasticsearch_ip"></a> [elasticsearch\_ip](#output\_elasticsearch\_ip) | Elasticsearch IP address and port |
| <a name="output_elasticsearch_password"></a> [elasticsearch\_password](#output\_elasticsearch\_password) | Elasticsearch password |
| <a name="output_kibana_ip"></a> [kibana\_ip](#output\_kibana\_ip) | Kibana IP address and port |
| <a name="output_kubernetes_cluster_host"></a> [kubernetes\_cluster\_host](#output\_kubernetes\_cluster\_host) | GKE Cluster Host |
| <a name="output_kubernetes_cluster_name"></a> [kubernetes\_cluster\_name](#output\_kubernetes\_cluster\_name) | GKE Cluster Name |
| <a name="output_project_id"></a> [project\_id](#output\_project\_id) | GCloud Project ID |
| <a name="output_region"></a> [region](#output\_region) | GCloud Region |
<!-- END_TF_DOCS -->