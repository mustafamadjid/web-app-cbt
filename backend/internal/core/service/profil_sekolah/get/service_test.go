package profil_sekolah_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/profil_sekolah"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/profil_sekolah"
	profil_sekolah_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/profil_sekolah/get"
	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/profil_sekolah/get/fake_test"
	"github.com/stretchr/testify/assert"
)

var _ out.ProfilSekolahRepository = (*fake_test.FakeProfilSekolahRepo)(nil)

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
		repo      *fake_test.FakeProfilSekolahRepo
		want      profil_sekolah.ProfilSekolah
		expectErr string
	}{
		{
			name:      "Branch 1 -> repo error",
			repo:      &fake_test.FakeProfilSekolahRepo{GetErr: errors.New("repo error")},
			expectErr: "repo error",
		},
		{
			name: "Branch 2 -> repo success",
			repo: &fake_test.FakeProfilSekolahRepo{
				GetResult: expected,
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
