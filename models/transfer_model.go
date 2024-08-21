package models

import (
	"Go-EventService/config"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransferModel struct {
	BlockNumber     string `bson:"block_number" json:"block_number" validate:"required"`
	ContractAddress string `bson:"contract_address" json:"contract_address" validate:"required"`
	ChainType       string `bson:"chain_type" json:"chain_type" validate:"required"`
	FromAddress     string `bson:"from_address" json:"from_address" validate:"required"`
	ToAddress       string `bson:"to_address" json:"to_address" validate:"required"`
	Amount          string `bson:"amount" json:"amount" validate:"required"`
	TransactionHash string `bson:"transaction_hash" json:"transaction_hash" validate:"required"`
	LogIndex        string `bson:"log_index" json:"log_index" validate:"required"`
}

func init() {
	var TransferCollection *mongo.Collection = config.GetCollection("TransferData")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Create index on email field
	compositeIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "transaction_hash", Value: 1},
			{Key: "contract_address", Value: 1},
			{Key: "chain_type", Value: 1},
			{Key: "log_index", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	if _, err := TransferCollection.Indexes().CreateOne(ctx, compositeIndex); err != nil {
		fmt.Println("UnConfirmed", err)
		log.Fatal(err)
	}
	fmt.Println("Transfer Index created successfully!")
}
