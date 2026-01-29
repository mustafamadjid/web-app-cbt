package auth_port_in

import (
	"context"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/service/auth_service"
)



type AuthUsecase interface {
	Login(ctx context.Context, cmd auth_service.LoginCmd) (auth_service.LoginRes, error)
	Logout(ctx context.Context, refreshtoken string, now time.Time) error
	RefreshAccessToken(ctx context.Context, refreshToken string, accessTTL time.Duration) (string,error)
}