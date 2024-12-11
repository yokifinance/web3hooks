package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"yoki.finance/common/rcommon"
	"yoki.finance/webhooks"
)

type WebhookResultArgs struct {
	WebhookJobId string
	Success      bool
	ErrorMessage string
}

const (
	resultJobQueueName = "resultWebhook_queue"
	resultJobTypeName  = "resultWebhook_jobType"
)

func main() {
	rcommon.Println("yoki-webhook-executor is running...")
	ctx, cancel := context.WithCancel(context.Background())

	// creating job queue to enqueue result webhook jobs
	resultJobQueue, err := webhooks.CreateWebhookJobQueue(resultJobQueueName, resultJobTypeName)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Create executor for result jobs – jobs that are used to inform about webhook call results
	resultJobsExecutor, err := webhooks.CreateQueueExecutor(
		ctx, resultJobQueueName, resultJobTypeName, func(execContext *webhooks.WebhookExecContext, args any, webhookRequestErr error) error {
			rcommon.Println("yoki-webhook-executor err - result job: url: %s job: %s: %s", execContext.WebhookUrl, execContext.WebhookJobId, webhookRequestErr.Error())
			return nil
		},
		func(execContext *webhooks.WebhookExecContext, twjd any, s struct{}) error {
			rcommon.Println("yoki-webhook-executor success - result job: url: %s job: %s", execContext.WebhookUrl, execContext.WebhookJobId)
			return nil
		})

	if err != nil {
		log.Fatalf(err.Error())
	}

	// if necessary created a result webhook job
	lazyCreateResultWebhookJob := func(execContext *webhooks.WebhookExecContext, success bool, errorMessage string) (resultJobString string) {
		if execContext.ResultWebhookUrl != "" {
			// if failed job has a result webhook URL then create a job to inform about the result.
			args := WebhookResultArgs{
				WebhookJobId: execContext.WebhookJobId,
				Success:      success,
				ErrorMessage: errorMessage,
			}
			// create job once, if not successful - we omit it to simplify code. This is just a result webhook job so its delivery is not critical
			resultJobId, err := resultJobQueue.EnqueueWebhookJobTx(nil, execContext.ResultWebhookUrl, 5, args)
			if err != nil {
				resultJobString = ". tried to create result job, failed: " + err.Error()
			} else {
				resultJobString = fmt.Sprintf(". created result job id: %s resultUrl: %s", resultJobId, execContext.ResultWebhookUrl)
			}
		}
		return resultJobString
	}

	// create executor for webhook service jobs (the ones put by webhook API)
	executor, err := webhooks.CreateQueueExecutor(
		ctx, "webhookService_queue", "webhookService_jobType", func(execContext *webhooks.WebhookExecContext, args any, webhookRequestErr error) error {
			resultJobString := lazyCreateResultWebhookJob(execContext, false, webhookRequestErr.Error())
			rcommon.Println("yoki-webhook-executor err: url: %s job: %s: %s %s", execContext.WebhookUrl, execContext.WebhookJobId, webhookRequestErr.Error(), resultJobString)
			return nil
		},
		func(execContext *webhooks.WebhookExecContext, twjd any, s struct{}) error {
			resultJobString := lazyCreateResultWebhookJob(execContext, true, "")
			rcommon.Println("yoki-webhook-executor success: url: %s job: %s %s", execContext.WebhookUrl, execContext.WebhookJobId, resultJobString)
			return nil
		})
	if err != nil {
		log.Fatalf(err.Error())
	}

	rcommon.Println("yoki-webhook-executor is running. Press Ctrl-C to stop")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan

	rcommon.Println("yoki-webhook-executor stop called")
	cancel()
	executor.WaitFinish()
	resultJobsExecutor.WaitFinish()
	rcommon.Println("yoki-webhook-executor stopped")
}
