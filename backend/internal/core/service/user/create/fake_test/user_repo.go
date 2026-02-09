package fake_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

// ===== minimal fake user repo =====

type FakeUserRepo struct {
	ExistByUsername bool
	UserExistErr    error

	CreateID  user.ID
	CreateErr error

	UserExistCalled bool
	CreateCalled    bool
}

func (r *FakeUserRepo) UserExistByUsername(ctx context.Context, username string) (bool, error) {
	r.UserExistCalled = true
	if r.UserExistErr != nil {
		return false, r.UserExistErr
	}
	return r.ExistByUsername, nil
}

func (r *FakeUserRepo) CreateUser(ctx context.Context, p user.Pengguna) (user.ID, error) {
	r.CreateCalled = true
	if r.CreateErr != nil {
		return 0, r.CreateErr
	}
	return r.CreateID, nil
}

// ---- methods below are NOT used by this test; keep minimal ----

func (r *FakeUserRepo) FindUserByID(ctx context.Context, id user.ID) (user.Pengguna, error) {
	panic("not used in this test")
}

func (r *FakeUserRepo) UpdateUser(ctx context.Context, idPengguna user.ID, pengguna outuser.UpdatePenggunaPatch) error {
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
