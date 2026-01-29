package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type SessionRepo struct {
	q Executor
}

func NewSessionRepo(q Executor) *SessionRepo {
	return &SessionRepo{q: q}
}


func (r *SessionRepo) CreateSession(ctx context.Context, userID user.ID, expiresAt time.Time) (string, error) {
	const query = `
		INSERT INTO sessions (id_pengguna, expires_at)
		VALUES ($1, $2)
		RETURNING id
	`

	var sessionID string
	if err := r.q.QueryRow(ctx, query, userID, expiresAt).Scan(&sessionID); err != nil {
		return "", err
	}

	return sessionID, nil
}

func (r *SessionRepo) GetSession(ctx context.Context, sessionID string) (session.Session, error) {
	const query = `
		SELECT id,
			id_pengguna,
			expires_at,
			revoked_at
		FROM sessions
		WHERE id = $1
	`

	var result session.Session
	var revokedAt *time.Time
	err := r.q.QueryRow(ctx, query, sessionID).Scan(
		&result.SessionID,
		&result.UserID,
		&result.ExpiresAt,
		&revokedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return session.Session{}, coreerror.ErrNotFound
	}
	if err != nil {
		return session.Session{}, err
	}

	result.Revoked = revokedAt != nil
	return result, nil
}

func (r *SessionRepo) RevokeSession(ctx context.Context, sessionID string) error {
	const query = `
		UPDATE sessions
		SET revoked_at = now()
		WHERE id = $1
			AND revoked_at IS NULL
	`

	result, err := r.q.Exec(ctx, query, sessionID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	return nil
}

func (r *SessionRepo) RevokeSessionAllbyUser(ctx context.Context, userID user.ID, now time.Time) error {
	const query = `
		UPDATE sessions
		SET revoked_at = $1
		WHERE id_pengguna = $2
			AND revoked_at IS NULL
			AND expires_at > $1
	`

	_, err := r.q.Exec(ctx, query, now, userID)
	return err
}
