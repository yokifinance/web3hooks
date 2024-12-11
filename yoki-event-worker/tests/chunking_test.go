package tests

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	yoki_common "yoki.finance/yoki-event-worker/common"
	"yoki.finance/yoki-event-worker/worker"
)

func chunkingTest(t *testing.T) {
	addressMap := map[common.Address][]*yoki_common.EventListener{}
	addressMap[common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")] = []*yoki_common.EventListener{{}}
	addressMap[common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")] = []*yoki_common.EventListener{{}}
	addressMap[common.HexToAddress("0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC")] = []*yoki_common.EventListener{{}}
	// map – 3 addresses, listener for each

	chunked := worker.ChunkAddresses(addressMap, 3)
	assert.Equal(t, 1, len(chunked), "must one one chunk only")
	assert.Equal(t, 3, len(chunked[0]))

	chunked = worker.ChunkAddresses(addressMap, 300)
	assert.Equal(t, 1, len(chunked), "must one one chunk only")
	assert.Equal(t, 3, len(chunked[0]))
	assert.Equal(t, []common.Address{
		common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
		common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"),
		common.HexToAddress("0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC")},
		chunked[0])

	chunked = worker.ChunkAddresses(addressMap, 2)
	assert.Equal(t, 2, len(chunked), "must be two chunks")
	assert.Equal(t, 2, len(chunked[0]))
	assert.Equal(t, []common.Address{
		common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
		common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")},
		chunked[0])

	assert.Equal(t, 1, len(chunked[1]))
	assert.Equal(t, []common.Address{
		common.HexToAddress("0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC")},
		chunked[1])

	chunked = worker.ChunkAddresses(map[common.Address][]*yoki_common.EventListener{}, 3)
	assert.Equal(t, 0, len(chunked), "no chunks, empty input")
}
