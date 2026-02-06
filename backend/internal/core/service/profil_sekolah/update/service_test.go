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
	getErr error
	profil profil_sekolah.ProfilSekolah

	updateErr error

	getCalled    bool
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
	f.getCalled = true
	if f.getErr != nil {
		return profil_sekolah.ProfilSekolah{}, f.getErr
	}
	return f.profil, nil
}

var _ out.ProfilSekolahRepository = (*fakeUpdateProfilSekolahRepo)(nil)

func TestUpdateProfilSekolahService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	makeCmd := func() profil_sekolah_service.UpdateProfilSekolahCmd {
		return profil_sekolah_service.UpdateProfilSekolahCmd{
			IDProfil:      1,
			EmailSekolah:  ptrString("admin@example.com"),
			NoTelpSekolah: ptrString("081234567"),
			KepalaSekolah: ptrString("Kepala"),
			WakaSekolah:   ptrString("Waka"),
			NamaSekolah:   ptrString("Sekolah"),
			AlamatSekolah: ptrString("Alamat"),
			LogoSekolah:   nil,
		}
	}

	tests := []struct {
		name           string
		getErr         error
		repoErr        error
		mutateCmd      func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd)
		expectErr      error
		expectGet      bool
		expectUpdate   bool
		expectLogo     *string
		expectNamaTrim string
		expectEmail    string
	}{
		{
			name: "Branch 1 -> id profil tidak valid",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				cmd.IDProfil = 0
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectGet:    false,
			expectUpdate: false,
		},
		{
			name: "Branch 2 -> request tanpa field update",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				cmd.EmailSekolah = nil
				cmd.NoTelpSekolah = nil
				cmd.KepalaSekolah = nil
				cmd.WakaSekolah = nil
				cmd.NamaSekolah = nil
				cmd.AlamatSekolah = nil
			},
			expectErr:    coreerror.ErrNoFieldToUpdate,
			expectGet:    false,
			expectUpdate: false,
		},
		{
			name: "Branch 3 -> email sekolah kosong",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				cmd.EmailSekolah = ptrString(" ")
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectGet:    true,
			expectUpdate: false,
		},
		{
			name: "Branch 4 -> no telp sekolah kosong",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				cmd.NoTelpSekolah = ptrString("")
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectGet:    true,
			expectUpdate: false,
		},
		{
			name: "Branch 5 -> kepala sekolah kosong",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				cmd.KepalaSekolah = ptrString("")
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectGet:    true,
			expectUpdate: false,
		},
		{
			name: "Branch 6 -> waka sekolah kosong",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				cmd.WakaSekolah = ptrString("")
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectGet:    true,
			expectUpdate: false,
		},
		{
			name: "Branch 7 -> nama sekolah kosong",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				cmd.NamaSekolah = ptrString("")
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectGet:    true,
			expectUpdate: false,
		},
		{
			name: "Branch 8 -> alamat sekolah kosong",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				cmd.AlamatSekolah = ptrString("")
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectGet:    true,
			expectUpdate: false,
		},
		{
			name: "Branch 9 -> logo sekolah kosong",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				blankLogo := "  "
				cmd.LogoSekolah = &blankLogo
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectGet:    true,
			expectUpdate: false,
		},
		{
			name:         "Branch 10 -> repo update gagal",
			getErr:       nil,
			repoErr:      errors.New("update failed"),
			expectErr:    errors.New("update failed"),
			expectGet:    true,
			expectUpdate: true,
		},
		{
			name:         "Branch 11 -> repo get gagal",
			getErr:       errors.New("get failed"),
			expectErr:    errors.New("get failed"),
			expectGet:    true,
			expectUpdate: false,
		},
		{
			name: "Branch 12 -> patch sebagian berhasil",
			mutateCmd: func(cmd *profil_sekolah_service.UpdateProfilSekolahCmd) {
				logo := " logo.png "
				cmd.LogoSekolah = &logo
				cmd.NamaSekolah = ptrString(" Sekolah Baru ")
				cmd.EmailSekolah = nil
			},
			expectUpdate:   true,
			expectGet:      true,
			expectNamaTrim: "Sekolah Baru",
			expectLogo:     ptrString("logo.png"),
			expectEmail:    "old@example.com",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeUpdateProfilSekolahRepo{
				updateErr: tc.repoErr,
				getErr:    tc.getErr,
				profil: profil_sekolah.ProfilSekolah{
					IDProfil:      1,
					EmailSekolah:  "old@example.com",
					NoTelpSekolah: "081111111",
					KepalaSekolah: "Kepala Lama",
					WakaSekolah:   "Waka Lama",
					NamaSekolah:   "Nama Lama",
					AlamatSekolah: "Alamat Lama",
					LogoSekolah:   ptrString("old-logo.png"),
				},
			}
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

			assert.Equal(t, tc.expectGet, repo.getCalled)
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
				if tc.expectEmail != "" {
					assert.Equal(t, tc.expectEmail, repo.lastProfil.EmailSekolah)
				}
			}
		})
	}
}

func ptrString(value string) *string {
	return &value
}
