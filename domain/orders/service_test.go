package orders

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/machado-br/order-service/domain/dtos"
	"github.com/machado-br/order-service/domain/entities"
	"github.com/machado-br/order-service/domain/orders/mocks"
	"github.com/pborman/uuid"
	"gotest.tools/assert"
)

func TestFind_Success(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	orderIdM := uuid.New()
	userIdM := uuid.New()
	productIdM := uuid.New()
	quantityM := 4
	deliveryAddressM := "Rua Dona Queridinha, 180, Itapoã, Belo Horizonte, Minas Gerais, Brasil"
	statusM := "pending"

	t.Run("Successful responses", func(t *testing.T) {

		repoM := mocks.NewMockRepository(ctrl)
		ordersService := NewService(repoM)

		want := dtos.Order{
			OrderId:         orderIdM,
			UserId:          userIdM,
			ProductId:       productIdM,
			Quantity:        quantityM,
			DeliveryAddress: deliveryAddressM,
			Status:          statusM,
		}

		returnM := entities.Order{
			OrderId:         orderIdM,
			UserId:          userIdM,
			ProductId:       productIdM,
			Quantity:        quantityM,
			DeliveryAddress: deliveryAddressM,
			Status:          statusM,
		}

		repoM.EXPECT().Find(ctx, orderIdM).Return(returnM, nil)

		order, err := ordersService.Find(ctx, orderIdM)
		if err != nil {
			t.Fatal(err)
		}

		assert.DeepEqual(t, want, order)
	})

	t.Run("Failure responses", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)

		orderIdM := uuid.New()

		repoM := mocks.NewMockRepository(ctrl)
		ordersService := NewService(repoM)

		errM := errors.New("mock error")

		repoM.EXPECT().Find(ctx, orderIdM).Return(entities.Order{}, errM)

		_, err := ordersService.Find(ctx, orderIdM)
		assert.Error(t, err, errM.Error())
	})
}

func TestFindAll_Success(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	orderIdM := uuid.New()
	userIdM := uuid.New()
	productIdM := uuid.New()
	quantityM := 4
	deliveryAddressM := "Rua Dona Queridinha, 180, Itapoã, Belo Horizonte, Minas Gerais, Brasil"
	statusM := "pending"

	t.Run("Successful responses", func(t *testing.T) {

		repoM := mocks.NewMockRepository(ctrl)
		ordersService := NewService(repoM)

		want := []dtos.Order{
			{
				OrderId:         orderIdM,
				UserId:          userIdM,
				ProductId:       productIdM,
				Quantity:        quantityM,
				DeliveryAddress: deliveryAddressM,
				Status:          statusM,
			},
		}

		returnM := []entities.Order{
			{
				OrderId:         orderIdM,
				UserId:          userIdM,
				ProductId:       productIdM,
				Quantity:        quantityM,
				DeliveryAddress: deliveryAddressM,
				Status:          statusM,
			},
		}

		repoM.EXPECT().FindAll(ctx).Return(returnM, nil)

		orders, err := ordersService.FindAll(ctx)
		if err != nil {
			t.Fatal(err)
		}

		assert.DeepEqual(t, want, orders)
	})

	t.Run("Failure responses", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)

		repoM := mocks.NewMockRepository(ctrl)
		ordersService := NewService(repoM)

		errM := errors.New("mock error")

		repoM.EXPECT().FindAll(ctx).Return([]entities.Order{}, errM)

		_, err := ordersService.FindAll(ctx)
		assert.Error(t, err, errM.Error())
	})
}

func TestCreate_Success(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	orderIdM := uuid.New()
	userIdM := uuid.New()
	productIdM := uuid.New()
	quantityM := 4
	deliveryAddressM := "Rua Dona Queridinha, 180, Itapoã, Belo Horizonte, Minas Gerais, Brasil"
	statusM := "pending"

	t.Run("Successful responses", func(t *testing.T) {

		repoM := mocks.NewMockRepository(ctrl)
		ordersService := NewService(repoM)

		payloadM := dtos.Order{
			OrderId:         orderIdM,
			UserId:          userIdM,
			ProductId:       productIdM,
			Quantity:        quantityM,
			DeliveryAddress: deliveryAddressM,
			Status:          statusM,
		}

		payloadRepoM := entities.Order{
			OrderId:         orderIdM,
			UserId:          userIdM,
			ProductId:       productIdM,
			Quantity:        quantityM,
			DeliveryAddress: deliveryAddressM,
			Status:          statusM,
		}

		returnM := "result"

		want := "result"

		repoM.EXPECT().Create(ctx, payloadRepoM).Return(returnM, nil)

		result, err := ordersService.Create(ctx, payloadM)
		if err != nil {
			t.Fatal(err)
		}

		assert.DeepEqual(t, want, result)
	})

	t.Run("Failure responses", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)

		repoM := mocks.NewMockRepository(ctrl)
		ordersService := NewService(repoM)

		payloadM := dtos.Order{
			OrderId:         orderIdM,
			UserId:          userIdM,
			ProductId:       productIdM,
			Quantity:        quantityM,
			DeliveryAddress: deliveryAddressM,
			Status:          statusM,
		}

		payloadRepoM := entities.Order{
			OrderId:         orderIdM,
			UserId:          userIdM,
			ProductId:       productIdM,
			Quantity:        quantityM,
			DeliveryAddress: deliveryAddressM,
			Status:          statusM,
		}

		errM := errors.New("mock error")

		repoM.EXPECT().Create(ctx, payloadRepoM).Return("", errM)

		_, err := ordersService.Create(ctx, payloadM)
		assert.Error(t, err, errM.Error())
	})
}
