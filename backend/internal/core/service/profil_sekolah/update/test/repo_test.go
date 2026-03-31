package profil_sekolah_service_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/profil_sekolah"
)

type FakeUpdateProfilSekolahRepo struct {
	GetErr error
	Profil profil_sekolah.ProfilSekolah

	UpdateErr error

	GetCalled    bool
	UpdateCalled bool
	LastID       profil_sekolah.IDProfil
	LastProfil   profil_sekolah.ProfilSekolah
}

func (f *FakeUpdateProfilSekolahRepo) UpdateProfilSekolah(ctx context.Context, idProfil profil_sekolah.IDProfil, profil profil_sekolah.ProfilSekolah) error {
	f.UpdateCalled = true
	f.LastID = idProfil
	f.LastProfil = profil
	if f.UpdateErr != nil {
		return f.UpdateErr
	}
	return nil
}

func (f *FakeUpdateProfilSekolahRepo) GetProfilSekolah(ctx context.Context) (profil_sekolah.ProfilSekolah, error) {
	f.GetCalled = true
	if f.GetErr != nil {
		return profil_sekolah.ProfilSekolah{}, f.GetErr
	}
	return f.Profil, nil
}
