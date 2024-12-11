package worker

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"os"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"yoki.finance/common/connectors/db"
	"yoki.finance/common/rcommon"
	"yoki.finance/webhooks"
	yoki_common "yoki.finance/yoki-event-worker/common"
	"yoki.finance/yoki-event-worker/data"
	"yoki.finance/yoki-event-worker/webhook"
)

// Webhooks settings
const (
	chainWorkerPanicDelay = 5 * time.Second

	// if any worker restarts more than this times, we will exit the program and let
	// the shell restart it again.
	unhandledCountLimit = 10

	maxProcessEventsRetries = 15
)

var (
	webhookMaxErrorCount = rcommon.GetParamIntOrDefault("WEBHOOK_MAX_ERROR_COUNT", 3)
)

// Block constants
const (
	// Number of last blocks that are ingored when handling events (in other words, number of block confirmation the worker waits)
	lastIgnoredBlocks = 3
)

var (
	/*
		It is max depth from lastIgnoredBlocks where worker still reaches logs.

		latestConsiderableBlock = headerBlock - lastIgnoredBlocks

		Example:
			maxBlocksDepth is 10 0008
			lastProcessedBlock (from Db) is 100 000
			latestConsiderableBlock is 200 000
			=> worker will start from block 190 000.

		As of Nov 2023, average block is produced every 2.5 seconds (Polygon, Optimism) and 12 seconds on Ethereum.
		So, 10 000 depth gives us (10000*2.5)/3600=7 hours depth
	*/
	maxBlocksDepth = int64(rcommon.GetParamIntOrDefault("MAX_BLOCKS_DEPTH", 300)) // about 300 , ~10 minutes by default

	// max number of addresses in eth_getLogs search query
	addressChunkSize = rcommon.GetParamIntOrDefault("ADDRESS_CHUNK_SIZE", 100)

	// max number of blocks in one chunk when requesting blockchain logs
	// This makes sense after first worker launch or after pause (e.g. restart) when worker "caches up" quickly
	maxBlockChunkSize = int64(rcommon.GetParamIntOrDefault("MAX_BLOCK_CHUNK_SIZE", 100))

	// minimum number of blocks in one chunk after which worker start processing
	// This makes sense when worker is in progress, so as soon MIN_BLOCK_CHUNK_SIZE blocks are available, worker starts processing.
	// 10 blocks is about 25 seconds delay on Polygon (~2.5 seconds per block)
	minBlockChunkSize = int64(rcommon.GetParamIntOrDefault("MIN_BLOCK_CHUNK_SIZE", 10))

	// Seconds to sleep if blocks chunk is incomplete.
	// The less this value is, the more requests to blockchain are made, but the faster worker launches chunk processing.
	incompleteChunkSleepSeconds = rcommon.GetParamIntOrDefault("INCOMPLETE_CHUNK_SLEEP_SECS", 10)
)

type ChainWorker struct {
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	client   EthClient
	jobQueue webhooks.WebhookJobQueue
	chain    int
}

type jobType struct {
	webhook.EventListenerResultDto
	WebhookUrl string
}

func Create(globalCtx context.Context, client EthClient, chain int) (*ChainWorker, error) {
	jobQueue, err := webhooks.CreateWebhookJobQueue(webhook.EventWebhookQueue, webhook.EventJobTypeWebhook)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(globalCtx)
	return &ChainWorker{
		client:   client,
		jobQueue: jobQueue,
		chain:    chain,
		ctx:      ctx,
		cancel:   cancel,
		wg:       sync.WaitGroup{},
	}, nil
}

func init() {
	assert.True(nil, maxBlockChunkSize < maxBlocksDepth)
}

func (w *ChainWorker) Run() {
	go w.chainWorker(0)
}

func (w *ChainWorker) Stop() {
	w.log("Stop called")
	w.cancel()
	w.wg.Wait()
	w.log("ChainWorker exited")
}

func (w *ChainWorker) Wait() {
	w.wg.Wait()
}

