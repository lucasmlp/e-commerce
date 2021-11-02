package main

import (
	"fmt"
	"os"

	"github.com/machado-br/order-service/helpers"

	"go.uber.org/cadence/worker"
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

	w := worker.New(workflowClient, domainName, "pocTasklist",
		worker.Options{
			Logger: logger,
		})

	err = w.Run()

	if err != nil {
		panic(err)
	}
}
