package fake_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

// ===== minimal fake profil siswa repo =====

type FakeProfilSiswaRepo struct {
	ExistsNisn   bool
	ExistNisnErr error

	CreateID  user.ID
	CreateErr error

	ExistCalled  bool
	CreateCalled bool
}

func (r *FakeProfilSiswaRepo) ExistByNISN(ctx context.Context, nisn string) (bool, error) {
	r.ExistCalled = true
	if r.ExistNisnErr != nil {
		return false, r.ExistNisnErr
	}
	return r.ExistsNisn, nil
}

func (r *FakeProfilSiswaRepo) CreateProfilSiswa(ctx context.Context, g user.ProfilSiswa) (user.ID, error) {
	r.CreateCalled = true
	if r.CreateErr != nil {
		return 0, r.CreateErr
	}
	return r.CreateID, nil
}

func (r *FakeProfilSiswaRepo) FindProfilSiswaByID(ctx context.Context, id user.ID) (user.DataSiswa, error) {
	panic("not used in this test")
}

func (r *FakeProfilSiswaRepo) UpdateProfilSiswa(ctx context.Context, idPengguna user.ID, profilSiswa updatepatch.ProfilSiswa) error {
	panic("not used in this test")
}
