package fake_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

// ===== minimal fake profil siswa repo =====

type FakeProfilSiswaRepo struct {
	UpdateErr error

	UpdateCalled bool
	LastID       user.ID
	LastPatch    updatepatch.ProfilSiswa
}

func (r *FakeProfilSiswaRepo) UpdateProfilSiswa(ctx context.Context, idPengguna user.ID, profilSiswa updatepatch.ProfilSiswa) error {
	r.UpdateCalled = true
	r.LastID = idPengguna
	r.LastPatch = profilSiswa
	if r.UpdateErr != nil {
		return r.UpdateErr
	}
	return nil
}

func (r *FakeProfilSiswaRepo) FindProfilSiswaByID(ctx context.Context, id user.ID) (user.DataSiswa, error) {
	panic("not used in this test")
}

func (r *FakeProfilSiswaRepo) ExistByNISN(ctx context.Context, nisn string) (bool, error) {
	panic("not used in this test")
}

func (r *FakeProfilSiswaRepo) CreateProfilSiswa(ctx context.Context, g user.ProfilSiswa) (user.ID, error) {
	panic("not used in this test")
}
