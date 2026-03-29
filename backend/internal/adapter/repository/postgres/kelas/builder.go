package kelasrepo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
)

func (r *KelasRepo) buildListKelasQuery(filter query.ListKelasFilter) (string, []any) {
	baseQuery := `
	SELECT 
	tk.id_kelas,
	tk.tingkat_kelas,
	k.id_nama_kelas,
	k.nama_kelas
	FROM kelas tk
	LEFT JOIN nama_kelas k ON tk.id_kelas = k.id_kelas
	`

	where := make([]string, 0, 2)
	args := make([]any, 0, 4)

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("(k.nama_kelas ILIKE $%d OR tk.tingkat_kelas::text ILIKE $%d)", idx, idx))
	}

	if filter.TingkatKelas != nil {
		args = append(args, *filter.TingkatKelas)
		where = append(where, fmt.Sprintf("tk.tingkat_kelas = $%d", len(args)))
	}

	if len(where) > 0 {
		baseQuery = fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(where, " AND "))
	}

	baseQuery = fmt.Sprintf("%s ORDER BY tk.tingkat_kelas ASC, k.nama_kelas ASC", baseQuery)

	if filter.Limit > 0 {
		args = append(args, filter.Limit)
		limitIndex := len(args)
		args = append(args, filter.Offset)
		offsetIndex := len(args)
		baseQuery = fmt.Sprintf("%s LIMIT $%d OFFSET $%d", baseQuery, limitIndex, offsetIndex)
	}

	return baseQuery, args
}

func (r *KelasRepo) scanKelasRows(ctx context.Context, op string, rows pgx.Rows) ([]kelas.FullKelasData, error) {
	var (
		itemsTingkat []kelas.TingkatKelas
		itemsNama    []kelas.NamaKelas
	)

	seenTingkat := make(map[int]bool)
	for rows.Next() {
		var (
			idTingkat int
			tingkat   int
			idNama    sql.NullInt64
			namaKelas sql.NullString
		)

		if err := rows.Scan(&idTingkat, &tingkat, &idNama, &namaKelas); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning kelas", "op", op, "err", err)
			return nil, err
		}

		if !seenTingkat[idTingkat] {
			seenTingkat[idTingkat] = true
			itemsTingkat = append(itemsTingkat, kelas.TingkatKelas{
				IdTingkatKelas: kelas.ID(idTingkat),
				TingkatKelas:   tingkat,
			})
		}

		if idNama.Valid && namaKelas.Valid {
			itemsNama = append(itemsNama, kelas.NamaKelas{
				IdNamaKelas:    kelas.ID(idNama.Int64),
				IdTingkatKelas: kelas.ID(idTingkat),
				NamaKelas:      namaKelas.String,
			})
		}
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating kelas", "op", op, "err", err)
		return nil, err
	}

	return []kelas.FullKelasData{{
		ItemsTingkatKelas: itemsTingkat,
		ItemsNamaKelas:    itemsNama,
	}}, nil
}

func scanKelasRow(row pgx.Row) (kelas.KelasData, error) {
	var (
		itemsTingkatKelas kelas.TingkatKelas
		itemsNamaKelas    kelas.NamaKelas
	)

	if err := row.Scan(
		&itemsTingkatKelas.IdTingkatKelas,
		&itemsTingkatKelas.TingkatKelas,
		&itemsNamaKelas.IdNamaKelas,
		&itemsNamaKelas.NamaKelas,
	); err != nil {
		return kelas.KelasData{}, err
	}

	itemsNamaKelas.IdTingkatKelas = itemsTingkatKelas.IdTingkatKelas

	return kelas.KelasData{
		ItemsTingkatKelas: itemsTingkatKelas,
		ItemsNamaKelas:    itemsNamaKelas,
	}, nil
}
