package main

import (
	"fmt"
	"log"
	"os"

	"github.com/machado-br/order-service/cadence/activities"
	"github.com/machado-br/order-service/cadence/helpers"
	"github.com/machado-br/order-service/cadence/workflows"
	"github.com/machado-br/order-service/domain/products"

	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/worker"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func main() {

	mongoUri := os.Getenv("MONGODB_URI")
	productsDatabaseName := os.Getenv("PRODUCTS_DATABASE_NAME")
	productsCollectionName := os.Getenv("PRODUCTS_COLLECTION_NAME")

	productsRepository, err := products.NewRepository(mongoUri, productsDatabaseName, productsCollectionName)
	if err != nil {
		log.Fatalln(err)
	}

	productsService := products.NewService(productsRepository)
	if err != nil {
		log.Fatalln(err)
	}

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

	storageWorkflow, err := workflows.NewStorageWorkflow(productsService)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = workflows.NewOrderWorkflow(storageWorkflow)
	if err != nil {
		log.Fatalln(err)
	}

	workflow.RegisterWithOptions(workflows.RunPayment, workflow.RegisterOptions{
		EnableShortName: true,
		Name:            "RunPayment",
	})
	workflow.RegisterWithOptions(workflows.RunShipment, workflow.RegisterOptions{
		EnableShortName: true,
		Name:            "RunShipment",
	})

	activity.Register(activities.GetOrder)
	activity.Register(storageWorkflow.ProductsService.Find)
	activity.Register(storageWorkflow.ProductsService.Update)

	w := worker.New(workflowClient, domainName, "order-tasklist",
		worker.Options{
			Logger: logger,
		})

	err = w.Run()

	if err != nil {
		panic(err)
	}
}
