package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/machado-br/e-commerce/cadence/activities"
	"github.com/machado-br/e-commerce/cadence/helpers"
	"github.com/machado-br/e-commerce/cadence/workflows"
	"github.com/machado-br/e-commerce/domain/orders"
	"github.com/machado-br/e-commerce/domain/products"
	"github.com/pborman/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.uber.org/cadence/client"
)

func buildMongoclient(ctx context.Context, databaseUri string) (*mongo.Client, error) {
	log.Println("cadence.main.buildMongoclient")

	if databaseUri == "" {
		log.Fatalln("You must set your 'MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/#environment-variable")
		return nil, errors.New("MongoDB URI not set in env file")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(databaseUri))
	if err != nil {
		log.Println("err mongo.Connect")
		log.Fatalln(err)
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return client, nil
}

func main() {
	serviceNameCadenceClient := os.Getenv("CADENCE_CLIENT_NAME")
	serviceNameCadenceFrontend := os.Getenv("CADENCE_FRONTEND_NAME")

	mongoUri := os.Getenv("MONGODB_URI")
	ordersDatabaseName := os.Getenv("ORDERS_DATABASE_NAME")
	ordersCollectionName := os.Getenv("ORDERS_COLLECTION_NAME")
	productsDatabaseName := os.Getenv("PRODUCTS_DATABASE_NAME")
	productsCollectionName := os.Getenv("PRODUCTS_COLLECTION_NAME")

	mongoClient, err := buildMongoclient(context.Background(), mongoUri)
	if err != nil {
		log.Fatalln(err)
	}

	ordersCollection := mongoClient.Database(ordersDatabaseName).Collection(ordersCollectionName)

	ordersRepository, err := orders.NewRepository(mongoClient, ordersCollection)
	if err != nil {
		log.Fatalln(err)
	}

	ordersService := orders.NewService(ordersRepository)
	if err != nil {
		log.Fatalln(err)
	}

	productsRepository, err := products.NewRepository(mongoUri, productsDatabaseName, productsCollectionName)
	if err != nil {
		log.Fatalln(err)
	}

	productsService := products.NewService(productsRepository)
	if err != nil {
		log.Fatalln(err)
	}

	activities, err := activities.NewActivities(ordersService, productsService)
	if err != nil {
		log.Fatalln(err)
	}

	storageWorkflow, err := workflows.NewStorageWorkflow(productsService, activities)
	if err != nil {
		log.Fatalln(err)
	}

	orderWorkflow, err := workflows.NewOrderWorkflow(storageWorkflow, activities)
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
		orderId := "8a786ba1-a9c9-45bb-a4c8-0d6f3a4b08eb"
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
