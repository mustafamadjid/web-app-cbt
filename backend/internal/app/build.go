package app


import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
)

type App struct {
	Infra  *InfraModule
	Tokens *TokenModule
	Auth   *AuthModule
	Users  *UserModule
	HTTP   *HTTPModule
}

func Build(ctx context.Context, cfg Config, dbURL string, hasher out.PasswordHasher) (*App, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	infra := BuildInfraModule(pool)
	tokens := BuildTokenModule(cfg)
	auth := BuildAuthModule(cfg, infra, tokens, hasher)
	users := BuildUserModule(cfg,infra, hasher)
	httpm := BuildHTTPModule(cfg, auth, users, tokens)

	return &App{
		Infra:  infra,
		Tokens: tokens,
		Auth:   auth,
		Users:  users,
		HTTP:   httpm,
	}, nil
}
