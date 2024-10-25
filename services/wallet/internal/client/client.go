package client

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

func NewETHClient(key string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(fmt.Sprintf("https://sepolia.infura.io/v3/%s", key))
	if err != nil {
		return nil, err
	}
	return client, nil
}
