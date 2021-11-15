package workflows

import (
	"errors"
	"time"

	"github.com/machado-br/order-service/activities"
	entities "github.com/machado-br/order-service/entities"
	"github.com/pborman/uuid"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

const (
	standardStepTimeout          = 90 * time.Second
	standardChildWorkflowTimeout = 30 * time.Second
)

func RunOrder(ctx workflow.Context, orderId string) error {
	logger := buildOrderWorkflowLogger(ctx, orderId)

	logger.Info("order workflow started")

	ao := workflow.ActivityOptions{
		StartToCloseTimeout:    time.Second * 5,
		ScheduleToStartTimeout: time.Second * 30,
	}

	var order entities.Order
	var err error

	ctx = workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(ctx, activities.GetOrder, orderId).Get(ctx, &order)
	if err != nil {
		return err
	}

	product, err := handleStorageCheckAndReservation(ctx, order.ProductId, order.Quantity)
	if err != nil {
		return err
	}

	err = handlePaymentProcess(ctx, order.Id, order.UserId, product.Price*float64(order.Quantity))
	if err != nil {
		return err
	}

	err = handleShipment(ctx, order.Id, order.DeliveryAddress)
	if err != nil {
		return err
	}

	logger.Info("order workflow finished succesfully")

	return nil
}

func buildOrderWorkflowLogger(ctx workflow.Context, orderId string) *zap.Logger {
	logger := workflow.GetLogger(ctx)
	workflowInfo := workflow.GetInfo(ctx)

	logger = logger.With(zap.String("WorkflowID", workflowInfo.WorkflowExecution.ID))
	logger = logger.With(zap.String("OrderID", orderId))

	return logger
}

func handleStorageCheckAndReservation(ctx workflow.Context, productId string, quantity int) (entities.Product, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info("started to check storage and reserve product(s)")

	cwo := workflow.ChildWorkflowOptions{
		WorkflowID:                   uuid.New(),
		ExecutionStartToCloseTimeout: time.Minute * 2,
	}
	ctx = workflow.WithChildOptions(ctx, cwo)

	var product entities.Product
	future := workflow.ExecuteChildWorkflow(ctx, RunStorage, productId, quantity)
	if err := future.Get(ctx, &product); err != nil {
		workflow.GetLogger(ctx).Error("storage workflow failed.", zap.Error(err))
		return entities.Product{}, err
	}

	return product, nil
}

func handlePaymentProcess(ctx workflow.Context, orderId string, userId string, orderValue float64) error {
	logger := workflow.GetLogger(ctx)

	logger.Info("started to process payment")

	cwo := workflow.ChildWorkflowOptions{
		WorkflowID:                   uuid.New(),
		ExecutionStartToCloseTimeout: time.Minute * 1,
	}
	ctx = workflow.WithChildOptions(ctx, cwo)

	var result string
	future := workflow.ExecuteChildWorkflow(ctx, RunPayment, orderId, userId, orderValue)
	if err := future.Get(ctx, &result); err != nil {
		workflow.GetLogger(ctx).Error("payment workflow failed.", zap.Error(err))
		return err
	}
	return nil
}

func handleShipment(ctx workflow.Context, orderId string, orderDeliveryAddress string) error {
	logger := workflow.GetLogger(ctx)

	logger.Info("shipment process started")

	cwo := workflow.ChildWorkflowOptions{
		WorkflowID:                   uuid.New(),
		ExecutionStartToCloseTimeout: time.Minute * 1,
	}
	ctx = workflow.WithChildOptions(ctx, cwo)

	var result string
	future := workflow.ExecuteChildWorkflow(ctx, RunShipment, orderId, orderDeliveryAddress)
	if err := future.Get(ctx, &result); err != nil {
		workflow.GetLogger(ctx).Error("shipment workflow failed.", zap.Error(err))
		return err
	}
	return nil
}

func handleStandardSignal(ctx workflow.Context, signalName string, receivedSignalMsg string, failedMsg string) error {
	logger := workflow.GetLogger(ctx)
	var signalVal string
	var receivedSignal bool
	signalChan := workflow.GetSignalChannel(ctx, signalName)
	timerCtx, cancelTimer := workflow.WithCancel(ctx)
	stepTimer := workflow.NewTimer(timerCtx, standardStepTimeout)
	s := workflow.NewSelector(ctx)
	s.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &signalVal)
		receivedSignal = true
		cancelTimer()
		logger.Info(receivedSignalMsg,
			zap.String("Signal value: ", signalVal))
	})

	s.AddFuture(stepTimer, func(f workflow.Future) {
		logger.Info("Step timed out")
	})

	s.Select(ctx)

	if receivedSignal && signalVal == "Success" {
		return nil
	} else {
		return errors.New(failedMsg)
	}
}
