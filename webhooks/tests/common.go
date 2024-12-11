package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"yoki.finance/common/connectors/db"
)

func clearDbData(t *testing.T) {
	_, err := db.Conn.Exec(`
		DELETE FROM gue_jobs;
		`)
	assert.NoError(t, err)
}
