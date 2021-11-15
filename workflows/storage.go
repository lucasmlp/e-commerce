package workflows

import (
	"time"

	"github.com/machado-br/order-service/activities"
	"go.uber.org/cadence/workflow"
)

func RunStorage(ctx workflow.Context, productId string, quantity int) (string, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info("storage workflow started")

	logger.Info("check storage availability started")
	ao := workflow.ActivityOptions{
		StartToCloseTimeout:    time.Second * 30,
		ScheduleToStartTimeout: time.Second * 120,
	}

	var units int
	ctx = workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(ctx, activities.GetProductUnits).Get(ctx, &units)
	if err != nil || units < 0 {
		return "Failure", err
	}

	logger.Info("reservation on product(s) unit(s) started")

	logger.Info("storage workflow finished succesfully")

	return "Success", nil
}
