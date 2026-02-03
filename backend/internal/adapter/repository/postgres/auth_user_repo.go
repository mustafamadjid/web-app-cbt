package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type AuthUserRepo struct {
	q      Executor
	logger corelog.Logger
}

func NewAuthUserRepo(q Executor, logger corelog.Logger) *AuthUserRepo {
	return &AuthUserRepo{q: q, logger: logger}
}

func (r *AuthUserRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *AuthUserRepo) FindByUsername(ctx context.Context, username string) (user.Pengguna, error) {
	const q = `
		SELECT p.id_pengguna, p.username, p.password, r.nama_role, p.status_akun
		FROM pengguna p
		JOIN role r ON p.id_role = r.id_role
		WHERE p.username = $1
	`

	var u user.Pengguna
	if err := r.q.QueryRow(ctx, q, username).Scan(&u.ID, &u.Username, &u.PasswordHashed, &u.Role, &u.StatusAkun); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.Pengguna{}, coreerror.ErrNotFound
		}
		r.loggerFor(ctx).Error(ctx, "failed finding auth user", "op", "auth_user_repo.find_by_username", "username", username, "err", err)
		return user.Pengguna{}, err
	}
	u.Role = user.Role(string(u.Role))
	u.StatusAkun = user.StatusAkun(string(u.StatusAkun))
	return u, nil
}
