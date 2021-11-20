provider "kubernetes" {
  host = "http://127.0.0.1:9874"
}

module "cassandra" {
  source = "./cassandra"
}