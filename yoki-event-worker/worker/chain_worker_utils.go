// In this file we put util, small functions for ChainWorker

package worker

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"yoki.finance/common/rcommon"
	yoki_common "yoki.finance/yoki-event-worker/common"
)

// Creates mapping of address and all its listeners
func (w *ChainWorker) mapListenersOnAddress(allEventListeners []yoki_common.EventListener) (res addressToListenersMapType) {
	res = addressToListenersMapType{}
	for _, listener := range allEventListeners {
		addr := common.HexToAddress(listener.Address)
		if res[addr] == nil {
			res[addr] = []*yoki_common.EventListener{}
		}
		res[addr] = append(res[addr], &listener)
	}
	return res
}

// Returns chunks of addresses not more than chunkSize in length
func ChunkAddresses(addressListenersMap addressToListenersMapType, chunkSize int) [][]common.Address {
	res := [][]common.Address{}
	if len(addressListenersMap) == 0 {
		return res
	}

	currentAddressChunk := []common.Address{}
	for address := range addressListenersMap {
		currentAddressChunk = append(currentAddressChunk, address)
		if len(currentAddressChunk) == chunkSize {
			res = append(res, currentAddressChunk)
			currentAddressChunk = []common.Address{}
		}
	}
	if len(currentAddressChunk) > 0 {
		res = append(res, currentAddressChunk)
	}

	return res
}

func (w *ChainWorker) log(msg string, a ...any) {
	rcommon.Println(fmt.Sprintf("chain %d: ", w.chain)+msg, a...)
}

func min(x, y *big.Int) *big.Int {
	if x.Cmp(y) < 0 {
		return x
	} else {
		return y
	}
}
