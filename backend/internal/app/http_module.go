package app

import (
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/cookie"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type HTTPModule struct {
	Handler http.Handler
	Server  *http.Server
}

func BuildHTTPModule(cfg Config, auth *AuthModule, users *UserModule, profilSekolah *ProfilSekolahModule, tokens *TokenModule) *HTTPModule {
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

	requireAccessRole := func(roles ...user.Role) func(next httprouter.Handle) httprouter.Handle {
		return func(next httprouter.Handle) httprouter.Handle {
			return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
				var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next(w, r, ps)
				})
				handler = middleware.RequireActorRole(handler, roles...)
				middleware.RequireValidAccessToken(handler, tokens.AccessTokenSvc, cookies).ServeHTTP(w, r)
			}
		}
	}

	requireAdmin := requireAccessRole(user.ADMIN)
	// requireAdminGuru := requireAccessRole(user.ADMIN, user.GURU)

	// requireSiswa := requireAccessRole(user.SISWA)
	// requireGuru := requireAccessRole(user.GURU)

	router.POST("/auth/login", auth.Handler.Login)
	router.POST("/auth/logout", auth.Handler.Logout)
	router.POST("/auth/refresh", auth.Handler.Refresh)

	// Admin
	router.GET("/admin/siswa", requireAdmin(users.GetSiswaHandler.ListSiswa))
	router.POST("/admin/siswa", requireAdmin(users.CreateHandler.CreateSiswa))
	router.PATCH("/admin/siswa/:id", requireAdmin(users.UpdateHandler.UpdateSiswa))

	router.GET("/admin/guru", requireAdmin(users.GetGuruHandler.ListGuru))
	router.POST("/admin", requireAdmin(users.CreateHandler.CreateGuru))
	router.PATCH("/admin/:id", requireAdmin(users.UpdateHandler.UpdateGuru))

	router.DELETE("admin/pengguna/:id", requireAdmin(users.DeleteHandler.DeleteUser))

	router.GET("/admin/profil-sekolah", requireAdmin(profilSekolah.GetHandler.GetProfilSekolah))
	router.PATCH("/admin/profil-sekolah", requireAdmin(profilSekolah.UpdateHandler.UpdateProfilSekolah))

	// Siswa

	// Guru

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
