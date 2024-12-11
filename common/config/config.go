package config

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"yoki.finance/common/rcommon"
)

// these parameters are set by env variables, otherwise default values (below) are used
var (
	DBhost = "localhost"
	DBport = "40001"
	DBuser = "postgres"
	DBpass = "postgrespw"
	DBname = "yoki"
)

var (
	AppChainIds         []int
	RpcEndpoints        = map[int]string{}
)

func init() {
	rcommon.LoadSecrets("/run/secrets/yoki_web3tasks_env")
	rcommon.LoadDefaultEnvFiles(IsInTests())

	// Initialize RPC API Urls
	anyRpcChain := false
	fmt.Println("RPC API for chains initializing...")
	for name, value := range rcommon.GetAllParams() {
		rpcUrlPrefix := "YOKI_RPC_URL_"
		if strings.HasPrefix(name, rpcUrlPrefix) {
			chainId, err := strconv.Atoi(strings.TrimPrefix(name, rpcUrlPrefix))
			if err != nil {
				log.Fatalf("%s key is not in correct format", name)
			}
			AppChainIds = append(AppChainIds, chainId)

			fmt.Printf("chain %d:\n", chainId)
			RpcEndpoints[chainId] = value
			slashIndex := strings.LastIndex(value, "/")
			if slashIndex != -1 {
				if len(value) > slashIndex+3 {
					obscuredUrl := value[:slashIndex+4]
					fmt.Printf("	RPC URL prefix: %s*\n", obscuredUrl)
				} else {
					panic("Not enough characters in RPC URL after '/'")
				}
			} else {
				panic("Incorrect RPC API URL")
			}

		
			anyRpcChain = true
		}
	}
	if !anyRpcChain {
		log.Fatalln("No RPC API config")
	}

	rcommon.SetParamStrOrLeaveDefault(&DBhost, "POSTGRES_HOST")
	rcommon.SetParamStrOrLeaveDefault(&DBport, "POSTGRES_PORT")
	rcommon.SetParamStrOrLeaveDefault(&DBuser, "POSTGRES_USER")
	rcommon.SetParamStrOrLeaveDefault(&DBpass, "POSTGRES_PASSWORD")
	rcommon.SetParamStrOrLeaveDefault(&DBname, "POSTGRES_DB")
}
