package workflows

import (
	"errors"
	"time"

	"github.com/pborman/uuid"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

const (
	standardStepTimeout          = 60 * time.Second
	standardChildWorkflowTimeout = 30 * time.Second
)

func RunOrder(ctx workflow.Context) error {
	logger := buildOrderWorkflowLogger(ctx, uuid.New(), "10568246624")

	logger.Info("order workflow started")
	var err error

	err = handleStorageCheckAndReservation(ctx)
	if err != nil {
		return err
	}

	err = handlePaymentProcess(ctx)
	if err != nil {
		return err
	}

	err = handleShipment(ctx)
	if err != nil {
		return err
	}

	logger.Info("order workflow finished succesfully")

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

func handleStorageCheckAndReservation(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)

	logger.Info("started to check storage and reserve product(s)")

	cwo := workflow.ChildWorkflowOptions{
		// Do not specify WorkflowID if you want Cadence to generate a unique ID for the child execution.
		WorkflowID:                   uuid.New(),
		ExecutionStartToCloseTimeout: time.Minute * 30,
	}
	ctx = workflow.WithChildOptions(ctx, cwo)

	var result string
	future := workflow.ExecuteChildWorkflow(ctx, RunStorage)
	if err := future.Get(ctx, &result); err != nil {
		workflow.GetLogger(ctx).Error("SimpleChildWorkflow failed.", zap.Error(err))
		return err
	}

	signalName := "storage-check-reservation-finished"

	err := handleStandardSignal(ctx, signalName, "recieved storage-check-reservation-finished signal", "failed to check storage and reserve product(s)")
	if err != nil {
		return err
	}
	return nil
}

func handlePaymentProcess(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)

	logger.Info("started to process payment")

	signalName := "payment-finished"

	err := handleStandardSignal(ctx, signalName, "recieved payment-finished signal", "failed to process product(s) payment")
	if err != nil {
		return err
	}
	return nil
}

func handleShipment(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)

	logger.Info("shipment process started")

	signalName := "shipment-finished"

	err := handleStandardSignal(ctx, signalName, "recieved shipment-finished signal", "failed to process product(s) shipment")
	if err != nil {
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
