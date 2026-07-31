package httpx

import (
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/stretchr/testify/assert"
)

func TestValidateUpdateUjianRequestFields(t *testing.T) {
	value, blank := " value ", " "
	tests := []struct {
		name    string
		req     UpdatePenjadwalanUjianRequest
		wantErr bool
	}{
		{name: "Path 1 -> omitted text fields are accepted"},
		{name: "Path 2 -> valid text fields are trimmed", req: UpdatePenjadwalanUjianRequest{NamaUjian: &value, Token: &value, DeskripsiUjian: &value, StatusUjian: &value}},
		{name: "Path 3 -> blank exam name is rejected", req: UpdatePenjadwalanUjianRequest{NamaUjian: &blank}, wantErr: true},
		{name: "Path 4 -> blank token is rejected", req: UpdatePenjadwalanUjianRequest{Token: &blank}, wantErr: true},
		{name: "Path 5 -> blank description is rejected", req: UpdatePenjadwalanUjianRequest{DeskripsiUjian: &blank}, wantErr: true},
		{name: "Path 6 -> blank status is rejected", req: UpdatePenjadwalanUjianRequest{StatusUjian: &blank}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateUpdateUjianRequestFields(tt.req)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestValidateUpdateUjianRequestIDs(t *testing.T) {
	positive, zero, negative := 1, 0, -1
	tests := []struct {
		name    string
		req     UpdatePenjadwalanUjianRequest
		wantErr error
	}{
		{name: "Path 1 -> omitted ids are accepted"},
		{name: "Path 2 -> positive ids and zero optional class-name id are accepted", req: UpdatePenjadwalanUjianRequest{IdBankSoal: &positive, IdKelas: &positive, IdNamaKelas: &zero, IdGuru: &positive, IdSesi: &positive, IdRuangan: &positive, IdPengawas: &positive}},
		{name: "Path 3 -> invalid bank id is rejected", req: UpdatePenjadwalanUjianRequest{IdBankSoal: &zero}, wantErr: coreerror.ErrMissingId},
		{name: "Path 4 -> invalid class id is rejected", req: UpdatePenjadwalanUjianRequest{IdKelas: &zero}, wantErr: coreerror.ErrMissingId},
		{name: "Path 5 -> negative class-name id is rejected", req: UpdatePenjadwalanUjianRequest{IdNamaKelas: &negative}, wantErr: coreerror.ErrMissingId},
		{name: "Path 6 -> invalid teacher id is rejected", req: UpdatePenjadwalanUjianRequest{IdGuru: &zero}, wantErr: coreerror.ErrMissingId},
		{name: "Path 7 -> invalid session id is rejected", req: UpdatePenjadwalanUjianRequest{IdSesi: &zero}, wantErr: coreerror.ErrMissingId},
		{name: "Path 8 -> invalid room id is rejected", req: UpdatePenjadwalanUjianRequest{IdRuangan: &zero}, wantErr: coreerror.ErrMissingId},
		{name: "Path 9 -> invalid supervisor id is rejected", req: UpdatePenjadwalanUjianRequest{IdPengawas: &zero}, wantErr: coreerror.ErrMissingId},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { assert.ErrorIs(t, ValidateUpdateUjianRequestIDs(tt.req), tt.wantErr) })
	}
}
