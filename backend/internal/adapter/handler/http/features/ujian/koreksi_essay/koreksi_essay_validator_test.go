package httpx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeAndValidateKoreksiEssayRequest(t *testing.T) {
	benar := true
	tests := []struct {
		name    string
		req     KoreksiEssayRequest
		wantErr bool
	}{
		{name: "Path 1 -> empty answer list is accepted"},
		{name: "Path 2 -> non positive answer id is rejected", req: KoreksiEssayRequest{Jawaban: []KoreksiEssayItemRequest{{IDJawaban: 0, EssayIsBenar: &benar}}}, wantErr: true},
		{name: "Path 3 -> missing essay result is rejected", req: KoreksiEssayRequest{Jawaban: []KoreksiEssayItemRequest{{IDJawaban: 1}}}, wantErr: true},
		{name: "Path 4 -> complete essay result is accepted", req: KoreksiEssayRequest{Jawaban: []KoreksiEssayItemRequest{{IDJawaban: 1, EssayIsBenar: &benar}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sanitizeAndValidateKoreksiEssayRequest(tt.req)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, tt.req, got)
			}
		})
	}
}
