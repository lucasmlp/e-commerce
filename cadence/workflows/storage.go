package workflows

import (
	"errors"
	"strconv"
	"time"

	"github.com/machado-br/order-service/cadence/activities"
	"github.com/machado-br/order-service/domain/entities"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func RunStorage(ctx workflow.Context, productId string, quantity int) (entities.Product, error) {
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

	var product entities.Product
	ctx = workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(ctx, activities.GetProduct, productId).Get(ctx, &product)
	if err != nil {
		return entities.Product{}, err
	}

	logger.Info("reservation on product units started")

	if product.Units < 0 || product.Units < quantity {
		logger.Error("not enough product units",
			zap.String("product units requested:", strconv.Itoa(quantity)),
			zap.String("product units available:", strconv.Itoa(product.Units)),
		)

		return entities.Product{}, errors.New("not enough product units")
	} else {
		var result string
		err := workflow.ExecuteActivity(ctx, activities.UpdateProductUnits, productId, product.Units-quantity).Get(ctx, &result)
		if err != nil {
			return entities.Product{}, err
		}
	}

	logger.Info("storage workflow finished succesfully")

	return product, nil
}
