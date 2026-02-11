package out

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type UserRepository interface {
	FindUserByID(ctx context.Context, id user.ID) (user.Pengguna, error)
	UserExistByUsername(ctx context.Context, username string) (bool, error)

	CreateUser(ctx context.Context, pengguna user.Pengguna) (user.ID, error)
	UpdateUser(ctx context.Context, idPengguna user.ID, pengguna updatepatch.Pengguna) error
	DeleteUser(ctx context.Context, id user.ID) error
	DeleteUsers(ctx context.Context, ids []user.ID) (int64, error)

	ListUser(ctx context.Context) ([]user.Pengguna, error)
}
