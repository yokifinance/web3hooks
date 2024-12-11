package tests

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"yoki.finance/common/gue_jobs"
	"yoki.finance/webhooks"
)

const (
	testWebhookQueue = "webhook_test"

	jobTypeWebhook = "webhook_job_test"
)

type tWebhookJobData struct {
	Id string
}

func successOnFirstTry(t *testing.T) {
	clearDbData(t)

	ctx, cancel := context.WithCancel(context.Background())

	// create server
	successServer := CreateWebhookTestServer(0)

	// create queue and insert a job
	jobQueue, err := webhooks.CreateWebhookJobQueue(testWebhookQueue, jobTypeWebhook)
	assert.NoError(t, err)
	jobId, err := jobQueue.EnqueueWebhookJobTx(nil, "http://localhost:8085/", 0, tWebhookJobData{"successOnFirstTry"})
	assert.NoError(t, err)

	// check job exists
	job, err := gue_jobs.SelectGueJob(jobId)
	assert.NoError(t, err)
	assert.Equal(t, 0, job.ErrorCount)

	successFuncChan := make(chan struct{}, 1)
	// create executor
	executor, err := webhooks.CreateQueueExecutor[tWebhookJobData, struct{}](
		ctx, testWebhookQueue, jobTypeWebhook, func(execContext *webhooks.WebhookExecContext, args tWebhookJobData, webhookRequestErr error) error {
			assert.FailNow(t, "must not be called")
			return nil
		},
		func(execContext *webhooks.WebhookExecContext, twjd tWebhookJobData, s struct{}) error {
			assert.Equal(t, execContext.WebhookUrl, "http://localhost:8085/")
			assert.NotEmpty(t, execContext.WebhookJobId)
			successFuncChan <- struct{}{}
			return nil
		})
	assert.NoError(t, err)

	// first call (will be successful)
	header := <-successServer.WebhookHeaderChan
	_, hasLastTryHeader := header["Webhook-Last-Try"]
	assert.True(t, hasLastTryHeader) // webhook must indicate that it will be last

	jobData := <-successServer.WebhookUrlOpenedChan
	receivedJobData := tWebhookJobData{}
	assert.NoError(t, json.Unmarshal(jobData, &receivedJobData))
	assert.Equal(t, "successOnFirstTry", receivedJobData.Id)

	// additionally making sure a success func is called
	<-successFuncChan

	// wait a bit – since job deletion is done after return from successFund – job must be removed
	time.Sleep(200 * time.Millisecond)
	job, err = gue_jobs.SelectGueJob(jobId)
	assert.NoError(t, err)
	assert.Nil(t, job, "job must be removed already")

	cancel()
	executor.WaitFinish()
	successServer.Stop()
}

/*
Server that always fails.
Webhook job that can tolerate only 2 error, then deleted.
*/
func alwaysFailingWebhook(t *testing.T) {
	clearDbData(t)

	ctx, cancel := context.WithCancel(context.Background())

	// create server
	alwaysFailingServer := CreateWebhookTestServer(999999)

	// create queue and insert a job
	jobQueue, err := webhooks.CreateWebhookJobQueue(testWebhookQueue, jobTypeWebhook)
	assert.NoError(t, err)
	jobId, err := jobQueue.EnqueueWebhookJobTx(nil, "http://localhost:8085/", 2, tWebhookJobData{"some_id"})
	assert.NoError(t, err)

	// check job exists
	job, err := gue_jobs.SelectGueJob(jobId)
	assert.NoError(t, err)
	assert.Equal(t, 0, job.ErrorCount)

	// create executor
	executor, err := webhooks.CreateQueueExecutor[tWebhookJobData, struct{}](ctx, testWebhookQueue, jobTypeWebhook, nil, nil)
	assert.NoError(t, err)

	// first call (will be failed)
	jobData := <-alwaysFailingServer.WebhookUrlOpenedChan
	receivedJobData := tWebhookJobData{}
	assert.NoError(t, json.Unmarshal(jobData, &receivedJobData))
	assert.Equal(t, "some_id", receivedJobData.Id)

	header := <-alwaysFailingServer.WebhookHeaderChan
	_, hasLastTryHeader := header["Webhook-Last-Try"]
	assert.False(t, hasLastTryHeader) // webhook must indicate that it is not last try

	// second call (will be also failed)
	<-alwaysFailingServer.WebhookUrlOpenedChan

	// check job (we are in the middle of second webhook call, so error count must be 1)
	job, err = gue_jobs.SelectGueJob(jobId)
	assert.NoError(t, err)
	assert.Equal(t, 1, job.ErrorCount)
	assert.NotNil(t, job.LastError)

	header = <-alwaysFailingServer.WebhookHeaderChan
	_, hasLastTryHeader = header["Webhook-Last-Try"]
	assert.True(t, hasLastTryHeader) // webhook must indicate that it is the last try

	// wait a bit – job must be removed
	time.Sleep(time.Second)
	job, err = gue_jobs.SelectGueJob(jobId)
	assert.NoError(t, err)
	assert.Nil(t, job, "job must be removed already")

	cancel()
	executor.WaitFinish()
	alwaysFailingServer.Stop()
}
