package workflows

import (
	"time"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func RunShipment(ctx workflow.Context, orderId string, orderDeliveryAddress string) (string, error) {
	logger := buildShipmentWorkflowLogger(ctx, orderId, orderDeliveryAddress)

	logger.Info("shipment process started")

	workflow.Sleep(ctx, time.Second*15)

	logger.Info("shipment workflow finished succesfully")
	return "Success", nil
}

func buildShipmentWorkflowLogger(ctx workflow.Context, orderId string, orderDeliveryAddress string) *zap.Logger {
	logger := workflow.GetLogger(ctx)
	workflowInfo := workflow.GetInfo(ctx)

	logger = logger.With(zap.String("WorkflowId", workflowInfo.WorkflowExecution.ID))
	logger = logger.With(zap.String("OrderId", orderId))
	logger = logger.With(zap.String("DeliveryAddress", orderDeliveryAddress))

	return logger
}
