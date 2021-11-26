package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	orders "github.com/machado-br/order-service/domain/order"
	"github.com/machado-br/order-service/entities"
)

func main() {
	router := gin.Default()

	root := router.Group("/orders")
	{
		root.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, "pong")
		})

		root.GET("", func(c *gin.Context) {
			orders, err := orders.GetOrders(c)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusOK, orders)
		})

		root.POST("", func(c *gin.Context) {

			var order entities.Order
			err := c.ShouldBindJSON(&order)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}

			result, err := orders.CreateOrder(c, order)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusOK, result)
		})
	}

	router.Run()
}
