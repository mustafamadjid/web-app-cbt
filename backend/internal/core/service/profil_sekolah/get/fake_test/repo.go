package fake_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/profil_sekolah"
)

type FakeProfilSekolahRepo struct {
	GetResult profil_sekolah.ProfilSekolah
	GetErr    error

	UpdateCalled bool
	UpdateErr    error
}

func (f *FakeProfilSekolahRepo) UpdateProfilSekolah(ctx context.Context, idProfil profil_sekolah.IDProfil, profil profil_sekolah.ProfilSekolah) error {
	f.UpdateCalled = true
	if f.UpdateErr != nil {
		return f.UpdateErr
	}
	return nil
}

func (f *FakeProfilSekolahRepo) GetProfilSekolah(ctx context.Context) (profil_sekolah.ProfilSekolah, error) {
	if f.GetErr != nil {
		return profil_sekolah.ProfilSekolah{}, f.GetErr
	}
	return f.GetResult, nil
}
