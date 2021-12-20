package workflows

import (
	"errors"
	"strings"
	"time"

	"github.com/machado-br/e-commerce/cadence/activities"
	"github.com/machado-br/e-commerce/domain/dtos"
	"github.com/pborman/uuid"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

type OrderWorkflow struct {
	StorageWorkflow StorageWorkflow
	Activities      activities.Activities
}

func NewOrderWorkflow(storageWorkflow StorageWorkflow, activities activities.Activities) (OrderWorkflow, error) {
	orderWorkflow := OrderWorkflow{
		StorageWorkflow: storageWorkflow,
		Activities:      activities,
	}

	workflow.RegisterWithOptions(orderWorkflow.RunOrder, workflow.RegisterOptions{
		EnableShortName: true,
		Name:            "RunOrder",
	})

	return orderWorkflow, nil
}

func (o OrderWorkflow) RunOrder(ctx workflow.Context, orderId string) error {
	logger := o.buildOrderWorkflowLogger(ctx, orderId)

	logger.Info("order workflow started")

	ao := workflow.ActivityOptions{
		StartToCloseTimeout:    time.Second * 5,
		ScheduleToStartTimeout: time.Second * 30,
	}

	var order dtos.Order
	var err error

	ctx = workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(ctx, o.Activities.GetOrder, orderId).Get(ctx, &order)
	if err != nil {
		return err
	}

	if strings.Compare(order.OrderId, "") == 0 || order.OrderId == string(uuid.NIL) {
		logger.Error("order not found",
			zap.String("order id:", order.OrderId),
		)

		return errors.New("order not found")
	}

	product, err := o.handleStorageCheckAndReservation(ctx, order.ProductId, order.Quantity)
	if err != nil {
		return err
	}

	err = o.handleUpdateOrderStatus(ctx, orderId, "payment pending")
	if err != nil {
		return err
	}

	err = o.handlePaymentProcess(ctx, order.OrderId, order.UserId, float64(product.Price*order.Quantity/100))
	if err != nil {
		return err
	}

	err = o.handleUpdateOrderStatus(ctx, orderId, "shipment sent")
	if err != nil {
		return err
	}

	err = o.handleShipment(ctx, order.OrderId, order.DeliveryAddress)
	if err != nil {
		return err
	}

	err = o.handleUpdateOrderStatus(ctx, orderId, "completed")
	if err != nil {
		return err
	}

	logger.Info("order workflow finished succesfully")

	return nil
}

func (o OrderWorkflow) buildOrderWorkflowLogger(ctx workflow.Context, orderId string) *zap.Logger {
	logger := workflow.GetLogger(ctx)
	workflowInfo := workflow.GetInfo(ctx)

	logger = logger.With(zap.String("WorkflowID", workflowInfo.WorkflowExecution.ID))
	logger = logger.With(zap.String("OrderID", orderId))

	return logger
}

func (o OrderWorkflow) handleStorageCheckAndReservation(ctx workflow.Context, productId string, quantity int) (dtos.Product, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info("started to check storage and reserve product(s)")

	cwo := workflow.ChildWorkflowOptions{
		WorkflowID:                   uuid.New(),
		ExecutionStartToCloseTimeout: time.Minute * 2,
	}
	ctx = workflow.WithChildOptions(ctx, cwo)

	var product dtos.Product
	future := workflow.ExecuteChildWorkflow(ctx, o.StorageWorkflow.RunStorage, productId, quantity)
	if err := future.Get(ctx, &product); err != nil {
		workflow.GetLogger(ctx).Error("storage workflow failed.", zap.Error(err))
		return dtos.Product{}, err
	}

	return product, nil
}

func (o OrderWorkflow) handlePaymentProcess(ctx workflow.Context, orderId string, userId string, orderValue float64) error {
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

func (o OrderWorkflow) handleShipment(ctx workflow.Context, orderId string, orderDeliveryAddress string) error {
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

func (o OrderWorkflow) handleUpdateOrderStatus(ctx workflow.Context, orderId string, status string) error {
	var orderWithStatusUpdated int
	future := workflow.ExecuteActivity(ctx, o.Activities.UpdateOrderStatus, orderId, status)
	err := future.Get(ctx, &orderWithStatusUpdated)
	if err != nil {
		return err
	}
	if orderWithStatusUpdated != 1 {
		return errors.New("failed to update order status")
	}
	return nil
}
