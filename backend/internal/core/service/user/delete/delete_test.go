package user_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/delete"
	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/delete/fake_test"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUserService_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := user.ID(123)

	t.Run("returns error when repository fails", func(t *testing.T) {
		repoErr := errors.New("delete failed")
		repo := &fake_test.FakeDeleteUserRepo{DeleteErr: repoErr}
		service := user_service.NewDeleteUserService(repo)

		err := service.Delete(ctx, userID)

		assert.ErrorIs(t, err, repoErr)
		assert.True(t, repo.DeleteCalled)
		assert.Equal(t, userID, repo.LastID)
	})

	t.Run("deletes user successfully", func(t *testing.T) {
		repo := &fake_test.FakeDeleteUserRepo{}
		service := user_service.NewDeleteUserService(repo)

		err := service.Delete(ctx, userID)

		assert.NoError(t, err)
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
		repo := &fake_test.FakeDeleteUserRepo{DeleteErr: repoErr}
		service := user_service.NewDeleteUserService(repo)

		affected, err := service.DeleteMany(ctx, userIDs)

		assert.ErrorIs(t, err, repoErr)
		assert.EqualValues(t, 0, affected)
		assert.True(t, repo.DeleteUsersCalled)
		assert.Equal(t, userIDs, repo.LastIDs)
	})

	t.Run("deletes users successfully", func(t *testing.T) {
		repo := &fake_test.FakeDeleteUserRepo{}
		service := user_service.NewDeleteUserService(repo)

		affected, err := service.DeleteMany(ctx, userIDs)

		assert.NoError(t, err)
		assert.EqualValues(t, len(userIDs), affected)
		assert.True(t, repo.DeleteUsersCalled)
		assert.Equal(t, userIDs, repo.LastIDs)
	})
}
