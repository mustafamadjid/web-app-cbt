package user_service

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

type DeleteUserService struct {
	users out.UserRepository
}

func NewDeleteUserService(user out.UserRepository) *DeleteUserService {
	return &DeleteUserService{users: user}
}

func (s *DeleteUserService) Delete(ctx context.Context, idPengguna user.ID) error {
	logger := corelog.FromContext(ctx)
	if err := s.users.DeleteUser(ctx, idPengguna); err != nil {
		logger.Error(ctx, "failed deleting user", "op", "user.delete", "user_id", idPengguna, "err", err)
		return err
	}

	return nil

}

func (s *DeleteUserService) DeleteMany(ctx context.Context, ids []user.ID) (int64, error) {
	logger := corelog.FromContext(ctx)
	affected, err := s.users.DeleteUsers(ctx, ids)
	if err != nil {
		logger.Error(ctx, "failed deleting users", "op", "user.delete_many", "err", err)
		return 0, err
	}

	return affected, nil
}
