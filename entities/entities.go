package entities

type Order struct {
	Id              string `json:"id" bson:"id" binding:"required"`
	UserId          string `json:"userId" bson:"userId" binding:"required"`
	ProductId       string `json:"productId" bson:"productId" binding:"required"`
	Quantity        int    `json:"quantity" bson:"quantity" binding:"required"`
	DeliveryAddress string `json:"deliveryAddress" bson:"deliveryAddress" binding:"required"`
}

type Product struct {
	Id    string  `json:"id" binding:"required"`
	Name  string  `json:"name" binding:"required"`
	Units int     `json:"units" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}
