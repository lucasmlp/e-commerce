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
	Delete(ctx context.Context, orderId string) error
	Replace(ctx context.Context, order entities.Order) (string, error)
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

	_, err = r.Collection.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}

	return order.OrderId, nil
}

func (r repository) Delete(ctx context.Context, orderId string) error {
	log.Println("repository.delete")

	filter := bson.D{primitive.E{Key: "orderId", Value: orderId}}

	_, err := r.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (r repository) Replace(ctx context.Context, order entities.Order) (string, error) {
	log.Println("repository.update")

	doc, err := bson.Marshal(order)
	if err != nil {
		return "", err
	}
	filter := bson.D{primitive.E{Key: "orderId", Value: order.OrderId}}

	_, err = r.Collection.ReplaceOne(ctx, filter, doc)
	if err != nil {
		log.Printf("err: %v\n", err)
		return "", err
	}

	return "", nil
}
