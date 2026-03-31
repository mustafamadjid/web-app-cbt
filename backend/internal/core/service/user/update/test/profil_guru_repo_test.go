package user_service_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

// ===== minimal fake profil guru repo =====

type FakeProfilGuruRepo struct {
	UpdateErr error

	UpdateCalled bool
	LastID       user.ID
	LastPatch    updatepatch.ProfilGuru
}

func (r *FakeProfilGuruRepo) UpdateProfilGuru(ctx context.Context, idPengguna user.ID, profilGuru updatepatch.ProfilGuru) error {
	r.UpdateCalled = true
	r.LastID = idPengguna
	r.LastPatch = profilGuru
	if r.UpdateErr != nil {
		return r.UpdateErr
	}
	return nil
}

func (r *FakeProfilGuruRepo) FindProfilGuruByID(ctx context.Context, id user.ID) (user.DataGuru, error) {
	panic("not used in this test")
}

func (r *FakeProfilGuruRepo) ExistByNIP(ctx context.Context, nip user.NIP) (bool, error) {
	panic("not used in this test")
}

func (r *FakeProfilGuruRepo) CreateProfilGuru(ctx context.Context, g user.ProfilGuru) (user.ID, error) {
	panic("not used in this test")
}
