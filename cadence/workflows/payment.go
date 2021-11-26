package workflows

import (
	"strconv"
	"time"

	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func RunPayment(ctx workflow.Context, orderId string, userId string, orderValue float64) (string, error) {
	logger := buildPaymentWorkflowLogger(ctx, userId, orderId, strconv.FormatFloat(orderValue, 'e', 2, 64))

	logger.Info("payment workflow started")

	workflow.Sleep(ctx, time.Second*15)

	logger.Info("payment workflow finished succesfully")
	return "Success", nil
}

func buildPaymentWorkflowLogger(ctx workflow.Context, orderId string, userId string, orderValue string) *zap.Logger {
	logger := workflow.GetLogger(ctx)
	workflowInfo := workflow.GetInfo(ctx)

	logger = logger.With(zap.String("WorkflowId", workflowInfo.WorkflowExecution.ID))
	logger = logger.With(zap.String("OrderId", orderId))
	logger = logger.With(zap.String("UserId", userId))
	logger = logger.With(zap.String("OrderValue", orderValue))

	return logger
}
