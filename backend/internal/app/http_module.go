package app

import (
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/cookie"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
)

type HTTPModule struct {
	Handler http.Handler
	Server  *http.Server
}

func BuildHTTPModule(cfg Config, auth *AuthModule, users *UserModule, tokens *TokenModule) *HTTPModule {
	cookies := cookie.CookieConfig{
		AccessName:  cfg.Cookie.AccessName,
		RefreshName: cfg.Cookie.RefreshName,
		Domain:      cfg.Cookie.Domain,
		Secure:      cfg.Cookie.Secure,
		SameSite:    cfg.Cookie.SameSite,
	}

	router := httprouter.New()

	requireAccess := func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next(w, r, ps)
			})
			middleware.RequireValidAccessToken(handler, tokens.AccessTokenSvc, cookies).ServeHTTP(w, r)
		}
	}

	router.POST("/auth/login", auth.Handler.Login)
	router.POST("/auth/logout", auth.Handler.Logout)
	router.POST("/auth/refresh", auth.Handler.Refresh)

	// Siswa
	router.GET("/siswa", requireAccess(users.GetSiswaHandler.ListSiswa))
	router.POST("/siswa", requireAccess(users.CreateHandler.CreateSiswa))
	router.PATCH("/siswa/:id", requireAccess(users.UpdateHandler.UpdateSiswa))

	// Guru
	router.GET("/guru", requireAccess(users.GetGuruHandler.ListGuru))
	router.POST("/guru", requireAccess(users.CreateHandler.CreateGuru))
	router.PATCH("/guru/:id", requireAccess(users.UpdateHandler.UpdateGuru))
	router.DELETE("/pengguna/:id", requireAccess(users.DeleteHandler.DeleteUser))

	router.GET("/uploads/*filepath", requireAccess(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		rel := ps.ByName("filepath")

		if rel == "/" || rel == "" {
			http.NotFound(w, r)
			return
		}
		clean := path.Clean(rel)
		if strings.Contains(clean, "..") {
			http.Error(w, "bad path", http.StatusBadRequest)
			return
		}

		full := filepath.Join("./public/uploads", clean)
		http.ServeFile(w, r, full)
	}))

	server := &http.Server{
		Addr:    cfg.HTTP.Addr,
		Handler: router,
	}

	return &HTTPModule{
		Handler: router,
		Server:  server,
	}
}
