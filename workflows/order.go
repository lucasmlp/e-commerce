package workflows

import (
	"time"

	"github.com/pborman/uuid"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func RunOrder(ctx workflow.Context) error {
	logger := buildOrderWorkflowLogger(ctx, uuid.New(), "10568246624")

	logger.Info("Order Workflow started")

	logger.Info("Stock check and reservation started")

	workflow.Sleep(ctx, 2*time.Second)

	logger.Info("Payment process started")

	workflow.Sleep(ctx, 10*time.Second)

	logger.Info("Shipment process started")

	workflow.Sleep(ctx, 10*time.Second)

	logger.Info("Order Workflow finished")

	return nil
}

func buildOrderWorkflowLogger(ctx workflow.Context, orderID string, reference string) *zap.Logger {
	logger := workflow.GetLogger(ctx)
	workflowInfo := workflow.GetInfo(ctx)

	logger = logger.With(zap.String("WorkflowID", workflowInfo.WorkflowExecution.ID))
	logger = logger.With(zap.String("OrderID", orderID))
	logger = logger.With(zap.String("Reference", reference))

	return logger
}
