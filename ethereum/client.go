package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	*ethclient.Client
}

func NewClient(url string) (*Client, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		fmt.Println("Error while configure instance")
	}
	return &Client{
		Client: client,
	}, nil
}