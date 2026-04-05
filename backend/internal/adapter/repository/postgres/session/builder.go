package sessionrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
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

func scanSessionWithUserRow(row pgx.Row) (session.SessionWithUser, error) {
	var (
		item         session.SessionWithUser
		revokedAt    *time.Time
		email        sql.NullString
		noHp         sql.NullString
		foto         sql.NullString
		jenisKelamin int16
		userRole     string
		statusAkun   string
	)

	if err := row.Scan(
		&item.Session.SessionID,
		&item.Session.UserID,
		&item.Session.Role,
		&item.Session.ExpiresAt,
		&revokedAt,
		&item.Pengguna.ID,
		&item.Pengguna.Username,
		&email,
		&item.Pengguna.NamaLengkap,
		&jenisKelamin,
		&noHp,
		&userRole,
		&statusAkun,
		&foto,
	); err != nil {
		return session.SessionWithUser{}, err
	}

	jenisKelaminValue, err := formatJenisKelaminForSessionRepo(jenisKelamin)
	if err != nil {
		return session.SessionWithUser{}, err
	}

	item.Session.Role = user.Role(string(item.Session.Role))
	item.Session.Revoked = revokedAt != nil

	item.Pengguna.Email = nullStringToUserEmailPtr(email)
	item.Pengguna.JenisKelamin = jenisKelaminValue
	item.Pengguna.NoHp = nullStringToPtr(noHp)
	item.Pengguna.Role = user.Role(userRole)
	item.Pengguna.StatusAkun = user.StatusAkun(statusAkun)
	if foto.Valid {
		item.Pengguna.Foto = foto.String
	}

	return item, nil
}

func (r *SessionRepo) scanSessionWithUserRows(ctx context.Context, op string, rows pgx.Rows) ([]session.SessionWithUser, error) {
	var results []session.SessionWithUser
	for rows.Next() {
		item, err := scanSessionWithUserRow(rows)
		if err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning active sessions", "op", op, "err", err)
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating active sessions", "op", op, "err", err)
		return nil, err
	}

	return results, nil
}

func formatJenisKelaminForSessionRepo(value int16) (string, error) {
	switch value {
	case 1:
		return "LAKI_LAKI", nil
	case 2:
		return "PEREMPUAN", nil
	default:
		return "", coreerror.ErrInvalidInput
	}
}

func nullStringToPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}

	v := value.String
	return &v
}

func nullStringToUserEmailPtr(value sql.NullString) *user.Email {
	if !value.Valid {
		return nil
	}

	email := user.Email(value.String)
	return &email
}
