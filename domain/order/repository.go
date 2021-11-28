package orders

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/machado-br/order-service/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func buildMongoclient(ctx context.Context) (*mongo.Client, error) {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/#environment-variable")
		return nil, errors.New("MongoDB URI not set in env file")
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetAll(ctx context.Context) ([]entities.Order, error) {
	fmt.Println("repository.getAll")

	client, err := buildMongoclient(ctx)
	if err != nil {
		log.Fatal(err)
		return []entities.Order{}, err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	db := client.Database("order-service")
	orderCollection := db.Collection("orders")

	filter := bson.D{{}}

	cursor, err := orderCollection.Find(ctx, filter)

	var order entities.Order
	var orders []entities.Order
	for cursor.Next(ctx) {
		if err := cursor.Decode(&order); err != nil {
			log.Fatal(err)
		}
		orders = append(orders, order)
		p, _ := json.MarshalIndent(order, "", "\t")
		fmt.Println(string(p))
	}

	return orders, nil
}

func Get(ctx context.Context, orderId string) (entities.Order, error) {
	fmt.Println("repository.get")
	client, err := buildMongoclient(ctx)
	if err != nil {
		log.Fatal(err)
		return entities.Order{}, err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	db := client.Database("order-service")
	productsCollection := db.Collection("orders")

	cursor, err := productsCollection.Find(
		ctx,
		bson.D{{"orderId", orderId}},
	)

	var order entities.Order
	for cursor.Next(ctx) {
		if err := cursor.Decode(&order); err != nil {
			log.Fatal(err)
		}
		p, _ := json.MarshalIndent(order, "", "\t")
		fmt.Println(string(p))
	}

	fmt.Printf("order: %v\n", order)

	return order, nil
}

func Create(ctx context.Context, order entities.Order) (string, error) {
	fmt.Println("service.create")

	client, err := buildMongoclient(ctx)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	db := client.Database("order-service")
	orderCollection := db.Collection("orders")

	order.Id = primitive.NewObjectID()

	doc, err := bson.Marshal(order)
	if err != nil {
		return "", err
	}

	result, err := orderCollection.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}

	fmt.Printf("result: %v\n", result)

	return order.OrderId, nil
}

func Delete(ctx context.Context, orderId string) error {
	fmt.Println("repository.delete")

	client, err := buildMongoclient(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	db := client.Database("order-service")
	orderCollection := db.Collection("orders")

	filter := bson.D{{"id", orderId}}
	result, err := orderCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	fmt.Printf("result: %v\n", result)

	return nil
}

func Update(ctx context.Context, order entities.Order) (string, error) {
	fmt.Println("repository.update")

	client, err := buildMongoclient(ctx)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	db := client.Database("order-service")
	orderCollection := db.Collection("orders")

	doc, err := bson.Marshal(order)
	if err != nil {
		return "", err
	}
	filter := bson.M{"$_id": order.Id}

	result, err := orderCollection.UpdateOne(ctx, doc, filter)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return "", err
	}

	fmt.Printf("result: %v\n", result)

	return "", nil
}
