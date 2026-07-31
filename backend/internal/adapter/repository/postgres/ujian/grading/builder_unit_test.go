package gradingrepo

import (
	"testing"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/stretchr/testify/assert"
)

func TestBuildStatistikSoalPayload(t *testing.T) {
	items := []ujian.StatistikSoal{{IDSoal: 1, IDUjian: 2}, {IDSoal: 3, IDUjian: 2}}
	tests := []struct {
		name                            string
		items                           []ujian.StatistikSoal
		correct                         bool
		wantLen, wantCorrect, wantWrong int
	}{
		{name: "Path 1 -> empty statistics returns nil"},
		{name: "Path 2 -> correct answers increment correct count", items: items, correct: true, wantLen: 2, wantCorrect: 1},
		{name: "Path 3 -> wrong answers increment wrong count", items: items, wantLen: 2, wantWrong: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildStatistikSoalPayload(tt.items, tt.correct)
			assert.Len(t, got, tt.wantLen)
			for i := range got {
				assert.Equal(t, tt.items[i].IDSoal, got[i].IDSoal)
				assert.Equal(t, tt.wantCorrect, got[i].JumlahJawabanBenar)
				assert.Equal(t, tt.wantWrong, got[i].JumlahJawabanSalah)
			}
		})
	}
}
