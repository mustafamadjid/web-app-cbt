package app

import (
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/token"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
)

type TokenModule struct {
	AccessTokenSvc out.AccessTokenService
	RefreshTokenSvc out.RefreshTokenService
}

func BuildTokenModule(cfg Config) *TokenModule {
	return &TokenModule{
		AccessTokenSvc: token.NewJWTAccessTokenService(cfg.JWT.AccessSecret, cfg.JWT.Issuer),
		RefreshTokenSvc: token.NewJWTRefreshTokenService(cfg.JWT.RefreshSecret,cfg.JWT.Issuer),
	}
}