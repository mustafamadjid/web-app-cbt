package fake_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

// ===== minimal fake user repo =====

type FakeUserRepo struct {
	UpdateErr error

	UpdateCalled bool
	LastID       user.ID
	LastPatch    outuser.UpdatePenggunaPatch
}

func (r *FakeUserRepo) UpdateUser(ctx context.Context, idPengguna user.ID, pengguna outuser.UpdatePenggunaPatch) error {
	r.UpdateCalled = true
	r.LastID = idPengguna
	r.LastPatch = pengguna
	if r.UpdateErr != nil {
		return r.UpdateErr
	}
	return nil
}

// ---- methods below are NOT used by this test; keep minimal ----

func (r *FakeUserRepo) FindUserByID(ctx context.Context, id user.ID) (user.Pengguna, error) {
	panic("not used in this test")
}

func (r *FakeUserRepo) UserExistByUsername(ctx context.Context, username string) (bool, error) {
	panic("not used in this test")
}

func (r *FakeUserRepo) CreateUser(ctx context.Context, p user.Pengguna) (user.ID, error) {
	panic("not used in this test")
}

func (r *FakeUserRepo) DeleteUser(ctx context.Context, id user.ID) error {
	panic("not used in this test")
}

func (r *FakeUserRepo) DeleteUsers(ctx context.Context, ids []user.ID) (int64, error) {
	panic("not used in this test")
}

func (r *FakeUserRepo) ListUser(ctx context.Context) ([]user.Pengguna, error) {
	panic("not used in this test")
}
