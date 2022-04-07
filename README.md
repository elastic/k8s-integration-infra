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
6. elasticsearch cluster reachable by gcp(e.g. deployed in elastic cloud)

### Bring up the cluster
In case the cluster is not running or an update is needed.
1. `cd infra`
2. `terraform init`
3. Update the `terraform.tfvars` file.
    * Fill in the correct elasticsearh url and credentials:
   ```
   es_host, es_user, es_password
   ```
    * Set the cluster name and the region:
   ```
   region, cluster_name
   ```
   
```
   project_id= "elastic-observability"
   region= "us-central1"
   cluster_name= ""
   gke_num_nodes = 1
   gke_max_num_nodes = 10
   es_password = ""
   es_user = "elastic"
   es_host = ""
   deployBeat =  true
   deployAgent = true
   imageTag = "8.2.0-SNAPSHOT"
   namespace = "kube-system"
```
4. `terraform apply`


The above command will bring up a regional gke cluster initially with `gke_num_nodes` * (number of zones in region) number of nodes with
autoscaling up to `gke_max_num_nodes` * (number of zones in region).
In case `deployBeat = true` it will also bring up metricbeat and filebeat using helm as well as some secrets with the es
credentials and host url.
If `deployAgent = true`, elastic-agent will also be deployed in the cluster. If set to false,
it won't.

NOTE.
The command may end with an error like this but everything should be successfully deployed:
```
Error: Kubernetes cluster unreachable: Get "https://35.239.222.162/version?timeout=32s": dial tcp 35.239.222.162:443: connect: connection refused
```
5. Configure `kubectl` by running `gcloud container clusters get-credentials <cluster-name> --zone us-central1-c --project elastic-observability`
The correct command can be obtained from Kubernetes Engine in GCP.
6. Check the cluster `kubectl get node`, `kubectl get pod -A`

### Put load on the cluster
1. `cd scripts`
2. `go build`
3. `./stress_test_k8s --kubeconfig=/Users/<username>/.kube/config --deployments=20 --namespaces=10 --podlabels=4 --podannotations=4`

The above command will create 10 namespaces and deploy one demo nginx deployment in each one with
as many 20 replicas as indicated in the `deployments` flag. Each pod will have 4 labels and 4 annotations,