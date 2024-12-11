package gue_jobs

import (
	"context"
	"fmt"

	"yoki.finance/common/connectors/db"
)

func SelectGueJob(jobId string) (*GueJob, error) {
	var jobs []GueJob

	if err := db.ORM.NewSelect().
		Model((*GueJob)(nil)).
		Where("job_id=?", jobId).
		Limit(1).
		Scan(context.Background(), &jobs); err != nil {

		return nil, fmt.Errorf("cannot select job %s: %s", jobId, err)
	}

	if len(jobs) == 0 {
		return nil, nil
	}

	return &jobs[0], nil
}
