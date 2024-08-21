package event

import "github.com/ethereum/go-ethereum/core/types"

type Event struct{
	Name string
	Data map[string]interface{}
	Raw types.Log
}