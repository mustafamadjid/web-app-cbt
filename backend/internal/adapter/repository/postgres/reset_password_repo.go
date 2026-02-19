package postgres

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type ResetPasswordRepo struct {
	q Executor
	logger corelog.Logger
}

func NewResetPasswordRepo(q Executor, logger corelog.Logger) *ResetPasswordRepo {
	return &ResetPasswordRepo{q: q, logger: logger}
}

func (r *ResetPasswordRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func(r *ResetPasswordRepo)ResetPassword(ctx context.Context, idPengguna user.ID, password string) error {
	const query = `
		UPDATE pengguna
		SET password = $1
		WHERE id_pengguna = $2
	`

	tag,err := r.q.Exec(ctx,query,password,idPengguna)
	if err != nil {
		r.loggerFor(ctx).Error(ctx,"failed reset password", "op", "reset_password_repo.reset", "user_id", idPengguna, "err", err)
		return err
	}
	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}
	return nil
}