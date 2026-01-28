package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
	userrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/user_repo"
)

type pgTx struct {
	ctx context.Context
	tx  *pgx.Tx
}

func (p *pgTx) Pengguna() outuser.UserRepository {
	return &userrepo.UserRepo{q: p.tx}
}