package httpx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeAndValidateAttemptUjianRequest(t *testing.T) {
	valid := AttemptUjianRequest{IdSiswa: 1, IdJadwalUjian: 2, TokenUjian: " token-abc ", WaktuMulai: time.Now()}
	tests := []struct {
		name      string
		mutate    func(*AttemptUjianRequest)
		wantErr   bool
		wantToken string
	}{
		{name: "Path 1 -> invalid student id is rejected", mutate: func(v *AttemptUjianRequest) { v.IdSiswa = 0 }, wantErr: true},
		{name: "Path 2 -> invalid schedule id is rejected", mutate: func(v *AttemptUjianRequest) { v.IdJadwalUjian = 0 }, wantErr: true},
		{name: "Path 3 -> zero start time is rejected", mutate: func(v *AttemptUjianRequest) { v.WaktuMulai = time.Time{} }, wantErr: true},
		{name: "Path 4 -> blank token is rejected", mutate: func(v *AttemptUjianRequest) { v.TokenUjian = " " }, wantErr: true},
		{name: "Path 5 -> valid token is trimmed and uppercased", wantToken: "TOKEN-ABC"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := valid
			if tt.mutate != nil {
				tt.mutate(&req)
			}
			got, err := sanitizeAndValidateAttemptUjianRequest(req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantToken, got.TokenUjian)
		})
	}
}
