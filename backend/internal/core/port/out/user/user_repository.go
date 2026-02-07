package out

import (
	"context"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type UpdatePenggunaPatch struct {
	NamaLengkap *string
	Email       *user.Email
	NoHp        *string
	Foto        *string
	StatusAkun  *user.StatusAkun
	Role        *user.Role
}

type UserRepository interface {
	FindUserByID(ctx context.Context, id user.ID) (user.Pengguna, error)
	UserExistByUsername(ctx context.Context, username string) (bool, error)

	CreateUser(ctx context.Context, pengguna user.Pengguna) (user.ID, error)
	UpdateUser(ctx context.Context, idPengguna user.ID, pengguna UpdatePenggunaPatch) error
	DeleteUser(ctx context.Context, id user.ID) error
	DeleteUsers(ctx context.Context, ids []user.ID) (int64, error)

	ListUser(ctx context.Context) ([]user.Pengguna, error)
}
