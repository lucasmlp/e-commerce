package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/machado-br/e-commerce/domain/dtos"
)

func (a api) getAllOrders(c *gin.Context) {
	log.Println("GET /orders")

	orders, err := a.OrdersService.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, orders)
}

func (a api) getOrderById(c *gin.Context) {
	log.Println("GET /orders/:id")

	orderId := c.Param("id")

	order, err := a.OrdersService.Find(c, orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, order)
}

func (a api) createOrder(c *gin.Context) {
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
}

func (a api) deleteOrder(c *gin.Context) {
	log.Println("DELETE /orders/:id")

	orderId := c.Param("id")
	err := a.OrdersService.Delete(c, orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.Status(http.StatusNoContent)
}

func (a api) updateOrder(c *gin.Context) {
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
}
