package app

import (
	httpx "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/auth"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/cookie"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/auth_service"
)

type AuthModule struct {
	Service *auth_service.AuthService
	Handler *httpx.AuthHandler
}


func BuildAuthModule(cfg Config, infra *InfraModule, tokens *TokenModule, hasher out.PasswordHasher) *AuthModule {
	svc := auth_service.NewAuthService(
		infra.AuthUsers,
		infra.users,
		hasher,
		infra.Sessions,
		tokens.AccessTokenSvc,
		tokens.RefreshTokenSvc,
	)

	cookies := cookie.CookieConfig{
		AccessName:  cfg.Cookie.AccessName,
		RefreshName: cfg.Cookie.RefreshName,
		Domain:      cfg.Cookie.Domain,
		Secure:      cfg.Cookie.Secure,
		SameSite:    cfg.Cookie.SameSite,
	}

	h := httpx.NewAuthHandler(svc, cookies, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)

	return &AuthModule{Service: svc, Handler: h}
}