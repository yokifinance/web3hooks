// Contains tests of different cases webhook can return upon URL call

package tests

import (
	"context"
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"yoki.finance/webhooks"
	wt "yoki.finance/webhooks/tests"
	"yoki.finance/yoki-event-worker/data"
	"yoki.finance/yoki-event-worker/webhook"
	"yoki.finance/yoki-event-worker/worker"
)

/*
Tests that webhook is returning body in good format
*/
func returnBodyInGoodFormatWebhookCall(t *testing.T) {
	clearGueJobs(t)
	t.Setenv("TEST_DEEPEST_CONSIDERABLE_BLOCK", "")

	ctx, cancel := context.WithCancel(context.Background())

	type resultType struct{ Value string }
	server := wt.CreateWebhookTestServer(0)

	webhookReturnBody, err := json.Marshal(resultType{Value: "test"}) // the object that will be returned by webhook
	assert.NoError(t, err)
	server.SetReturnBody(webhookReturnBody)

	queueExecutor, err := webhooks.CreateQueueExecutor[webhook.EventListenerResultDto, resultType](
		ctx, webhook.EventWebhookQueue, webhook.EventJobTypeWebhook, nil, nil)
	assert.NoError(t, err)

	chain := 137
	data.SetLastProcessedBlockNumberTx(nil, chain, big.NewInt(49))

	client := &TestEthClient{HeaderNumber: 100}
	worker, err := worker.Create(ctx, client, chain)
	assert.NoError(t, err)
	worker.Run()

	// wait till success execution of webhook
	res := <-server.WebhookSuccessChan
	var webhookInput webhook.EventListenerResultDto
	var webhookOutput resultType

	err = json.Unmarshal(res.ResponseBody, &webhookOutput)
	assert.NoError(t, err)
	assert.Equal(t, "test", webhookOutput.Value)

	err = json.Unmarshal(res.RequestBody, &webhookInput)
	assert.NoError(t, err)
	// check that it was this listener (see TestEthClient)
	assert.Equal(t, uint64(50), webhookInput.BlockNumber)
	assert.Equal(t, common.HexToAddress("0xe467fab1e5ddA1AAf51Ad3d4e10a2667e9efF2c3"), webhookInput.Address)

	jobs := selectAllGueJobs(t)
	assert.Equal(t, 0, len(jobs)) // we test success webhook call, so must not be any jobs
	// error not "unexpected end of JSON input"

	cancel() // cancel everything

	// wait till worker and executor stopped
	worker.Wait()
	queueExecutor.WaitFinish()

	server.Stop()
}

func returnBodyInBadFormatWebhookCall(t *testing.T) {
	clearGueJobs(t)
	t.Setenv("TEST_DEEPEST_CONSIDERABLE_BLOCK", "")

	ctx, cancel := context.WithCancel(context.Background())

	server := wt.CreateWebhookTestServer(0)

	server.SetReturnBody([]byte("bad format"))
	queueExecutor, err := webhooks.CreateQueueExecutor[webhook.EventListenerResultDto, struct{ Value string }](
		ctx, webhook.EventWebhookQueue, webhook.EventJobTypeWebhook, nil, nil)
	assert.NoError(t, err)

	chain := 137
	data.SetLastProcessedBlockNumberTx(nil, chain, big.NewInt(49))

	client := &TestEthClient{HeaderNumber: 100}
	worker, err := worker.Create(ctx, client, chain)
	assert.NoError(t, err)
	worker.Run()

	// wait till success execution of webhook
	res := <-server.WebhookSuccessChan
	var webhookInput webhook.EventListenerResultDto
	err = json.Unmarshal(res.RequestBody, &webhookInput)
	assert.NoError(t, err)
	// check that it was this listener (see TestEthClient)
	assert.Equal(t, uint64(50), webhookInput.BlockNumber)
	assert.Equal(t, common.HexToAddress("0xe467fab1e5ddA1AAf51Ad3d4e10a2667e9efF2c3"), webhookInput.Address)

	jobs := selectAllGueJobs(t)
	assert.Equal(t, 1, len(jobs)) // must be error
	assert.NotNil(t, jobs[0].LastError)
	assert.True(t, strings.HasPrefix(*jobs[0].LastError, "webhook returned body not in expected format. return empty body or conform to expected format"))

	cancel() // cancel everything

	// wait till worker and executor stopped
	worker.Wait()
	queueExecutor.WaitFinish()

	server.Stop()
}
