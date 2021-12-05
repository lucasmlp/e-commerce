package main

import (
	"log"
	"os"
	"strconv"

	"github.com/machado-br/e-commerce/api"
	"github.com/machado-br/e-commerce/domain/orders"
	"github.com/machado-br/e-commerce/domain/products"
)

func main() {
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Fatalln(err)
	}

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

	api, err := api.NewApi(debug, ordersService, productsService)
	if err != nil {
		log.Fatalln(err)
	}

	api.Run()
}
