package ujianlistrepo

import (
	"context"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func (r *ListSoalUjianRepo) GetOpsiPilihanGandaByBankSoal(ctx context.Context, idBankSoalAktif int) ([]ujian.OpsiPilganUjian, error) {
	const query = `
		SELECT
			op.id_pilihan_ganda,
			op.id_soal,
			op.isi_pilihan,
			op.is_benar
		FROM opsi_pilihan_ganda op
		JOIN isi_soal s
			ON s.id_soal = op.id_soal
		WHERE s.id_bank_soal_version = $1
		ORDER BY op.id_soal, op.id_pilihan_ganda
	`

	rows, err := r.q.Query(ctx, query, idBankSoalAktif)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get opsi pilihan ganda by bank soal", "layer", "repo.db", "op", "ujian.list_opsi_pilgan", "err", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanOpsiPilganRows(ctx, "ujian.list_opsi_pilgan", rows)
}
