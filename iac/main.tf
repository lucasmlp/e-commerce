terraform {
  backend "gcs" {
    bucket = "iac-cadence"
    # credentials = "./cadence-poc-b055cfd0672c.json"
  }
}

provider "kubernetes" {
  host = "http://127.0.0.1:9874"
}

module "cassandra" {
  source = "./cadence"
}