package jawabanujian_repo

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func (r *JawabanUjianRepo) ListHasilJawabanUjianByIdAttempt(ctx context.Context, idAttempt ujian.ID) ([]ujian.HasilJawabanUjian, error) {
	const query = `
		SELECT
			s.id_soal,
			s.id_bank_soal_version,
			s.tipe_soal,
			s.pertanyaan,
			s.gambar,
			s.bobot_soal,
			s.no_urut_soal,
			op.id_pilihan_ganda,
			op.isi_pilihan,
			op.is_benar,
			jus.id_jawaban,
			jus.id_attempt,
			jus.id_pilihan,
			jus.jawaban_essay,
			jus.waktu_jawab,
			jus.essay_is_benar
		FROM attempt_ujian au
		JOIN peserta_ujian pu
			ON pu.id_peserta_ujian = au.id_peserta_ujian
		JOIN jadwal_ujian ju
			ON ju.id_jadwal_ujian = pu.id_jadwal_ujian
		JOIN ujian u
			ON u.id_ujian = ju.id_ujian
		JOIN bank_soal bs
			ON bs.id_bank_soal = u.id_bank_soal
		JOIN isi_soal s
			ON s.id_bank_soal_version = bs.id_bank_soal_version_aktif
		LEFT JOIN opsi_pilihan_ganda op
			ON op.id_soal = s.id_soal
		LEFT JOIN jawaban_ujian_siswa jus
			ON jus.id_attempt = au.id_attempt
			AND jus.id_soal = s.id_soal
		WHERE au.id_attempt = $1
		ORDER BY s.no_urut_soal ASC, op.id_pilihan_ganda ASC
	`

	rows, err := r.q.Query(ctx, query, idAttempt)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed list hasil jawaban ujian by attempt id", "layer", "adapter.repository", "op", "ujian.jawaban.list_hasil_by_attempt_id", "err", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanHasilJawabanRows(ctx, "ujian.jawaban.list_hasil_by_attempt_id", rows)
}

func (r *JawabanUjianRepo) scanHasilJawabanRows(ctx context.Context, op string, rows pgx.Rows) ([]ujian.HasilJawabanUjian, error) {
	var (
		results        []ujian.HasilJawabanUjian
		orderedSoalIDs []ujian.ID
		itemsBySoalID  = make(map[ujian.ID]*ujian.HasilJawabanUjian)
	)

	for rows.Next() {
		var (
			idSoal            ujian.ID
			idBankSoalVersion ujian.ID
			tipeSoal          string
			pertanyaan        string
			gambar            sql.NullString
			bobotSoal         float64
			noUrutSoal        int
			idPilgan          sql.NullInt64
			isiPilgan         sql.NullString
			isBenar           sql.NullBool
			idJawaban         sql.NullInt64
			idAttempt         sql.NullInt64
			idPilihan         sql.NullInt64
			jawabanEssay      sql.NullString
			waktuJawab        sql.NullTime
			essayIsBenar      sql.NullBool
		)

		if err := rows.Scan(
			&idSoal,
			&idBankSoalVersion,
			&tipeSoal,
			&pertanyaan,
			&gambar,
			&bobotSoal,
			&noUrutSoal,
			&idPilgan,
			&isiPilgan,
			&isBenar,
			&idJawaban,
			&idAttempt,
			&idPilihan,
			&jawabanEssay,
			&waktuJawab,
			&essayIsBenar,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scan hasil jawaban ujian by attempt id", "layer", "adapter.repository", "op", op, "err", err)
			return nil, err
		}

		item, exists := itemsBySoalID[idSoal]
		if !exists {
			item = &ujian.HasilJawabanUjian{
				SoalUjianSiswa: ujian.SoalUjianSiswa{
					IdSoal:            idSoal,
					IdBankSoalVersion: idBankSoalVersion,
					TipeSoal:          tipeSoal,
					Pertanyaan:        pertanyaan,
					BobotSoal:         bobotSoal,
					NoUrutSoal:        noUrutSoal,
				},
				JawabanSiswa: ujian.JawabanUjian{
					IdSoal:       idSoal,
					IdPilihan:    nullInt64ToUjianIDPtr(idPilihan),
					JawabanEssay: nullStringToPtr(jawabanEssay),
				},
			}

			if gambar.Valid {
				item.SoalUjianSiswa.Gambar = gambar.String
			}
			if idJawaban.Valid {
				item.JawabanSiswa.IdJawaban = ujian.ID(idJawaban.Int64)
			}
			if idAttempt.Valid {
				item.JawabanSiswa.IdAttempt = ujian.ID(idAttempt.Int64)
			}
			if waktuJawab.Valid {
				item.JawabanSiswa.WaktuJawab = &waktuJawab.Time
			}
			if essayIsBenar.Valid {
				v := essayIsBenar.Bool
				item.JawabanSiswa.EssayIsBenar = &v
			}

			itemsBySoalID[idSoal] = item
			orderedSoalIDs = append(orderedSoalIDs, idSoal)
		}

		if idPilgan.Valid {
			opsi := ujian.OpsiPilganUjian{
				IdPilihanGanda: ujian.ID(idPilgan.Int64),
				IdSoal:         idSoal,
			}
			if isiPilgan.Valid {
				opsi.IsiPilihan = isiPilgan.String
			}
			if isBenar.Valid {
				opsi.IsBenar = isBenar.Bool
			}
			item.SoalUjianSiswa.OpsiJawaban = append(item.SoalUjianSiswa.OpsiJawaban, opsi)
		}
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterate hasil jawaban ujian by attempt id", "layer", "adapter.repository", "op", op, "err", err)
		return nil, err
	}

	for _, idSoal := range orderedSoalIDs {
		results = append(results, *itemsBySoalID[idSoal])
	}

	return results, nil
}
