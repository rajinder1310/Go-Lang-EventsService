package services

import (
	"Go-EventService/config"
	"Go-EventService/contracts"
	"Go-EventService/ethereum"
	"Go-EventService/event"
	"Go-EventService/models"
	"context"
	"fmt"
	"math/big"
	"strconv"

	eth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EventService struct {
	Client       *ethereum.Client
	Contract     *contracts.Contract
	ContractName string
	AbiContract  string
	ChainId      *big.Int
	StartBlock   *big.Int
}

func NewEventService(rpcURl string, contractAddress string, contractName string, chainId *big.Int, abi string, startBlock *big.Int) (*EventService, error) {
	client, err := ethereum.NewClient(rpcURl)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ethereum client: %w", err)
	}
	contract, err := contracts.ContractInstance(contractAddress, abi)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract instance: %w", err)
	}
	return &EventService{
		Client:       client,
		Contract:     contract,
		ContractName: contractName,
		ChainId:      chainId,
		StartBlock:   startBlock,
	}, nil
}

func (s *EventService) decodeEvent(log types.Log) (*event.Event, map[string]interface{}, error) {
	eventABI, err := s.Contract.ABI.EventByID(log.Topics[0])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get event by ID: %w", err)
	}

	data := make(map[string]interface{})
	if err := s.Contract.ABI.UnpackIntoMap(data, eventABI.Name, log.Data); err != nil {
		return nil, nil, fmt.Errorf("failed to unpack event data: %w", err)
	}
	j := 1
	// Iterate over the inputs to handle indexed parameters
	for _, input := range eventABI.Inputs {
		if input.Indexed && j < len(log.Topics) {
			value, err := parseTopicValue(input.Type, log.Topics[j])
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse indexed parameter %s: %w", input.Name, err)
			}
			data[input.Name] = value
			j++
		}
	}

	return &event.Event{
		Name: eventABI.Name,
		Data: data,
		Raw:  log,
	}, data, nil
}

func (s *EventService) FetchAndProcessEvents() error {
	fmt.Println("******************************************************************************************************", s.ContractName)
	chainID, err := s.Client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("error getting chain ID: %w", err)
	}

	latestBlock, err := s.Client.BlockNumber(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching latest block number: %w", err)
	}
	// Convert latestBlock from uint64 to *big.Int
	latestBlockBigInt := new(big.Int).SetUint64(latestBlock)
	blockInfo, err := GetLastBlockInfo(s.StartBlock, latestBlockBigInt, s.Contract.Address.String(), s.ChainId)
	if err != nil {
		return fmt.Errorf("error fetching block info: %w", err)
	}
	fmt.Println("Start Block is ", blockInfo.start_block)
	fmt.Println("Last Block is ", blockInfo.last_block)
	if blockInfo.start_block.Cmp(blockInfo.last_block) >= 0 {
		return nil
	}

	query := eth.FilterQuery{
		FromBlock: blockInfo.start_block,
		ToBlock:   blockInfo.last_block,
		Addresses: []common.Address{s.Contract.Address},
	}

	fmt.Println("Start Block is ***", blockInfo.start_block)
	fmt.Println("Last Block is *** ", blockInfo.last_block)

	logs, err := s.Client.FilterLogs(context.Background(), query)
	if err != nil {
		return fmt.Errorf("error filtering logs: %w", err)
	}
	var events []*event.Event
	for _, logDetail := range logs {
		event, _, err := s.decodeEvent(logDetail)
		if err != nil {
			return fmt.Errorf("error decoding event: %w", err)
		}
		transactionHash := logDetail.TxHash.Hex()
		blockNumber := logDetail.BlockNumber
		event.Data["transactionHash"] = transactionHash
		event.Data["blockNumber"] = blockNumber
		event.Data["chainID"] = chainID
		event.Data["logIndex"] = logDetail.Index
		events = append(events, event)
	}
	if err := s.processEvents(events); err != nil {
		return fmt.Errorf("error processing events: %w", err)
	}
	if err = s.updateBlockInfo(blockInfo.last_block); err != nil {
		return fmt.Errorf("error updating block info: %w", err)
	}
	return nil
}

func convertToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint64: // Handle uint64 type
		return strconv.FormatUint(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10) // Handle uint type
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case *big.Int:
		return v.String()
	case common.Address:
		return v.Hex()
	default:
		fmt.Printf("*******Unhandled type: %T", v)
		return ""
	}
}

func (s *EventService) processEvents(events []*event.Event) error {
	configDetail, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error while loading env in web3service: %w", err)
	}

	var stakingEvents []interface{}
	var tokenEvents []interface{}
	for _, event := range events {
		switch s.Contract.Address.Hex() {
		case configDetail.ContractAddress:
			if event.Name == "Deposit" {
				stakingData := models.StakingModel{}
				// Convert the string to a common.Address
				user := convertToString(event.Data["user"])
				stakingData.User = user

				transactionHash := convertToString(event.Data["transactionHash"])
				stakingData.TransactionHash = transactionHash

				stakeType := convertToString(event.Data["stakeType"])
				stakingData.StakeType = stakeType

				stakeAmount := convertToString(event.Data["stakeAmount"])
				stakingData.StakeAmount = stakeAmount

				stakePeriod := convertToString(event.Data["period"])
				stakingData.StakePeriod = stakePeriod

				time := convertToString(event.Data["blockTimestamp"])
				stakingData.Time = time

				totalStakedInPool := convertToString(event.Data["totalStakedInPool"])
				stakingData.TotalStakedInPool = totalStakedInPool

				blockNumber := convertToString(event.Data["blockNumber"])
				stakingData.BlockNumber = blockNumber

				contractAddress := convertToString(configDetail.ContractAddress)
				stakingData.ContractAddress = contractAddress

				chainId := convertToString(event.Data["chainID"])
				stakingData.ChainType = chainId

				logIndex := convertToString(event.Data["logIndex"])
				stakingData.LogIndex = logIndex

				stakingEvents = append(stakingEvents, stakingData)
			}

		case configDetail.TokenContract:
			var transferModel models.TransferModel
			if event.Name == "Transfer" {
				transactionHash := convertToString(event.Data["transactionHash"])
				transferModel.TransactionHash = transactionHash

				contractAddress := convertToString(configDetail.ContractAddress)
				transferModel.ContractAddress = contractAddress

				chainId := convertToString(event.Data["chainID"])
				transferModel.ChainType = chainId

				blockNumber := convertToString(event.Data["blockNumber"])
				transferModel.BlockNumber = blockNumber

				from := convertToString(event.Data["from"])
				transferModel.FromAddress = from

				to := convertToString(event.Data["to"])
				transferModel.ToAddress = to

				amount := convertToString(event.Data["value"])
				transferModel.Amount = amount

				logIndex := convertToString(event.Data["logIndex"])
				transferModel.LogIndex = logIndex

				tokenEvents = append(tokenEvents, transferModel)

			}

		default:
			fmt.Printf("Unhandled contract address: %s\n", s.Contract.Address)
		}
	}

	if len(stakingEvents) > 0 {
		if err := s.insertBulkData(stakingEvents, StakingCollection); err != nil {
			return fmt.Errorf("failed to insert staking data: %w", err)
		}
	}
	if len(tokenEvents) > 0 {
		if err := s.insertBulkData(tokenEvents, TransferCollection); err != nil {
			return fmt.Errorf("failed to insert Transfer data: %w", err)
		}
	}
	return nil
}

func parseTopicValue(t abi.Type, topic common.Hash) (interface{}, error) {
	switch t.T {
	case abi.IntTy, abi.UintTy:
		return topic.Big(), nil
	case abi.BoolTy:
		return topic.Big().Cmp(common.Big0) != 0, nil
	case abi.AddressTy:
		return common.HexToAddress(topic.Hex()), nil
	case abi.HashTy:
		return topic, nil
	default:
		return nil, fmt.Errorf("unsupported indexed parameter type: %v", t.String())
	}
}
