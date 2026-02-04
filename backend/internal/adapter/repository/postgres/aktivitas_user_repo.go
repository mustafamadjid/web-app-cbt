package postgres

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)
type AktivitasUserRepo struct {
	q Executor
	logger corelog.Logger
}

func NewAktivitasUserRepo(q Executor, logger corelog.Logger) *AktivitasUserRepo {
	return &AktivitasUserRepo{q: q, logger: logger}
}

func (r *AktivitasUserRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func(r *AktivitasUserRepo)CreateAktivitasUser(ctx context.Context, aktivitasUser aktivitas_user.AktivitasUser)error{
	const query = `
		INSERT INTO aktivitas_user (id_pengguna, action, description, ip_address) VALUES ($1, $2, $3, $4)
	`

	_, err := r.q.Exec(ctx, query, aktivitasUser.IdPengguna, aktivitasUser.Action, aktivitasUser.Description, aktivitasUser.IpAddress)

	if err != nil {
		r.loggerFor(ctx).Error(ctx,"Failed creating aktivitas user","op","aktivitas_user_repo.insert","err",err)
	}

	r.loggerFor(ctx).Info(ctx,"Successfully created aktivitas user","op","aktivitas_user_repo.insert")
	return err
}

func (r *AktivitasUserRepo)GetAktivitasUser(ctx context.Context) ([]aktivitas_user.AktivitasUser, error){
	const query = `
		SELECT id_aktivitas, id_pengguna, action, description, ip_address, created_at FROM aktivitas_user
	`

	rows, err := r.q.Query(ctx, query)
	if err != nil {
		r.loggerFor(ctx).Error(ctx,"Failed listing aktivitas user","op","aktivitas_user_repo.list","err",err)
		return nil, err
	}
	defer rows.Close()

	var results []aktivitas_user.AktivitasUser
	for rows.Next() {
		var item aktivitas_user.AktivitasUser
		if err := rows.Scan(
			&item.IdAktivitas,
			&item.IdPengguna,
			&item.Action,
			&item.Description,
			&item.IpAddress,
			&item.CreatedAt,
		); err != nil {
			r.loggerFor(ctx).Error(ctx,"Failed scanning aktivitas user","op","aktivitas_user_repo.list_scan","err",err)
			return nil, err
		}
		results = append(results, item)
	}

	r.loggerFor(ctx).Info(ctx,"Successfully listed aktivitas user","op","aktivitas_user_repo.list")
	return results, nil
}