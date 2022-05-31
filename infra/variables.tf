variable "gke_username" {
  default     = ""
  description = "gke username"
}

variable "gke_password" {
  default     = ""
  description = "gke password"
}

variable "gke_num_nodes" {
  default     = 2
  description = "number of gke nodes"
}

variable "gke_max_num_nodes" {
  default     = 10
  description = "max number of gke nodes"
}

variable "machine_type" {
  default     = "e2-standard-4"
  description = "machine type for gke cluster"
}

variable "deployAgent" {
  type        = bool
  default     = false
  description = "deploy elastic-agent"
}

variable "deployBeat" {
  type        = bool
  default     = true
  description = "deploy metricbeat and filebeat"
}

variable "imageTag" {
  type        = string
  default     = "8.2.0-SNAPSHOT"
  description = "The beats and agent image version"
}

variable "es_password" {
  default     = "changeme"
  type        = string
  description = "The password of elasticsearch."
}

variable "es_user" {
  default     = "elastic"
  type        = string
  description = "The username of elasticsearch."
}

variable "es_host" {
  default     = "elasticsearch:9200"
  type        = string
  description = "The url:port of elasticsearch."
}

variable "namespace" {
  default     = "default"
  type        = string
  description = "The namespace to deploy beats and agent"
}

variable "provision_elasticsearch" {
  type        = bool
  default     = false
  description = "provision an elasticsearch cluster on gke with eck"
}

variable "project_id" {
  description = "project id"
}

variable "region" {
  description = "region"
}

variable "cluster_name" {
  description = "cluster name"
}