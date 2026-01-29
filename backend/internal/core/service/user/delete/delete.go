package user_service

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

type DeleteUserService struct {
	users out.UserRepository
}


func NewDeleteUserService(user out.UserRepository) *DeleteUserService {
	return &DeleteUserService{users: user}
}

func (s *DeleteUserService)Delete(ctx context.Context, idPengguna user.ID) error {
	
	if err := s.users.DeleteUser(ctx, idPengguna); err != nil {
		return err
	}

	return nil

}