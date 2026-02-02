package profil_sekolah_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/profil_sekolah"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/profil_sekolah"
	profil_sekolah_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/profil_sekolah/update"
	"github.com/stretchr/testify/assert"
)

type fakeUpdateProfilSekolahRepo struct {
	updateErr error

	updateCalled bool
	lastID       profil_sekolah.IDProfil
	lastProfil   profil_sekolah.ProfilSekolah
}

func (f *fakeUpdateProfilSekolahRepo) UpdateProfilSekolah(ctx context.Context, idProfil profil_sekolah.IDProfil, profil profil_sekolah.ProfilSekolah) error {
	f.updateCalled = true
	f.lastID = idProfil
	f.lastProfil = profil
	if f.updateErr != nil {
		return f.updateErr
	}
	return nil
}

func (f *fakeUpdateProfilSekolahRepo) GetProfilSekolah(ctx context.Context) (profil_sekolah.ProfilSekolah, error) {
	panic("not used in this test")
}

var _ out.ProfilSekolahRepository = (*fakeUpdateProfilSekolahRepo)(nil)

func TestUpdateProfilSekolahService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	makeCmd := func() profil_sekolah_service.UpdateProfilSekolahCmd {
		return profil_sekolah_service.UpdateProfilSekolahCmd{
			IDProfil:      1,
			EmailSekolah:  "admin@example.com",
			NoTelpSekolah: "081234567",
			KepalaSekolah: "Kepala",
			WakaSekolah:   "Waka",
			NamaSekolah:   "Sekolah",
			AlamatSekolah: "Alamat",
			LogoSekolah:   nil,
		}
	}

	tests := []struct {
		name           string
		repoErr        error
		mutateCmd      func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd)
		expectErr      error
		expectUpdate   bool
		expectLogo     *string
		expectNamaTrim string
	}{
		{
			name: "Branch 1 -> id profil tidak valid",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				cmd.IDProfil = 0
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectUpdate: false,
		},
		{
			name: "Branch 2 -> email sekolah kosong",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				cmd.EmailSekolah = " "
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectUpdate: false,
		},
		{
			name: "Branch 3 -> no telp sekolah kosong",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				cmd.NoTelpSekolah = ""
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectUpdate: false,
		},
		{
			name: "Branch 4 -> kepala sekolah kosong",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				cmd.KepalaSekolah = ""
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectUpdate: false,
		},
		{
			name: "Branch 5 -> waka sekolah kosong",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				cmd.WakaSekolah = ""
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectUpdate: false,
		},
		{
			name: "Branch 6 -> nama sekolah kosong",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				cmd.NamaSekolah = ""
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectUpdate: false,
		},
		{
			name: "Branch 7 -> alamat sekolah kosong",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				cmd.AlamatSekolah = ""
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectUpdate: false,
		},
		{
			name: "Branch 8 -> logo sekolah kosong",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				blankLogo := "  "
				cmd.LogoSekolah = &blankLogo
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectUpdate: false,
		},
		{
			name:         "Branch 9 -> repo update gagal",
			repoErr:      errors.New("update failed"),
			expectErr:    errors.New("update failed"),
			expectUpdate: true,
		},
		{
			name: "Branch 10 -> semua patch berhasil",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				logo := " logo.png "
				cmd.LogoSekolah = &logo
				cmd.NamaSekolah = " Sekolah Baru "
			},
			expectUpdate:   true,
			expectNamaTrim: "Sekolah Baru",
			expectLogo:     ptrString("logo.png"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeUpdateProfilSekolahRepo{updateErr: tc.repoErr}
			svc := profil_sekolah_service.NewUpdateProfilSekolahService(repo)

			cmd := makeCmd()
			if tc.mutateCmd != nil {
				tc.mutateCmd(&cmd)
			}

			err := svc.UpdateProfilSekolah(ctx, cmd)

			if tc.expectErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectUpdate, repo.updateCalled)

			if tc.expectUpdate && tc.expectErr == nil {
				assert.Equal(t, cmd.IDProfil, repo.lastID)
				if tc.expectNamaTrim != "" {
					assert.Equal(t, tc.expectNamaTrim, repo.lastProfil.NamaSekolah)
				}
				if tc.expectLogo != nil {
					assert.NotNil(t, repo.lastProfil.LogoSekolah)
					assert.Equal(t, *tc.expectLogo, *repo.lastProfil.LogoSekolah)
				}
			}
		})
	}
}

func ptrString(value string) *string {
	return &value
}
