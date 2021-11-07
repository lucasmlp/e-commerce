package activities

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoURI = "mongodb://localhost:27017"
)

func GetProductUnits(ctx context.Context) (int, error) {
	// Open Connection
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatal(err)
	}

	// End Open Connection Code

	// Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	fmt.Println("Connected to MongoDB!")
	return 1, nil
}

func PrintCurrentTime(ctx context.Context) error {
	loc, err := time.LoadLocation("Asia/Tokyo")

	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	tokyoTime := time.Now().In(loc)

	fmt.Printf("tokyoTime: %v\n", tokyoTime)

	loc, err = time.LoadLocation("America/Sao_Paulo")

	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}

	saoPauloTime := time.Now().In(loc)
	fmt.Printf("saoPauloTime: %v\n", saoPauloTime)

	return nil
}

func ActivityA(data string) (string, error) {
	return data + " antigo", nil
}

func ActivityB(data string) (string, error) {
	return data + " + final", nil
}

func ActivityC(data string) (string, error) {
	return data + " novo", nil
}
