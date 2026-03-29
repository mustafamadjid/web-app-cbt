package sesirepo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
)

type sesiScanner interface {
	Scan(dest ...any) error
}

func (r *SesiRepo) buildListSesiQuery(filter query.ListSesiFilter) (string, []any) {
	queryText := `
		SELECT
			id_sesi,
			kode_sesi,
			nama_sesi
		FROM sesi_ujian	
	`

	where := make([]string, 0, 1)
	args := make([]any, 0, 3)

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("(nama_sesi ILIKE $%d OR kode_sesi ILIKE $%d)", idx, idx))
	}

	if len(where) > 0 {
		queryText = fmt.Sprintf("%s WHERE %s", queryText, strings.Join(where, " AND "))
	}

	queryText = fmt.Sprintf("%s ORDER BY id_sesi ASC", queryText)

	if filter.Limit > 0 {
		args = append(args, filter.Limit)
		limitIndex := len(args)
		args = append(args, filter.Offset)
		offsetIndex := len(args)
		queryText = fmt.Sprintf("%s LIMIT $%d OFFSET $%d", queryText, limitIndex, offsetIndex)
	}

	return queryText, args
}

func scanSesiRow(row sesiScanner) (sesi.Sesi, error) {
	var item sesi.Sesi
	if err := row.Scan(
		&item.IdSesi,
		&item.KodeSesi,
		&item.NamaSesi,
	); err != nil {
		return sesi.Sesi{}, err
	}

	return item, nil
}

func (r *SesiRepo) scanSesiRows(ctx context.Context, op string, rows pgx.Rows) ([]sesi.Sesi, error) {
	var results []sesi.Sesi
	for rows.Next() {
		item, err := scanSesiRow(rows)
		if err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning sesi", "layer", "repo.db", "op", op, "err", err)
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating sesi", "layer", "repo.db", "op", op, "err", err)
		return nil, err
	}

	return results, nil
}
