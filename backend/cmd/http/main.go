package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/security/bcrypt"
	"github.com/mustafamadjid/web-app-cbt/internal/app"
	"github.com/mustafamadjid/web-app-cbt/internal/infra/db"
	"github.com/mustafamadjid/web-app-cbt/internal/infra/logging"
)

func main() {
	// 1) Root ctx: cancel on SIGINT/SIGTERM
	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := logging.NewLogger(os.Getenv("ENV"))

	cookieSecure := true
	if os.Getenv("ENV") == "dev" || os.Getenv("ENV") == "" {
		cookieSecure = false
	}

	appDir := ResolveAppDir()
	uploadDir := ResolveUploadDir(appDir)

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	imageUploadDir := filepath.Join(uploadDir, "image")
	documentUploadDir := filepath.Join(uploadDir, "document")

	cfg := app.Config{
		HTTP: app.HTTPConfig{Addr: ":8080"},
		JWT: app.JWTConfig{
			Issuer:        "web-app-cbt",
			AccessSecret:  os.Getenv("JWT_ACCESS_SECRET"),
			RefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
			AccessTTL:     15 * time.Minute,
			RefreshTTL:    2 * 24 * time.Hour,
		},
		Cookie: app.CookieConfig{
			AccessName:  "access_token",
			RefreshName: "refresh_token",
			Domain:      "",
			Secure:      cookieSecure,
			SameSite:    http.SameSiteLaxMode,
		},
		ImageStore: app.ImageStoreConfig{
			Dir:      imageUploadDir,
			BaseURL:  baseURL,
			Route:    "/uploads/image",
			MaxBytes: 5 << 20,
		},
		DocumentStore: app.DocumentStoreConfig{
			Dir:      documentUploadDir,
			BaseURL:  baseURL,
			Route:    "/uploads/document",
			MaxBytes: 10 << 20,
		},
	}

	// 2) DB uses rootCtx (so can be cancelled on shutdown)
	pool, err := db.OpenPgxPool(rootCtx, db.PgxConfig{
		DbURL:           os.Getenv("POSTGRES_DBURL"),
		MaxConns:        20,
		MinConns:        5,
		MaxConnLifetime: 30 * time.Minute,
		MaxConnIdleTime: 5 * time.Minute,
		HealthTimeout:   3 * time.Second,
	})
	if err != nil {
		log.Fatalf("DB init failed: %v", err)
	}
	log.Println("DB init success")
	defer pool.Close()

	infra := app.BuildInfraModule(pool, logger)
	tokens := app.BuildTokenModule(cfg)
	hasher := bcrypt.NewHasher(0)
	deleteFileSystem := app.BuildDeleteFileModule(uploadDir)

	aktivitasUserMod := app.BuildAktivitasUserModule(infra)
	authMod := app.BuildAuthModule(cfg, infra, tokens, hasher, aktivitasUserMod)
	userMod := app.BuildUserModule(cfg, infra, hasher, aktivitasUserMod, deleteFileSystem)
	profilSekolahMod := app.BuildProfilSekolahModule(cfg, infra)
	kelasMod := app.BuildKelasModule(infra, aktivitasUserMod)
	mapelMod := app.BuildMataPelajaranModule(infra)
	ruangUjianMod := app.BuildRuangUjianModule(infra)
	ujianMod := app.BuildUjianModule(infra)
	sesiMod := app.BuildSesiModule(infra)
	pengumumanMod := app.BuildPengumumanModule(cfg, infra, deleteFileSystem)
	bankSoalMod := app.BuildBankSoalModule(infra)
	resetPasswordMod := app.BuildResetPasswordModule(infra, hasher)
	importSoalMod := app.BuildImportSoalModule(infra, cfg, logger)

	httpMod := app.BuildHTTPModule(app.HTTPDeps{
		Config:        cfg,
		Logger:        logger,
		Auth:          authMod,
		Users:         userMod,
		ProfilSekolah: profilSekolahMod,
		AktivitasUser: aktivitasUserMod,
		Kelas:         kelasMod,
		MataPelajaran: mapelMod,
		RuangUjian:    ruangUjianMod,
		Ujian:         ujianMod,
		Sesi:          sesiMod,
		Pengumuman:    pengumumanMod,
		BankSoal:      bankSoalMod,
		ResetPassword: resetPasswordMod,
		ImportSoal:    importSoalMod,
		Tokens:        tokens,
		Infra:         infra,
	})

	// 3) Start worker
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		importSoalMod.Worker.Start(rootCtx)
	}()

	// 4) Start HTTP server
	srv := httpMod.Server
	go func() {
		log.Println("Listening on", cfg.HTTP.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	// 5) Wait for shutdown signal
	<-rootCtx.Done()
	log.Println("Shutdown signal received")

	// 6) Graceful shutdown (HTTP)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown error: %v", err)
	}

	// 7) Wait worker exit (depends on ctx-aware processJobs!)
	wg.Wait()
	log.Println("Shutdown complete")
}
