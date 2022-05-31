# k8s-integration-infra

The purpose of this repo is to:
1. Automate the provision of a Kubernetes cluster in GKE only, using Terraform.
2. Provision elastic stack(Elasticsearch, Kibana) using ECK(if user requests it) on the same K8s cluster. 
3. Deploy metricbeat/filebeat/standalone_agent on the K8s cluster.
4. Stress test the cluster by deploying multiple pods in various namespaces using a cli tool.
5. Take statistics of elasticsearch target indices, storage size, docs counts and execution time of [query 12](https://gist.github.com/ChrsMark/f4292c388879eeb5368218068d09d40c#12-top-cpu-intensive-pods) by executing es_bench script

## Prerequisites
1. configured gcloud SDK
2. kubectl >= 1.7.0
3. terraform >= 0.14
4. helm > 2.4.1
5. golang >= 1.17.0
6. jq
7. elasticsearch cluster reachable by gcp(only in case `provision_elasticsearch` is set to `false`)

### Bring up the cluster for the first time
1. `cd infra`
2. `terraform init`
3. Update the `terraform.tfvars` file.
    * Set the google cloud `project_id`, k8s `cluster_name`, k8s nodes `machine_type` and the cluster `region`.
   For `project_id`, `region` and `machine_type` defaults can be used. `cluster_name` has to be unique.

   * Configure elasticsearch url, credentials and version. There are two options available:
         
     a. In case user wants a new elasticsearch cluster to be provisioned using ECK then `provision_elasticsearch` must be set to `true`. In that case variables `es_password` and `es_host` can be left empty. `es_user` should keep the default value and `imageTag` should be set to the version required.
   
     b. In case user already has an elasticsearch cluster deployed and reachable by gcp (elastic cloud) then set `provision_elasticsearch = false` and set the right values to variables  
      
        ```
         es_host, es_user, es_password, imageTag
         ```
   * Set the size of the cluster by setting the required nodes number in variables `gke_num_nodes` and `gke_max_num_nodes`. As the cluster is regional with 3 zones per region, the value set in those variables will result to 3x number of nodes created.
   * User can select if they want their cluster to be monitored by either metricbeat, filebeat or elastic-agent in standalone mode by setting the appropriate value in variables `deployBeat`, `deployAgent`. Both options can be used.

   * Example configuration:
      ```
         project_id              = "elastic-observability"
         region                  = "europe-west1"
         cluster_name            = "test-k8s-cluster"
         machine_type            = "e2-standard-4"
         gke_num_nodes           = 1
         gke_max_num_nodes       = 2
         provision_elasticsearch = true
         es_password             = ""
         es_user                 = "elastic"
         es_host                 = ""
         deployBeat              = true
         deployAgent             = false
         imageTag                = "8.3.0-SNAPSHOT"
         namespace               = "kube-system"
      ```
4. `terraform apply`


The above command will bring up a regional gke cluster initially with `gke_num_nodes` * (number of zones in region) number of nodes with
autoscaling up to `gke_max_num_nodes` * (number of zones in region).
In case `provision_elasticsearch = true` it will also provision an elasticsearch cluster using ECK in the same kubernetes cluster.
If `provision_elasticsearch = false` it will not bring up an elasticsearch cluster and will rely on `es_host`, `es_user`, `es_password` variables' values.
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
as many 20 replicas as indicated in the `deployments` flag. Each pod will have 4 labels and 4 annotations.

### Helper es_bench script (TSDB use case only)
Helper Script: scripts/es_bench
The helper script executes [q12](https://gist.github.com/ChrsMark/f4292c388879eeb5368218068d09d40c#12-top-cpu-intensive-pods) from the list of queries.
Execute the helper script to get the stats from the  ES instances one simple and one using TSDB:

TSDB_ES_URL="https://35.157.42.42:9200/" TSDB_ES_PASS="passpasstsdb" SIMPLE_ES_URL="https://104.199.42.42:9200/" SIMPLE_ES_PASS="passpasssimple" TSDB_INDEX=".ds-metricbeat-tsdb-8.3.0-2022.05.24-000001" SIMPLE_INDEX=".ds-metricbeat-8.3.0-2022.05.24-000001" go run main.go