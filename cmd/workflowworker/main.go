package main

import (
	"fmt"
	"os"

	"github.com/machado-br/order-service/cadence/activities"
	"github.com/machado-br/order-service/cadence/helpers"
	"github.com/machado-br/order-service/cadence/workflows"

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

	activity.Register(activities.GetOrder)
	activity.Register(activities.GetProduct)
	activity.Register(activities.UpdateProductUnits)

	w := worker.New(workflowClient, domainName, "order-tasklist",
		worker.Options{
			Logger: logger,
		})

	err = w.Run()

	if err != nil {
		panic(err)
	}
}
