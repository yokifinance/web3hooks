package tests

import (
	"testing"
)

func TestUtils(t *testing.T) {
	eventWorkerTestInit(t)

	t.Run("lastProcessedBlockTest", lastProcessedBlockTest)
}

func TestSimple(t *testing.T) {
	eventWorkerTestInit(t)

	// must be first – clean Db
	t.Run("simpleTestFromHeader", simpleTestFromHeader)

	t.Run("simpleTestSpecificLogs", simpleTestSpecificLogs)
	t.Run("simpleTestWebhookCall", simpleTestWebhookCall)

	t.Run("simpleMultipleAddressesTest", simpleMultipleAddressesTest)
}

func TestWebhookReturns(t *testing.T) {
	eventWorkerTestInit(t)

	t.Run("returnBodyInGoodFormatWebhookCall", returnBodyInGoodFormatWebhookCall)
	t.Run("returnBodyInBadFormatWebhookCall", returnBodyInBadFormatWebhookCall)
}

func TestChunking(t *testing.T) {
	t.Run("chunkingTest", chunkingTest)
}
