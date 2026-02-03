package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type App struct {
	Infra         *InfraModule
	Tokens        *TokenModule
	Auth          *AuthModule
	Users         *UserModule
	ProfilSekolah *ProfilSekolahModule
	HTTP          *HTTPModule
}

func Build(ctx context.Context, cfg Config, dbURL string, hasher out.PasswordHasher, logger corelog.Logger) (*App, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	infra := BuildInfraModule(pool, logger)
	tokens := BuildTokenModule(cfg)
	auth := BuildAuthModule(cfg, infra, tokens, hasher)
	users := BuildUserModule(cfg, infra, hasher)
	profilSekolah := BuildProfilSekolahModule(cfg, infra)
	httpm := BuildHTTPModule(cfg, auth, users, profilSekolah, tokens, logger)

	return &App{
		Infra:         infra,
		Tokens:        tokens,
		Auth:          auth,
		Users:         users,
		ProfilSekolah: profilSekolah,
		HTTP:          httpm,
	}, nil
}
