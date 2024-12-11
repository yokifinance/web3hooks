package worker

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ProdEthClient struct {
	client *ethclient.Client
}

func NewEthClient(rawurl string) (EthClient, error) {
	client, err := ethclient.Dial(rawurl)
	if err != nil {
		return nil, err
	}

	return &ProdEthClient{client: client}, nil
}

func (c *ProdEthClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	return c.client.FilterLogs(ctx, q)
}

func (c *ProdEthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return c.client.HeaderByNumber(ctx, number)
}
