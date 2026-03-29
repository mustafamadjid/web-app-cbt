package userrepo

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type nullString = sql.NullString

type penggunaScanner interface {
	Scan(dest ...any) error
}

func scanPenggunaRow(row penggunaScanner) (user.Pengguna, error) {
	var (
		result       user.Pengguna
		jenisKelamin int16
		roleName     string
		status       string
		email        sql.NullString
		noHp         sql.NullString
		foto         sql.NullString
	)

	err := row.Scan(
		&result.ID,
		&result.Username,
		&email,
		&result.PasswordHashed,
		&result.NamaLengkap,
		&jenisKelamin,
		&noHp,
		&roleName,
		&status,
		&foto,
	)
	if err != nil {
		return user.Pengguna{}, err
	}

	jenisKelaminValue, err := formatJenisKelamin(jenisKelamin)
	if err != nil {
		return user.Pengguna{}, err
	}

	result.Email = nullStringToUserEmailPtr(email)
	result.JenisKelamin = jenisKelaminValue
	result.Role = user.Role(roleName)
	result.StatusAkun = user.StatusAkun(status)
	result.NoHp = nullStringToPtr(noHp)
	if foto.Valid {
		result.Foto = foto.String
	}

	return result, nil
}

func (r *UserRepo) scanPenggunaRows(ctx context.Context, op string, rows pgx.Rows) ([]user.Pengguna, error) {
	var results []user.Pengguna
	for rows.Next() {
		item, err := scanPenggunaRow(rows)
		if err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning users", "op", op, "err", err)
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating users", "op", op, "err", err)
		return nil, err
	}

	return results, nil
}
