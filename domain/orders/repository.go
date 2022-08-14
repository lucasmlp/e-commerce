package orders

import (
	"context"
	"log"

	"github.com/machado-br/e-commerce/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository struct {
	MongoCollection *mongo.Collection
	MongoClient     *mongo.Client
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
	mongoClient *mongo.Client,
	mongoCollection *mongo.Collection,
) (Repository, error) {

	return repository{
		MongoCollection: mongoCollection,
		MongoClient:     mongoClient,
	}, nil
}

func (r repository) FindAll(ctx context.Context) ([]entities.Order, error) {
	log.Println("ordersRepository.getAll")

	filter := bson.D{{}}

	cursor, err := r.MongoCollection.Find(ctx, filter)
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
	log.Println("ordersRepository.get")
	log.Println("orderId: ", orderId)

	filter := bson.D{primitive.E{Key: "orderId", Value: orderId}}

	cursor, err := r.MongoCollection.Find(ctx, filter)
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
	log.Println("ordersRepository.create")

	order.Id = primitive.NewObjectID()

	doc, err := bson.Marshal(order)
	if err != nil {
		return "", err
	}

	result, err := r.MongoCollection.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}

	insertedId := result.InsertedID.(primitive.ObjectID).String()
	return insertedId, nil
}

func (r repository) Delete(ctx context.Context, orderId string) (int, error) {
	log.Println("ordersRepository.delete")

	filter := bson.D{primitive.E{Key: "orderId", Value: orderId}}

	result, err := r.MongoCollection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(result.DeletedCount), nil
}

func (r repository) Replace(ctx context.Context, order entities.Order) (int, error) {
	log.Println("ordersRepository.update")

	doc, err := bson.Marshal(order)
	if err != nil {
		return 0, err
	}
	filter := bson.D{primitive.E{Key: "orderId", Value: order.OrderId}}

	result, err := r.MongoCollection.ReplaceOne(ctx, filter, doc)
	if err != nil {
		log.Printf("err: %v\n", err)
		return 0, err
	}

	return int(result.ModifiedCount), nil
}

func (r repository) UpdateStatus(ctx context.Context, orderId string, status string) (int, error) {
	log.Println("ordersRepository.updateStatus")

	filter := bson.D{primitive.E{Key: "orderId", Value: orderId}}
	update := bson.D{
		{"$set", bson.D{{"status", status}}},
	}

	result, err := r.MongoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("err: %v\n", err)
		return 0, err
	}

	return int(result.ModifiedCount), nil
}
