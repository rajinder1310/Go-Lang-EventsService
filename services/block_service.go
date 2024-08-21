package services

import (
	"Go-EventService/config"
	"Go-EventService/models"
	"context"
	"log"
	"math/big"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GetBlockService struct {
	start_block *big.Int
	last_block  *big.Int
}

var (
	BlockCollection    *mongo.Collection = config.GetCollection("blocks")
	StakingCollection  *mongo.Collection = config.GetCollection("staking")
	TransferCollection *mongo.Collection = config.GetCollection("TransferData")
)

func GetLastBlockInfo(startBlockNumber *big.Int, _latestBlock *big.Int, contractAddress string, chainID *big.Int) (*GetBlockService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var Blocks models.BlocksModel
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return nil, err
	}

	err = BlockCollection.FindOne(ctx, bson.D{
		{Key: "contract_address", Value: contractAddress},
		{Key: "chain_id", Value: chainID},
	}).Decode(&Blocks)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			startBlock := startBlockNumber
			lastBlock := new(big.Int).Add(startBlockNumber, big.NewInt(int64(cfg.BLOCKRANGE)))
			if lastBlock.Cmp(_latestBlock) > 0 {
				lastBlock = _latestBlock
			}
			return &GetBlockService{
				start_block: startBlock,
				last_block:  lastBlock,
			}, nil
		}
		log.Printf("Error finding document: %v", err)
		return nil, err
	}

	startBlock := big.NewInt(int64(Blocks.LastBlockNumber + 1))
	lastBlock := new(big.Int).Add(startBlock, big.NewInt(int64(cfg.BLOCKRANGE)))

	if lastBlock.Cmp(_latestBlock) > 0 {
		lastBlock = _latestBlock
	}

	return &GetBlockService{
		start_block: startBlock,
		last_block:  lastBlock,
	}, nil
}

func (e *EventService) updateBlockInfo(lastBlockNumber *big.Int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{
		{Key: "contract_address", Value: e.Contract.Address.String()},
		{Key: "chain_id", Value: e.ChainId},
		{Key: "contract_name", Value: e.ContractName},
	}
	opts := options.Update().SetUpsert(true)
	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "last_block_number", Value: lastBlockNumber.Int64()}}},
	}

	_, err := BlockCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("Error while updating block number in db: %v", err)
		return err
	}
	return nil
}
