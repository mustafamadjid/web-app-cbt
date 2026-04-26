package importsoal

import "time"

type JobStatus string

const (
	StatusPending    JobStatus = "pending"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

type ImportSoalJob struct {
	IDJob      int64
	IDBankSoal int64
	IDPengguna int64
	Status     JobStatus
	FilePath   string
	ErrorMsg   string
	WarningMsg string
	TotalSoal  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
