package importsoaljobrepo

import (
	"context"

	"github.com/jackc/pgx/v5"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
)

type importSoalJobScanner interface {
	Scan(dest ...any) error
}

func scanImportSoalJobRow(row importSoalJobScanner) (importsoal.ImportSoalJob, error) {
	var item importsoal.ImportSoalJob
	if err := row.Scan(
		&item.IDJob,
		&item.IDBankSoal,
		&item.IDPengguna,
		&item.Status,
		&item.FilePath,
		&item.ErrorMsg,
		&item.TotalSoal,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return importsoal.ImportSoalJob{}, err
	}

	return item, nil
}

func (r *ImportSoalJobRepo) scanImportSoalJobRows(ctx context.Context, op string, rows pgx.Rows) ([]importsoal.ImportSoalJob, error) {
	var items []importsoal.ImportSoalJob
	for rows.Next() {
		item, err := scanImportSoalJobRow(rows)
		if err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning import soal job", "layer", "repo.db", "op", op, "err", err)
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating import soal job", "layer", "repo.db", "op", op, "err", err)
		return nil, err
	}

	return items, nil
}
