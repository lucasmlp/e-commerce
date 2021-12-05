resource "kubernetes_config_map" "e-commerce-workflow-worker-configmap" {
  metadata {
    name = "e-commerce-workflow-worker-config-map"
    labels = {
      name                         = "e-commerce-workflow-worker"
    }
  }

  data = {
    GO_ENV                                               = "production"
    LOG_LEVEL                                            = "Info"
    CADENCE_DOMAIN_NAME                                  = "e-commerce"
    CADENCE_FRONTEND_NAME                                = "cadence-frontend"
    CADENCE_CLIENT_NAME                                  = "cadence-client"
    CADENCE_HOST_ADDRESS_AND_PORT                        = "cadence-frontend-headless:7933"
  }
}

resource "kubernetes_secret" "e-commerce-workflow-worker-secret" {
  metadata {
    name = "e-commerce-workflow-worker-secret"
    labels = {
      name                         = "e-commerce-workflow-worker"
    }
  }
}

resource "kubernetes_deployment" "e-commerce-workflow-worker-deployment" {

  metadata {
    name = "e-commerce-workflow-worker"
    labels = {
      name                         = "e-commerce-workflow-worker"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        name = "e-commerce-workflow-worker"
      }
    }

    template {
      metadata {
        name = "e-commerce-workflow-worker"
        labels = {
          "name"                       = "e-commerce-workflow-worker"
        }
      }

      spec {
        restart_policy = "Always"

        toleration {
          key      = "node-pool"
          operator = "Equal"
          value    = "cadence"
          effect   = "NoSchedule"
        }

        node_selector = {
          "cloud.google.com/gke-nodepool" = "cadence"
        }

        container {
          name              = "e-commerce-workflow-worker"
          image             = "gcr.io/e-commerce-332709/e-commerce-workflow-worker:latest"
          image_pull_policy = "IfNotPresent"

          resources {
            requests = {
              cpu    = "20m"
              memory = "128Mi"
            }

            limits = {
              cpu    = "100m"
              memory = "512Mi"
            }
          }

          port {
            container_port = 9090
            name           = "metrics"
            protocol       = "TCP"
          }

          port {
            container_port = 3004
            name           = "cdnc-metrics"
            protocol       = "TCP"
          }

          env {
            name = "GO_ENV"
            value_from {
              config_map_key_ref {
                key  = "GO_ENV"
                name = "e-commerce-workflow-worker-config-map"
              }
            }
          }

          env {
            name = "LOG_LEVEL"
            value_from {
              config_map_key_ref {
                key  = "LOG_LEVEL"
                name = "e-commerce-workflow-worker-config-map"
              }
            }
          }

          env {
            name = "CADENCE_DOMAIN_NAME"
            value_from {
              config_map_key_ref {
                key  = "CADENCE_DOMAIN_NAME"
                name = "e-commerce-workflow-worker-config-map"
              }
            }
          }

          env {
            name = "CADENCE_FRONTEND_NAME"
            value_from {
              config_map_key_ref {
                key  = "CADENCE_FRONTEND_NAME"
                name = "e-commerce-workflow-worker-config-map"
              }
            }
          }

          env {
            name = "CADENCE_CLIENT_NAME"
            value_from {
              config_map_key_ref {
                key  = "CADENCE_CLIENT_NAME"
                name = "e-commerce-workflow-worker-config-map"
              }
            }
          }

          env {
            name = "CADENCE_HOST_ADDRESS_AND_PORT"
            value_from {
              config_map_key_ref {
                key  = "CADENCE_HOST_ADDRESS_AND_PORT"
                name = "e-commerce-workflow-worker-config-map"
              }
            }
          }

        }
      }
    }
  }
}