package dtos

import "github.com/pborman/uuid"

type Order struct {
	OrderId         string `json:"orderId" binding:"required"`
	UserId          string `json:"userId" binding:"required"`
	ProductId       string `json:"productId" binding:"required"`
	Quantity        int    `json:"quantity" binding:"required"`
	DeliveryAddress string `json:"deliveryAddress" binding:"required"`
	Status          string `json:"status" binding:"required"`
}

type Product struct {
	ProductId uuid.UUID `bson:"productId" binding:"required"`
	Name      string    `bson:"name" binding:"required"`
	Units     int       `bson:"units" binding:"required"`
	Price     int       `bson:"price" binding:"required"`
}
