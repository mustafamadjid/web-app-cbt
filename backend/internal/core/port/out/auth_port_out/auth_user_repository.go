package auth_port_out

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type AuthUserrepository interface {
	FindByUsername(ctx context.Context, username string) (user.Pengguna, error)
}