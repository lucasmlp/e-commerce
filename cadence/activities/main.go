package activities

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/machado-br/order-service/domain/dtos"
	"github.com/machado-br/order-service/domain/entities"
	"github.com/machado-br/order-service/domain/orders"
	"github.com/machado-br/order-service/domain/products"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Activities struct {
	OrderService    orders.Service
	ProductsService products.Service
}

func NewActivities(orderService orders.Service, productService products.Service) (Activities, error) {
	activities := Activities{
		OrderService:    orderService,
		ProductsService: productService,
	}
	return activities, nil
}

func (a Activities) GetOrder(ctx context.Context, orderId string) (dtos.Order, error) {
	order, err := a.OrderService.Find(ctx, orderId)
	if err != nil {
		return dtos.Order{}, err
	}

	return order, nil
}

func GetProduct(ctx context.Context, productId string) (entities.Product, error) {
	client, err := buildMongoclient(ctx)
	if err != nil {
		log.Fatal(err)
		return entities.Product{}, err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	db := client.Database("storage-service")

	productsCollection := db.Collection("products")

	cursor, err := productsCollection.Find(
		ctx,
		bson.D{{"id", productId}},
	)
	if err != nil {
		return entities.Product{}, err
	}

	var product entities.Product
	for cursor.Next(ctx) {
		if err := cursor.Decode(&product); err != nil {
			log.Fatal(err)
		}
		p, _ := json.MarshalIndent(product, "", "\t")
		fmt.Println(string(p))
	}

	return product, nil
}

func UpdateProductUnits(ctx context.Context, productId string, units int) error {

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

	coll := client.Database("storage-service").Collection("products")
	filter := bson.D{{"productId", productId}}
	update := bson.D{{"$set", bson.D{{"units", units}}}}
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		println(err)
		return err
	}

	if result.ModifiedCount <= 0 {
		println(result.ModifiedCount)
		return errors.New("failed to update product units")
	}

	return nil
}

func PingMongo(ctx context.Context) (int, error) {
	client, err := buildMongoclient(ctx)
	if err != nil {
		log.Fatal(err)
		return -1, err
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	return 1, nil
}

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
