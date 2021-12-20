package products

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/machado-br/e-commerce/domain/dtos"
	"github.com/machado-br/e-commerce/domain/entities"
	"github.com/machado-br/e-commerce/domain/products/mocks"
	"github.com/pborman/uuid"
	"gotest.tools/assert"
)

func TestFind(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	productIdM := uuid.NewUUID()

	t.Run("Successful responses", func(t *testing.T) {

		repoM := mocks.NewMockRepository(ctrl)
		productsService := NewService(repoM)

		want := dtos.Product{
			ProductId: productIdM,
		}

		returnM := entities.Product{
			ProductId: productIdM,
		}

		repoM.EXPECT().Find(ctx, productIdM).Return(returnM, nil)

		product, err := productsService.Find(ctx, productIdM)
		if err != nil {
			t.Fatal(err)
		}

		assert.DeepEqual(t, want, product)
	})

	t.Run("Failure responses", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)

		productIdM := uuid.NewUUID()

		repoM := mocks.NewMockRepository(ctrl)
		productsService := NewService(repoM)

		errM := errors.New("mock error")

		repoM.EXPECT().Find(ctx, productIdM).Return(entities.Product{}, errM)

		_, err := productsService.Find(ctx, productIdM)
		assert.Error(t, err, errM.Error())
	})
}

func TestFindAll(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	productIdM := uuid.NewUUID()

	t.Run("Successful responses", func(t *testing.T) {

		repoM := mocks.NewMockRepository(ctrl)
		productsService := NewService(repoM)

		want := []dtos.Product{
			{
				ProductId: productIdM,
			},
		}

		returnM := []entities.Product{
			{
				ProductId: productIdM,
			},
		}

		repoM.EXPECT().FindAll(ctx).Return(returnM, nil)

		products, err := productsService.FindAll(ctx)
		if err != nil {
			t.Fatal(err)
		}

		assert.DeepEqual(t, want, products)
	})

	t.Run("Failure responses", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)

		repoM := mocks.NewMockRepository(ctrl)
		productsService := NewService(repoM)

		errM := errors.New("mock error")

		repoM.EXPECT().FindAll(ctx).Return([]entities.Product{}, errM)

		_, err := productsService.FindAll(ctx)
		assert.Error(t, err, errM.Error())
	})
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	productIdM := uuid.NewUUID()

	t.Run("Successful responses", func(t *testing.T) {

		repoM := mocks.NewMockRepository(ctrl)
		productsService := NewService(repoM)

		payloadM := dtos.Product{
			ProductId: productIdM,
		}

		payloadRepoM := entities.Product{
			ProductId: productIdM,
		}

		returnM := "result"

		want := "result"

		repoM.EXPECT().Create(ctx, payloadRepoM).Return(returnM, nil)

		result, err := productsService.Create(ctx, payloadM)
		if err != nil {
			t.Fatal(err)
		}

		assert.DeepEqual(t, want, result)
	})

	t.Run("Failure responses", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)

		repoM := mocks.NewMockRepository(ctrl)
		productsService := NewService(repoM)

		payloadM := dtos.Product{
			ProductId: productIdM,
		}

		payloadRepoM := entities.Product{
			ProductId: productIdM,
		}

		errM := errors.New("mock error")

		repoM.EXPECT().Create(ctx, payloadRepoM).Return("", errM)

		_, err := productsService.Create(ctx, payloadM)
		assert.Error(t, err, errM.Error())
	})
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	productIdM := uuid.NewUUID()

	t.Run("Successful responses", func(t *testing.T) {

		repoM := mocks.NewMockRepository(ctrl)
		productsService := NewService(repoM)

		repoM.EXPECT().Delete(ctx, productIdM).Return(1, nil)

		err := productsService.Delete(ctx, productIdM)

		assert.NilError(t, err)
	})

	t.Run("Failure responses", func(t *testing.T) {
		ctx := context.Background()
		ctrl := gomock.NewController(t)

		productIdM := uuid.NewUUID()

		repoM := mocks.NewMockRepository(ctrl)
		productsService := NewService(repoM)

		errM := errors.New("mock error")

		repoM.EXPECT().Delete(ctx, productIdM).Return(0, errM)

		err := productsService.Delete(ctx, productIdM)
		assert.Error(t, err, errM.Error())
	})
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	productIdM := uuid.NewUUID()

	t.Run("Successful responses", func(t *testing.T) {

		repoM := mocks.NewMockRepository(ctrl)
		productsService := NewService(repoM)

		payloadM := dtos.Product{
			ProductId: productIdM,
		}

		payloadRepoM := entities.Product{
			ProductId: productIdM,
		}

		returnFindM := entities.Product{
			ProductId: productIdM,
		}

		repoM.EXPECT().Find(ctx, productIdM).Return(returnFindM, nil)

		repoM.EXPECT().Replace(ctx, payloadRepoM).Return(1, nil)

		err := productsService.Update(ctx, payloadM)
		if err != nil {
			t.Fatal(err)
		}

		assert.NilError(t, err)
	})

	t.Run("Failure responses", func(t *testing.T) {
		ctx := context.Background()

		payloadM := dtos.Product{
			ProductId: productIdM,
		}

		payloadRepoM := entities.Product{
			ProductId: productIdM,
		}

		returnFindM := entities.Product{
			ProductId: productIdM,
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

					repoM.EXPECT().Find(ctx, productIdM).Return(entities.Product{}, errM)

					return repoM
				},
			},
			{
				name: "Fail when replacing mongo document",
				err:  errM,
				repoM: func(ctx context.Context, ctrl *gomock.Controller) Repository {
					repoM := mocks.NewMockRepository(ctrl)

					repoM.EXPECT().Find(ctx, productIdM).Return(returnFindM, nil)

					repoM.EXPECT().Replace(ctx, payloadRepoM).Return(0, errM)

					return repoM
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)

				productsService := NewService(tc.repoM(ctx, ctrl))

				err := productsService.Update(ctx, payloadM)
				assert.DeepEqual(t, err.Error(), tc.err.Error())
			})
		}

	})
}
