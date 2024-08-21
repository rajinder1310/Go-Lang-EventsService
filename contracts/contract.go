package contracts

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type Contract struct {
	Address common.Address
	ABI abi.ABI
}

func ContractInstance(address string, abiJSON string) (*Contract, error) {
		contractABI, err := abi.JSON(strings.NewReader(abiJSON))
		if err != nil {
			return nil, err
		}

		return &Contract{
			Address: common.HexToAddress(address),
			ABI: contractABI,
		}, nil
}