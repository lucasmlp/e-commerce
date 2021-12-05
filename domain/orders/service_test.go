package orders

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/machado-br/e-commerce/domain/dtos"
	"github.com/machado-br/e-commerce/domain/entities"
	"github.com/machado-br/e-commerce/domain/orders/mocks"
	"github.com/pborman/uuid"
	"gotest.tools/assert"
)

func TestFind(t *testing.T) {
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

func TestFindAll(t *testing.T) {
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

func TestCreate(t *testing.T) {
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

func TestDelete(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	orderIdM := uuid.New()

	t.Run("Successful responses", func(t *testing.T) {

		repoM := mocks.NewMockRepository(ctrl)
		ordersService := NewService(repoM)

		repoM.EXPECT().Delete(ctx, orderIdM).Return(1, nil)

		err := ordersService.Delete(ctx, orderIdM)

		assert.NilError(t, err)
	})

	t.Run("Failure responses", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)

		orderIdM := uuid.New()

		repoM := mocks.NewMockRepository(ctrl)
		ordersService := NewService(repoM)

		errM := errors.New("mock error")

		repoM.EXPECT().Delete(ctx, orderIdM).Return(0, errM)

		err := ordersService.Delete(ctx, orderIdM)
		assert.Error(t, err, errM.Error())
	})
}

func TestUpdate(t *testing.T) {
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

		returnFindM := entities.Order{
			OrderId:         orderIdM,
			UserId:          userIdM,
			ProductId:       productIdM,
			Quantity:        quantityM,
			DeliveryAddress: deliveryAddressM,
			Status:          statusM,
		}

		returnM := 1

		repoM.EXPECT().Find(ctx, orderIdM).Return(returnFindM, nil)

		repoM.EXPECT().Replace(ctx, payloadRepoM).Return(returnM, nil)

		err := ordersService.Update(ctx, payloadM)
		if err != nil {
			t.Fatal(err)
		}

		assert.NilError(t, err)
	})

	t.Run("Failure responses", func(t *testing.T) {
		ctx := context.Background()

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

		returnFindM := entities.Order{
			OrderId:         orderIdM,
			UserId:          userIdM,
			ProductId:       productIdM,
			Quantity:        quantityM,
			DeliveryAddress: deliveryAddressM,
			Status:          statusM,
		}

		errM := errors.New("mock error")

		tt := []struct {
			name  string
			err   error
			repoM func(context.Context, *gomock.Controller) Repository
		}{
			{
				name: "Fail when finding document",
				err:  errM,
				repoM: func(ctx context.Context, ctrl *gomock.Controller) Repository {
					repoM := mocks.NewMockRepository(ctrl)

					repoM.EXPECT().Find(ctx, orderIdM).Return(entities.Order{}, errM)

					return repoM
				},
			},
			{
				name: "Fail when replacing mongo document",
				err:  errM,
				repoM: func(ctx context.Context, ctrl *gomock.Controller) Repository {
					repoM := mocks.NewMockRepository(ctrl)

					repoM.EXPECT().Find(ctx, orderIdM).Return(returnFindM, nil)

					repoM.EXPECT().Replace(ctx, payloadRepoM).Return(0, errM)

					return repoM
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)

				ordersService := NewService(tc.repoM(ctx, ctrl))

				err := ordersService.Update(ctx, payloadM)
				assert.DeepEqual(t, err.Error(), tc.err.Error())
			})
		}

	})
}

func TestUpdateStatus(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	orderIdM := uuid.New()
	statusM := "payment pending"
	statusPayloadM := "payment pending"

	repoM := mocks.NewMockRepository(ctrl)
	ordersService := NewService(repoM)

	t.Run("Successful responses", func(t *testing.T) {

		repoM.EXPECT().UpdateStatus(ctx, orderIdM, statusM).Return(1, nil)

		ordersUpdated, err := ordersService.UpdateStatus(ctx, orderIdM, statusPayloadM)
		if err != nil {
			t.Fatal(err)
		}

		assert.DeepEqual(t, 1, ordersUpdated)
	})

	t.Run("Failure responses", func(t *testing.T) {

		errM := errors.New("mock error")

		repoM.EXPECT().UpdateStatus(ctx, orderIdM, statusM).Return(0, errM)

		_, err := ordersService.UpdateStatus(ctx, orderIdM, statusPayloadM)
		assert.Error(t, err, errM.Error())
	})
}
