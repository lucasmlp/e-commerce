package main

import (
	"fmt"
	"os"

	"github.com/machado-br/order-service/activities"
	"github.com/machado-br/order-service/helpers"
	"github.com/machado-br/order-service/workflows"

	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/worker"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func main() {
	serviceNameCadenceClient := os.Getenv("CADENCE_CLIENT_NAME")
	serviceNameCadenceFrontend := os.Getenv("CADENCE_FRONTEND_NAME")
	domainName := os.Getenv("CADENCE_DOMAIN_NAME")

	fmt.Printf("serviceNameCadenceClient: %v\n", serviceNameCadenceClient)

	workflowClient, err := helpers.NewWorkflowClient(serviceNameCadenceClient, serviceNameCadenceFrontend)
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	workflow.RegisterWithOptions(workflows.RunOrder, workflow.RegisterOptions{
		EnableShortName: true,
		Name:            "RunOrder",
	})
	workflow.RegisterWithOptions(workflows.RunStorage, workflow.RegisterOptions{
		EnableShortName: true,
		Name:            "RunStorage",
	})
	workflow.RegisterWithOptions(workflows.RunPayment, workflow.RegisterOptions{
		EnableShortName: true,
		Name:            "RunPayment",
	})
	workflow.RegisterWithOptions(workflows.RunShipment, workflow.RegisterOptions{
		EnableShortName: true,
		Name:            "RunShipment",
	})

	activity.Register(activities.GetProductUnits)
	activity.Register(activities.GetOrder)

	w := worker.New(workflowClient, domainName, "order-tasklist",
		worker.Options{
			Logger: logger,
		})

	err = w.Run()

	if err != nil {
		panic(err)
	}
}
