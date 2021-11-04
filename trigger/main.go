package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/machado-br/order-service/helpers"
	"github.com/machado-br/order-service/workflows"
	"github.com/pborman/uuid"
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

	workflowID := uuid.New()

	switch name := action; name {
	case "RunOrder":
		_, err = triggerClient.StartWorkflow(context.Background(), client.StartWorkflowOptions{
			ID:                           workflowID,
			TaskList:                     "order-tasklist",
			ExecutionStartToCloseTimeout: 1 * time.Second,
		}, workflows.RunOrder)
	}

	if err != nil {
		panic(err)
	}

	fmt.Println("Started Action: ", workflowID)
	fmt.Println("Action: ", action)
}
