package pengumumanrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
)

type pengumumanScanner interface {
	Scan(dest ...any) error
}

func scanPengumumanRow(row pengumumanScanner) (pengumuman.Pengumuman, error) {
	var (
		item           pengumuman.Pengumuman
		tanggalRilis   time.Time
		tanggalSelesai time.Time
		dokumen        sql.NullString
	)

	if err := row.Scan(
		&item.IdPengumuman,
		&item.IdPengguna,
		&item.JudulPengumuman,
		&item.IsiPengumuman,
		&tanggalRilis,
		&tanggalSelesai,
		&dokumen,
	); err != nil {
		return pengumuman.Pengumuman{}, err
	}

	item.TanggalRilisPengumuman = tanggalRilis.Format("2006-01-02")
	item.TanggalSelesaiPengumuman = tanggalSelesai.Format("2006-01-02")
	if dokumen.Valid {
		item.DokumenPengumuman = dokumen.String
	}

	return item, nil
}

func (r *PengumumanRepo) scanPengumumanRows(ctx context.Context, op string, rows pgx.Rows) ([]pengumuman.Pengumuman, error) {
	var results []pengumuman.Pengumuman
	for rows.Next() {
		item, err := scanPengumumanRow(rows)
		if err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning pengumuman", "layer", "repo.db", "op", op, "err", err)
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating pengumuman", "layer", "repo.db", "op", op, "err", err)
		return nil, err
	}

	return results, nil
}
