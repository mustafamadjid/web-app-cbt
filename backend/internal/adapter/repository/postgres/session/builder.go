package sessionrepo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func scanSessionRow(row pgx.Row) (session.Session, error) {
	var (
		item      session.Session
		revokedAt *time.Time
	)

	if err := row.Scan(
		&item.SessionID,
		&item.UserID,
		&item.Role,
		&item.ExpiresAt,
		&revokedAt,
	); err != nil {
		return session.Session{}, err
	}

	item.Role = user.Role(string(item.Role))
	item.Revoked = revokedAt != nil
	return item, nil
}

func (r *SessionRepo) scanSessionRows(ctx context.Context, op string, rows pgx.Rows) ([]session.Session, error) {
	var results []session.Session
	for rows.Next() {
		item, err := scanSessionRow(rows)
		if err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning sessions", "op", op, "err", err)
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating sessions", "op", op, "err", err)
		return nil, err
	}

	return results, nil
}
