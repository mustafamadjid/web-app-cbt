package out

import (
	"context"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type UserRepository interface {
	FindByUsername(ctx context.Context,username string) (user.Pengguna,error) 
	FindByEmail(ctx context.Context,email string) (user.Pengguna,error)
	ExistByUsername(ctx context.Context,username string) (bool,error)

	Create(ctx context.Context,pengguna user.Pengguna) (user.ID,error)
}

