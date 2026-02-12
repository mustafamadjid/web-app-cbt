package main

import (
	"context"
	"net/http"

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

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./public/uploads"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

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
			Dir:      uploadDir,
			BaseURL:  baseURL,
			Route:    "/uploads",
			MaxBytes: 5 << 20,
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

	aktivitasUserMod := app.BuildAktivitasUserModule(infra)
	authMod := app.BuildAuthModule(cfg, infra, tokens, hasher, aktivitasUserMod)
	userMod := app.BuildUserModule(cfg, infra, hasher, aktivitasUserMod)
	profilSekolahMod := app.BuildProfilSekolahModule(cfg, infra)
	kelasMod := app.BuildKelasModule(infra, aktivitasUserMod)
	mapelMod := app.BuildMataPelajaranModule(infra)

	httpMod := app.BuildHTTPModule(cfg, authMod, userMod, profilSekolahMod, aktivitasUserMod, kelasMod, mapelMod, tokens, logger)

	log.Println("Listening on", cfg.HTTP.Addr)
	if err := httpMod.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
