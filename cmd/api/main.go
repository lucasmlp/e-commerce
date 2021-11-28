package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/machado-br/order-service/domain/dtos"
	orders "github.com/machado-br/order-service/domain/order"
)

func main() {
	router := gin.Default()

	root := router.Group("/orders")
	{
		root.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, "pong")
		})

		root.GET("/:id", func(c *gin.Context) {
			fmt.Println("/orders/:id")

			orderId := c.Param("id")
			fmt.Printf("orderId: %v\n", orderId)

			order, err := orders.GetOrder(c, orderId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusOK, order)
		})

		root.GET("", func(c *gin.Context) {
			fmt.Println("/orders")
			orders, err := orders.GetOrders(c)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusOK, orders)
		})

		root.POST("", func(c *gin.Context) {

			var order dtos.Order
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

		root.DELETE(":id", func(c *gin.Context) {
			orderId := c.Param("id")
			err := orders.Delete(c, orderId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusNoContent, "")
		})

		root.PUT("", func(c *gin.Context) {

			var order dtos.Order
			err := c.ShouldBindJSON(&order)
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, err)
			}

			result, err := orders.UpdateOrder(c, order)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
			}
			c.JSON(http.StatusOK, result)
		})
	}

	router.Run()
}
