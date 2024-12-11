package tests

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"yoki.finance/common/gue_jobs"
	"yoki.finance/webhooks"
	"yoki.finance/yoki-event-worker/webhook"
)

func unwrapArgs(t *testing.T, job *gue_jobs.GueJob) (args webhook.EventListenerResultDto, wrapper webhooks.JobArgsWrapper) {
	assert.NoError(t, json.Unmarshal(job.Args, &wrapper))
	assert.NoError(t, json.Unmarshal(wrapper.Args, &args))
	return args, wrapper
}
