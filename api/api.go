package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/machado-br/order-service/domain/dtos"
	orders "github.com/machado-br/order-service/domain/order"
)

type api struct {
	Debug         bool
	OrdersService orders.Service
}

func NewApi(
	debug bool,
	orderService orders.Service,
) (api, error) {
	return api{
		Debug:         debug,
		OrdersService: orderService,
	}, nil
}

func (a api) Engine() *gin.Engine {
	router := gin.Default()

	root := router.Group("/orders")
	{
		root.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, "pong")
		})

		root.GET("/:id", func(c *gin.Context) {
			log.Println("GET /orders/:id")

			orderId := c.Param("id")

			order, err := a.OrdersService.Find(c, orderId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusOK, order)
		})

		root.GET("", func(c *gin.Context) {
			log.Println("GET /orders")

			orders, err := a.OrdersService.FindAll(c)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusOK, orders)
		})

		root.POST("", func(c *gin.Context) {
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

		root.DELETE(":id", func(c *gin.Context) {
			log.Println("DELETE /orders")

			orderId := c.Param("id")
			err := a.OrdersService.Delete(c, orderId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusNoContent, "")
		})

		root.PUT("", func(c *gin.Context) {
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

	return router
}

func (a api) Run() {

	router := a.Engine()
	router.Run()
}
