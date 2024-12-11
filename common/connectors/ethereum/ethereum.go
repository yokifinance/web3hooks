package ethereum

import (
	"context"
	"log"
	"math/big"

	"github.com/pkg/errors"

	"yoki.finance/common/config"
	"yoki.finance/common/connectors/ethereum/token"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	clients = make(map[int]*ethclient.Client)
)

func GetBalance(ctx context.Context, chainId int, accountAddr string, blockNum int64) (balance *big.Int, err error) {
	rpcUrl, ok := config.RpcEndpoints[chainId]
	if !ok {
		return nil, errors.New("chainId is unsupported")
	}

	client, ok := clients[chainId]
	if !ok {
		if client, err = ethclient.Dial(rpcUrl); err != nil {
			return nil, errors.Wrap(err, "ethclient dial")
		}
		clients[chainId] = client
	}

	account := common.HexToAddress(accountAddr)
	var blockNumber *big.Int
	if blockNum > 0 {
		blockNumber.SetInt64(blockNum)
	}
	balance, err = client.BalanceAt(ctx, account, blockNumber)
	if err != nil {
		return nil, err
	}
	log.Println("chainId:", chainId, accountAddr, "balance:", balance)

	return balance, nil
}

func GetBalanceERC20(ctx context.Context, chainId int, accountAddr, tokenAddr string) (balance *big.Int, err error) {
	instance, err := getToken(ctx, chainId, tokenAddr)
	if err != nil {
		return nil, err
	}

	bal, err := instance.BalanceOf(&bind.CallOpts{}, common.HexToAddress(accountAddr))
	if err != nil {
		return nil, err
	}
	return bal, nil
}

func GetAllowanceERC20(ctx context.Context, chainId int, tokenAddr, ownerAddr, spenderAddr string) (*big.Int, error) {
	instance, err := getToken(ctx, chainId, tokenAddr)
	if err != nil {
		return nil, err
	}

	allowance, err := instance.Allowance(&bind.CallOpts{}, common.HexToAddress(ownerAddr), common.HexToAddress(spenderAddr))
	if err != nil {
		return nil, err
	}

	return allowance, nil
}

func getToken(ctx context.Context, chainId int, tokenAddr string) (instance *token.Token, err error) {
	rpcUrl, ok := config.RpcEndpoints[chainId]
	if !ok {
		return nil, errors.New("chainId is unsupported")
	}

	client, ok := clients[chainId]
	if !ok {
		if client, err = ethclient.Dial(rpcUrl); err != nil {
			return nil, errors.Wrap(err, "ethclient dial")
		}
		clients[chainId] = client
	}

	instance, err = token.NewToken(common.HexToAddress(tokenAddr), client)
	if err != nil {
		return nil, err
	}
	return instance, nil
}
