package db

import (
	"context"
	"database/sql"
	"regexp"
	"strconv"

	"github.com/uptrace/bun"
	"github.com/vgarvardt/gue/v5/adapter"
)

/*
This is an object that allows to execute Bun and gue operations in one transaction.
The idea is to create an instance of WrappedBunTx and use it in EnqueueTx function of gue.
Internally WrappedBunTx contains bun.Tx.
Also, Exec function is implemented.
*/
type WrappedBunTx struct {
	adapter.Tx

	bunTx *bun.Tx
}

func CreateWrappedBunTx(bunTx bun.Tx) WrappedBunTx {
	return WrappedBunTx{bunTx: &bunTx}
}

func (tx WrappedBunTx) Exec(ctx context.Context, query string, args ...any) (adapter.CommandTag, error) {
	query = queryFormatFromSqlPgToBun(query)
	res, err := tx.bunTx.ExecContext(ctx, query, args...)
	return aCommandTag{res}, err
}

/*
Converts query from pg/sql args format like "%1, %2, .." to bun format "?0, ?1, .." (%->? and zero-based indexes)
*/
func queryFormatFromSqlPgToBun(input string) string {
	re := regexp.MustCompile(`\$(\d+)`)
	output := re.ReplaceAllStringFunc(input, func(s string) string {
		numStr := s[1:]
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return s
		}
		return "?" + strconv.Itoa(num-1)
	})
	return output
}

type aCommandTag struct {
	ct sql.Result
}

// RowsAffected implements adapter.CommandTag.RowsAffected() using github.com/lib/pq
func (ct aCommandTag) RowsAffected() int64 {
	ra, err := ct.ct.RowsAffected()
	if err != nil {
		return 0
	}
	return ra
}
