package webhooks

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/vgarvardt/backoff"
	"github.com/vgarvardt/gue/v5/adapter/libpq"
	"golang.org/x/sync/errgroup"
	"yoki.finance/common/config"
	"yoki.finance/common/connectors/db"
	"yoki.finance/common/rcommon"

	"github.com/vgarvardt/gue/v5"
)

type WebhookExecContext struct {
	WebhookUrl       string
	WebhookJobId     string
	ResultWebhookUrl string
}

type QueueExecutor[JobArgsType any, JobResultType any] struct {
	ctx         context.Context
	errGroup    *errgroup.Group
	gc          *gue.Client
	errorFunc   func(*WebhookExecContext, JobArgsType, error) error
	successFunc func(*WebhookExecContext, JobArgsType, JobResultType) error
}

func CreateQueueExecutor[JobArgsType any, JobResultType any](
	ctx context.Context,
	webhookQueue,
	jobTypeWebhook string,

	/* called when maxErrorCount exceeded. errorFunc can potentially "prolong" execution returning non-null error (but it is not
	recommended to do). In this case webhook will be called again.
	*/
	errorFunc func(*WebhookExecContext, JobArgsType, error) error,

	/*
		Called when webhook successfully called.  can potentially "prolong" execution returning non-null error (but it is not
		recommended to do). In this case webhook will be called again.
	*/
	successFunc func(*WebhookExecContext, JobArgsType, JobResultType) error,
) (*QueueExecutor[JobArgsType, JobResultType], error) {
	rcommon.Println("CreateQueueExecutor called: queue %s, jobType %s", webhookQueue, jobTypeWebhook)

	var backoffConf backoff.Config
	if !config.IsInTests() {
		backoffConf = backoff.Config{
			BaseDelay:  10.0 * time.Second,
			Multiplier: 2.9,
			Jitter:     0.1,
			MaxDelay:   1.0 * time.Hour,
		}
	} else {
		backoffConf = backoff.DefaultConfig
	}
	bo := gue.WithClientBackoff(gue.NewExponentialBackoff(backoffConf))

	gc, err := gue.NewClient(libpq.NewConnPool(db.Conn), bo)
	if err != nil {
		return nil, err
	}

	errGroup, _ := errgroup.WithContext(ctx)
	executor := &QueueExecutor[JobArgsType, JobResultType]{
		ctx:         ctx,
		errGroup:    errGroup,
		gc:          gc,
		errorFunc:   errorFunc,
		successFunc: successFunc,
	}

	wm := gue.WorkMap{
		jobTypeWebhook: func(ctx context.Context, j *gue.Job) error {
			return workerFunc(executor, j)
		},
	}
	// poll interval set to 1 second, but default is 5
	workers, err := gue.NewWorkerPool(gc, wm, 2, gue.WithPoolQueue(webhookQueue), gue.WithPoolPollInterval(time.Second))
	if err != nil {
		log.Fatal(err)
	}

	errGroup.Go(func() error {
		err := workers.Run(ctx)
		if err != nil {
			// In a real-world applications, use a better way to shut down
			// application on unrecoverable error. E.g. fx.Shutdowner from
			// go.uber.org/fx module.
			log.Fatal(err)
		}
		return err
	})

	rcommon.Println("CreateQueueExecutor completed: queue %s, jobType %s", webhookQueue, jobTypeWebhook)

	return executor, nil
}

func (qe *QueueExecutor[JobArgsType, JobResultType]) WaitFinish() {
	if err := qe.errGroup.Wait(); err != nil {
		log.Fatal(err)
	}
}

func workerFunc[JobArgsType any, JobResultType any](executor *QueueExecutor[JobArgsType, JobResultType], j *gue.Job) error {
	rcommon.Println("QueueExecutor: job %s picked up", j.ID)

	var argsWrapper JobArgsWrapper
	if err := json.Unmarshal(j.Args, &argsWrapper); err != nil {
		rcommon.Println("QueueExecutor: job %s: unmarshal wrapper: %s", j.ID, err.Error())
		j.Delete(executor.ctx)
		return nil
	}

	var args JobArgsType
	if err := json.Unmarshal(argsWrapper.Args, &args); err != nil {
		rcommon.Println("QueueExecutor: job %s: unmarshal args: %s", j.ID, err.Error())
		j.Delete(executor.ctx)
		return nil
	}

	execContext := &WebhookExecContext{
		WebhookUrl:       argsWrapper.WebhookUrl,
		WebhookJobId:     j.ID.String(),
		ResultWebhookUrl: argsWrapper.ResultWebhookUrl,
	}
	rcommon.Println("QueueExecutor: job %s calling %s", j.ID, argsWrapper.WebhookUrl)
	requester := &HttpRequester[JobArgsType, JobResultType]{}
	lastTry := j.ErrorCount+1 >= int32(argsWrapper.MaxErrorCount)
	response, webhookRequestErr := requester.Request(argsWrapper.WebhookUrl, &args, lastTry)
	if webhookRequestErr != nil { // unsuccessful webhook call

		if lastTry {
			rcommon.Println("QueueExecutor: job %s failed %d/%d, calling errorFunc: %s", j.ID, j.ErrorCount+1, argsWrapper.MaxErrorCount, webhookRequestErr)
			if executor.errorFunc == nil {
				j.Delete(executor.ctx)
				return nil
			}

			errorFuncErr := executor.errorFunc(execContext, args, webhookRequestErr)
			if errorFuncErr != nil {
				// if error is returned then we don't delete a job and just return this error
				// This is useful when we need to update some "final" state in the DB and this didn't work out, so we'll retry
				return errorFuncErr
			}

			j.Delete(executor.ctx)
			return nil
		}

		rcommon.Println("QueueExecutor: job %s failed %d/%d: %s", j.ID, j.ErrorCount+1, argsWrapper.MaxErrorCount, webhookRequestErr)
		return webhookRequestErr
	}

	rcommon.Println("QueueExecutor: job %s: webhook success", j.ID)

	if executor.successFunc != nil {
		err := executor.successFunc(execContext, args, *response)
		msg := ""
		if err != nil {
			msg = err.Error()
		} else {
			msg = "success"
		}
		rcommon.Println("QueueExecutor: job %s: successFunc called: %s", j.ID, msg)
		return err
	}

	rcommon.Println("QueueExecutor: job %s: completed", j.ID)
	return nil
}
