package workflows

import (
	"go.uber.org/cadence/workflow"
)

func RunStorage(ctx workflow.Context) (string, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info("storage workflow started")

	logger.Info("check storage availability started")

	logger.Info("reservation on product(s) unit(s) started")

	logger.Info("storage workflow finished succesfully")

	return "Success", nil
}
