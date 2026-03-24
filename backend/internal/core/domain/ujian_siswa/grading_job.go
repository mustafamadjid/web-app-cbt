package ujian

type JobStatus string

const (
	StatusQueued     JobStatus = "queued"
	StatusProcessing JobStatus = "processing"
	StatusDone  JobStatus = "done"
	StatusFailed     JobStatus = "failed"
)

type GradingJob struct {
	IDgradingJob int
	IDAttempt    int
	Status       JobStatus
	RetryCount   int
	MaxRetries   int
	AvailableAt  string
	LockedAt     string
	ErrorCode    string
	ErrorMessage string
}