package entities

type Order struct {
	Id              string
	UserId          string
	ProductId       string
	Quantity        int
	DeliveryAddress string
}

type Product struct {
	Id    string
	Name  string
	Units int
	Price float64
}
