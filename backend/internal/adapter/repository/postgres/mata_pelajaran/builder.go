package matapelajaranrepo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
)

type mapelScanner interface {
	Scan(dest ...any) error
}

func (r *MapelRepo) buildListMapelQuery(filter query.ListMapelFilter) (string, []any) {
	queryText := `
		SELECT 
			id_mapel,
			id_kelas,
			kode_mapel,
			nama_mapel,
			deskripsi
		FROM mata_pelajaran
	`

	where := make([]string, 0, 2)
	args := make([]any, 0, 4)

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("(nama_mapel ILIKE $%d OR deskripsi ILIKE $%d OR kode_mapel ILIKE $%d)", idx, idx, idx))
	}

	if filter.NamaMapel != nil {
		args = append(args, *filter.NamaMapel)
		where = append(where, fmt.Sprintf("nama_mapel = $%d", len(args)))
	}

	if filter.TingkatKelas != nil {
		args = append(args, *filter.TingkatKelas)
		where = append(where, fmt.Sprintf("id_kelas = $%d", len(args)))
	}

	if len(where) > 0 {
		queryText = fmt.Sprintf("%s WHERE %s", queryText, strings.Join(where, " AND "))
	}

	queryText = fmt.Sprintf("%s ORDER BY created_at ASC", queryText)

	if filter.Limit > 0 {
		args = append(args, filter.Limit)
		limitIndex := len(args)
		args = append(args, filter.Offset)
		offsetIndex := len(args)
		queryText = fmt.Sprintf("%s LIMIT $%d OFFSET $%d", queryText, limitIndex, offsetIndex)
	}

	return queryText, args
}

func scanMapelRow(row mapelScanner) (matapelajaran.MataPelajaran, error) {
	var item matapelajaran.MataPelajaran
	if err := row.Scan(
		&item.IdMapel,
		&item.IdKelas,
		&item.KodeMapel,
		&item.NamaMapel,
		&item.Deskripsi,
	); err != nil {
		return matapelajaran.MataPelajaran{}, err
	}

	return item, nil
}

func (r *MapelRepo) scanMapelRows(ctx context.Context, op string, rows pgx.Rows) ([]matapelajaran.MataPelajaran, error) {
	var results []matapelajaran.MataPelajaran
	for rows.Next() {
		item, err := scanMapelRow(rows)
		if err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning mapel", "layer", "core.service", "op", op, "err", err)
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating mapel", "layer", "core.service", "op", op, "err", err)
		return nil, err
	}

	return results, nil
}
