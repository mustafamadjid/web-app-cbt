package httpx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeAndValidateUpdateAttemptUjianRequest(t *testing.T) {
	submitted, expired, blank, invalid := " SUBMITTED ", "expired", " ", "active"
	zero, now := time.Time{}, time.Now()
	tests := []struct {
		name       string
		req        UpdateAttemptUjianRequest
		wantErr    bool
		wantStatus string
	}{
		{name: "Path 1 -> empty patch is rejected", wantErr: true},
		{name: "Path 2 -> blank status is rejected", req: UpdateAttemptUjianRequest{StatusAttempt: &blank}, wantErr: true},
		{name: "Path 3 -> unsupported status is rejected", req: UpdateAttemptUjianRequest{StatusAttempt: &invalid}, wantErr: true},
		{name: "Path 4 -> submitted status is normalized", req: UpdateAttemptUjianRequest{StatusAttempt: &submitted}, wantStatus: "submitted"},
		{name: "Path 5 -> expired status is accepted", req: UpdateAttemptUjianRequest{StatusAttempt: &expired}, wantStatus: "expired"},
		{name: "Path 6 -> zero submit time is rejected", req: UpdateAttemptUjianRequest{WaktuSubmit: &zero}, wantErr: true},
		{name: "Path 7 -> valid submit time is accepted", req: UpdateAttemptUjianRequest{WaktuSubmit: &now}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sanitizeAndValidateUpdateAttemptUjianRequest(tt.req)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantStatus != "" && assert.NotNil(t, got.StatusAttempt) {
				assert.Equal(t, tt.wantStatus, *got.StatusAttempt)
			}
		})
	}
}
