package aktivitasuserrepo

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func (r *AktivitasUserRepo) scanAktivitasUserRows(ctx context.Context, op string, rows pgx.Rows) ([]aktivitas_user.AktivitasUser, error) {
	var results []aktivitas_user.AktivitasUser
	for rows.Next() {
		var item aktivitas_user.AktivitasUser
		var description sql.NullString
		var ipAddress sql.NullString
		var roleName string

		if err := rows.Scan(
			&item.IdAktivitas,
			&item.IdPengguna,
			&item.Username,
			&roleName,
			&item.Action,
			&description,
			&ipAddress,
			&item.CreatedAt,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning aktivitas user", "op", op, "err", err)
			return nil, err
		}

		if description.Valid {
			item.Description = description.String
		}
		if ipAddress.Valid {
			item.IpAddress = ipAddress.String
		}
		item.Role = user.Role(roleName)

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating aktivitas user", "op", op, "err", err)
		return nil, err
	}

	return results, nil
}
