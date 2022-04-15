# High availability Elasticsearch on Kubernetes with ECK and GKE

## Create GKE cluster

Creating a Kubernetes cluster using GKE is very straightforward. Navigate to the Kubernetes Engine page
and select Create Cluster. To ensure high-availability and prevent data loss, you want to create a cluster
with nodes that go across three availability zones in a region, so select Regional under Location Type and select
`europe-west-1` as the Region.
Also choose `e2-standard-4` as machine type for the Nodes to ensure the required resources.
This will create Kubernetes nodes across multiple zones in the region selected.


## Install the Operator

`kubectl create -f https://download.elastic.co/downloads/eck/2.1.0/crds.yaml`

`kubectl apply -f https://download.elastic.co/downloads/eck/2.1.0/operator.yaml`


## Exposing the service

`kubectl create -f storage_class.yml`

`kubectl patch storageclass standard -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"false"}}}'`

`kubectl patch storageclass zone-storage -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'`


## Deploy Elasticsearch using ECK

`kubectl apply -f es.yml`


## Update strategy

`kubectl apply -f es2.yml`


## Deploy Kibana

`kubectl apply -f kb.yml`


## Access the stack

`kubectl get svc ha-es-http`

`PASSWORD=$(kubectl get secret ha-es-elastic-user -o go-template='{{.data.elastic | base64decode}}') `
and `curl` using the EXTERNAL_IP ie:
`curl -k -u "elastic:$PASSWORD" https://34.76.163.44:9200/_cat/indices?v=true&s=index`

And for Kibana:
`kubectl get svc ha-kb-http`
And then visit using the EXTERNAL_IP ie: `https://104.199.52.165:5601/`

## Ship Metrics using the following config:

```yaml
output.elasticsearch:
  hosts: ['https://34.76.163.44:9200'] # replace with the proper EXTERNAL_IP
  username: 'elastic'
  password: 'um5k4C632Kx6cpdmE0785Ayx' # use proper PASSWORD
  ssl.verification_mode: none
```


