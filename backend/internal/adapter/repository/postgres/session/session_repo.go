package sessionrepo

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type SessionRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewSessionRepo(q pg.Executor, logger corelog.Logger) *SessionRepo {
	return &SessionRepo{q: q, logger: logger}
}

func (r *SessionRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *SessionRepo) CreateSession(ctx context.Context, userID user.ID, role user.Role, expiresAt time.Time) (string, error) {
	const query = `
		INSERT INTO sessions (id_pengguna, role, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var sessionID string
	if err := r.q.QueryRow(ctx, query, userID, role, expiresAt).Scan(&sessionID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return "", coreerror.ErrHasSession
			}
		}

		r.loggerFor(ctx).Error(ctx, "failed creating session", "op", "session_repo.create", "user_id", userID, "err", err)
		return "", err
	}

	return sessionID, nil
}

func (r *SessionRepo) GetSession(ctx context.Context, sessionID string) (session.Session, error) {
	const query = `
		SELECT id,
			id_pengguna,
			role,
			expires_at,
			revoked_at
		FROM sessions
		WHERE id = $1
	`

	result, err := scanSessionRow(r.q.QueryRow(ctx, query, sessionID))
	if errors.Is(err, pgx.ErrNoRows) {
		return session.Session{}, coreerror.ErrNotFound
	}
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed getting session", "op", "session_repo.get", "session_id", sessionID, "err", err)
		return session.Session{}, err
	}
	return result, nil
}

func (r *SessionRepo) GetSessionByUserId(ctx context.Context, userID user.ID) (session.Session, error) {
	const query = `
		SELECT id,
			id_pengguna,
			role,
			expires_at,
			revoked_at
		FROM sessions
		WHERE id_pengguna = $1
		AND revoked_at IS NULL
		AND expires_at > now()
	`

	result, err := scanSessionRow(r.q.QueryRow(ctx, query, userID))

	if errors.Is(err, pgx.ErrNoRows) {
		return session.Session{}, coreerror.ErrNotFound
	}
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed getting session by user", "op", "session_repo.get_by_user", "user_id", userID, "err", err)
		return session.Session{}, err
	}

	return result, nil
}

func (r *SessionRepo) GetAllActiveSession(ctx context.Context) ([]session.SessionWithUser, error) {
	const query = `
	SELECT s.id,
		   s.id_pengguna,
		   s.role,
		   s.expires_at,
		   s.revoked_at,
		   p.id_pengguna,
		   p.username,
		   p.email,
		   p.nama_lengkap,
		   p.jenis_kelamin,
		   p.no_hp,
		   r.nama_role,
		   p.status_akun,
		   p.foto
	FROM sessions s
	JOIN pengguna p ON s.id_pengguna = p.id_pengguna
	JOIN role r ON p.id_role = r.id_role
	WHERE s.revoked_at IS NULL
	AND s.expires_at > now()
	ORDER BY s.expires_at DESC
	`

	rows, err := r.q.Query(ctx, query)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed listing sessions", "op", "session_repo.list_active", "err", err)
		return nil, err
	}

	defer rows.Close()

	return r.scanSessionWithUserRows(ctx, "session_repo.list_active", rows)
}

func (r *SessionRepo) HasActiveSession(ctx context.Context, userID user.ID) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM sessions
			WHERE id_pengguna = $1
				AND revoked_at IS NULL
		)
	`

	var exists bool
	if err := r.q.QueryRow(ctx, query, userID).Scan(&exists); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed checking active session", "op", "session_repo.has_active", "user_id", userID, "err", err)
		return false, err
	}

	return exists, nil
}

func (r *SessionRepo) RevokeExpiredSessions(ctx context.Context, userID user.ID) (bool, error) {
	const q = `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE id_pengguna = $1
		  AND revoked_at IS NULL
		  AND expires_at <= NOW()
	`
	tag, err := r.q.Exec(ctx, q, userID)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed revoking expired sessions", "op", "session_repo.revoke_expired", "user_id", userID, "err", err)
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func (r *SessionRepo) RevokeSession(ctx context.Context, sessionID string) error {
	const query = `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE id = $1
			AND revoked_at IS NULL
	`

	result, err := r.q.Exec(ctx, query, sessionID)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed revoking session", "op", "session_repo.revoke", "session_id", sessionID, "err", err)
		return err
	}
	if result.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	return nil
}

func (r *SessionRepo) RevokeSessionAllbyUser(ctx context.Context, userID user.ID) error {
	const query = `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE id_pengguna = $1
			AND revoked_at IS NULL
	`

	ct, err := r.q.Exec(ctx, query, userID)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed revoking sessions for user", "op", "session_repo.revoke_by_user", "user_id", userID, "err", err)
		return err
	}

	if ct.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}
	return nil
}
