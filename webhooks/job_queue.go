package webhooks

import (
	"context"
	"encoding/json"

	"github.com/uptrace/bun"
	"github.com/vgarvardt/gue/v5"
	"github.com/vgarvardt/gue/v5/adapter/libpq"
	"yoki.finance/common/connectors/db"
)

type WebhookJobQueue interface {
	EnqueueWebhookJobTx(tx *bun.Tx, webhookUrl string, maxErrorCount int32, jobArgs any) (jobId string, err error)
}

type JobArgsWrapper struct {
	WebhookUrl       string
	MaxErrorCount    int32
	Args             []byte
	ResultWebhookUrl string
}

type JobQueue struct {
	gc             *gue.Client
	webhookQueue   string
	jobTypeWebhook string
}

func CreateWebhookJobQueue(webhookQueue, jobTypeWebhook string) (WebhookJobQueue, error) {
	gc, err := gue.NewClient(libpq.NewConnPool(db.Conn))
	if err != nil {
		return nil, err
	}
	return &JobQueue{
		gc:             gc,
		webhookQueue:   webhookQueue,
		jobTypeWebhook: jobTypeWebhook,
	}, nil
}

func (q *JobQueue) EnqueueWebhookJobTx(tx *bun.Tx, webhookUrl string, maxErrorCount int32, jobArgs any) (jobId string, err error) {
	args, err := json.Marshal(jobArgs)
	if err != nil {
		return "", err
	}

	wrappedArgs, err := json.Marshal(
		JobArgsWrapper{
			WebhookUrl:    webhookUrl,
			MaxErrorCount: maxErrorCount,
			Args:          args,
		})
	if err != nil {
		return "", err
	}

	j := &gue.Job{
		Type:  q.jobTypeWebhook,
		Queue: q.webhookQueue,
		Args:  wrappedArgs,
	}

	if tx != nil {
		err = q.gc.EnqueueTx(context.Background(), j, db.CreateWrappedBunTx(*tx))
	} else {
		err = q.gc.Enqueue(context.Background(), j)
	}
	if err != nil {
		return "", err
	}

	id := j.ID.String()
	return id, nil
}
