package sesi_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	sesi_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/create"
	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/create/fake_test"
	"github.com/stretchr/testify/assert"
)

func TestCreateSesiService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	existErr := errors.New("exist kode sesi error")
	createErr := errors.New("create sesi error")

	tests := []struct {
		name       string
		repo       *fake_test.FakeSesiRepo
		input      sesi.Sesi
		wantErr    error
		wantCreate bool
	}{
		{
			name: "Path 1 -> NamaSesi kosong",
			repo: &fake_test.FakeSesiRepo{},
			input: sesi.Sesi{
				NamaSesi: "   ",
				KodeSesi: "SESI01",
			},
			wantErr:    coreerror.ErrMissingField,
			wantCreate: false,
		},
		{
			name: "Path 2 -> KodeSesi kosong",
			repo: &fake_test.FakeSesiRepo{},
			input: sesi.Sesi{
				NamaSesi: "Sesi 1",
				KodeSesi: "   ",
			},
			wantErr:    coreerror.ErrMissingField,
			wantCreate: false,
		},
		{
			name: "Path 3 -> ExistByKodeSesi error",
			repo: &fake_test.FakeSesiRepo{ExistByKodeSesiErr: existErr},
			input: sesi.Sesi{
				NamaSesi: "Sesi 1",
				KodeSesi: "SESI01",
			},
			wantErr:    existErr,
			wantCreate: false,
		},
		{
			name: "Path 4 -> kode sesi sudah exist",
			repo: &fake_test.FakeSesiRepo{ExistByKodeSesiRet: true},
			input: sesi.Sesi{
				NamaSesi: "Sesi 1",
				KodeSesi: "SESI01",
			},
			wantErr:    coreerror.ErrSesiUjianExist,
			wantCreate: false,
		},
		{
			name: "Path 5 -> CreateSesi error",
			repo: &fake_test.FakeSesiRepo{CreateSesiErr: createErr},
			input: sesi.Sesi{
				NamaSesi: "Sesi 1",
				KodeSesi: "SESI01",
			},
			wantErr:    createErr,
			wantCreate: true,
		},
		{
			name: "Path 6 -> berhasil create sesi",
			repo: &fake_test.FakeSesiRepo{},
			input: sesi.Sesi{
				NamaSesi: "  Sesi 1  ",
				KodeSesi: "  sesi01  ",
			},
			wantErr:    nil,
			wantCreate: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := sesi_service.NewCreateSesiService(tt.repo)
			err := svc.CreateSesiService(ctx, tt.input)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantCreate, tt.repo.CreateSesiCalled)
		})
	}
}
