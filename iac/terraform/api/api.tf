resource "kubernetes_config_map" "e-commerce-api-configmap" {
  metadata {
    name = "e-commerce-api-config-map"
    labels = {
      name                         = "e-commerce-api"
    }
  }

  data = {
    GO_ENV                                               = "development"
    LOG_LEVEL                                            = "Info"
    MONGODB_URI                                  = "mongodb://localhost:27017"
    ORDERS_DATABASE_NAME                                = "order-service"
    ORDERS_COLLECTION_NAME                                  = "orders"
    PRODUCTS_DATABASE_NAME                                = "product-service"
    PRODUCTS_COLLECTION_NAME                                  = "products"
    
  }
}

resource "kubernetes_secret" "e-commerce-api-secret" {
  metadata {
    name = "e-commerce-api-secret"
    labels = {
      name                         = "e-commerce-api"
    }
  }
}

resource "kubernetes_deployment" "e-commerce-api-deployment" {

  metadata {
    name = "e-commerce-api"
    labels = {
      name                         = "e-commerce-api"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        name = "e-commerce-api"
      }
    }

    template {
      metadata {
        name = "e-commerce-api"
        labels = {
          "name"                       = "e-commerce-api"
        }
      }

      spec {
        restart_policy = "Always"

        toleration {
          key      = "node-pool"
          operator = "Equal"
          value    = "api"
          effect   = "NoSchedule"
        }

        node_selector = {
          # change
          "cloud.google.com/gke-nodepool" = "api"
        }

        container {
          name              = "e-commerce-api"
          # change
          image             = "gcr.io/e-commerce-332709/e-commerce-api:latest"
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

          env {
            name = "GO_ENV"
            value_from {
              config_map_key_ref {
                key  = "GO_ENV"
                name = "e-commerce-api-config-map"
              }
            }
          }

          env {
            name = "LOG_LEVEL"
            value_from {
              config_map_key_ref {
                key  = "LOG_LEVEL"
                name = "e-commerce-api-config-map"
              }
            }
          }

          env {
            name = "MONGODB_URI"
            value_from {
              config_map_key_ref {
                key  = "MONGODB_URI"
                name = "e-commerce-api-config-map"
              }
            }
          }

          env {
            name = "ORDERS_DATABASE_NAME"
            value_from {
              config_map_key_ref {
                key  = "ORDERS_DATABASE_NAME"
                name = "e-commerce-api-config-map"
              }
            }
          }

          env {
            name = "ORDERS_COLLECTION_NAME"
            value_from {
              config_map_key_ref {
                key  = "ORDERS_COLLECTION_NAME"
                name = "e-commerce-api-config-map"
              }
            }
          }

          env {
            name = "PRODUCTS_DATABASE_NAME"
            value_from {
              config_map_key_ref {
                key  = "PRODUCTS_DATABASE_NAME"
                name = "e-commerce-api-config-map"
              }
            }
          }

          env {
            name = "PRODUCTS_COLLECTION_NAME"
            value_from {
              config_map_key_ref {
                key  = "PRODUCTS_COLLECTION_NAME"
                name = "e-commerce-api-config-map"
              }
            }
          }

        }
      }
    }
  }
}