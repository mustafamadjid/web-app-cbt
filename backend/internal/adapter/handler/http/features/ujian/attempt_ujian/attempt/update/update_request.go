package httpx

import "time"

type UpdateAttemptUjianRequest struct {
	StatusAttempt *string    `json:"status_attempt"`
	WaktuSubmit   *time.Time `json:"waktu_submit"`
}
