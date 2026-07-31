package isisoalbatchrepo

import (
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	"github.com/stretchr/testify/assert"
)

func TestValidateSingleCorrectOption(t *testing.T) {
	correct := importsoal.ParsedOpsi{IsBenar: true}
	wrong := importsoal.ParsedOpsi{}
	tests := []struct {
		name    string
		items   []importsoal.ParsedSoal
		wantErr bool
	}{
		{name: "Branch 1 -> essay skips option validation", items: []importsoal.ParsedSoal{{TipeSoal: "essay"}}},
		{name: "Branch 2 -> multiple choice with exactly one correct option passes", items: []importsoal.ParsedSoal{{TipeSoal: "pilihan_ganda", Opsi: []importsoal.ParsedOpsi{correct, wrong}}}},
		{name: "Branch 3 -> multiple choice without correct option fails", items: []importsoal.ParsedSoal{{TipeSoal: "pilihan_ganda", Opsi: []importsoal.ParsedOpsi{wrong}}}, wantErr: true},
		{name: "Branch 4 -> multiple choice with multiple correct options fails", items: []importsoal.ParsedSoal{{TipeSoal: "pilihan_ganda", Opsi: []importsoal.ParsedOpsi{correct, correct}}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSingleCorrectOption(tt.items)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantErr {
				assert.ErrorIs(t, err, coreerror.ErrInvalidInput)
			}
		})
	}
}

func TestBuilderUtilities(t *testing.T) {
	tests := []struct {
		name, kind, image string
		a, b              int
		want              any
	}{
		{name: "Branch 1 -> blank image normalizes to nil", kind: "image", image: "  ", want: nil},
		{name: "Branch 2 -> nonblank image remains unchanged", kind: "image", image: " image.png ", want: " image.png "},
		{name: "Branch 3 -> min returns first smaller integer", kind: "min", a: 1, b: 2, want: 1},
		{name: "Branch 4 -> min returns second integer when equal or smaller", kind: "min", a: 2, b: 1, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got any
			if tt.kind == "image" {
				got = normalizeSoalImage(tt.image)
			} else {
				got = minInt(tt.a, tt.b)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
