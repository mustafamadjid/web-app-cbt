package user_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/delete"
	"github.com/stretchr/testify/assert"
)

type fakeDeleteFileRepo struct {
	DeleteErr error

	DeleteCalled bool
	LastPath     string
}

func (f *fakeDeleteFileRepo) DeleteFile(ctx context.Context, filePath string) error {
	f.DeleteCalled = true
	f.LastPath = filePath
	if f.DeleteErr != nil {
		return f.DeleteErr
	}
	return nil
}

func TestDeleteUserService_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := user.ID(123)

	t.Run("returns error when id is invalid", func(t *testing.T) {
		repo := &FakeDeleteUserRepo{}
		deleteFile := &fakeDeleteFileRepo{}
		service := user_service.NewDeleteUserService(repo, deleteFile)

		err := service.Delete(ctx, 0)

		assert.ErrorIs(t, err, coreerror.ErrMissingId)
		assert.False(t, repo.FindCalled)
		assert.False(t, deleteFile.DeleteCalled)
		assert.False(t, repo.DeleteCalled)
	})

	t.Run("returns error when find user fails", func(t *testing.T) {
		findErr := errors.New("find failed")
		repo := &FakeDeleteUserRepo{FindErr: findErr}
		deleteFile := &fakeDeleteFileRepo{}
		service := user_service.NewDeleteUserService(repo, deleteFile)

		err := service.Delete(ctx, userID)

		assert.ErrorIs(t, err, findErr)
		assert.True(t, repo.FindCalled)
		assert.Equal(t, userID, repo.LastFindID)
		assert.False(t, deleteFile.DeleteCalled)
		assert.False(t, repo.DeleteCalled)
	})

	t.Run("returns error when delete file fails", func(t *testing.T) {
		deleteFileErr := errors.New("delete file failed")
		repo := &FakeDeleteUserRepo{
			FindResult: user.Pengguna{ID: userID, Foto: "foto-lama.png"},
		}
		deleteFile := &fakeDeleteFileRepo{DeleteErr: deleteFileErr}
		service := user_service.NewDeleteUserService(repo, deleteFile)

		err := service.Delete(ctx, userID)

		assert.ErrorIs(t, err, deleteFileErr)
		assert.True(t, repo.FindCalled)
		assert.Equal(t, userID, repo.LastFindID)
		assert.True(t, deleteFile.DeleteCalled)
		assert.Equal(t, "foto-lama.png", deleteFile.LastPath)
		assert.False(t, repo.DeleteCalled)
	})

	t.Run("returns error when repository delete fails", func(t *testing.T) {
		repoErr := errors.New("delete failed")
		repo := &FakeDeleteUserRepo{
			DeleteErr:  repoErr,
			FindResult: user.Pengguna{ID: userID, Foto: "foto-lama.png"},
		}
		deleteFile := &fakeDeleteFileRepo{}
		service := user_service.NewDeleteUserService(repo, deleteFile)

		err := service.Delete(ctx, userID)

		assert.ErrorIs(t, err, repoErr)
		assert.True(t, repo.FindCalled)
		assert.Equal(t, userID, repo.LastFindID)
		assert.True(t, deleteFile.DeleteCalled)
		assert.Equal(t, "foto-lama.png", deleteFile.LastPath)
		assert.True(t, repo.DeleteCalled)
		assert.Equal(t, userID, repo.LastID)
	})

	t.Run("deletes user successfully", func(t *testing.T) {
		repo := &FakeDeleteUserRepo{
			FindResult: user.Pengguna{ID: userID, Foto: "foto-lama.png"},
		}
		deleteFile := &fakeDeleteFileRepo{}
		service := user_service.NewDeleteUserService(repo, deleteFile)

		err := service.Delete(ctx, userID)

		assert.NoError(t, err)
		assert.True(t, repo.FindCalled)
		assert.Equal(t, userID, repo.LastFindID)
		assert.True(t, deleteFile.DeleteCalled)
		assert.Equal(t, "foto-lama.png", deleteFile.LastPath)
		assert.True(t, repo.DeleteCalled)
		assert.Equal(t, userID, repo.LastID)
	})
}

func TestDeleteUserService_DeleteMany(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userIDs := []user.ID{101, 202, 303}

	t.Run("returns error when repository fails", func(t *testing.T) {
		repoErr := errors.New("delete many failed")
		repo := &FakeDeleteUserRepo{DeleteErr: repoErr}
		deleteFile := &fakeDeleteFileRepo{}
		service := user_service.NewDeleteUserService(repo, deleteFile)

		affected, err := service.DeleteMany(ctx, userIDs)

		assert.ErrorIs(t, err, repoErr)
		assert.EqualValues(t, 0, affected)
		assert.True(t, repo.DeleteUsersCalled)
		assert.Equal(t, userIDs, repo.LastIDs)
		assert.False(t, deleteFile.DeleteCalled)
	})

	t.Run("deletes users successfully", func(t *testing.T) {
		repo := &FakeDeleteUserRepo{}
		deleteFile := &fakeDeleteFileRepo{}
		service := user_service.NewDeleteUserService(repo, deleteFile)

		affected, err := service.DeleteMany(ctx, userIDs)

		assert.NoError(t, err)
		assert.EqualValues(t, len(userIDs), affected)
		assert.True(t, repo.DeleteUsersCalled)
		assert.Equal(t, userIDs, repo.LastIDs)
		assert.False(t, deleteFile.DeleteCalled)
	})
}
