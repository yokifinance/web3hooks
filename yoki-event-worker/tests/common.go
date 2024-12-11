package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"yoki.finance/common/connectors/db"
	"yoki.finance/common/test"
)

func eventWorkerTestInit(t *testing.T) {
	t.Setenv("INCOMPLETE_CHUNK_SLEEP_SECS", "0") // doesn't sleep in tests
	clearEventWorkerDbData(t)
	test.CreateTestDbChainClientData(t)
	addDummyEventListener(t)
}

// func eventWorkerTeardown(t *testing.T) {
// 	manage.StopEventWorker()
// }

func addDummyEventListener(t *testing.T) {
	_, err := db.Conn.Exec(`

		INSERT INTO public.event_listeners 
		 (chain, client_id, address, webhook_url, created_timestamp, active)
		VALUES
		 (137, '221fabd2-9df0-44e7-98d4-4eda556c4143', '0xe467fab1e5ddA1AAf51Ad3d4e10a2667e9efF2c3', 
		 	$1, NOW(), TRUE);
		`, "http://localhost:8085/")
	assert.NoError(t, err)
}

// clears table last_processed_blocks
func clearEventWorkerDbData(t *testing.T) {
	_, err := db.Conn.Exec(`
		DELETE FROM last_processed_blocks;
		DELETE FROM event_listeners;
		DELETE FROM gue_jobs;
		`)
	assert.NoError(t, err)
}

func clearGueJobs(t *testing.T) {
	_, err := db.Conn.Exec(`
		DELETE FROM gue_jobs;
		`)
	assert.NoError(t, err)
}
