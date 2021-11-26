package orders

import (
	"context"

	"github.com/machado-br/order-service/entities"
)

func GetOrders(ctx context.Context) ([]entities.Order, error) {
	orders, err := GetAll(ctx)
	if err != nil {
		return []entities.Order{}, err
	}
	return orders, nil
}

func CreateOrder(ctx context.Context, order entities.Order) (string, error) {
	result, err := Create(ctx, order)
	if err != nil {
		return "", err
	}
	return result, nil
}
