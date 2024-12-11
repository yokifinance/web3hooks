package tests

import (
	"testing"
)

func TestSimple(t *testing.T) {
	t.Run("successOnFirstTry", successOnFirstTry)

	t.Run("alwaysFailingWebhook", alwaysFailingWebhook)
}
