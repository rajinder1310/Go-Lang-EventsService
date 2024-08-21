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
	// "Go-EventService/config"
	// "Go-EventService/services"
	// "context"
	// "fmt"
	// "log"
	// "time"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

type StakingModel struct {
	User              string `bson:"user" json:"user" validate:"required"`
	TransactionHash   string `bson:"transaction_hash" json:"transaction_hash" validate:"required"`
	StakeType         string `bson:"stake_type" json:"stake_type" validate:"required"`
	StakeAmount       string `bson:"stake_amount" json:"stake_amount" validate:"required"`
	StakePeriod       string `bson:"stake_period" json:"stake_period" validate:"required"`
	Time              string `bson:"time" json:"time" validate:"required"`
	TotalStakedInPool string `bson:"total_staked_in_pool" json:"total_staked_in_pool" validate:"required"`
	BlockNumber       string `bson:"block_number" json:"block_number" validate:"required"`
	ContractAddress   string `bson:"contract_address" json:"contract_address" validate:"required"`
	ChainType         string `bson:"chain_type" json:"chain_type" validate:"required"`
	LogIndex          string `bson:"log_index" json:"log_index" validate:"required"`
}

func init() {
	var stakingCollection *mongo.Collection = config.GetCollection("staking")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Create index on email field
	compositeIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "transaction_hash", Value: 1},
			{Key: "user", Value: 1},
			{Key: "contract_address", Value: 1},
			{Key: "chain_type", Value: 1},
			{Key: "log_index", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	if _, err := stakingCollection.Indexes().CreateOne(ctx, compositeIndex); err != nil {
		fmt.Println("UnConfirmed", err)
		log.Fatal(err)
	}

	fmt.Println("Staking Index created successfully!")
}
