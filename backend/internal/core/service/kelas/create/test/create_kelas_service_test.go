package kelas_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	kelas_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/create"
	"github.com/stretchr/testify/assert"
)

func TestCreateKelasService_CreateTingkatKelas(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	tingkatKelas := 10
	existErr := errors.New("exist tingkat kelas error")
	createErr := errors.New("create tingkat kelas error")

	tests := []struct {
		name       string
		repo       *FakeKelasRepo
		wantErr    error
		wantCreate bool
	}{
		{
			name:       "Branch 1 -> tingkat kelas exist",
			repo:       &FakeKelasRepo{ExistTingkatKelasRet: true},
			wantErr:    coreerror.ErrTingkatKelasExist,
			wantCreate: false,
		},
		{
			name:       "Branch 2 -> gagal cek tingkat kelas",
			repo:       &FakeKelasRepo{ExistTingkatKelasErr: existErr},
			wantErr:    existErr,
			wantCreate: false,
		},
		{
			name:       "Branch 3 -> gagal create tingkat kelas",
			repo:       &FakeKelasRepo{CreateTingkatKelasErr: createErr},
			wantErr:    createErr,
			wantCreate: true,
		},
		{
			name:       "Branch 4 -> berhasil create tingkat kelas",
			repo:       &FakeKelasRepo{},
			wantErr:    nil,
			wantCreate: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := kelas_service.NewCreateKelasService(tt.repo)
			err := svc.CreateTingkatKelas(ctx, kelas_service.CreateTingkatKelasCmd{TingkatKelas: tingkatKelas})

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.True(t, tt.repo.ExistTingkatKelasCalled)
			assert.Equal(t, tingkatKelas, tt.repo.GotExistTingkatKelas)
			assert.Equal(t, tt.wantCreate, tt.repo.CreateTingkatKelasCalled)
			if tt.repo.CreateTingkatKelasCalled {
				assert.Equal(t, tingkatKelas, tt.repo.GotCreateTingkatKelas)
			}
		})
	}
}

func TestCreateKelasService_CreateNamaKelas(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	inputNamaKelas := "  IPA 1  "
	trimmedNamaKelas := "IPA 1"
	idTingkatKelas := kelas.ID(7)
	existErr := errors.New("exist nama kelas error")
	createErr := errors.New("create nama kelas error")

	tests := []struct {
		name               string
		repo               *FakeKelasRepo
		wantErr            error
		wantCreate         bool
		wantExistNamaKelas string
	}{
		{
			name:               "Branch 1 -> nama kelas exist",
			repo:               &FakeKelasRepo{ExistNamaKelasRet: true},
			wantErr:            coreerror.ErrNamaKelasExist,
			wantCreate:         false,
			wantExistNamaKelas: trimmedNamaKelas,
		},
		{
			name:               "Branch 2 -> gagal cek nama kelas",
			repo:               &FakeKelasRepo{ExistNamaKelasErr: existErr},
			wantErr:            existErr,
			wantCreate:         false,
			wantExistNamaKelas: trimmedNamaKelas,
		},
		{
			name:               "Branch 3 -> gagal create nama kelas",
			repo:               &FakeKelasRepo{CreateNamaKelasErr: createErr},
			wantErr:            createErr,
			wantCreate:         true,
			wantExistNamaKelas: trimmedNamaKelas,
		},
		{
			name:               "Branch 4 -> berhasil create nama kelas",
			repo:               &FakeKelasRepo{},
			wantErr:            nil,
			wantCreate:         true,
			wantExistNamaKelas: trimmedNamaKelas,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := kelas_service.NewCreateKelasService(tt.repo)
			err := svc.CreateNamaKelas(ctx, kelas_service.CreateNamaKelasCmd{
				IdTingkatKelas: idTingkatKelas,
				NamaKelas:      inputNamaKelas,
			})

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.True(t, tt.repo.ExistNamaKelasCalled)
			assert.Equal(t, tt.wantExistNamaKelas, tt.repo.GotExistNamaKelas)
			assert.Equal(t, tt.wantCreate, tt.repo.CreateNamaKelasCalled)
			if tt.repo.CreateNamaKelasCalled {
				assert.Equal(t, kelas.NamaKelas{
					IdTingkatKelas: idTingkatKelas,
					NamaKelas:      trimmedNamaKelas,
				}, tt.repo.GotCreateNamaKelas)
			}
		})
	}
}
