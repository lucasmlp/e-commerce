package main

import (
	"log"
	"os"

	"github.com/machado-br/e-commerce/cadence/activities"
	"github.com/machado-br/e-commerce/cadence/helpers"
	"github.com/machado-br/e-commerce/cadence/workflows"
	"github.com/machado-br/e-commerce/domain/orders"
	"github.com/machado-br/e-commerce/domain/products"

	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/worker"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func main() {

	mongoUri := os.Getenv("MONGODB_URI")
	ordersDatabaseName := os.Getenv("ORDERS_DATABASE_NAME")
	ordersCollectionName := os.Getenv("ORDERS_COLLECTION_NAME")
	productsDatabaseName := os.Getenv("PRODUCTS_DATABASE_NAME")
	productsCollectionName := os.Getenv("PRODUCTS_COLLECTION_NAME")

	ordersRepository, err := orders.NewRepository(mongoUri, ordersDatabaseName, ordersCollectionName)
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

	serviceNameCadenceClient := os.Getenv("CADENCE_CLIENT_NAME")
	serviceNameCadenceFrontend := os.Getenv("CADENCE_FRONTEND_NAME")
	domainName := os.Getenv("CADENCE_DOMAIN_NAME")

	workflowClient, err := helpers.NewWorkflowClient(serviceNameCadenceClient, serviceNameCadenceFrontend)
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	activities, err := activities.NewActivities(ordersService, productsService)
	if err != nil {
		log.Fatalln(err)
	}

	storageWorkflow, err := workflows.NewStorageWorkflow(productsService, activities)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = workflows.NewOrderWorkflow(storageWorkflow, activities)
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
	activity.Register(activities.GetProduct)
	activity.Register(activities.UpdateProduct)
	activity.Register(activities.UpdateOrderStatus)

	w := worker.New(workflowClient, domainName, "order-tasklist",
		worker.Options{
			Logger: logger,
		})

	err = w.Run()

	if err != nil {
		panic(err)
	}
}
