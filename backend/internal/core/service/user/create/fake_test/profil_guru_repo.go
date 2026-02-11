package fake_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

// ===== minimal fake profil guru repo =====

type FakeProfilGuruRepo struct {
	ExistsNip bool
	ExistErr  error

	CreateID  user.ID
	CreateErr error

	ExistCalled  bool
	CreateCalled bool
}

func (r *FakeProfilGuruRepo) ExistByNIP(ctx context.Context, nip user.NIP) (bool, error) {
	r.ExistCalled = true
	if r.ExistErr != nil {
		return false, r.ExistErr
	}
	return r.ExistsNip, nil
}

func (r *FakeProfilGuruRepo) CreateProfilGuru(ctx context.Context, g user.ProfilGuru) (user.ID, error) {
	r.CreateCalled = true
	if r.CreateErr != nil {
		return 0, r.CreateErr
	}
	return r.CreateID, nil
}

func (r *FakeProfilGuruRepo) FindProfilGuruByID(ctx context.Context, id user.ID) (user.DataGuru, error) {
	panic("not used in this test")
}

func (r *FakeProfilGuruRepo) UpdateProfilGuru(ctx context.Context, idPengguna user.ID, profilGuru updatepatch.ProfilGuru) error {
	panic("not used in this test")
}
