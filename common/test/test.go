package test

import (
	"testing"

	"yoki.finance/common/connectors/db"
)

func CreateTestDbChainClientData(t *testing.T) {
	_, err := db.Conn.Exec(`
	INSERT INTO clients (id, name, secret_key) SELECT '221fabd2-9df0-44e7-98d4-4eda556c4143', 'Test client', 'secret' 
	WHERE NOT EXISTS (SELECT 1 FROM clients WHERE id = '221fabd2-9df0-44e7-98d4-4eda556c4143');

	INSERT INTO supported_chains (chain, name, rpc_url) SELECT 137, 'Polygon', 'RPC' 
	WHERE NOT EXISTS (SELECT 1 FROM supported_chains WHERE chain = 137);

	INSERT INTO supported_chains (chain, name, rpc_url) SELECT 80001, 'Polygon Mumbai Testnet', 'RPC' 
	WHERE NOT EXISTS (SELECT 1 FROM supported_chains WHERE chain = 80001);
	`)
	if err != nil {
		t.Fatal(err)
	}
}
