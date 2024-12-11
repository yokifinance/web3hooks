package tests

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"yoki.finance/common/connectors/db"
	"yoki.finance/yoki-event-worker/data"
)

func lastProcessedBlockTest(t *testing.T) {
	res, err := data.GetLastProcessedBlockNumber(137)
	assert.NoError(t, err)
	assert.Equal(t, new(big.Int), res)

	tx, _ := db.ORM.Begin()
	defer tx.Rollback()
	err = data.SetLastProcessedBlockNumberTx(&tx, 137, new(big.Int).SetInt64(1985))
	assert.NoError(t, err)
	tx.Commit()

	res, _ = data.GetLastProcessedBlockNumber(137)
	assert.Equal(t, big.NewInt(1985), res)
}
