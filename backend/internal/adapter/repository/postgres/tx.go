package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	pgprofilgururepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/profil_guru"
	pgprofilsiswarepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/profil_siswa"
	pguserrepo "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/user"

	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

type pgTx struct {
	ctx    context.Context
	tx     pgx.Tx
	logger corelog.Logger
}

func (p *pgTx) Pengguna() outuser.UserRepository {
	return pguserrepo.NewUserRepo(p.tx, p.logger)
}

func (p *pgTx) ProfilGuru() outuser.ProfilGuruRepository {
	return pgprofilgururepo.NewProfilgGuruRepo(p.tx, p.logger)
}

func (p *pgTx) ProfilSiswa() outuser.ProfilSiswaRepository {
	return pgprofilsiswarepo.NewProfilSiswaRepo(p.tx, p.logger)
}

func (p *pgTx) Commit() error {
	return p.tx.Commit(p.ctx)
}

func (p *pgTx) Rollback() error {
	return p.tx.Rollback(p.ctx)
}
