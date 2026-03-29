package banksoalrepo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
)

const bankSoalSelectColumns = `
	SELECT
		b.id_bank_soal,
		b.id_mapel,
		b.id_kelas,
		b.id_pengguna,
		b.nama_bank_soal,
		b.deskripsi,
		b.materi,
		b.created_at,
		(b.id_bank_soal_version_aktif IS NOT NULL) AS soal_uploaded,
		k.tingkat_kelas,
		m.nama_mapel,
		p.nama_lengkap
	FROM bank_soal b
	JOIN kelas k ON b.id_kelas = k.id_kelas
	JOIN mata_pelajaran m ON b.id_mapel = m.id_mapel
	JOIN pengguna p ON b.id_pengguna = p.id_pengguna
`

func (r *BankSoalRepo) buildListBankSoalQuery(filter query.BankSoalFilter, uploadedOnly bool) (string, []any) {
	baseQuery := bankSoalSelectColumns
	where := make([]string, 0, 4)
	args := make([]any, 0, 4)

	if uploadedOnly {
		where = append(where, "b.id_bank_soal_version_aktif IS NOT NULL")
	}

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("(b.nama_bank_soal ILIKE $%d OR b.deskripsi ILIKE $%d OR b.materi ILIKE $%d)", idx, idx, idx))
	}

	if filter.TingkatKelas != nil {
		args = append(args, *filter.TingkatKelas)
		where = append(where, fmt.Sprintf("b.id_kelas = $%d", len(args)))
	}

	if filter.Mapel != nil {
		args = append(args, *filter.Mapel)
		where = append(where, fmt.Sprintf("b.id_mapel = $%d", len(args)))
	}

	if len(where) > 0 {
		baseQuery = fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(where, " AND "))
	}

	baseQuery = fmt.Sprintf("%s ORDER BY b.created_at ASC", baseQuery)

	if filter.Limit > 0 {
		args = append(args, filter.Limit)
		limitIndex := len(args)
		args = append(args, filter.Offset)
		offsetIndex := len(args)
		baseQuery = fmt.Sprintf("%s LIMIT $%d OFFSET $%d", baseQuery, limitIndex, offsetIndex)
	}

	return baseQuery, args
}

func (r *BankSoalRepo) scanBankSoalRows(ctx context.Context, op string, rows pgx.Rows) ([]bank_soal.BankSoal, error) {
	var results []bank_soal.BankSoal
	for rows.Next() {
		var item bank_soal.BankSoal
		if err := rows.Scan(
			&item.IdBankSoal,
			&item.IdMapel,
			&item.IdKelas,
			&item.IdPengguna,
			&item.NamaBankSoal,
			&item.Deskripsi,
			&item.Materi,
			&item.CreatedAt,
			&item.SoalUploaded,
			&item.TingkatKelas,
			&item.Mapel,
			&item.GuruPembuat,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scan bank soal", "layer", "repo.db", "op", op, "err", err)
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating bank soal rows", "layer", "repo.db", "op", op, "err", err)
		return nil, err
	}

	return results, nil
}