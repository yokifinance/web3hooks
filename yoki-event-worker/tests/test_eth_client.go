package tests

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TestEthClient struct {
	HeaderNumber int64
}

func (c *TestEthClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	if q.FromBlock.Cmp(big.NewInt(50)) == 0 {
		return []types.Log{{
			Address:     common.HexToAddress("0xe467fab1e5ddA1AAf51Ad3d4e10a2667e9efF2c3"),
			BlockNumber: 50,
		}}, nil
	}

	if q.FromBlock.Cmp(big.NewInt(5000)) == 0 {
		return []types.Log{
			{
				Address:     common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"),
				BlockNumber: 5000,
			}}, nil
	}

	return []types.Log{}, nil // always empty logs
}

func (c *TestEthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	c.HeaderNumber++
	return &types.Header{Number: big.NewInt(c.HeaderNumber)}, nil
}
