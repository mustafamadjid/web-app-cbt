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

func (r *SessionRepo)GetSessionByUserId(ctx context.Context, userID user.ID) (session.Session, error) {
	const query = `
		SELECT id,
			id_pengguna,
			expires_at,
			revoked_at
		FROM sessions
		WHERE id_pengguna = $1
		AND revoked_at IS NULL
		AND expires_at > now()
	`

	var result session.Session
	var revokedAt *time.Time
	err := r.q.QueryRow(ctx, query, userID).Scan(
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

func (r *SessionRepo) GetAllActiveSession(ctx context.Context)([]session.Session,error){
	const query = `
	SELECT id,
		   id_pengguna,
		   expires_at,
		   revoked_at
	FROM sessions
	WHERE revoked_at IS NULL
	AND expires_at > now()
	`

	rows, err := r.q.Query(ctx,query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []session.Session
	for rows.Next() {
		var item session.Session
		var revokedAt *time.Time
		if err := rows.Scan(
			&item.SessionID,
			&item.UserID,
			&item.ExpiresAt,
			&revokedAt,
		); err != nil {
			return nil, err
		}
		item.Revoked = revokedAt != nil
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
    return nil, err
}
	return results, nil
}

func (r *SessionRepo) HasActiveSession(ctx context.Context, userID user.ID) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM sessions
			WHERE id_pengguna = $1
				AND revoked_at IS NULL
				AND expires_at > now()
		)
	`

	var exists bool
	if err := r.q.QueryRow(ctx, query, userID).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
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
