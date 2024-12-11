package gue_jobs

type GueJob struct {
	JobId      string
	ErrorCount int
	LastError  *string
	Args       []byte
}
