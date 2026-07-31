package httpx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCreateUjianRequestFields(t *testing.T) {
	description := " description "
	valid := CreatePenjadwalanUjianRequest{NamaUjian: " exam ", Token: " token ", StatusUjian: " draft ", DeskripsiUjian: &description}
	tests := []struct {
		name    string
		mutate  func(*CreatePenjadwalanUjianRequest)
		wantErr bool
	}{
		{name: "Path 1 -> blank exam name is rejected", mutate: func(v *CreatePenjadwalanUjianRequest) { v.NamaUjian = " " }, wantErr: true},
		{name: "Path 2 -> blank token is rejected", mutate: func(v *CreatePenjadwalanUjianRequest) { v.Token = " " }, wantErr: true},
		{name: "Path 3 -> blank provided description is rejected", mutate: func(v *CreatePenjadwalanUjianRequest) { s := " "; v.DeskripsiUjian = &s }, wantErr: true},
		{name: "Path 4 -> blank status is rejected", mutate: func(v *CreatePenjadwalanUjianRequest) { v.StatusUjian = " " }, wantErr: true},
		{name: "Path 5 -> valid fields are trimmed"},
		{name: "Path 6 -> omitted description is accepted", mutate: func(v *CreatePenjadwalanUjianRequest) { v.DeskripsiUjian = nil }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := valid
			if tt.mutate != nil {
				tt.mutate(&req)
			}
			got, err := ValidateCreateUjianRequestFields(req)
			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, "exam", got.NamaUjian)
				assert.Equal(t, "token", got.Token)
				assert.Equal(t, "draft", got.StatusUjian)
			}
		})
	}
}
