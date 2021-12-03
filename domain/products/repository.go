package products

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

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
	GetAll(ctx context.Context) ([]entities.Product, error)
	Get(ctx context.Context, productId string) (entities.Product, error)
	Create(ctx context.Context, product entities.Product) (string, error)
	Delete(ctx context.Context, productId string) error
	Replace(ctx context.Context, product entities.Product) (string, error)
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
	return client, nil
}

func (r repository) GetAll(ctx context.Context) ([]entities.Product, error) {
	log.Println("repository.getAll")

	filter := bson.D{{}}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return []entities.Product{}, err
	}

	var product entities.Product
	var products []entities.Product
	for cursor.Next(ctx) {
		if err := cursor.Decode(&product); err != nil {
			log.Fatalln(err)
			log.Fatal(err)
		}
		products = append(products, product)
	}

	return products, nil
}

func (r repository) Get(ctx context.Context, productId string) (entities.Product, error) {
	log.Println("repository.get")

	filter := bson.D{primitive.E{Key: "productId", Value: productId}}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return entities.Product{}, err
	}

	var product entities.Product
	for cursor.Next(ctx) {
		if err := cursor.Decode(&product); err != nil {
			log.Fatal(err)
		}
	}

	return product, nil
}

func (r repository) Create(ctx context.Context, product entities.Product) (string, error) {
	log.Println("repository.create")

	product.Id = primitive.NewObjectID()

	doc, err := bson.Marshal(product)
	if err != nil {
		return "", err
	}

	_, err = r.Collection.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}

	return product.ProductId, nil
}

func (r repository) Delete(ctx context.Context, productId string) error {
	log.Println("repository.delete")

	filter := bson.D{primitive.E{Key: "productId", Value: productId}}

	_, err := r.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (r repository) Replace(ctx context.Context, product entities.Product) (string, error) {
	log.Println("repository.update")

	doc, err := bson.Marshal(product)
	if err != nil {
		return "", err
	}

	fmt.Printf("product: %v\n", product)

	filter := bson.D{primitive.E{Key: "productId", Value: product.ProductId}}

	result, err := r.Collection.ReplaceOne(ctx, filter, doc)
	if err != nil {
		log.Printf("err: %v\n", err)
		return "", err
	}

	return strconv.FormatInt(result.ModifiedCount, 32), nil
}
