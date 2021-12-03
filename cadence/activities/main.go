package activities

import (
	"context"

	"github.com/machado-br/order-service/domain/dtos"
	"github.com/machado-br/order-service/domain/orders"
	"github.com/machado-br/order-service/domain/products"
)

type Activities struct {
	OrderService    orders.Service
	ProductsService products.Service
}

func NewActivities(orderService orders.Service, productService products.Service) (Activities, error) {
	activities := Activities{
		OrderService:    orderService,
		ProductsService: productService,
	}
	return activities, nil
}

func (a Activities) GetOrder(ctx context.Context, orderId string) (dtos.Order, error) {

	order, err := a.OrderService.Find(ctx, orderId)
	if err != nil {
		return dtos.Order{}, err
	}
	return order, nil
}

func (a Activities) GetProduct(ctx context.Context, productId string) (dtos.Product, error) {

	product, err := a.ProductsService.Find(ctx, productId)
	if err != nil {
		return dtos.Product{}, err
	}
	return product, nil
}

func (a Activities) UpdateProduct(ctx context.Context, product dtos.Product) (string, error) {

	result, err := a.ProductsService.Update(ctx, product)
	if err != nil {
		return "", err
	}
	return result, nil
}
