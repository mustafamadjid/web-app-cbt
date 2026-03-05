package httpmodule

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/rate_limit"
	"github.com/mustafamadjid/web-app-cbt/internal/app/http/routes"
	"golang.org/x/time/rate"
)

type HTTPModule struct {
	Handler http.Handler
	Server  *http.Server
}

func BuildHTTPModule(deps HTTPDeps) *HTTPModule {
	router := httprouter.New()

	standardLimiter := rate_limit.NewMemoryTokenBucket(rate.Limit(10), 15, 5*time.Minute)
	authLimiter := rate_limit.NewMemoryTokenBucket(rate.Limit(1), 3, 5*time.Minute)

	mw := NewMiddlewares(
		deps.Logger,
		deps.AccessTokenSvc,
		deps.RefreshTokenSvc,
		deps.Sessions,
		deps.CookieConfig,
		standardLimiter,
		authLimiter,
	)

	routes.RegisterAuthRoutes(router, deps.Auth, mw)
	routes.RegisterUserRoutes(router, deps.Users, deps.ResetPassword, mw)
	routes.RegisterAktivitasUserRoutes(router, deps.AktivitasUser, mw)
	routes.RegisterProfilSekolahRoutes(router, deps.ProfilSekolah, mw)
	routes.RegisterKelasRoutes(router, deps.Kelas, mw)
	routes.RegisterMataPelajaranRoutes(router, deps.MataPelajaran, mw)
	routes.RegisterRuangUjianRoutes(router, deps.RuangUjian, mw)
	routes.RegisterUjianRoutes(router, deps.Ujian, mw)
	routes.RegisterSesiRoutes(router, deps.Sesi, mw)
	routes.RegisterPengumumanRoutes(router, deps.Pengumuman, mw)
	routes.RegisterImportSoalRoutes(router, deps.ImportSoal, mw)
	routes.RegisterBankSoalRoutes(router, deps.BankSoal, mw)

	documentPengumumanRoute := strings.TrimRight(deps.DocumentStoreRoute, "/") + "/pengumuman/*filepath"
	imageUploadRoute := strings.TrimRight(deps.ImageStoreRoute, "/") + "/*filepath"

	documentPengumumanHandler := ServeStaticFrom(filepath.Join(deps.DocumentStoreDir, "pengumuman"))
	imageUploadHandler := ServeStaticFrom(deps.ImageStoreDir)

	routes.RegisterStaticRoutes(
		router,
		documentPengumumanRoute,
		imageUploadRoute,
		documentPengumumanHandler,
		imageUploadHandler,
		mw,
	)

	corsHandler := middleware.CORSPolicy(router)
	protectCSRFHandler := middleware.PreventCSRF(corsHandler)

	server := &http.Server{
		Addr:    deps.HTTPAddr,
		Handler: protectCSRFHandler,
	}

	return &HTTPModule{
		Handler: corsHandler,
		Server:  server,
	}
}
