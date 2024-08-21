package models

type BlocksModel struct {
	ChainType       string `bson:"chain_type" json:"chain_type" validate:"required"`
	ContractAddress string `bson:"contract_address" json:"contract_address" validate:"required"`
	ContractName    string `bson:"contract_name" json:"contract_name" validate:"required"`
	LastBlockNumber int    `bson:"last_block_number" json:"last_block_number" validate:"required"`
	// CronInProgress  bool   `bson:"cron_in_progress" json:"cron_in_progress" validate:"required"`
	Status bool `bson:"status" json:"status" validate:"required"`
}
