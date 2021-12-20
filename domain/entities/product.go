package entities

import (
	"errors"

	"github.com/pborman/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id        primitive.ObjectID `bson:"_id"`
	ProductId uuid.UUID          `bson:"productId" binding:"required"`
	Name      string             `bson:"name" binding:"required"`
	Units     int                `bson:"units" binding:"required"`
	Price     int                `bson:"price" binding:"required"`
}

func (p Product) NewProduct(productId string, name string, units int, price float32) (Product, error) {

	pId := uuid.Parse(productId)
	if pId == nil {
		return Product{}, errors.New("productId must be a uuid")
	}

	err := p.validateProduct()
	if err != nil {
		return Product{}, err
	}

	priceInt := int(price * 100)

	id := primitive.NewObjectID()

	product := Product{
		Id:        id,
		ProductId: pId,
		Name:      name,
		Units:     units,
		Price:     priceInt,
	}

	return product, nil
}

func (p Product) validateProduct() error {
	if p.Price <= 0 {
		return errors.New("product price cannot be less or equal to zero")
	}

	if p.Units <= 0 {
		return errors.New("product units cannot be less or equal to zero")
	}
	return nil
}
