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
6. jq >= 1.6
7. elasticsearch cluster reachable by gcp(only in case `provision_elasticsearch` is set to `false`)

### Bring up the cluster for the first time
1. `cd infra`
2. `terraform init`
3. Set the google cloud `project_id`, k8s `cluster_name`, k8s nodes `machine_type` and the cluster `region` in `terraform.tfvars` file.
    * For `project_id`, `region` and `machine_type` defaults can be used. `cluster_name` has to be unique.
4. Configure ElasticSearch Cluster in `terraform.tfvars` file. There are two options available:
   1. In case user wants a new elasticsearch cluster to be provisioned using ECK then `provision_elasticsearch` must be set to `true`.
      In that case variables `es_password` and `es_host` can be left empty. `es_user` should keep the default value and `imageTag` should be set to the version required.
      
   2. In case user already has an elasticsearch cluster deployed and reachable by gcp (elastic cloud) then `provision_elasticsearch = false`  must be set as well as the right values to variables  
      
        ```
         es_host, es_user, es_password, imageTag
         ```
5. Set the size of the cluster by setting the required nodes number in variables `gke_num_nodes` and `gke_max_num_nodes` inside `terraform.tfvars` file.
   As the cluster is regional with 3 zones per region, the value set in those variables will result to 3x number of nodes created (`gke_num_nodes` * (number of zones in region)). 
   `gke_max_num_nodes` enables cluster autoscaling in case there is a need for more resources.
6. Configure Monitoring.
   User can select if they want their cluster to be monitored by either metricbeat/filebeat or elastic-agent in standalone mode by setting the appropriate values in variables `deployBeat`, `deployAgent`.
   Both options can be used.
7. `terraform apply`
8. Configure `kubectl` by running `gcloud container clusters get-credentials <cluster-name> --zone europe-west1 --project elastic-obs-integrations-dev`
   The correct command can be obtained from Kubernetes Engine in GCP.
9. Check the cluster `kubectl get node`, `kubectl get pod -A`

### Examples of different setups:
1. Bring up a kubernetes cluster with 3 nodes and no autoscaling, without provisioning Elasticsearch and without monitoring
* Example configuration:
   ```
      project_id              = "elastic-obs-integrations-dev"
      region                  = "europe-west1"
      cluster_name            = "test-k8s-cluster-simple"
      machine_type            = "e2-standard-4"
      gke_num_nodes           = 1
      gke_max_num_nodes       = 1
      provision_elasticsearch = false
      es_password             = ""
      es_user                 = "elastic"
      es_host                 = ""
      deployBeat              = false
      deployAgent             = false
      imageTag                = ""
      namespace               = "kube-system"
   ```

2. Bring up a kubernetes cluster with 3 nodes and autoscaling up to 18 nodes, without provisioning Elasticsearch and with Beats monitoring version 8.3.0.
   Prerequisite is the existence of an elastic stack. 
* Example configuration:
   ```
      project_id              = "elastic-obs-integrations-dev"
      region                  = "europe-west1"
      cluster_name            = "test-k8s-cluster-autoscaling-beats"
      machine_type            = "e2-standard-4"
      gke_num_nodes           = 1
      gke_max_num_nodes       = 6
      provision_elasticsearch = false
      es_password             = "mypassword"
      es_user                 = "elastic"
      es_host                 = "https://bxxxxxed.europe-west1.gcp.cloud.es.io:9243"
      deployBeat              = true
      deployAgent             = false
      imageTag                = "8.3.0"
      namespace               = "kube-system"
   ```

3. Bring up a kubernetes cluster with 3 nodes and autoscaling up to 18 nodes, with Elasticsearch provisioning and with elastic-agent monitoring version 8.3.0.
* Example configuration:
   ```
      project_id              = "elastic-obs-integrations-dev"
      region                  = "europe-west1"
      cluster_name            = "test-k8s-cluster-autoscaling-elasticsearch-agent"
      machine_type            = "e2-standard-4"
      gke_num_nodes           = 1
      gke_max_num_nodes       = 6
      provision_elasticsearch = true
      es_password             = ""
      es_user                 = "elastic"
      es_host                 = ""
      deployBeat              = false
      deployAgent             = true
      imageTag                = "8.3.0"
      namespace               = "kube-system"
   ```

NOTE.
The command may end with an error like this but everything should be successfully deployed:
```
Error: Kubernetes cluster unreachable: Get "https://35.239.222.162/version?timeout=32s": dial tcp 35.239.222.162:443: connect: connection refused
```


### Put load on the cluster
1. `cd scripts`
2. `go build`
3. `./stress_test_k8s --kubeconfig=/Users/<username>/.kube/config --deployments=20 --namespaces=10 --podlabels=4 --podannotations=4`

The above command will create 10 namespaces and deploy one demo nginx deployment in each one with
as many 20 replicas as indicated in the `deployments` flag. Each pod will have 4 labels and 4 annotations.

### Helper es_bench script (TSDB use case only)
####Prerequisite: Existence of 2 Elasticsearch Clusters. One with metricbeat index TSDB enabled and one without.

In order to get a quick estimation of the status of the 2 Elasticsearch indices(one simple and one TSDB enabled) one can execute `scripts/es_bench`.  By now the script can only be executed manually. More specifically the script provides the following information about the cluster:
* `pri.store.size`
* `docs.count`
  This information would be available through `_cat/indices?v=true&s=index` API.
  In addition to this, the script also executes [q12](https://gist.github.com/ChrsMark/f4292c388879eeb5368218068d09d40c#12-top-cpu-intensive-pods) which is considered as "expensive" for our use case. The query is executed 20 times sequentially for each ES cluster and provides the median of the execution times.

Execution example:
TSDB_ES_URL="https://35.157.42.42:9200/" TSDB_ES_PASS="passpasstsdb" SIMPLE_ES_URL="https://104.199.42.42:9200/" SIMPLE_ES_PASS="passpasssimple" TSDB_INDEX=".ds-metricbeat-tsdb-8.3.0-2022.05.24-000001" SIMPLE_INDEX=".ds-metricbeat-8.3.0-2022.05.24-000001" go run main.go
Example output:
*************************************
Executing against new ES cluster
*************************************
Client: 8.2.0
Server: 8.3.0-SNAPSHOT
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
index name: .ds-metricbeat-tsdb-8.3.0-2022.05.24-000001
pri.store.size: 5.8gb
docs.count: 25635493
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
median query time is: 2ms
*************************************
Executing against new ES cluster
*************************************
Client: 8.2.0
Server: 8.3.0-SNAPSHOT
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
index name: .ds-metricbeat-8.3.0-2022.05.24-000001
pri.store.size: 23gb
docs.count: 39051417
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
median query time is: 333ms