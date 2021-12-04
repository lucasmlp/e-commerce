package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	Id              primitive.ObjectID `bson:"_id"`
	OrderId         string             `bson:"orderId" binding:"required"`
	UserId          string             `bson:"userId" binding:"required"`
	ProductId       string             `bson:"productId" binding:"required"`
	Quantity        int                `bson:"quantity" binding:"required"`
	DeliveryAddress string             `bson:"deliveryAddress" binding:"required"`
	Status          string             `bson:"status" binding:"required"`
}

type Product struct {
	Id        primitive.ObjectID `bson:"_id"`
	ProductId string             `bson:"productId" binding:"required"`
	Name      string             `bson:"name" binding:"required"`
	Units     int                `bson:"units" binding:"required"`
	Price     float64            `bson:"price" binding:"required"`
}
