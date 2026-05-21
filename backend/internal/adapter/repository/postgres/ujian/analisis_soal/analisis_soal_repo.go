package analisissoal_repo

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	pgshared "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/shared"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type AnalisisSoalRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewAnalisisSoalRepo(q pg.Executor, logger corelog.Logger) *AnalisisSoalRepo {
	return &AnalisisSoalRepo{q: q, logger: logger}
}

func (r *AnalisisSoalRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *AnalisisSoalRepo) GetListAnalisisSoal(ctx context.Context, idJadwalUjian ujian.ID) ([]ujian.AnalisisSoal, error) {
	const query = `
		SELECT
			s.id_soal,
			s.id_bank_soal_version,
			s.tipe_soal,
			s.pertanyaan,
			s.pertanyaan_content,
			s.gambar,
			s.bobot_soal,
			s.no_urut_soal,
			COALESCE(ss.jumlah_jawaban_benar, 0) AS jumlah_jawaban_benar,
			COALESCE(ss.jumlah_jawaban_salah, 0) AS jumlah_jawaban_salah,
			op.id_pilihan_ganda,
			op.isi_pilihan,
			op.isi_pilihan_content,
			op.is_benar
		FROM jadwal_ujian ju
		JOIN ujian u
			ON u.id_ujian = ju.id_ujian
		JOIN bank_soal bs
			ON bs.id_bank_soal = u.id_bank_soal
		JOIN isi_soal s
			ON s.id_bank_soal_version = bs.id_bank_soal_version_aktif
		LEFT JOIN statistik_soal ss
			ON ss.id_ujian = u.id_ujian
			AND ss.id_soal = s.id_soal
		LEFT JOIN opsi_pilihan_ganda op
			ON op.id_soal = s.id_soal
		WHERE ju.id_jadwal_ujian = $1
		ORDER BY s.no_urut_soal ASC, op.id_pilihan_ganda ASC
	`

	rows, err := r.q.Query(ctx, query, idJadwalUjian)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed list analisis soal", "layer", "repo.db", "op", "ujian.analisis_soal.list", "jadwal_id", idJadwalUjian, "err", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanAnalisisSoalRows(ctx, "ujian.analisis_soal.list", rows)
}

func (r *AnalisisSoalRepo) scanAnalisisSoalRows(ctx context.Context, op string, rows pgx.Rows) ([]ujian.AnalisisSoal, error) {
	var (
		results        []ujian.AnalisisSoal
		orderedSoalIDs []ujian.ID
		itemsBySoalID  = make(map[ujian.ID]*ujian.AnalisisSoal)
	)

	for rows.Next() {
		var (
			idSoal             ujian.ID
			idBankSoalVersion  ujian.ID
			tipeSoal           string
			pertanyaan         string
			pertanyaanContent  []byte
			gambar             sql.NullString
			bobotSoal          float64
			noUrutSoal         int
			jumlahJawabanBenar int
			jumlahJawabanSalah int
			idPilgan           sql.NullInt64
			isiPilgan          sql.NullString
			isiPilihanContent  []byte
			isBenar            sql.NullBool
		)

		if err := rows.Scan(
			&idSoal,
			&idBankSoalVersion,
			&tipeSoal,
			&pertanyaan,
			&pertanyaanContent,
			&gambar,
			&bobotSoal,
			&noUrutSoal,
			&jumlahJawabanBenar,
			&jumlahJawabanSalah,
			&idPilgan,
			&isiPilgan,
			&isiPilihanContent,
			&isBenar,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scan analisis soal", "layer", "repo.db", "op", op, "err", err)
			return nil, err
		}

		item, exists := itemsBySoalID[idSoal]
		if !exists {
			contentValue, err := pgshared.UnmarshalRichContent(pertanyaanContent)
			if err != nil {
				r.loggerFor(ctx).Error(ctx, "failed unmarshal pertanyaan content", "layer", "repo.db", "op", op, "err", err)
				return nil, err
			}

			item = &ujian.AnalisisSoal{
				Soal: ujian.SoalUjianSiswa{
					IdSoal:            idSoal,
					IdBankSoalVersion: idBankSoalVersion,
					TipeSoal:          tipeSoal,
					Pertanyaan:        pertanyaan,
					PertanyaanContent: contentValue,
					BobotSoal:         bobotSoal,
					NoUrutSoal:        noUrutSoal,
				},
				JumlahJawabanBenar: jumlahJawabanBenar,
				JumlahJawabanSalah: jumlahJawabanSalah,
			}
			if gambar.Valid {
				item.Soal.Gambar = gambar.String
			}

			itemsBySoalID[idSoal] = item
			orderedSoalIDs = append(orderedSoalIDs, idSoal)
		}

		if idPilgan.Valid {
			optionContent, err := pgshared.UnmarshalRichContent(isiPilihanContent)
			if err != nil {
				r.loggerFor(ctx).Error(ctx, "failed unmarshal opsi content", "layer", "repo.db", "op", op, "err", err)
				return nil, err
			}

			opsi := ujian.OpsiPilganUjian{
				IdPilihanGanda:    ujian.ID(idPilgan.Int64),
				IdSoal:            idSoal,
				IsiPilihanContent: optionContent,
			}
			if isiPilgan.Valid {
				opsi.IsiPilihan = isiPilgan.String
			}
			if isBenar.Valid {
				opsi.IsBenar = isBenar.Bool
			}

			item.Soal.OpsiJawaban = append(item.Soal.OpsiJawaban, opsi)
		}
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterate analisis soal", "layer", "repo.db", "op", op, "err", err)
		return nil, err
	}

	for _, idSoal := range orderedSoalIDs {
		results = append(results, *itemsBySoalID[idSoal])
	}

	return results, nil
}
