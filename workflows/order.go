package workflows

import (
	"errors"

	"github.com/pborman/uuid"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func RunOrder(ctx workflow.Context) error {
	logger := buildOrderWorkflowLogger(ctx, uuid.New(), "10568246624")

	logger.Info("order workflow started")
	var err error

	err = handleStockCheckAndReservation(ctx)
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

func handleStockCheckAndReservation(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)

	logger.Info("Started to check stock and reserve product(s)")

	var signalVal string
	signalName := "stock-check-reservation-finished"
	signalChan := workflow.GetSignalChannel(ctx, signalName)
	s := workflow.NewSelector(ctx)

	s.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &signalVal)
		logger.Info("Recieved stock-check-reservation-finished signal",
			zap.String("Signal value: ", signalVal))
	})

	s.Select(ctx)

	if signalVal == "Success" {
		return nil
	} else {
		return errors.New("Failed to check stock and reserve product(s)")
	}
}

func handlePaymentProcess(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)

	logger.Info("Started process payment")

	var signalVal string
	signalName := "payment-finished"
	signalChan := workflow.GetSignalChannel(ctx, signalName)
	s := workflow.NewSelector(ctx)

	s.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &signalVal)
		logger.Info("Recieved payment-finished signal",
			zap.String("Signal value: ", signalVal))
	})

	s.Select(ctx)

	if signalVal == "Success" {
		return nil
	} else {
		return errors.New("Failed to process product(s) payment")
	}
}

func handleShipment(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)

	logger.Info("Shipment process started")

	var signalVal string
	signalName := "shipment-finished"
	signalChan := workflow.GetSignalChannel(ctx, signalName)
	s := workflow.NewSelector(ctx)

	s.AddReceive(signalChan, func(c workflow.Channel, more bool) {
		c.Receive(ctx, &signalVal)
		logger.Info("Recieved shipment-finished signal",
			zap.String("Signal value: ", signalVal))
	})

	s.Select(ctx)

	if signalVal == "Success" {
		return nil
	} else {
		return errors.New("Failed to process product(s) shipment")
	}
}
