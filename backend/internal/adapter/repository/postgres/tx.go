package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

type pgTx struct {
	ctx context.Context
	tx  pgx.Tx
}

func (p *pgTx) Pengguna() outuser.UserRepository {
	return NewUserRepo(p.tx)
}

func (p *pgTx) ProfilGuru() outuser.ProfilGuruRepository {
	return NewProfilgGuruRepo(p.tx)
}

func (p *pgTx) ProfilSiswa() outuser.ProfilSiswaRepository {
	return NewProfilSiswaRepo(p.tx)
}

func (p *pgTx) Commit() error {
	return p.tx.Commit(p.ctx)
}

func (p *pgTx) Rollback() error {
	return p.tx.Rollback(p.ctx)
}
