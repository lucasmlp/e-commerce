package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	Id              primitive.ObjectID `bson:"_id"`
	OrderId         string             `json:"orderId" bson:"orderId" binding:"required"`
	UserId          string             `json:"userId" bson:"userId" binding:"required"`
	ProductId       string             `json:"productId" bson:"productId" binding:"required"`
	Quantity        int                `json:"quantity" bson:"quantity" binding:"required"`
	DeliveryAddress string             `json:"deliveryAddress" bson:"deliveryAddress" binding:"required"`
}

type Product struct {
	Id        primitive.ObjectID `bson:"_id"`
	ProductId string             `json:"productId" bson:"productId" binding:"required"`
	Name      string             `json:"name" bson:"name" binding:"required"`
	Units     int                `json:"units" bson:"units" binding:"required"`
	Price     float64            `json:"price" bson:"price" binding:"required"`
}
