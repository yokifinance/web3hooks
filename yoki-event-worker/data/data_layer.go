package data

import (
	"context"
	"math/big"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bunbig"
	"yoki.finance/common/connectors/db"
	"yoki.finance/yoki-event-worker/common"
)

func GetLastProcessedBlockNumber(chain int) (block *big.Int, err error) {
	var res []common.LastProcessedBlock

	if err = db.ORM.NewSelect().
		Model((*common.LastProcessedBlock)(nil)).
		Where("chain=?", chain).
		Scan(context.Background(), &res); err != nil {
		return nil, err
	}

	if len(res) == 0 {
		block = new(big.Int)
	} else {
		block = res[0].LastProcessedBlockNumber.ToMathBig()
	}

	return block, nil
}

func SetLastProcessedBlockNumberTx(tx *bun.Tx, chain int, blockNumber *big.Int) error {
	blockObj := common.LastProcessedBlock{
		Chain:                    chain,
		LastProcessedBlockNumber: *bunbig.FromMathBig(blockNumber),
	}

	var q *bun.InsertQuery
	if tx != nil {
		q = tx.NewInsert()
	} else {
		q = db.ORM.NewInsert()
	}
	// insert or update chain
	if _, err := q.
		Model(&blockObj).
		On("CONFLICT (chain) DO UPDATE").
		Set("last_processed_block_number=EXCLUDED.last_processed_block_number").
		Exec(context.Background()); err != nil {
		return err
	}
	return nil
}

func SelectEventListeners(chain int) (listeners []common.EventListener, err error) {
	if err = db.ORM.NewSelect().
		Model((*common.EventListener)(nil)).
		Where("chain=?", chain).
		Scan(context.Background(), &listeners); err != nil {
		return nil, err
	}
	return listeners, nil
}
