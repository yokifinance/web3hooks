package tests

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"yoki.finance/common/config"
	"yoki.finance/common/connectors/db"
	"yoki.finance/webhooks"
	wt "yoki.finance/webhooks/tests"
	"yoki.finance/yoki-event-worker/data"
	"yoki.finance/yoki-event-worker/webhook"
	"yoki.finance/yoki-event-worker/worker"
)

/*
tests for two events of contract 0xe467fab1e5ddA1AAf51Ad3d4e10a2667e9efF2c3 (see test for more details)
As a result, two jobs are created. We don't even create executor for them – just check they're added.
*/
func simpleTestSpecificLogs(t *testing.T) {
	clearGueJobs(t)

	// set worker to consider blocks starting from this
	t.Setenv("TEST_DEEPEST_CONSIDERABLE_BLOCK", "50169056")

	// catch logs for test
	//	-- Unsubscribe	50169057	https://polygonscan.com/tx/0xca0c73f34a177711a4bd627eac38578da86cc5f03bd6be62b83d701ca84d7afc
	//	-- Subscribe	50169082	https://polygonscan.com/tx/0x2b129348479b3eba287b7df66cd6ba86b26d6816bc4d60b2535160e6f763f110
	//	Difference is 25 blocks (~10 seconds)
	data.SetLastProcessedBlockNumberTx(nil, 137, new(big.Int).SetInt64(50169056))

	chain := 137

	client, err := worker.NewEthClient(config.RpcEndpoints[chain])
	assert.NoError(t, err)
	worker, err := worker.Create(context.Background(), client, chain)
	assert.NoError(t, err)
	worker.Run()
	time.Sleep(time.Second * 2)
	worker.Stop()

	jobs := selectAllGueJobs(t)
	assert.Equal(t, 2, len(jobs))
	args1, wrp1 := unwrapArgs(t, &jobs[0])
	args2, wrp2 := unwrapArgs(t, &jobs[1])

	assert.NotEmpty(t, wrp1.WebhookUrl)
	assert.Equal(t, int32(3), wrp1.MaxErrorCount)
	assert.NotEmpty(t, wrp2.WebhookUrl)
	assert.Equal(t, int32(3), wrp2.MaxErrorCount)

	assert.Equal(t, 137, args1.Chain)
	assert.Equal(t, common.HexToAddress("0xe467fab1e5ddA1AAf51Ad3d4e10a2667e9efF2c3"), args1.Address)
	assert.Equal(t, uint64(50169057), args1.BlockNumber)
	assert.Equal(t, common.HexToHash("0x18d0621149ecf449d154c57f283b7368ecea3d64cb5ef5a91d716facc48f86b9"), args1.BlockHash)
	assert.Equal(t, "000000000000000000000000e467fab1e5dda1aaf51ad3d4e10a2667e9eff2c30000000000000000000000002d9a8be931f1eab82abfcb9697023424e440cd43", args1.Data)
	assert.Equal(t, 1, len(args1.Topics))
	assert.Equal(t, common.HexToHash("0x7773c30acd0762ed6b4b92a9aa2c6b3c074e29ad93b334cbed8ba807c596f13a"), args1.Topics[0])

	assert.Equal(t, 137, args2.Chain)
	assert.Equal(t, common.HexToAddress("0xe467fab1e5ddA1AAf51Ad3d4e10a2667e9efF2c3"), args2.Address)
	assert.Equal(t, uint64(50169082), args2.BlockNumber)
	assert.Equal(t, common.HexToHash("0x6f6f2152e683461bbb1f4e129b02b916889e80a713dcf5de2a6ff8070c8f175f"), args2.BlockHash)
	assert.Equal(t, "000000000000000000000000e467fab1e5dda1aaf51ad3d4e10a2667e9eff2c300000000000000000000000000000000000000000000000000000000000000a00000000000000000000000002d9a8be931f1eab82abfcb9697023424e440cd4300000000000000000000000000000000000000000000000000000000655b35ba00000000000000000000000000000000000000000000000000000000655b35f60000000000000000000000000000000000000000000000000000000000000015596f6b692d6578616d706c652d6d65726368616e740000000000000000000000", args2.Data)
	assert.Equal(t, 1, len(args2.Topics))
	assert.Equal(t, common.HexToHash("0xb2afd60ec89ad71deb13ea0be4a196179313144abff030e53bcbaba1ea4856da"), args2.Topics[0])
}

