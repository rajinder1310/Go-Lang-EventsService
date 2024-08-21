package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	RPCURL          string `mapstructure:"RPC_URL"`
	MongoURI        string `mapstructure:"MONGO_URI"`
	DATABASENAME    string `mapstructure:"DATABASE_NAME"`
	ContractAddress string `mapstructure:"STAKING_CONTRACT_ADDRESS"`
	TokenContract   string `mapstructure:"TOKEN_CONTRACT_ADDRESS"`
	CRONDURATION    string `mapstructure:"CRON_DURATION"`
	BLOCKRANGE      int    `mapstructure:"BLOCK_RANGE"`
	STARTBLOCK      int    `mapstructure:"START_BLOCK_NUMBER"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error reading config file", err)
		return nil, err
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Println("Error unable to decode into struct", err)
		return nil, err
	}
	return &config, nil
}
