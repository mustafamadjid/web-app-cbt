package fake_test

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type FakeDeleteUserRepo struct {
	DeleteErr error

	DeleteCalled bool
	LastID       user.ID

	DeleteUsersCalled bool
	LastIDs           []user.ID
}

func (r *FakeDeleteUserRepo) DeleteUser(ctx context.Context, id user.ID) error {
	r.DeleteCalled = true
	r.LastID = id
	if r.DeleteErr != nil {
		return r.DeleteErr
	}
	return nil
}

func (r *FakeDeleteUserRepo) DeleteUsers(ctx context.Context, ids []user.ID) (int64, error) {
	r.DeleteUsersCalled = true
	r.LastIDs = ids
	if r.DeleteErr != nil {
		return 0, r.DeleteErr
	}
	return int64(len(ids)), nil
}

func (r *FakeDeleteUserRepo) FindUserByID(ctx context.Context, id user.ID) (user.Pengguna, error) {
	panic("not used in this test")
}

func (r *FakeDeleteUserRepo) UserExistByUsername(ctx context.Context, username string) (bool, error) {
	panic("not used in this test")
}

func (r *FakeDeleteUserRepo) CreateUser(ctx context.Context, pengguna user.Pengguna) (user.ID, error) {
	panic("not used in this test")
}

func (r *FakeDeleteUserRepo) UpdateUser(ctx context.Context, idPengguna user.ID, pengguna updatepatch.Pengguna) error {
	panic("not used in this test")
}

func (r *FakeDeleteUserRepo) ListUser(ctx context.Context) ([]user.Pengguna, error) {
	panic("not used in this test")
}
