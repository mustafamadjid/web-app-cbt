package fake_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

// ===== minimal fake profil guru repo =====

type FakeProfilGuruRepo struct {
	UpdateErr error

	UpdateCalled bool
	LastID       user.ID
	LastPatch    outuser.UpdateProfilGuruPatch
}

func (r *FakeProfilGuruRepo) UpdateProfilGuru(ctx context.Context, idPengguna user.ID, profilGuru outuser.UpdateProfilGuruPatch) error {
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
