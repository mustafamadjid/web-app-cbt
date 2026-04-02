package gradingrepo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

type listGradingUjianScanner interface {
	Scan(dest ...any) error
}

func (r *ListGradingRepo) buildListUjianEssayUngradedQuery(filter query.ListUjianEssayUngradedFilter) (string, []any) {
	baseQuery := `
		SELECT
			u.id_ujian,
			u.id_bank_soal,
			u.id_guru,
			u.nama_ujian,
			pg.username AS pembuat_username,
			u.id_kelas,
			u.id_nama_kelas,
			k.tingkat_kelas,
			nk.nama_kelas,
			ju.id_jadwal_ujian,
			ju.tanggal_ujian,
			ju.waktu_mulai,
			ju.waktu_selesai,
			ju.status_ujian,
			ju.id_pengawas,
			p.nama_lengkap,
			p.username AS pengawas_username,
			ju.id_sesi,
			s.nama_sesi,
			ju.id_ruangan,
			ru.nama_ruangan
		FROM ujian u
		JOIN jadwal_ujian ju
			ON ju.id_ujian = u.id_ujian
		JOIN bank_soal bs
			ON bs.id_bank_soal = u.id_bank_soal
		JOIN mata_pelajaran mp
			ON mp.id_mapel = bs.id_mapel
		JOIN kelas k
			ON k.id_kelas = u.id_kelas
		LEFT JOIN nama_kelas nk
			ON nk.id_nama_kelas = u.id_nama_kelas
		JOIN pengguna pg
			ON pg.id_pengguna = u.id_guru
		JOIN pengguna p
			ON p.id_pengguna = ju.id_pengawas
		JOIN sesi_ujian s
			ON s.id_sesi = ju.id_sesi
		JOIN ruang_ujian ru
			ON ru.id_ruangan = ju.id_ruangan
	`

	where := make([]string, 0, 10)
	args := make([]any, 0, 10)

	where = append(where, `
		EXISTS (
			SELECT 1
			FROM peserta_ujian pu
			JOIN attempt_ujian au
				ON au.id_peserta_ujian = pu.id_peserta_ujian
			JOIN hasil_ujian hu
				ON hu.id_attempt = au.id_attempt
			WHERE pu.id_jadwal_ujian = ju.id_jadwal_ujian
				AND hu.essay_graded = FALSE
		)
	`)
	where = append(where, `
		EXISTS (
			SELECT 1
			FROM isi_soal is2
			WHERE is2.id_bank_soal_version = bs.id_bank_soal_version_aktif
				AND is2.tipe_soal = 'essay'
		)
	`)

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("(u.nama_ujian ILIKE $%d OR mp.nama_mapel ILIKE $%d OR COALESCE(nk.nama_kelas, '') ILIKE $%d OR p.nama_lengkap ILIKE $%d OR s.nama_sesi ILIKE $%d)", idx, idx, idx, idx, idx))
	}

	if filter.TanggalUjian != nil {
		args = append(args, *filter.TanggalUjian)
		where = append(where, fmt.Sprintf("ju.tanggal_ujian = $%d::date", len(args)))
	}

	if filter.Tahun != nil {
		args = append(args, *filter.Tahun)
		where = append(where, fmt.Sprintf("EXTRACT(YEAR FROM ju.tanggal_ujian)::int = $%d::int", len(args)))
	}

	if filter.Bulan != nil {
		args = append(args, *filter.Bulan)
		where = append(where, fmt.Sprintf("EXTRACT(MONTH FROM ju.tanggal_ujian)::int = $%d::int", len(args)))
	}

	if filter.TingkatKelasID != nil {
		args = append(args, *filter.TingkatKelasID)
		where = append(where, fmt.Sprintf("u.id_kelas = $%d", len(args)))
	}

	if filter.NamaKelasID != nil {
		args = append(args, *filter.NamaKelasID)
		where = append(where, fmt.Sprintf("u.id_nama_kelas = $%d", len(args)))
	}

	if filter.MapelID != nil {
		args = append(args, *filter.MapelID)
		where = append(where, fmt.Sprintf("bs.id_mapel = $%d", len(args)))
	}

	if filter.SesiID != nil {
		args = append(args, *filter.SesiID)
		where = append(where, fmt.Sprintf("ju.id_sesi = $%d", len(args)))
	}

	baseQuery = fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(where, " AND "))
	baseQuery = fmt.Sprintf("%s ORDER BY ju.tanggal_ujian ASC, ju.waktu_mulai ASC, ju.id_jadwal_ujian ASC", baseQuery)

	if filter.Limit > 0 {
		args = append(args, filter.Limit)
		limitIdx := len(args)
		args = append(args, filter.Offset)
		offsetIdx := len(args)
		baseQuery = fmt.Sprintf("%s LIMIT $%d OFFSET $%d", baseQuery, limitIdx, offsetIdx)
	}

	return baseQuery, args
}

func scanListUjianEssayUngradedRow(row listGradingUjianScanner) (ujian.ListUjian, error) {
	var (
		item        ujian.ListUjian
		idNamaKelas sql.NullInt64
		namaKelas   sql.NullString
		statusUjian sql.NullString
	)

	if err := row.Scan(
		&item.IdUjian,
		&item.IdBankSoal,
		&item.IdGuru,
		&item.NamaUjian,
		&item.PembuatUsername,
		&item.IdKelas,
		&idNamaKelas,
		&item.TingkatKelas,
		&namaKelas,
		&item.IdJadwalUjian,
		&item.TanggalUjian,
		&item.WaktuMulai,
		&item.WaktuSelesai,
		&statusUjian,
		&item.IdPengawas,
		&item.NamaPengawas,
		&item.PengawasUsername,
		&item.IdSesi,
		&item.NamaSesi,
		&item.IdRuangan,
		&item.NamaRuangan,
	); err != nil {
		return ujian.ListUjian{}, err
	}

	item.IdNamaKelas = nullInt64ToUjianIDPtr(idNamaKelas)
	item.NamaKelas = nullStringToPtr(namaKelas)
	item.StatusUjian = nullStringToStatusUjianPtr(statusUjian)

	return item, nil
}

func (r *ListGradingRepo) scanListUjianEssayUngradedRows(ctx context.Context, op string, rows pgx.Rows) ([]ujian.ListUjian, error) {
	var results []ujian.ListUjian
	for rows.Next() {
		item, err := scanListUjianEssayUngradedRow(rows)
		if err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning list essay ungraded ujian", "layer", "repo.db", "op", op, "err", err)
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating list essay ungraded ujian", "layer", "repo.db", "op", op, "err", err)
		return nil, err
	}

	return results, nil
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

func nullStringToStatusUjianPtr(v sql.NullString) *ujian.StatusUjian {
	if !v.Valid {
		return nil
	}

	status := ujian.StatusUjian(v.String)
	return &status
}
