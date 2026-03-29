package aktivitasuserrepo

import (
	"context"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type AktivitasUserRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewAktivitasUserRepo(q pg.Executor, logger corelog.Logger) *AktivitasUserRepo {
	return &AktivitasUserRepo{q: q, logger: logger}
}

func (r *AktivitasUserRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *AktivitasUserRepo) CreateAktivitasUser(ctx context.Context, aktivitasUser aktivitas_user.AktivitasUser) error {
	const query = `
		INSERT INTO aktivitas_user (id_pengguna, action, description, ip_address) VALUES ($1, $2, $3, $4)
	`

	_, err := r.q.Exec(ctx, query, aktivitasUser.IdPengguna, aktivitasUser.Action, aktivitasUser.Description, aktivitasUser.IpAddress)

	if err != nil {
		r.loggerFor(ctx).Error(ctx, "Failed creating aktivitas user", "op", "aktivitas_user_repo.insert", "err", err)
	}
	return err
}

func (r *AktivitasUserRepo) GetAktivitasUser(ctx context.Context) ([]aktivitas_user.AktivitasUser, error) {
	const query = `
		SELECT
			au.id_aktivitas,
			au.id_pengguna,
			p.username,
			r.nama_role,
			au.action,
			au.description,
			au.ip_address,
			au.created_at
		FROM aktivitas_user au
		JOIN pengguna p ON au.id_pengguna = p.id_pengguna
		JOIN role r ON p.id_role = r.id_role
		ORDER BY au.created_at DESC
		LIMIT 30
	`

	rows, err := r.q.Query(ctx, query)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "Failed listing aktivitas user", "op", "aktivitas_user_repo.list", "err", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanAktivitasUserRows(ctx, "aktivitas_user_repo.list_scan", rows)
}
