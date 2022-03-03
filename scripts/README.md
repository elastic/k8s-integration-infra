1. `gcloud container clusters get-credentials cloud-native-integrations-team-gke --zone us-central1-c --project elastic-observability`
2. `go build`
3. `./stress_test_k8s --kubeconfig=<path> --deployments=4 --namespaces=2`