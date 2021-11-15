package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/machado-br/order-service/helpers"
	"github.com/machado-br/order-service/workflows"

	"go.uber.org/cadence/client"
)

func main() {
	serviceNameCadenceClient := os.Getenv("CADENCE_CLIENT_NAME")
	serviceNameCadenceFrontend := os.Getenv("CADENCE_FRONTEND_NAME")

	action := os.Args[1]

	workflowClient, err := helpers.NewWorkflowClient(serviceNameCadenceClient, serviceNameCadenceFrontend)

	if err != nil {
		panic(err)
	}

	triggerClient := helpers.NewCadenceClient(workflowClient)

	switch name := action; name {
	case "RunOrder":
		orderId := "543dc1a5-bf19-49f1-93de-89bb3bf3785f"
		userId := "462cf3f0-c1db-4d0c-b70d-483b35f441d1"
		workflowId := orderId + ":" + userId

		_, err = triggerClient.StartWorkflow(context.Background(), client.StartWorkflowOptions{
			ID:                           workflowId,
			TaskList:                     "order-tasklist",
			ExecutionStartToCloseTimeout: 5 * time.Minute,
		}, workflows.RunOrder, orderId)

		fmt.Println("Started order workflow. Order ID: ", orderId)
	case "StorageCheckReservationFinished":
		err = triggerClient.SignalWorkflow(context.Background(), os.Args[2], "", "storage-check-reservation-finished", "Success")
	case "PaymentFinished":
		err = triggerClient.SignalWorkflow(context.Background(), os.Args[2], "", "payment-finished", "Success")
	case "ShipmentFinished":
		err = triggerClient.SignalWorkflow(context.Background(), os.Args[2], "", "shipment-finished", "Success")
	}

	if err != nil {
		panic(err)
	}

}
