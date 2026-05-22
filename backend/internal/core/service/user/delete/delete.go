package user_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	delete_file_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/delete_file_system"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

type DeleteUserService struct {
	users out.UserRepository
	deleteFile delete_file_repo.DeleteFileRepo
}

func NewDeleteUserService(user out.UserRepository, deleteFile delete_file_repo.DeleteFileRepo) *DeleteUserService {
	return &DeleteUserService{users: user, deleteFile: deleteFile}
}

func (s *DeleteUserService) Delete(ctx context.Context, idPengguna user.ID) error {
	logger := corelog.FromContext(ctx)

	if idPengguna <= 0 {
		logger.Error(ctx, "failed deleting user", "layer", "core.service", "op", "user.delete", "user_id", idPengguna, "err", coreerror.ErrMissingId)
		return coreerror.ErrMissingId
	}

	u,err := s.users.FindUserByID(ctx,idPengguna)
	if err != nil {
		logger.Error(ctx, "failed deleting user", "layer", "core.service", "op", "user.delete", "user_id", idPengguna, "err", err)
		return err
	}
	
	if err := s.deleteFile.DeleteFile(ctx,u.Foto); err != nil {
		logger.Error(ctx, "failed deleting user", "layer", "core.service", "op", "user.delete_file_foto", "user_id", idPengguna, "err", err)
	}

	if err := s.users.DeleteUser(ctx, idPengguna); err != nil {
		logger.Error(ctx, "failed deleting user", "layer", "core.service", "op", "user.delete", "user_id", idPengguna, "err", err)
		return err
	}

	return nil

}

func (s *DeleteUserService) DeleteMany(ctx context.Context, ids []user.ID) (int64, error) {
	logger := corelog.FromContext(ctx)
	affected, err := s.users.DeleteUsers(ctx, ids)
	if err != nil {
		logger.Error(ctx, "failed deleting users", "layer", "core.service", "op", "user.delete_many", "err", err)
		return 0, err
	}

	return affected, nil
}
