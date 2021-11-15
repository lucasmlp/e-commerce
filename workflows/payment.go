package workflows

import (
	"strconv"
	"time"

	"github.com/pborman/uuid"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func RunPayment(ctx workflow.Context, userId uuid.UUID, orderId uuid.UUID, orderValue float64) {
	logger := buildPaymentWorkflowLogger(ctx, userId.String(), orderId.String(), strconv.FormatFloat(orderValue, 'e', 2, 1))
	logger.Info("payment workflow started")

	workflow.Sleep(ctx, time.Second*15)

	logger.Info("payment workflow finished succesfully")
}

func buildPaymentWorkflowLogger(ctx workflow.Context, userId string, orderId string, orderValue string) *zap.Logger {
	logger := workflow.GetLogger(ctx)
	workflowInfo := workflow.GetInfo(ctx)

	logger = logger.With(zap.String("WorkflowId", workflowInfo.WorkflowExecution.ID))
	logger = logger.With(zap.String("OrderId", orderId))
	logger = logger.With(zap.String("UserId", userId))
	logger = logger.With(zap.String("OrderValue", orderValue))

	return logger
}
