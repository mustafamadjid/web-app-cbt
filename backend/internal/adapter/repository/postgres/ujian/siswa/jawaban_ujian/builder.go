package jawabanujian_repo

import (
	"context"
	"database/sql"
	"strings"

	"github.com/jackc/pgx/v5"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func splitSaveJawabanItems(jawaban []ujian.JawabanUjian) ([]ujian.JawabanUjian, []int64) {
	upsertItems := make([]ujian.JawabanUjian, 0, len(jawaban))
	clearSoalIDs := make([]int64, 0, len(jawaban))

	for _, item := range jawaban {
		normalizedItem := item
		if normalizedItem.JawabanEssay != nil {
			trimmedEssay := strings.TrimSpace(*normalizedItem.JawabanEssay)
			if trimmedEssay == "" {
				normalizedItem.JawabanEssay = nil
			} else {
				normalizedItem.JawabanEssay = &trimmedEssay
			}
		}

		if normalizedItem.IdPilihan == nil && normalizedItem.JawabanEssay == nil {
			clearSoalIDs = append(clearSoalIDs, int64(normalizedItem.IdSoal))
			continue
		}

		upsertItems = append(upsertItems, normalizedItem)
	}

	return upsertItems, clearSoalIDs
}

func (r *JawabanUjianRepo) scanJawabanUjianRows(ctx context.Context, op string, rows pgx.Rows) ([]ujian.JawabanUjian, error) {
	var jawaban []ujian.JawabanUjian
	for rows.Next() {
		var (
			itemJawaban  ujian.JawabanUjian
			idPilihan    sql.NullInt64
			jawabanEssay sql.NullString
			waktuJawab   sql.NullTime
		)

		if err := rows.Scan(
			&itemJawaban.IdJawaban,
			&itemJawaban.IdSoal,
			&idPilihan,
			&jawabanEssay,
			&waktuJawab,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scan jawaban ujian by attempt id", "layer", "adapter.repository", "op", op, "err", err)
			return nil, err
		}

		itemJawaban.IdPilihan = nullInt64ToUjianIDPtr(idPilihan)
		itemJawaban.JawabanEssay = nullStringToPtr(jawabanEssay)
		if waktuJawab.Valid {
			itemJawaban.WaktuJawab = &waktuJawab.Time
		}

		jawaban = append(jawaban, itemJawaban)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get jawaban ujian by attempt id", "layer", "adapter.repository", "op", op, "err", err)
		return nil, err
	}

	return jawaban, nil
}

func nullInt64ToUjianIDPtr(v sql.NullInt64) *ujian.ID {
	if !v.Valid {
		return nil
	}
	id := ujian.ID(v.Int64)
	return &id
}

func nullStringToPtr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	s := v.String
	return &s
}
