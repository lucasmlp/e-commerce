package products

import (
	"context"
	"errors"
	"log"

	"github.com/machado-br/e-commerce/domain/dtos"
	"github.com/machado-br/e-commerce/domain/entities"
)

type service struct {
	repo Repository
}
type Service interface {
	FindAll(ctx context.Context) ([]dtos.Product, error)
	Find(ctx context.Context, id string) (dtos.Product, error)
	Create(ctx context.Context, product dtos.Product) (string, error)
	Delete(ctx context.Context, productId string) error
	Update(ctx context.Context, product dtos.Product) error
}

func NewService(
	repository Repository,
) Service {
	return service{
		repo: repository,
	}
}

func (s service) FindAll(ctx context.Context) ([]dtos.Product, error) {
	log.Println("service.getProducts")

	products, err := s.repo.FindAll(ctx)
	if err != nil {
		return []dtos.Product{}, err
	}

	var result []dtos.Product

	for i := range products {
		dto, err := mapToDto(ctx, products[i])
		if err != nil {
			return []dtos.Product{}, err
		}
		result = append(result, dto)
	}
	return result, nil
}

func (s service) Find(ctx context.Context, id string) (dtos.Product, error) {
	log.Println("service.getProduct")

	product, err := s.repo.Find(ctx, id)
	if err != nil {
		return dtos.Product{}, err
	}

	dto, err := mapToDto(ctx, product)
	if err != nil {
		return dtos.Product{}, err
	}

	return dto, nil
}

func (s service) Create(ctx context.Context, product dtos.Product) (string, error) {
	log.Println("service.createProduct")

	entity, err := mapToEntity(ctx, product)
	if err != nil {
		return "", err
	}

	result, err := s.repo.Create(ctx, entity)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s service) Delete(ctx context.Context, productId string) error {
	log.Println("service.deleteProduct")

	deletedCount, err := s.repo.Delete(ctx, productId)
	if err != nil {
		return err
	}

	if deletedCount == 0 {
		return errors.New("unable to delete document")
	}

	return nil
}

func (s service) Update(ctx context.Context, product dtos.Product) error {
	log.Println("service.updateProduct")

	productEntity, err := s.repo.Find(ctx, product.ProductId)
	if err != nil {
		return err
	}

	log.Printf("productEntity: %v\n", productEntity)

	entity, err := mapToEntity(ctx, product)
	if err != nil {
		return err
	}

	entity.Id = productEntity.Id

	log.Printf("entity: %v\n", entity)

	updatedCount, err := s.repo.Replace(ctx, entity)
	if err != nil {
		return err
	}

	if updatedCount == 0 {
		return errors.New("unable to update document")
	}

	return nil
}

func mapToEntity(ctx context.Context, productDto dtos.Product) (entities.Product, error) {
	return entities.Product{
		ProductId: productDto.ProductId,
		Name:      productDto.Name,
		Units:     productDto.Units,
		Price:     productDto.Price,
	}, nil
}

func mapToDto(ctx context.Context, productEntity entities.Product) (dtos.Product, error) {
	return dtos.Product{
		ProductId: productEntity.ProductId,
		Name:      productEntity.Name,
		Units:     productEntity.Units,
		Price:     productEntity.Price,
	}, nil
}
