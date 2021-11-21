resource "kubernetes_config_map" "order-service-workflow-worker-configmap" {
  metadata {
    name = "order-service-workflow-worker-config-map"
    labels = {
      name                         = "order-service-workflow-worker"
    }
  }

  data = {
    GO_ENV                                               = "production"
    LOG_LEVEL                                            = "Info"
    CADENCE_DOMAIN_NAME                                  = "order-service"
    CADENCE_FRONTEND_NAME                                = "cadence-frontend"
    CADENCE_CLIENT_NAME                                  = "cadence-client"
    CADENCE_HOST_ADDRESS_AND_PORT                        = "cadence-frontend-headless:7933"
  }
}

resource "kubernetes_secret" "order-service-workflow-worker-secret" {
  metadata {
    name = "order-service-workflow-worker-secret"
    labels = {
      name                         = "order-service-workflow-worker"
    }
  }
}

resource "kubernetes_deployment" "order-service-workflow-worker-deployment" {

  metadata {
    name = "order-service-workflow-worker"
    labels = {
      name                         = "order-service-workflow-worker"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        name = "order-service-workflow-worker"
      }
    }

    template {
      metadata {
        name = "order-service-workflow-worker"
        labels = {
          "name"                       = "order-service-workflow-worker"
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
          name              = "order-service-workflow-worker"
          image             = "gcr.io/e-commerce-332709/order-service-workflow-worker:latest"
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
                name = "order-service-workflow-worker-config-map"
              }
            }
          }

          env {
            name = "LOG_LEVEL"
            value_from {
              config_map_key_ref {
                key  = "LOG_LEVEL"
                name = "order-service-workflow-worker-config-map"
              }
            }
          }

          env {
            name = "CADENCE_DOMAIN_NAME"
            value_from {
              config_map_key_ref {
                key  = "CADENCE_DOMAIN_NAME"
                name = "order-service-workflow-worker-config-map"
              }
            }
          }

          env {
            name = "CADENCE_FRONTEND_NAME"
            value_from {
              config_map_key_ref {
                key  = "CADENCE_FRONTEND_NAME"
                name = "order-service-workflow-worker-config-map"
              }
            }
          }

          env {
            name = "CADENCE_CLIENT_NAME"
            value_from {
              config_map_key_ref {
                key  = "CADENCE_CLIENT_NAME"
                name = "order-service-workflow-worker-config-map"
              }
            }
          }

          env {
            name = "CADENCE_HOST_ADDRESS_AND_PORT"
            value_from {
              config_map_key_ref {
                key  = "CADENCE_HOST_ADDRESS_AND_PORT"
                name = "order-service-workflow-worker-config-map"
              }
            }
          }

        }
      }
    }
  }
}