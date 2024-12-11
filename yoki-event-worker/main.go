package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"yoki.finance/common/config"
	"yoki.finance/common/rcommon"
	"yoki.finance/webhooks"
	"yoki.finance/yoki-event-worker/webhook"
	"yoki.finance/yoki-event-worker/worker"
)

var (
	queueExecutor *webhooks.QueueExecutor[webhook.EventListenerResultDto, struct{}]
	ctx           context.Context
	cancel        context.CancelFunc
	workers       = []*worker.ChainWorker{}
)

func main() {
	rcommon.Println("yoki-event-worker is starting...")
	ctx, cancel = context.WithCancel(context.Background())

	var err error
	queueExecutor, err = webhooks.CreateQueueExecutor[webhook.EventListenerResultDto, struct{}](ctx, webhook.EventWebhookQueue, webhook.EventJobTypeWebhook, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Init workers for chains
	for _, chain := range config.AppChainIds {
		client, err := worker.NewEthClient(config.RpcEndpoints[chain])
		if err != nil {
			log.Printf("NewEthClient failed: chain %d: ethclient.Dial: %s", chain, err)
			continue
		}

		worker, err := worker.Create(ctx, client, chain)
		if err != nil {
			log.Printf("worker.Create failed: chain %d: %s", chain, err)
			continue
		}
		workers = append(workers, worker)
		worker.Run()
	}

	rcommon.Println("yoki-event-worker is running. Press Ctrl-C to stop")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan

	rcommon.Println("yoki-event-worker stop called")
	cancel()
	queueExecutor.WaitFinish()
	for _, w := range workers {
		w.Wait()
	}
	rcommon.Println("yoki-event-worker stopped")
}
