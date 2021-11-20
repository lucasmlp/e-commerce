resource "kubernetes_storage_class" "cassandra_storage_class" {
  metadata {
    name = "cadence-cassandra-sc"
  }
  storage_provisioner    = "kubernetes.io/gce-pd"
  reclaim_policy         = "Retain"
  allow_volume_expansion = true
  parameters = {
    type = "pd-ssd"
  }
}