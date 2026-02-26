package app

import (
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/time/rate"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/cookie"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/rate_limit"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type HTTPModule struct {
	Handler http.Handler
	Server  *http.Server
}

func BuildHTTPModule(cfg Config, auth *AuthModule, users *UserModule, profilSekolah *ProfilSekolahModule, aktivitasUser *AktivitasUserModule, kelas *KelasModule, mapel *MataPelajaranModule, ruangUjian *RuangUjianModule, sesi *SesiModule, pengumuman *PengumumanModule, resetPassword *ResetPasswordModule, importSoal *ImportSoalModule, tokens *TokenModule, infra *InfraModule, logger corelog.Logger) *HTTPModule {
	cookies := cookie.CookieConfig{
		AccessName:  cfg.Cookie.AccessName,
		RefreshName: cfg.Cookie.RefreshName,
		Domain:      cfg.Cookie.Domain,
		Secure:      cfg.Cookie.Secure,
		SameSite:    cfg.Cookie.SameSite,
	}

	router := httprouter.New()

	// Rate limiter
	standardLimiter := rate_limit.NewMemoryTokenBucket(rate.Limit(10), 15, 5*time.Minute)
	authLimiter := rate_limit.NewMemoryTokenBucket(rate.Limit(1), 3, 5*time.Minute)

	withRequestLogger := func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next(w, r, ps)
			})
			middleware.RequestLogger(handler, logger).ServeHTTP(w, r)
		}
	}

	requireAccess := func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next(w, r, ps)
			})
			handler = middleware.RequestLogger(handler, logger)
			middleware.RequireValidTokenAndSession(handler, tokens.AccessTokenSvc, tokens.RefreshTokenSvc, infra.Sessions, cookies).ServeHTTP(w, r)
		}
	}

	requireAccessRole := func(roles ...user.Role) func(next httprouter.Handle) httprouter.Handle {
		return func(next httprouter.Handle) httprouter.Handle {
			return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
				var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next(w, r, ps)
				})
				handler = middleware.RequireActorRole(handler, roles...)
				handler = middleware.RequestLogger(handler, logger)
				middleware.RequireValidTokenAndSession(handler, tokens.AccessTokenSvc, tokens.RefreshTokenSvc, infra.Sessions, cookies).ServeHTTP(w, r)
			}
		}
	}

	rateLimiterAuth := func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			// adapt httprouter.Handle -> http.Handler
			var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next(w, r, ps)
			})

			handler = middleware.LoginRateLimit(authLimiter, handler)
			handler.ServeHTTP(w, r)
		}
	}
	rateLimitStandard := func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next(w, r, ps)
			})
			handler = middleware.StandardRateLimit(standardLimiter, handler)
			handler.ServeHTTP(w, r)
		}
	}

	requireAdmin := requireAccessRole(user.ADMIN)
	requireAdminGuru := requireAccessRole(user.ADMIN, user.GURU)
	documentPengumumanRoute := strings.TrimRight(cfg.DocumentStore.Route, "/") + "/pengumuman/*filepath"
	imageUploadRoute := strings.TrimRight(cfg.ImageStore.Route, "/") + "/*filepath"

	// requireSiswa := requireAccessRole(user.SISWA)
	// requireGuru := requireAccessRole(user.GURU)

	router.POST("/auth/login", rateLimiterAuth(withRequestLogger(auth.Handler.Login)))
	router.POST("/auth/logout", withRequestLogger(rateLimitStandard(auth.Handler.Logout)))
	router.POST("/auth/refresh", withRequestLogger(rateLimitStandard(auth.Handler.Refresh)))
	router.GET("/auth/me", requireAccess(rateLimitStandard(auth.Handler.AuthMe)))

	// SISWA
	router.GET("/admin/siswa", requireAdmin(rateLimitStandard(users.GetSiswaHandler.ListSiswa)))
	router.GET("/admin/siswa/:id", requireAdmin(rateLimitStandard(users.GetSiswaHandler.GetSiswaByID)))
	router.POST("/admin/siswa", requireAdmin(rateLimitStandard(users.CreateHandler.CreateSiswa)))
	router.PATCH("/admin/siswa/:id", requireAdmin(rateLimitStandard(users.UpdateHandler.UpdateSiswa)))

	// GURU
	router.GET("/admin/guru", requireAdmin(rateLimitStandard(users.GetGuruHandler.ListGuru)))
	router.GET("/admin/guru/:id", requireAdmin(rateLimitStandard(users.GetGuruHandler.GetGuruByID)))
	router.POST("/admin/guru", requireAdmin(rateLimitStandard(users.CreateHandler.CreateGuru)))
	router.PATCH("/admin/guru/:id", requireAdmin(rateLimitStandard(users.UpdateHandler.UpdateGuru)))

	// PENGGUNA
	router.DELETE("/admin/pengguna", requireAdmin(rateLimitStandard(users.DeleteHandler.DeleteUsers)))
	router.DELETE("/admin/pengguna/:id", requireAdmin(rateLimitStandard(users.DeleteHandler.DeleteUser)))

	// RESET PASSWORD
	router.PUT("/admin/pengguna/:idPengguna/reset-password", requireAdmin(rateLimitStandard(resetPassword.Handler.ResetPasswordHandler)))

	// AKTIVITAS USER
	router.GET("/admin/aktivitas-user", requireAdmin(rateLimitStandard(aktivitasUser.GetHandler.GetAktivitasUser)))

	// PROFIL SEKOLAH
	router.GET("/admin/profil-sekolah", requireAdmin(rateLimitStandard(profilSekolah.GetHandler.GetProfilSekolah)))
	router.PATCH("/admin/profil-sekolah", requireAdmin(rateLimitStandard(profilSekolah.UpdateHandler.UpdateProfilSekolah)))

	// KELAS
	router.GET("/admin/kelas", requireAdmin(rateLimitStandard(kelas.GetHandler.ListKelas)))
	router.GET("/admin/kelas/:idTingkatKelas/:idNamaKelas", requireAdmin(rateLimitStandard(kelas.GetHandler.GetKelasByID)))
	router.POST("/admin/kelas/tingkat-kelas", requireAdmin(rateLimitStandard(kelas.CreateHandler.CreateTingkatKelas)))
	router.POST("/admin/kelas/nama-kelas", requireAdmin(rateLimitStandard(kelas.CreateHandler.CreateNamaKelas)))
	router.PATCH("/admin/kelas/nama-kelas/:idNamaKelas", requireAdmin(rateLimitStandard(kelas.UpdateHandler.UpdateNamaKelas)))
	router.DELETE("/admin/kelas/nama-kelas/:idNamaKelas", requireAdmin(rateLimitStandard(kelas.DeleteHandler.DeleteKelas)))

	// MATA PELAJARAN
	router.GET("/admin/mata-pelajaran", requireAdmin(rateLimitStandard(mapel.GetHandler.ListMapel)))
	router.GET("/admin/mata-pelajaran/:idMapel", requireAdmin(rateLimitStandard(mapel.GetHandler.GetMapelByID)))
	router.POST("/admin/mata-pelajaran", requireAdmin(rateLimitStandard(mapel.CreateHandler.CreateMapel)))
	router.PATCH("/admin/mata-pelajaran/:idMapel", requireAdmin(rateLimitStandard(mapel.UpdateHandler.UpdateMapel)))
	router.DELETE("/admin/mata-pelajaran/:idMapel", requireAdmin(rateLimitStandard(mapel.DeleteHandler.DeleteMapel)))

	// RUANG UJIAN
	router.GET("/admin/ruang-ujian", requireAdmin(rateLimitStandard(ruangUjian.GetHandler.GetRuangUjian)))
	router.GET("/admin/ruang-ujian/id/:IdRuangan", requireAdmin(rateLimitStandard(ruangUjian.GetHandler.GetRuangUjianByID)))
	router.GET("/admin/ruang-ujian/kode/:KodeRuang", requireAdmin(rateLimitStandard(ruangUjian.GetHandler.GetRuangUjianByKode)))
	router.POST("/admin/ruang-ujian", requireAdmin(rateLimitStandard(ruangUjian.CreateHandler.CreateRuangUian)))
	router.PATCH("/admin/ruang-ujian/:idRuangan", requireAdmin(rateLimitStandard(ruangUjian.UpdateHandler.UpdateRuangUjian)))
	router.DELETE("/admin/ruang-ujian/:idRuangan", requireAdmin(rateLimitStandard(ruangUjian.DeleteHandler.DeleteRuangUjian)))

	// SESI
	router.GET("/admin/sesi", requireAdmin(rateLimitStandard(sesi.GetHandler.ListSesi)))
	router.GET("/admin/sesi/sesi-id/:idSesi", requireAdmin(rateLimitStandard(sesi.GetHandler.GetSesiByID)))
	router.GET("/admin/sesi/kode/:kodeSesi", requireAdmin(rateLimitStandard(sesi.GetHandler.GetSesiByKode)))
	router.POST("/admin/sesi", requireAdmin(rateLimitStandard(sesi.CreateHandler.CreateSesiHandler)))
	router.PATCH("/admin/sesi/:idSesi", requireAdmin(rateLimitStandard(sesi.UpdateHandler.UpdateSesi)))
	router.DELETE("/admin/sesi/:idSesi", requireAdmin(rateLimitStandard(sesi.DeleteHandler.DeleteSesi)))

	// PENGUMUMAN
	router.GET("/pengumuman/active", requireAdminGuru(rateLimitStandard(pengumuman.GetHandler.GetPengumumanActive)))
	router.GET("/pengumuman/non-active", requireAdminGuru(rateLimitStandard(pengumuman.GetHandler.GetPengumumanNonActive)))
	router.GET("/pengumuman/incoming", requireAdminGuru(rateLimitStandard(pengumuman.GetHandler.GetPengumumanIncoming)))
	router.GET("/pengumuman/id/:idPengumuman", requireAdminGuru(rateLimitStandard(pengumuman.GetHandler.GetPengumumanByID)))
	router.POST("/pengumuman", requireAdminGuru(rateLimitStandard(pengumuman.CreateHandler.CreatePengumuman)))
	router.PATCH("/pengumuman/:idPengumuman", requireAdminGuru(rateLimitStandard(pengumuman.UpdateHandler.UpdatePengumuman)))
	router.DELETE("/pengumuman/:idPengumuman", requireAdminGuru(rateLimitStandard(pengumuman.DeleteHandler.DeletePengumuman)))

	// IMPORT SOAL
	router.POST("/admin/bank-soal/import/:idBankSoal", requireAdminGuru(rateLimitStandard(importSoal.ImportHandler.ImportSoal)))
	router.GET("/admin/bank-soal/import-job/:idJob", requireAdminGuru(rateLimitStandard(importSoal.GetJobHandler.GetJobByID)))
	router.GET("/admin/bank-soal/import-jobs/:idBankSoal", requireAdminGuru(rateLimitStandard(importSoal.GetJobHandler.GetJobsByBankSoal)))

	// Siswa

	// Guru

	router.GET(documentPengumumanRoute, requireAdminGuru(rateLimitStandard(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

		clean = strings.TrimPrefix(clean, "/")
		if clean == "." || clean == "" {
			http.NotFound(w, r)
			return
		}

		full := filepath.Join(cfg.DocumentStore.Dir, "pengumuman", clean)
		http.ServeFile(w, r, full)
	})))

	router.GET(imageUploadRoute, rateLimitStandard(requireAccess(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

		clean = strings.TrimPrefix(clean, "/")
		if clean == "." || clean == "" {
			http.NotFound(w, r)
			return
		}

		full := filepath.Join(cfg.ImageStore.Dir, clean)
		http.ServeFile(w, r, full)
	})))

	corsHandler := middleware.CORSPolicy(router)

	protectCSRFHandler := middleware.PreventCSRF(corsHandler)

	server := &http.Server{
		Addr:    cfg.HTTP.Addr,
		Handler: protectCSRFHandler,
	}

	return &HTTPModule{
		Handler: corsHandler,
		Server:  server,
	}
}
