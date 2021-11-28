package orders

import (
	"context"
	"fmt"

	"github.com/machado-br/order-service/domain/dtos"
	"github.com/machado-br/order-service/domain/entities"
)

func GetOrders(ctx context.Context) ([]dtos.Order, error) {
	fmt.Println("service.getOrders")

	orders, err := GetAll(ctx)
	if err != nil {
		return []dtos.Order{}, err
	}

	var result []dtos.Order

	for i, _ := range orders {
		dto, err := mapToDto(ctx, orders[i])
		if err != nil {
			return []dtos.Order{}, err
		}
		result = append(result, dto)
	}
	return result, nil
}

func GetOrder(ctx context.Context, id string) (dtos.Order, error) {
	fmt.Println("service.getOrder")

	order, err := Get(ctx, id)
	if err != nil {
		return dtos.Order{}, err
	}

	dto, err := mapToDto(ctx, order)
	if err != nil {
		return dtos.Order{}, err
	}

	return dto, nil
}

func CreateOrder(ctx context.Context, order dtos.Order) (string, error) {
	fmt.Println("service.createOrder")

	entity, err := mapToEntity(ctx, order)
	if err != nil {
		return "", err
	}

	result, err := Create(ctx, entity)
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

func UpdateOrder(ctx context.Context, order dtos.Order) (string, error) {
	fmt.Println("service.updateOrder")

	orderEntity, err := Get(ctx, order.OrderId)
	if err != nil {
		return "", err
	}

	entity, err := mapToEntity(ctx, order)
	if err != nil {
		return "", err
	}

	entity.Id = orderEntity.Id

	result, err := Update(ctx, entity)
	if err != nil {
		return "", err
	}
	return result, nil
}

func mapToEntity(ctx context.Context, orderDto dtos.Order) (entities.Order, error) {
	return entities.Order{
		OrderId:         orderDto.OrderId,
		UserId:          orderDto.UserId,
		ProductId:       orderDto.ProductId,
		Quantity:        orderDto.Quantity,
		DeliveryAddress: orderDto.DeliveryAddress,
	}, nil
}

func mapToDto(ctx context.Context, orderEntity entities.Order) (dtos.Order, error) {
	return dtos.Order{
		OrderId:         orderEntity.OrderId,
		UserId:          orderEntity.UserId,
		ProductId:       orderEntity.ProductId,
		Quantity:        orderEntity.Quantity,
		DeliveryAddress: orderEntity.DeliveryAddress,
	}, nil
}
