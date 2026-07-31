package httpx

import (
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestParseAndValidateSubmitUjianRequest(t *testing.T) {
	tests := []struct {
		name, raw      string
		direct         int
		parse, wantErr bool
		want           int
	}{
		{name: "Path 1 -> missing route id is rejected", parse: true, wantErr: true},
		{name: "Path 2 -> non numeric route id is rejected", raw: "abc", parse: true, wantErr: true},
		{name: "Path 3 -> zero route id fails validation", raw: "0", parse: true, wantErr: true},
		{name: "Path 4 -> whitespace is trimmed from valid route id", raw: " 42 ", parse: true, want: 42},
		{name: "Path 5 -> negative direct id is rejected", direct: -1, wantErr: true},
		{name: "Path 6 -> positive direct id is accepted", direct: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.parse {
				got, err := parseSubmitUjianRequest(httprouter.Params{{Key: "idAttempt", Value: tt.raw}})
				assert.Equal(t, tt.wantErr, err != nil)
				assert.Equal(t, tt.want, got.IDAttempt)
				return
			}
			assert.Equal(t, tt.wantErr, ValidateSubmitUjianRequest(SubmitUjianRequest{IDAttempt: tt.direct}) != nil)
		})
	}
}
