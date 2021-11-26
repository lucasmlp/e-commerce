package orders

import (
	"context"
	"fmt"

	"github.com/machado-br/order-service/domain/entities"
)

func GetOrders(ctx context.Context) ([]entities.Order, error) {
	fmt.Println("service.getOrders")

	orders, err := GetAll(ctx)
	if err != nil {
		return []entities.Order{}, err
	}
	return orders, nil
}

func GetOrder(ctx context.Context, id string) (entities.Order, error) {
	fmt.Println("service.getOrder")

	order, err := Get(ctx, id)
	if err != nil {
		return entities.Order{}, err
	}

	fmt.Printf("order: %v\n", order)

	return order, nil
}

func CreateOrder(ctx context.Context, order entities.Order) (string, error) {
	fmt.Println("service.createOrder")

	result, err := Create(ctx, order)
	if err != nil {
		return "", err
	}
	return result, nil
}

func DeleteOrder(ctx context.Context, orderId string) error {
	fmt.Println("service.deleteOrder")

	err := Delete(ctx, orderId)
	if err != nil {
		return err
	}

	return nil
}

func UpdateOrder(ctx context.Context, order entities.Order) (string, error) {
	fmt.Println("service.updateOrder")

	result, err := Update(ctx, order)
	if err != nil {
		return "", err
	}
	return result, nil
}
