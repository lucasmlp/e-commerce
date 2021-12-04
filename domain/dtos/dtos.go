package dtos

type Order struct {
	OrderId         string `json:"orderId" binding:"required"`
	UserId          string `json:"userId" binding:"required"`
	ProductId       string `json:"productId" binding:"required"`
	Quantity        int    `json:"quantity" binding:"required"`
	DeliveryAddress string `json:"deliveryAddress" binding:"required"`
	Status          string `json:"status" binding:"required"`
}

type Product struct {
	ProductId string  `json:"productId" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Units     int     `json:"units" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
}
