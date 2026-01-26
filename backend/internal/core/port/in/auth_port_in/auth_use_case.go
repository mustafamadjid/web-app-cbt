package auth_port_in

import (
	"context"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/login"
)



type AuthUsecase interface {
	Login(ctx context.Context, cmd login.LoginCmd) (login.LoginRes, error)
	Logout(ctx context.Context, refreshtoken string, now time.Time) error
}