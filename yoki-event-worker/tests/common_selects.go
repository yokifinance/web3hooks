package tests

import (
	"context"
	"testing"

	"yoki.finance/common/connectors/db"
	"yoki.finance/common/gue_jobs"
)

func selectAllGueJobs(t *testing.T) []gue_jobs.GueJob {
	var jobs []gue_jobs.GueJob

	if err := db.ORM.NewSelect().
		Order("created_at").
		Model((*gue_jobs.GueJob)(nil)).
		Scan(context.Background(), &jobs); err != nil {
		t.Fatalf("cannot select jobs: %s", err)
	}

	return jobs
}
