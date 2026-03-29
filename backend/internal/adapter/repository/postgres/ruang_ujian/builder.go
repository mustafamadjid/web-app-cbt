package ruangujianrepo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ruang_ujian"
)

type ruangUjianScanner interface {
	Scan(dest ...any) error
}

func (r *RuangUjianRepo) buildListRuangUjianQuery(filter query.ListRuangUjianFilter) (string, []any) {
	baseQuery := `
		SELECT
			id_ruangan,
			nama_ruangan,
			kode_ruang
		FROM ruang_ujian
	`

	where := make([]string, 0, 1)
	args := make([]any, 0, 3)

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("(nama_ruangan ILIKE $%d OR kode_ruang ILIKE $%d)", idx, idx))
	}

	if len(where) > 0 {
		baseQuery = fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(where, " AND "))
	}

	baseQuery = fmt.Sprintf("%s ORDER BY id_ruangan ASC", baseQuery)

	if filter.Limit > 0 {
		args = append(args, filter.Limit)
		limitIndex := len(args)
		args = append(args, filter.Offset)
		offsetIndex := len(args)
		baseQuery = fmt.Sprintf("%s LIMIT $%d OFFSET $%d", baseQuery, limitIndex, offsetIndex)
	}

	return baseQuery, args
}

func scanRuangUjianRow(row ruangUjianScanner) (ruangujian.RuangUjian, error) {
	var item ruangujian.RuangUjian
	if err := row.Scan(
		&item.IdRuangan,
		&item.NamaRuangan,
		&item.KodeRuang,
	); err != nil {
		return ruangujian.RuangUjian{}, err
	}

	return item, nil
}

func (r *RuangUjianRepo) scanRuangUjianRows(ctx context.Context, op string, rows pgx.Rows) ([]ruangujian.RuangUjian, error) {
	var results []ruangujian.RuangUjian
	for rows.Next() {
		item, err := scanRuangUjianRow(rows)
		if err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning ruang ujian", "layer", "repo.db", "op", op, "err", err)
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating ruang ujian", "layer", "repo.db", "op", op, "err", err)
		return nil, err
	}

	return results, nil
}
