package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	entities "github.com/machado-br/order-service/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURI  = "mongodb://localhost:27017"
	productId = "49f2cea1-7648-4e88-947f-0b8db5cb845a"
)

func GetOrder(ctx context.Context, orderId string) (entities.Order, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
		return entities.Order{}, err
	}

	db := client.Database("order-service")
	productsCollection := db.Collection("orders")

	cursor, err := productsCollection.Find(
		ctx,
		bson.D{{"id", orderId}},
	)

	var order entities.Order
	for cursor.Next(ctx) {
		if err := cursor.Decode(&order); err != nil {
			log.Fatal(err)
		}
		p, _ := json.MarshalIndent(order, "", "\t")
		fmt.Println(string(p))
	}

	client.Disconnect(ctx)
	if err != nil {
		return entities.Order{}, err
	}

	return order, nil
}

func GetProduct(ctx context.Context) (entities.Product, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
		return entities.Product{}, err
	}

	db := client.Database("storage-service")

	productsCollection := db.Collection("products")

	cursor, err := productsCollection.Find(
		ctx,
		bson.D{{"id", productId}},
	)
	var product entities.Product
	for cursor.Next(ctx) {
		if err := cursor.Decode(&product); err != nil {
			log.Fatal(err)
		}
		p, _ := json.MarshalIndent(product, "", "\t")
		fmt.Println(string(p))
	}

	client.Disconnect(ctx)
	if err != nil {
		return entities.Product{}, err
	}

	return product, nil
}

func PingMongo(ctx context.Context) (int, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	fmt.Println("Connected to MongoDB!")

	client.Disconnect(ctx)
	if err != nil {
		return -1, err
	}

	fmt.Println("Disconnected from MongoDB!")
	return 1, nil
}
