package orders

import (
	"context"
	"errors"
	"log"

	"github.com/machado-br/order-service/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repository struct {
	DatabaseUri string
	Collection  *mongo.Collection
	MongoClient *mongo.Client
}
type Repository interface {
	FindAll(ctx context.Context) ([]entities.Order, error)
	Find(ctx context.Context, orderId string) (entities.Order, error)
	Create(ctx context.Context, order entities.Order) (string, error)
	Delete(ctx context.Context, orderId string) (int, error)
	Replace(ctx context.Context, order entities.Order) (int, error)
	UpdateStatus(ctx context.Context, orderId string, status string) (int, error)
}

func NewRepository(
	databaseUri string,
	databaseName string,
	collectionName string,
) (Repository, error) {

	mongoClient, err := buildMongoclient(context.Background(), databaseUri)
	if err != nil {
		return nil, err
	}

	mongoCollection := mongoClient.Database(databaseName).Collection(collectionName)

	return repository{
		DatabaseUri: databaseUri,
		Collection:  mongoCollection,
		MongoClient: mongoClient,
	}, nil
}

func buildMongoclient(ctx context.Context, databaseUri string) (*mongo.Client, error) {
	log.Println("repository.buildMongoclient")

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

func (r repository) FindAll(ctx context.Context) ([]entities.Order, error) {
	log.Println("repository.getAll")

	filter := bson.D{{}}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return []entities.Order{}, err
	}

	var order entities.Order
	var orders []entities.Order
	for cursor.Next(ctx) {
		if err := cursor.Decode(&order); err != nil {
			log.Fatalln(err)
			log.Fatal(err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r repository) Find(ctx context.Context, orderId string) (entities.Order, error) {
	log.Println("repository.get")
	log.Println("orderId: ", orderId)

	filter := bson.D{primitive.E{Key: "orderId", Value: orderId}}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return entities.Order{}, err
	}

	var order entities.Order
	for cursor.Next(ctx) {
		if err := cursor.Decode(&order); err != nil {
			log.Fatal(err)
		}
	}

	return order, nil
}

func (r repository) Create(ctx context.Context, order entities.Order) (string, error) {
	log.Println("repository.create")

	order.Id = primitive.NewObjectID()

	doc, err := bson.Marshal(order)
	if err != nil {
		return "", err
	}

	result, err := r.Collection.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}

	insertedId := result.InsertedID.(primitive.ObjectID).String()
	return insertedId, nil
}

func (r repository) Delete(ctx context.Context, orderId string) (int, error) {
	log.Println("repository.delete")

	filter := bson.D{primitive.E{Key: "orderId", Value: orderId}}

	result, err := r.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(result.DeletedCount), nil
}

func (r repository) Replace(ctx context.Context, order entities.Order) (int, error) {
	log.Println("repository.update")

	doc, err := bson.Marshal(order)
	if err != nil {
		return 0, err
	}
	filter := bson.D{primitive.E{Key: "orderId", Value: order.OrderId}}

	result, err := r.Collection.ReplaceOne(ctx, filter, doc)
	if err != nil {
		log.Printf("err: %v\n", err)
		return 0, err
	}

	return int(result.ModifiedCount), nil
}

func (r repository) UpdateStatus(ctx context.Context, orderId string, status string) (int, error) {
	log.Println("repository.updateStatus")

	filter := bson.D{primitive.E{Key: "orderId", Value: orderId}}
	update := bson.D{
		{"$set", bson.D{{"status", status}}},
	}

	result, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("err: %v\n", err)
		return 0, err
	}

	return int(result.ModifiedCount), nil
}