func (w *ChainWorker) chainWorker(restartCount int) {
	defer func() {
		if r := recover(); r != nil {
			w.log("chain %d: chainWorker: %s stacktrace: %s", w.chain, r, strings.ReplaceAll(string(debug.Stack()), "\n", `\n`))

			if restartCount+1 < unhandledCountLimit {
				w.log("chain %d: chainWorker restarting (%d/%d)", w.chain, restartCount+1, unhandledCountLimit)

				time.Sleep(chainWorkerPanicDelay)
				go w.chainWorker(restartCount + 1)
			} else {
				w.log("chain %d: chainWorker error limit exceeded (%d/%d). Shutting down the app", w.chain, unhandledCountLimit, unhandledCountLimit)
				panic(r)
			}
		}
	}()
	defer w.wg.Done()
	w.wg.Add(1)

	// simple exponential backoff (https://en.wikipedia.org/wiki/Exponential_backoff) algorithm, params are mine.
	errorCount := 0
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
			err := w.nextBlockRange()
			if err != nil {
				errorCount++
				delay := math.Pow(1.3, float64(errorCount))
				sleepDuration := time.Duration(delay) * time.Second

				text := ""
				willShutdown := errorCount >= maxProcessEventsRetries
				if willShutdown {
					text = "Shutting down worker process"
				} else {
					text = fmt.Sprintf("Will sleep for %s", sleepDuration)
				}

				w.log("processEvents error (%d/%d): %s: %s", errorCount, maxProcessEventsRetries, err, text)
				if willShutdown {
					os.Exit(1)
				}

				time.Sleep(sleepDuration)
			} else {
				errorCount = 0
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

func (w *ChainWorker) nextBlockRange() error {
	// for each new block we select all listeners – for simplicity let's leave it for now.
	allEventListeners, err := data.SelectEventListeners(w.chain)
	if err != nil {
		return err
	}

	return w.processBlockRangeLogs(allEventListeners)
}

// processes blocks for the next available range
func (w *ChainWorker) processBlockRangeLogs(allEventListeners []yoki_common.EventListener) error {
	if len(allEventListeners) == 0 {
		return nil
	}
	lastProcessedBlockNumber, err := data.GetLastProcessedBlockNumber(w.chain)
	if err != nil {
		return err
	}
	header, err := w.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	w.log("blocks: header %s", header.Number)

	latestConsiderableBlock := new(big.Int).Sub(header.Number, big.NewInt(lastIgnoredBlocks))
	fromBlock, toBlock, err := w.calculateNextBlockRange(lastProcessedBlockNumber, latestConsiderableBlock)
	if err != nil {
		return err
	}

	toBlock = min(toBlock, latestConsiderableBlock)

	rangeBlockCount := new(big.Int).Add( // length of the blocks range
		new(big.Int).Sub(toBlock, fromBlock),
		big.NewInt(1),
	)

	// sometimes block range length is less than 0, so we set it to 0 for better display
	if rangeBlockCount.Cmp(big.NewInt(0)) < 0 {
		rangeBlockCount = big.NewInt(0)
	}

	if rangeBlockCount.Cmp(big.NewInt(minBlockChunkSize)) < 0 { // don't have enough blocks to process
		w.log("Block chunk not full: %s of %d. Sleeping..", rangeBlockCount, minBlockChunkSize)
		time.Sleep(time.Duration(incompleteChunkSleepSeconds) * time.Second)
		return nil
	}

	// this is the map of all addresses of all contracts that we are listening to indexes of event listeners
	addressListenersMap := w.mapListenersOnAddress(allEventListeners)
	// we chunk addresses not to exceed eth_getLogs limitations
	chunkedAddresses := ChunkAddresses(addressListenersMap, addressChunkSize)

	// then we load logs chunk by chunk for the same block range
	blockRangeLogs, err := w.queryLogsForBlockRange(chunkedAddresses, fromBlock, toBlock)
	if err != nil {
		return err
	}

	// here we got all logs loaded for the range fromBlock - toBlock for all our smart contracts.
	// let's transform them to jobs
	blockRangeJobs := w.createJobsFromLogs(blockRangeLogs, addressListenersMap)

	w.log("processing %s-%s, listeners: %d, address chunks: %d, addresses chunk size: %d, logs received: %d, jobs created: %d",
		fromBlock, toBlock, len(allEventListeners), len(chunkedAddresses), addressChunkSize, len(blockRangeLogs), len(blockRangeJobs))

	if err := w.enqueueJobs(blockRangeJobs); err != nil {
		return err
	}

	// set last processed block to the end of interval – we checked inside of it already.
	err = data.SetLastProcessedBlockNumberTx(nil, w.chain, toBlock)

	return err
}

// calculates the next block range to process logs
func (w *ChainWorker) calculateNextBlockRange(lastProcessedBlockNumber, latestConsiderableBlock *big.Int) (fromBlock, toBlock *big.Int, err error) {
	// if set by test env, we use this deepest considerable block.
	// This is to test "old" logs of smart contract to test how event worker is performing.
	var firstConsiderableBlock *big.Int
	testFirstConsiderableBlock := rcommon.GetParamStrOrDefault("TEST_DEEPEST_CONSIDERABLE_BLOCK", "")
	if testFirstConsiderableBlock != "" {
		firstConsiderableBlock, _ = new(big.Int).SetString(testFirstConsiderableBlock, 10)
	} else {
		firstConsiderableBlock = new(big.Int).Sub(latestConsiderableBlock, big.NewInt(int64(maxBlocksDepth)))
	}

	if // no records about last processed logs
	lastProcessedBlockNumber.Cmp(big.NewInt(0)) == 0 ||
		// or information hasn't been updated for a long time
		lastProcessedBlockNumber.Cmp(firstConsiderableBlock) < 0 {

		// we just take first block range (right after deepestConsiderableBlock)
		fromBlock = firstConsiderableBlock
		toBlock = new(big.Int).Add(fromBlock, big.NewInt(maxBlockChunkSize-1))
	} else {
		fromBlock = new(big.Int).Add(lastProcessedBlockNumber, big.NewInt(1)) // start with the next after lastProcessed
		toBlock = new(big.Int).Add(fromBlock, big.NewInt(maxBlockChunkSize-1))
	}

	return fromBlock, toBlock, nil
}

type addressToListenersMapType map[common.Address][]*yoki_common.EventListener

// load all logs for specified addresses (in chunks) for from-to range blocks.
// Also sort them by blockNumber
func (w *ChainWorker) queryLogsForBlockRange(chunkedAddresses [][]common.Address, fromBlock, toBlock *big.Int) ([]types.Log, error) {
	blockRangeLogs := []types.Log{}
	for _, addressChunk := range chunkedAddresses {
		chunkLogs, err := w.queryLogsForAddressChunk(addressChunk, fromBlock, toBlock)
		if err != nil {
			return nil, err
		}
		blockRangeLogs = append(blockRangeLogs, chunkLogs...)
	}

	return blockRangeLogs, nil
}

// creates jobs from logs considering listeners.
// I.e., for one log we can have any number of listeners.
func (w *ChainWorker) createJobsFromLogs(logs []types.Log, addressListenersMap addressToListenersMapType) []jobType {
	blockRangeJobs := []jobType{}
	for _, log := range logs { // iterate over all received logs
		addressEventListeners := addressListenersMap[log.Address]

		// iterave over all listeners for that address and create a job for each
		for _, listener := range addressEventListeners {
			job := jobType{
				EventListenerResultDto: webhook.EventListenerResultDto{
					EventListenerId: listener.Id,
					Chain:           w.chain,
					Address:         log.Address,
					BlockHash:       log.BlockHash,
					BlockNumber:     log.BlockNumber,
					LogIndex:        log.Index,
					TxHash:          log.TxHash,
					TxIndex:         log.TxIndex,
					Data:            common.Bytes2Hex(log.Data),
					Topics:          log.Topics,
				},
				WebhookUrl: listener.WebhookUrl,
			}
			blockRangeJobs = append(blockRangeJobs, job)
		}
	}
	return blockRangeJobs
}

// enqueues jobs block by block
func (w *ChainWorker) enqueueJobs(jobs []jobType) error {
	if len(jobs) == 0 {
		return nil
	}
	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].BlockNumber < jobs[j].BlockNumber
	})

	if jobs[0].BlockNumber > jobs[len(jobs)-1].BlockNumber {
		panic("uncorrect sorting")
	}

	lastBlockNumber := jobs[0].BlockNumber
	blockJobs := []jobType{jobs[0]}
	for i := 1; i <= len(jobs); i++ {
		if i == len(jobs) || jobs[i].BlockNumber != lastBlockNumber { // block changed or jobs completed
			// here we process jobs for <lastBlockNumber>
			tx, err := db.ORM.Begin()
			if err != nil {
				return err
			}
			defer tx.Rollback()

			for _, job := range blockJobs {
				w.jobQueue.EnqueueWebhookJobTx(&tx, job.WebhookUrl, int32(webhookMaxErrorCount), job.EventListenerResultDto)
			}

			if err := data.SetLastProcessedBlockNumberTx(&tx, w.chain, new(big.Int).SetUint64(lastBlockNumber)); err != nil {
				return err
			}
			if err := tx.Commit(); err != nil {
				return err
			}
			w.log("enqueueJobs: block %d: %d jobs added", lastBlockNumber, len(blockJobs))

			if i != len(jobs) {
				blockJobs = []jobType{jobs[i]} // restart chunk again, block has changed
				lastBlockNumber = jobs[i].BlockNumber
			}
		} else {
			blockJobs = append(blockJobs, jobs[i])
		}
	}
	return nil
}

func (w *ChainWorker) queryLogsForAddressChunk(currentAddressChunk []common.Address, fromBlock, toBlock *big.Int) ([]types.Log, error) {
	query := ethereum.FilterQuery{Addresses: currentAddressChunk, FromBlock: fromBlock, ToBlock: toBlock}
	logs, err := w.client.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
