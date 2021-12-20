package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/machado-br/e-commerce/domain/dtos"
	"github.com/machado-br/e-commerce/domain/orders"
	"github.com/machado-br/e-commerce/domain/products"
	"github.com/pborman/uuid"
)

type api struct {
	Debug           bool
	OrdersService   orders.Service
	ProductsService products.Service
}

func NewApi(
	debug bool,
	orderService orders.Service,
	productsService products.Service,
) (api, error) {
	return api{
		Debug:           debug,
		OrdersService:   orderService,
		ProductsService: productsService,
	}, nil
}

func (a api) Engine() *gin.Engine {
	router := gin.Default()

	root := router.Group("")
	{
		root.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, "pong")
		})
	}
	ordersRoot := router.Group("/orders")
	{
		ordersRoot.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, "pong")
		})

		ordersRoot.GET("/:id", func(c *gin.Context) {
			log.Println("GET /orders/:id")

			orderId := c.Param("id")

			order, err := a.OrdersService.Find(c, orderId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusOK, order)
		})

		ordersRoot.GET("", func(c *gin.Context) {
			log.Println("GET /orders")

			orders, err := a.OrdersService.FindAll(c)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusOK, orders)
		})

		ordersRoot.POST("", func(c *gin.Context) {
			log.Println("POST /orders")

			var order dtos.Order
			err := c.ShouldBindJSON(&order)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}

			result, err := a.OrdersService.Create(c, order)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusOK, result)
		})

		ordersRoot.DELETE(":id", func(c *gin.Context) {
			log.Println("DELETE /orders")

			orderId := c.Param("id")
			err := a.OrdersService.Delete(c, orderId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.Status(http.StatusNoContent)
		})

		ordersRoot.PUT("", func(c *gin.Context) {
			log.Println("PUT /orders")

			var order dtos.Order
			err := c.ShouldBindJSON(&order)
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, err)
			}

			err = a.OrdersService.Update(c, order)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.Status(http.StatusNoContent)
		})
	}

	productsRoot := router.Group("/products")
	{
		productsRoot.GET("/:id", func(c *gin.Context) {
			log.Println("GET /products/:id")

			productId := uuid.Parse(c.Param("id"))
			if productId == nil {
				c.JSON(http.StatusUnprocessableEntity, "productId must be a uuid")
			} else if strings.Compare(productId.String(), uuid.NIL.String()) == 0 {
				c.JSON(http.StatusUnprocessableEntity, "productId must cannot be an empty uuid")
			} else {
				product, err := a.ProductsService.Find(c, productId)
				if err != nil {
					c.JSON(http.StatusInternalServerError, err)
				}
				c.JSON(http.StatusOK, product)
			}
		})

		productsRoot.GET("", func(c *gin.Context) {
			log.Println("GET /products")

			products, err := a.ProductsService.FindAll(c)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusOK, products)
		})

		productsRoot.POST("", func(c *gin.Context) {
			log.Println("POST /products")

			var product dtos.Product
			err := c.ShouldBindJSON(&product)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			} else if strings.Compare(product.ProductId.String(), uuid.NIL.String()) == 0 {
				c.JSON(http.StatusUnprocessableEntity, "productId must cannot be an empty uuid")
			} else {
				result, err := a.ProductsService.Create(c, product)
				if err != nil {
					c.JSON(http.StatusInternalServerError, err)
				}
				c.JSON(http.StatusOK, result)
			}
		})

		productsRoot.DELETE(":id", func(c *gin.Context) {
			log.Println("DELETE /products")

			productId := uuid.Parse(c.Param("id"))
			if productId == nil {
				c.JSON(http.StatusUnprocessableEntity, "productId must be a uuid")
			} else if strings.Compare(productId.String(), uuid.NIL.String()) == 0 {
				c.JSON(http.StatusUnprocessableEntity, "productId must cannot be an empty uuid")
			} else {
				err := a.ProductsService.Delete(c, productId)
				if err != nil {
					c.JSON(http.StatusInternalServerError, err)
				} else {
					c.Status(http.StatusNoContent)
				}
			}
		})

		productsRoot.PUT("", func(c *gin.Context) {
			log.Println("PUT /products")

			var product dtos.Product
			err := c.ShouldBindJSON(&product)
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, err)
			} else if strings.Compare(product.ProductId.String(), uuid.NIL.String()) == 0 {
				c.JSON(http.StatusUnprocessableEntity, "productId must cannot be an empty uuid")
			} else {
				err = a.ProductsService.Update(c, product)
				if err != nil {
					c.JSON(http.StatusInternalServerError, err)
				} else {
					c.Status(http.StatusNoContent)
				}
			}
		})
	}

	return router
}

func (a api) Run() {

	router := a.Engine()
	router.Run()
}
