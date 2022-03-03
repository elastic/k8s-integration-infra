# k8s-integration-infra

The purpose of this repo is to:
1. Automate the provision of a Kubernetes cluster in GKE, monitored by filebeat and metricbeat.
2. Stress test the cluster by deploying multiple pods in various namespaces using a cli tool.

## Prerequisites
1. configured gcloud SDK
2. kubectl
3. terraform
4. helm > 2.4.1
5. golang

### Bring up the cluster
In case the cluster is not running or an update is needed.
1. `cd infra`
2. `terraform init`
3. Update the `terraform.tfvars` file with the correct values for 
```
   gke_num_nodes = 2
   gke_max_num_nodes = 10
   es_password = ""
   es_user = "elastic"
   es_host = ""
```
4. `terraform apply`

The above command will bring up a gke cluster initially with `gke_num_nodes` number of nodes with
autoscaling up to `gke_max_num_nodes`.
It will also bring up metricbeat and filebeat using helm as well as some secrets with the es
credentials and host url.

5. Configure `kubectl` by running `gcloud container clusters get-credentials cloud-native-integrations-team-gke --zone us-central1-c --project elastic-observability`
6. Check the cluster `kubectl get node`, `kubectl get pod -A`