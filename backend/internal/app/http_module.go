package app

import (
	"net/http"
	"github.com/julienschmidt/httprouter"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/cookie"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
)

type HTTPModule struct {
	Handler http.Handler
	Server *http.Server
}

func BuildHTTPModule(cfg Config, auth *AuthModule, users *UserModule,tokens *TokenModule) *HTTPModule {
	cookies := cookie.CookieConfig{
		AccessName:  cfg.Cookie.AccessName,
		RefreshName: cfg.Cookie.RefreshName,
		Domain:      cfg.Cookie.Domain,
		Secure:      cfg.Cookie.Secure,
		SameSite:    cfg.Cookie.SameSite,
	}

	protected := middleware.RequireValidAccessToken(http.NewServeMux(),tokens.AccessTokenSvc,cookies)

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/guru/register",users.CreateHandler.CreateGuru)
}
	