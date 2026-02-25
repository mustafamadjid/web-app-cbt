package fake_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

// ===== minimal fake user repo =====

type FakeUserRepo struct {
	UpdateErr error
	FindErr   error

	UpdateCalled bool
	LastID       user.ID
	LastPatch    updatepatch.Pengguna

	FindCalled bool
	LastFindID user.ID
	FindResult user.Pengguna
}

func (r *FakeUserRepo) UpdateUser(ctx context.Context, idPengguna user.ID, pengguna updatepatch.Pengguna) error {
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
	r.FindCalled = true
	r.LastFindID = id
	if r.FindErr != nil {
		return user.Pengguna{}, r.FindErr
	}
	if r.FindResult.ID == 0 {
		r.FindResult = user.Pengguna{ID: id, Foto: "foto-lama.png"}
	}
	return r.FindResult, nil
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
