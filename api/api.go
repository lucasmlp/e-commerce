package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/machado-br/order-service/domain/dtos"
	"github.com/machado-br/order-service/domain/orders"
	"github.com/machado-br/order-service/domain/products"
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
			c.JSON(http.StatusNoContent, "")
		})

		ordersRoot.PUT("", func(c *gin.Context) {
			log.Println("PUT /orders")

			var order dtos.Order
			err := c.ShouldBindJSON(&order)
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, err)
			}

			result, err := a.OrdersService.Update(c, order)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusOK, result)
		})
	}

	productsRoot := router.Group("/products")
	{
		productsRoot.GET("/:id", func(c *gin.Context) {
			log.Println("GET /products/:id")

			productId := c.Param("id")

			product, err := a.ProductsService.Find(c, productId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusOK, product)
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
			}

			result, err := a.ProductsService.Create(c, product)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusOK, result)
		})

		productsRoot.DELETE(":id", func(c *gin.Context) {
			log.Println("DELETE /products")

			productId := c.Param("id")
			err := a.ProductsService.Delete(c, productId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusNoContent, "")
		})

		productsRoot.PUT("", func(c *gin.Context) {
			log.Println("PUT /products")

			var product dtos.Product
			err := c.ShouldBindJSON(&product)
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, err)
			}

			result, err := a.ProductsService.Update(c, product)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusOK, result)
		})
	}

	return router
}

func (a api) Run() {

	router := a.Engine()
	router.Run()
}
