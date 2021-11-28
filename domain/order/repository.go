package orders

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/machado-br/order-service/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repository struct {
	DatabaseUri string
}
type Repository interface {
	GetAll(ctx context.Context) ([]entities.Order, error)
	Get(ctx context.Context, orderId string) (entities.Order, error)
	Create(ctx context.Context, order entities.Order) (string, error)
	Delete(ctx context.Context, orderId string) error
	Update(ctx context.Context, order entities.Order) (string, error)
}

func NewRepository(
	databaseUri string,
) (Repository, error) {
	return repository{
		DatabaseUri: databaseUri,
	}, nil
}
func (r repository) buildMongoclient(ctx context.Context) (*mongo.Client, error) {
	log.Println("repository.buildMongoclient")

	if r.DatabaseUri == "" {
		log.Fatalln("You must set your 'MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/#environment-variable")
		return nil, errors.New("MongoDB URI not set in env file")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(r.DatabaseUri))
	if err != nil {
		log.Println("err mongo.Connect")
		log.Fatalln(err)
		return nil, err
	}
	return client, nil
}

func (r repository) GetAll(ctx context.Context) ([]entities.Order, error) {
	log.Println("repository.getAll")

	client, err := r.buildMongoclient(ctx)
	if err != nil {
		log.Fatalln(err)
		return []entities.Order{}, err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalln(err)
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
			log.Fatalln(err)
			log.Fatal(err)
		}
		orders = append(orders, order)
		p, _ := json.MarshalIndent(order, "", "\t")
		fmt.Println(string(p))
	}

	return orders, nil
}

func (r repository) Get(ctx context.Context, orderId string) (entities.Order, error) {
	log.Println("repository.get")

	client, err := r.buildMongoclient(ctx)
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

	log.Printf("order: %v\n", order)

	return order, nil
}

func (r repository) Create(ctx context.Context, order entities.Order) (string, error) {
	log.Println("service.create")

	client, err := r.buildMongoclient(ctx)
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

	log.Printf("result: %v\n", result)

	return order.OrderId, nil
}

func (r repository) Delete(ctx context.Context, orderId string) error {
	log.Println("repository.delete")

	client, err := r.buildMongoclient(ctx)
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

	log.Printf("result: %v\n", result)

	return nil
}

func (r repository) Update(ctx context.Context, order entities.Order) (string, error) {
	log.Println("repository.update")

	client, err := r.buildMongoclient(ctx)
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
	filter := bson.M{"_id": order.Id}

	result, err := orderCollection.ReplaceOne(ctx, filter, doc)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return "", err
	}

	log.Printf("result: %v\n", result)

	return "", nil
}
