package httpx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeAndValidateSaveJawabanRequest(t *testing.T) {
	choice, badChoice := 2, 0
	essay, blankEssay, invalidUTF8 := "  jawaban  ", "   ", string([]byte{0xff})
	now, zero := time.Now(), time.Time{}
	valid := SaveJawabanRequest{IDAttempt: 1, Jawaban: []SaveJawabanItemRequest{{IDSoal: 1, IDPilihan: &choice, WaktuJawab: &now}}}
	tests := []struct {
		name      string
		req       SaveJawabanRequest
		wantErr   bool
		wantEssay *string
	}{
		{name: "Path 1 -> invalid attempt id is rejected", req: SaveJawabanRequest{}, wantErr: true},
		{name: "Path 2 -> invalid question id is rejected", req: SaveJawabanRequest{IDAttempt: 1, Jawaban: []SaveJawabanItemRequest{{IDSoal: 0, IDPilihan: &choice}}}, wantErr: true},
		{name: "Path 3 -> invalid choice id is rejected", req: SaveJawabanRequest{IDAttempt: 1, Jawaban: []SaveJawabanItemRequest{{IDSoal: 1, IDPilihan: &badChoice}}}, wantErr: true},
		{name: "Path 4 -> blank essay becomes nil then fails exactly-one rule", req: SaveJawabanRequest{IDAttempt: 1, Jawaban: []SaveJawabanItemRequest{{IDSoal: 1, JawabanEssay: &blankEssay}}}, wantErr: true},
		{name: "Path 5 -> invalid UTF-8 essay is rejected", req: SaveJawabanRequest{IDAttempt: 1, Jawaban: []SaveJawabanItemRequest{{IDSoal: 1, JawabanEssay: &invalidUTF8}}}, wantErr: true},
		{name: "Path 6 -> choice and essay together are rejected", req: SaveJawabanRequest{IDAttempt: 1, Jawaban: []SaveJawabanItemRequest{{IDSoal: 1, IDPilihan: &choice, JawabanEssay: &essay}}}, wantErr: true},
		{name: "Path 7 -> zero answer time is rejected", req: SaveJawabanRequest{IDAttempt: 1, Jawaban: []SaveJawabanItemRequest{{IDSoal: 1, IDPilihan: &choice, WaktuJawab: &zero}}}, wantErr: true},
		{name: "Path 8 -> valid choice answer is accepted", req: valid},
		{name: "Path 9 -> valid essay is trimmed", req: SaveJawabanRequest{IDAttempt: 1, Jawaban: []SaveJawabanItemRequest{{IDSoal: 1, JawabanEssay: &essay}}}, wantEssay: stringPointer("jawaban")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sanitizeAndValidateSaveJawabanRequest(tt.req)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantEssay != nil && assert.NotNil(t, got.Jawaban[0].JawabanEssay) {
				assert.Equal(t, *tt.wantEssay, *got.Jawaban[0].JawabanEssay)
			}
		})
	}
}

func stringPointer(value string) *string { return &value }
