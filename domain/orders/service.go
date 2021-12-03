package orders

import (
	"context"
	"log"

	"github.com/machado-br/order-service/domain/dtos"
	"github.com/machado-br/order-service/domain/entities"
)

type service struct {
	repo Repository
}
type Service interface {
	FindAll(ctx context.Context) ([]dtos.Order, error)
	Find(ctx context.Context, id string) (dtos.Order, error)
	Create(ctx context.Context, order dtos.Order) (string, error)
	Delete(ctx context.Context, orderId string) error
	Update(ctx context.Context, order dtos.Order) (string, error)
}

func NewService(
	repository Repository,
) Service {
	return service{
		repo: repository,
	}
}

func (s service) FindAll(ctx context.Context) ([]dtos.Order, error) {
	log.Println("service.getOrders")

	orders, err := s.repo.GetAll(ctx)
	if err != nil {
		return []dtos.Order{}, err
	}

	var result []dtos.Order

	for i := range orders {
		dto, err := mapToDto(ctx, orders[i])
		if err != nil {
			return []dtos.Order{}, err
		}
		result = append(result, dto)
	}
	return result, nil
}

func (s service) Find(ctx context.Context, id string) (dtos.Order, error) {
	log.Println("service.getOrder")

	order, err := s.repo.Get(ctx, id)
	if err != nil {
		return dtos.Order{}, err
	}

	dto, err := mapToDto(ctx, order)
	if err != nil {
		return dtos.Order{}, err
	}

	return dto, nil
}

func (s service) Create(ctx context.Context, order dtos.Order) (string, error) {
	log.Println("service.createOrder")

	entity, err := mapToEntity(ctx, order)
	if err != nil {
		return "", err
	}

	result, err := s.repo.Create(ctx, entity)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s service) Delete(ctx context.Context, orderId string) error {
	log.Println("service.deleteOrder")

	err := s.repo.Delete(ctx, orderId)
	if err != nil {
		return err
	}

	return nil
}

func (s service) Update(ctx context.Context, order dtos.Order) (string, error) {
	log.Println("service.updateOrder")

	orderEntity, err := s.repo.Get(ctx, order.OrderId)
	if err != nil {
		return "", err
	}

	entity, err := mapToEntity(ctx, order)
	if err != nil {
		return "", err
	}

	entity.Id = orderEntity.Id

	result, err := s.repo.Replace(ctx, entity)
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
