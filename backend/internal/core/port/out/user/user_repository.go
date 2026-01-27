package out

import (
	"context"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)



type UserRepository interface {
	FindUserByID(ctx context.Context,id user.ID) (user.Pengguna,error)
	UserExistByUsername(ctx context.Context,username string) (bool,error)

	CreateUser(ctx context.Context,pengguna user.Pengguna) (user.ID,error)
	UpdateUser(ctx context.Context,pengguna user.Pengguna) (user.ID,error)
	DeleteUser(ctx context.Context,id user.ID) error

	ListUser(ctx context.Context) ([]user.Pengguna,error)
}



