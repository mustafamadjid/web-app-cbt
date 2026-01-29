package user_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/delete"
	"github.com/stretchr/testify/assert"
)

type fakeDeleteUserRepo struct {
	deleteErr error

	deleteCalled bool
	lastID       user.ID
}

func (r *fakeDeleteUserRepo) DeleteUser(ctx context.Context, id user.ID) error {
	r.deleteCalled = true
	r.lastID = id
	if r.deleteErr != nil {
		return r.deleteErr
	}
	return nil
}

func (r *fakeDeleteUserRepo) FindUserByID(ctx context.Context, id user.ID) (user.Pengguna, error) {
	panic("not used in this test")
}

func (r *fakeDeleteUserRepo) UserExistByUsername(ctx context.Context, username string) (bool, error) {
	panic("not used in this test")
}

func (r *fakeDeleteUserRepo) CreateUser(ctx context.Context, pengguna user.Pengguna) (user.ID, error) {
	panic("not used in this test")
}

func (r *fakeDeleteUserRepo) UpdateUser(ctx context.Context, idPengguna user.ID, pengguna outuser.UpdatePenggunaPatch) error {
	panic("not used in this test")
}

func (r *fakeDeleteUserRepo) ListUser(ctx context.Context) ([]user.Pengguna, error) {
	panic("not used in this test")
}

func TestDeleteUserService_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := user.ID(123)

	t.Run("returns error when repository fails", func(t *testing.T) {
		repoErr := errors.New("delete failed")
		repo := &fakeDeleteUserRepo{deleteErr: repoErr}
		service := user_service.NewDeleteUserService(repo)

		err := service.Delete(ctx, userID)

		assert.ErrorIs(t, err, repoErr)
		assert.True(t, repo.deleteCalled)
		assert.Equal(t, userID, repo.lastID)
	})

	t.Run("deletes user successfully", func(t *testing.T) {
		repo := &fakeDeleteUserRepo{}
		service := user_service.NewDeleteUserService(repo)

		err := service.Delete(ctx, userID)

		assert.NoError(t, err)
		assert.True(t, repo.deleteCalled)
		assert.Equal(t, userID, repo.lastID)
	})
}
