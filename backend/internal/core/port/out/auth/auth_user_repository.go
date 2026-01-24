package auth

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type AuthUserrepository interface {
	FindByUsername(ctx context.Context, username string) (user.Pengguna, error)
}