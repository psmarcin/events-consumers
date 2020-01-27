variable "project" {
  type = string
  default = "events-consumer"
}

variable "region" {
  type = string
  default = "europe-west1"
}

variable "zone" {
  type = string
  default = "europe-west1-c"
}

provider "google" {
  project = var.project
  region  = var.region
  zone    = var.zone
}
