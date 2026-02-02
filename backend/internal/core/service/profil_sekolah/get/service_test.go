package profil_sekolah_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/profil_sekolah"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/profil_sekolah"
	profil_sekolah_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/profil_sekolah/get"
	"github.com/stretchr/testify/assert"
)

type fakeProfilSekolahRepo struct {
	getResult profil_sekolah.ProfilSekolah
	getErr    error

	updateCalled bool
	updateErr    error
}

func (f *fakeProfilSekolahRepo) UpdateProfilSekolah(ctx context.Context, idProfil profil_sekolah.IDProfil, profil profil_sekolah.ProfilSekolah) error {
	f.updateCalled = true
	if f.updateErr != nil {
		return f.updateErr
	}
	return nil
}

func (f *fakeProfilSekolahRepo) GetProfilSekolah(ctx context.Context) (profil_sekolah.ProfilSekolah, error) {
	if f.getErr != nil {
		return profil_sekolah.ProfilSekolah{}, f.getErr
	}
	return f.getResult, nil
}

var _ out.ProfilSekolahRepository = (*fakeProfilSekolahRepo)(nil)

func TestGetProfilSekolahService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expected := profil_sekolah.ProfilSekolah{
		IDProfil:      1,
		EmailSekolah:  "admin@example.com",
		NoTelpSekolah: "081234",
		KepalaSekolah: "Kepala",
		WakaSekolah:   "Waka",
		NamaSekolah:   "Sekolah",
		AlamatSekolah: "Alamat",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	tests := []struct {
		name      string
		repo      *fakeProfilSekolahRepo
		want      profil_sekolah.ProfilSekolah
		expectErr string
	}{
		{
			name:      "Branch 1 -> repo error",
			repo:      &fakeProfilSekolahRepo{getErr: errors.New("repo error")},
			expectErr: "repo error",
		},
		{
			name: "Branch 2 -> repo success",
			repo: &fakeProfilSekolahRepo{
				getResult: expected,
			},
			want: expected,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			svc := profil_sekolah_service.NewGetProfilSekolahService(tc.repo)

			result, err := svc.GetProfilSekolah(ctx)

			if tc.expectErr != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want, result)
		})
	}
}
