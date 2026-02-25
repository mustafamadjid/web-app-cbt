package main

import (
	"context"
	"net/http"
	"path/filepath"

	// "fmt"
	"log"
	"os"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/security/bcrypt"
	"github.com/mustafamadjid/web-app-cbt/internal/app"
	"github.com/mustafamadjid/web-app-cbt/internal/infra/db"
	"github.com/mustafamadjid/web-app-cbt/internal/infra/logging"
)

func main() {
	ctx := context.Background()
	// env, ok := os.LookupEnv("POSTGRES_DBURL")
	// if !ok {
	// 	log.Fatal("POSTGRES_DBURL tidak ada")
	// }
	// fmt.Println(env)

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
		HTTP: app.HTTPConfig{
			Addr: ":8080",
		},
		JWT: app.JWTConfig{
			Issuer:        "web-app-cbt",
			AccessSecret:  os.Getenv("JWT_ACCESS_SECRET"),
			RefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
			AccessTTL:     15 * time.Minute,
			// AccessTTL:     15 * time.Second,
			RefreshTTL: 2 * 24 * time.Hour,
			// RefreshTTL:    30 * time.Second,
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

	pool, err := db.OpenPgxPool(ctx, db.PgxConfig{
		DbURL:           os.Getenv("POSTGRES_DBURL"),
		MaxConns:        20,
		MinConns:        5,
		MaxConnLifetime: 30 * time.Minute,
		MaxConnIdleTime: 5 * time.Minute,
		HealthTimeout:   3 * time.Second,
	})

	if err != nil {
		log.Fatalf("DB init failed: %v", err)
		log.Fatal("Exiting....")
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
	sesiMod := app.BuildSesiModule(infra)
	pengumumanMod := app.BuildPengumumanModule(cfg, infra, deleteFileSystem)
	resetPasswordMod := app.BuildResetPasswordModule(infra, hasher)

	httpMod := app.BuildHTTPModule(cfg, authMod, userMod, profilSekolahMod, aktivitasUserMod, kelasMod, mapelMod, ruangUjianMod, sesiMod, pengumumanMod, resetPasswordMod, tokens, infra, logger)

	log.Println("Listening on", cfg.HTTP.Addr)
	if err := httpMod.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

