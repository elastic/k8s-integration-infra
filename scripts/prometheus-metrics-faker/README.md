gcloud auth configure-docker


docker build --tag prometheus-metrics-faker .
docker tag prometheus-metrics-faker eu.gcr.io/elastic-obs-integrations-dev/prometheus-metrics-faker:0.1

docker push eu.gcr.io/elastic-obs-integrations-dev/prometheus-metrics-faker:0.1


add "https://www.googleapis.com/auth/devstorage.read_only" scope to the nodepool (require recreation of the node pool)


docker run -it -p 127.0.0.1:9000:9000 prometheus-metrics-faker --gauge=100 --labels=4
