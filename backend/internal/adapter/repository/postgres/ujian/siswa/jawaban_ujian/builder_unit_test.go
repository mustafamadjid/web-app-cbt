package jawabanujian_repo

import (
	"database/sql"
	"testing"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	"github.com/stretchr/testify/assert"
)

func TestSplitSaveJawabanItemsBasisPaths(t *testing.T) {
	choice := ujian.ID(8)
	essay, blank := " answer ", "  "
	tests := []struct {
		name                  string
		items                 []ujian.JawabanUjian
		wantUpsert, wantClear int
		wantEssay             string
	}{
		{name: "Path 1 -> unanswered question is selected for clearing", items: []ujian.JawabanUjian{{IdSoal: 1}}, wantClear: 1},
		{name: "Path 2 -> blank essay becomes unanswered and is cleared", items: []ujian.JawabanUjian{{IdSoal: 1, JawabanEssay: &blank}}, wantClear: 1},
		{name: "Path 3 -> choice answer is selected for upsert", items: []ujian.JawabanUjian{{IdSoal: 1, IdPilihan: &choice}}, wantUpsert: 1},
		{name: "Path 4 -> essay answer is trimmed and selected for upsert", items: []ujian.JawabanUjian{{IdSoal: 1, JawabanEssay: &essay}}, wantUpsert: 1, wantEssay: "answer"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upsert, clear := splitSaveJawabanItems(tt.items)
			assert.Len(t, upsert, tt.wantUpsert)
			assert.Len(t, clear, tt.wantClear)
			if tt.wantEssay != "" {
				assert.Equal(t, tt.wantEssay, *upsert[0].JawabanEssay)
			}
		})
	}
}

func TestNullableConverters(t *testing.T) {
	tests := []struct {
		name, kind  string
		valid       bool
		intValue    int64
		stringValue string
	}{
		{name: "Path 1 -> invalid nullable id becomes nil", kind: "id"},
		{name: "Path 2 -> valid nullable id becomes typed pointer", kind: "id", valid: true, intValue: 42},
		{name: "Path 3 -> invalid nullable string becomes nil", kind: "string"},
		{name: "Path 4 -> valid nullable string becomes pointer", kind: "string", valid: true, stringValue: "answer"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.kind == "id" {
				got := nullInt64ToUjianIDPtr(sql.NullInt64{Int64: tt.intValue, Valid: tt.valid})
				if !tt.valid {
					assert.Nil(t, got)
				} else {
					assert.Equal(t, ujian.ID(tt.intValue), *got)
				}
				return
			}
			got := nullStringToPtr(sql.NullString{String: tt.stringValue, Valid: tt.valid})
			if !tt.valid {
				assert.Nil(t, got)
			} else {
				assert.Equal(t, tt.stringValue, *got)
			}
		})
	}
}
