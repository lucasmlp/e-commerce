package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/machado-br/order-service/cadence/helpers"
	"github.com/machado-br/order-service/cadence/workflows"
	"github.com/machado-br/order-service/domain/products"
	"github.com/pborman/uuid"

	"go.uber.org/cadence/client"
)

func main() {
	serviceNameCadenceClient := os.Getenv("CADENCE_CLIENT_NAME")
	serviceNameCadenceFrontend := os.Getenv("CADENCE_FRONTEND_NAME")
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

	storageWorkflow, err := workflows.NewStorageWorkflow(productsService)
	if err != nil {
		log.Fatalln(err)
	}

	orderWorkflow, err := workflows.NewOrderWorkflow(storageWorkflow)
	if err != nil {
		log.Fatalln(err)
	}

	action := os.Args[1]

	workflowClient, err := helpers.NewWorkflowClient(serviceNameCadenceClient, serviceNameCadenceFrontend)

	if err != nil {
		panic(err)
	}

	triggerClient := helpers.NewCadenceClient(workflowClient)

	switch name := action; name {
	case "RunOrder":
		orderId := "e49ebcbf-8a46-4a0c-801e-dea11c05cd83"
		userId := uuid.New()
		workflowId := orderId + ":" + userId

		_, err = triggerClient.StartWorkflow(context.Background(), client.StartWorkflowOptions{
			ID:                           workflowId,
			TaskList:                     "order-tasklist",
			ExecutionStartToCloseTimeout: 5 * time.Minute,
		}, orderWorkflow.RunOrder, orderId)

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
