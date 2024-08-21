package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	configDetails, _ := LoadConfig()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(configDetails.MongoURI))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected To MongoDB Database")
	return client
}

var DB *mongo.Client = ConnectDB()

func GetCollection(CollectionName string) *mongo.Collection {
	configDetails, _ := LoadConfig()
	collection := DB.Database(configDetails.DATABASENAME).Collection(CollectionName)
	return collection
}