func simpleTestFromHeader(t *testing.T) {
	clearGueJobs(t)
	t.Setenv("TEST_DEEPEST_CONSIDERABLE_BLOCK", "")

	lastBlock, err := data.GetLastProcessedBlockNumber(137)
	assert.NoError(t, err)
	assert.Equal(t, *big.NewInt(0), *lastBlock) // must be nothing processed

	chain := 137
	client := &TestEthClient{
		HeaderNumber: 10000,
	}

	worker, err := worker.Create(context.Background(), client, chain)
	assert.NoError(t, err)
	worker.Run()
	time.Sleep(time.Second * 5)
	client.HeaderNumber = 100000

	time.Sleep(time.Second * 3)

	worker.Stop()
}

/*
Tests that webhook is called and completed with no error
Dummy TestEthClient is used.
*/
func simpleTestWebhookCall(t *testing.T) {
	clearGueJobs(t)
	t.Setenv("TEST_DEEPEST_CONSIDERABLE_BLOCK", "")

	ctx, cancel := context.WithCancel(context.Background())

	server := wt.CreateWebhookTestServer(0)

	queueExecutor, err := webhooks.CreateQueueExecutor[webhook.EventListenerResultDto, struct{}](
		ctx, webhook.EventWebhookQueue, webhook.EventJobTypeWebhook, nil, nil)
	assert.NoError(t, err)

	chain := 137
	data.SetLastProcessedBlockNumberTx(nil, chain, big.NewInt(49))

	client := &TestEthClient{
		HeaderNumber: 100,
	}

	worker, err := worker.Create(ctx, client, chain)
	assert.NoError(t, err)
	worker.Run()

	// wait till success execution of webhook
	res := <-server.WebhookSuccessChan
	var listenerResult webhook.EventListenerResultDto
	err = json.Unmarshal(res.RequestBody, &listenerResult)
	assert.NoError(t, err)
	// check that it was this listener (see TestEthClient)
	assert.Equal(t, uint64(50), listenerResult.BlockNumber)
	assert.Equal(t, common.HexToAddress("0xe467fab1e5ddA1AAf51Ad3d4e10a2667e9efF2c3"), listenerResult.Address)

	jobs := selectAllGueJobs(t)
	assert.Equal(t, 0, len(jobs)) // we test success webhook call, so must not be any jobs

	cancel() // cancel everything

	// wait till worker and executor stopped
	worker.Wait()
	queueExecutor.WaitFinish()

	server.Stop()
}

// Just multiple listeners with different addresses
func simpleMultipleAddressesTest(t *testing.T) {
	createEventListener137(t, "0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	createEventListener137(t, "0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
	createEventListener137(t, "0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC")

	clearGueJobs(t)
	t.Setenv("TEST_DEEPEST_CONSIDERABLE_BLOCK", "")

	ctx, cancel := context.WithCancel(context.Background())
	server := wt.CreateWebhookTestServer(0)

	queueExecutor, err := webhooks.CreateQueueExecutor[webhook.EventListenerResultDto, struct{}](
		ctx, webhook.EventWebhookQueue, webhook.EventJobTypeWebhook, nil, nil)
	assert.NoError(t, err)

	chain := 137
	data.SetLastProcessedBlockNumberTx(nil, chain, big.NewInt(4999))

	client := &TestEthClient{
		HeaderNumber: 5050,
	}

	worker, err := worker.Create(ctx, client, chain)
	assert.NoError(t, err)
	worker.Run()

	// wait till success execution of webhook
	res := <-server.WebhookSuccessChan
	var listenerResult webhook.EventListenerResultDto
	err = json.Unmarshal(res.RequestBody, &listenerResult)
	assert.NoError(t, err)
	// check that it was this listener (see TestEthClient)
	assert.Equal(t, uint64(5000), listenerResult.BlockNumber)
	assert.Equal(t, common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"), listenerResult.Address)

	jobs := selectAllGueJobs(t)
	assert.Equal(t, 0, len(jobs)) // we test success webhook call, so must not be any jobs

	cancel() // cancel everything
	// wait till worker and executor stopped
	worker.Wait()
	queueExecutor.WaitFinish()

	server.Stop()
}

func createEventListener137(t *testing.T, address string) {
	_, err := db.Conn.Exec(`INSERT INTO public.event_listeners 
		 (chain, client_id, address, webhook_url, created_timestamp, active)
		VALUES
		 (137, '221fabd2-9df0-44e7-98d4-4eda556c4143', $1, 
		 	$2, NOW(), TRUE);
		`, address, "http://localhost:8085/")
	assert.NoError(t, err)
}
