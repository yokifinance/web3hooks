package webhook

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type EventListenerResultDto struct {
	EventListenerId uuid.UUID      `json:"eventListenerId"`
	Chain           int            `json:"chain"`
	Address         common.Address `json:"address"`
	BlockHash       common.Hash    `json:"blockHash"`
	BlockNumber     uint64         `json:"blockNumber"`
	LogIndex        uint           `json:"logIndex"`
	TxHash          common.Hash    `json:"txHash"`
	TxIndex         uint           `json:"txIndex"`
	Data            string         `json:"data"`
	Topics          []common.Hash  `json:"topics"`
}
