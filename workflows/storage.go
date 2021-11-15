package workflows

import (
	"time"

	"github.com/machado-br/order-service/activities"
	"github.com/machado-br/order-service/entities"
	"go.uber.org/cadence/workflow"
)

func RunStorage(ctx workflow.Context, productId string, quantity int) (entities.Product, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info("storage workflow started")

	logger.Info("check storage availability started")
	ao := workflow.ActivityOptions{
		StartToCloseTimeout:    time.Second * 30,
		ScheduleToStartTimeout: time.Second * 120,
	}

	var product entities.Product
	ctx = workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(ctx, activities.GetProduct).Get(ctx, &product)
	if err != nil || product.Units < 0 {
		return entities.Product{}, err
	}

	logger.Info("reservation on product(s) unit(s) started")

	logger.Info("storage workflow finished succesfully")

	return product, nil
}
