package workflows

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/machado-br/e-commerce/cadence/activities"
	"github.com/machado-br/e-commerce/domain/dtos"
	"github.com/machado-br/e-commerce/domain/products"
	"github.com/pborman/uuid"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

type StorageWorkflow struct {
	ProductsService products.Service
	Activities      activities.Activities
}

func NewStorageWorkflow(productsService products.Service, activities activities.Activities) (StorageWorkflow, error) {
	storageWorkflow := StorageWorkflow{
		ProductsService: productsService,
		Activities:      activities,
	}
	workflow.RegisterWithOptions(storageWorkflow.RunStorage, workflow.RegisterOptions{
		EnableShortName: true,
		Name:            "RunStorage",
	})

	return storageWorkflow, nil
}

func (s StorageWorkflow) RunStorage(ctx workflow.Context, productId string, quantity int) (dtos.Product, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info("storage workflow started",
		zap.String("ProductId:", productId),
		zap.String("Quantity:", strconv.Itoa(quantity)),
	)

	logger.Info("check storage availability started")
	ao := workflow.ActivityOptions{
		StartToCloseTimeout:    time.Second * 30,
		ScheduleToStartTimeout: time.Second * 120,
	}

	var product dtos.Product
	ctx = workflow.WithActivityOptions(ctx, ao)
	future := workflow.ExecuteActivity(ctx, s.Activities.GetProduct, productId)
	if err := future.Get(ctx, &product); err != nil {
		workflow.GetLogger(ctx).Error("check storage availability failed.", zap.Error(err))
		return dtos.Product{}, err
	}

	if strings.Compare(product.ProductId.String(), uuid.NIL.String()) == 0 || product.ProductId.String() == uuid.NIL.String() {
		logger.Error("product not found",
			zap.String("product id:", product.ProductId.String()),
		)

		return dtos.Product{}, errors.New("product not found")
	}

	logger.Info("reservation on product units started")

	if product.Units < 0 || product.Units < quantity {
		logger.Error("not enough product units",
			zap.String("product units requested:", strconv.Itoa(quantity)),
			zap.String("product units available:", strconv.Itoa(product.Units)),
		)

		return dtos.Product{}, errors.New("not enough product units")
	} else {
		product.Units = product.Units - quantity
		var result int
		future := workflow.ExecuteActivity(ctx, s.Activities.UpdateProduct, product)
		if err := future.Get(ctx, &result); err != nil {
			workflow.GetLogger(ctx).Error("reservation on product units failed.", zap.Error(err))
			return dtos.Product{}, err
		}
	}

	logger.Info("storage workflow finished succesfully")

	return product, nil
}
