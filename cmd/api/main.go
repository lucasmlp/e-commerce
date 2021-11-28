package main

import (
	"log"
	"os"
	"strconv"

	"github.com/machado-br/order-service/api"
	orders "github.com/machado-br/order-service/domain/order"
)

func main() {
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Fatalln(err)
	}

	mongoUri := os.Getenv("MONGODB_URI")
	databaseName := os.Getenv("DATABASE_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")

	ordersRepository, err := orders.NewRepository(mongoUri, databaseName, collectionName)
	if err != nil {
		log.Fatalln(err)
	}

	ordersService := orders.NewService(ordersRepository)
	if err != nil {
		log.Fatalln(err)
	}

	ordersApi, err := api.NewApi(debug, ordersService)
	if err != nil {
		log.Fatalln(err)
	}

	ordersApi.Run()
}
