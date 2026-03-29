package authuserrepo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type AuthUserRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewAuthUserRepo(q pg.Executor, logger corelog.Logger) *AuthUserRepo {
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

	u, err := scanAuthUserRow(r.q.QueryRow(ctx, q, username))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.Pengguna{}, coreerror.ErrNotFound
		}
		r.loggerFor(ctx).Error(ctx, "failed finding auth user", "op", "auth_user_repo.find_by_username", "username", username, "err", err)
		return user.Pengguna{}, err
	}

	return u, nil
}
