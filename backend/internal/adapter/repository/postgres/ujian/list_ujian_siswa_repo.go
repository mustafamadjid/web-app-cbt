package ujianrepo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

func (r *SiswaUjianRepo) buildListUjianSiswaQuery(idSiswa int, filter query.ListUjianFilter) (string, []any) {
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
			ru.nama_ruangan,
			u.deskripsi_ujian,
			ju.token,
			u.acak_soal
		FROM peserta_ujian pu
		JOIN jadwal_ujian ju
			ON ju.id_jadwal_ujian = pu.id_jadwal_ujian
		JOIN ujian u
			ON u.id_ujian = ju.id_ujian
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

	args = append(args, idSiswa)
	where = append(where, fmt.Sprintf("pu.id_siswa = $%d", len(args)))

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("(u.nama_ujian ILIKE $%d OR COALESCE(nk.nama_kelas, '') ILIKE $%d OR p.nama_lengkap ILIKE $%d OR s.nama_sesi ILIKE $%d OR ru.nama_ruangan ILIKE $%d OR mp.nama_mapel ILIKE $%d)", idx, idx, idx, idx, idx, idx))
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
	} else if filter.TingkatKelas != nil {
		args = append(args, *filter.TingkatKelas)
		where = append(where, fmt.Sprintf("k.tingkat_kelas = $%d", len(args)))
	}

	if filter.RuangUjian != nil {
		args = append(args, *filter.RuangUjian)
		where = append(where, fmt.Sprintf("ju.id_ruangan = $%d", len(args)))
	}

	if filter.IDMapel != nil {
		args = append(args, *filter.IDMapel)
		where = append(where, fmt.Sprintf("bs.id_mapel = $%d", len(args)))
	}

	switch filter.KategoriUjian {
	case query.MENDATANG:
		where = append(where, "ju.waktu_mulai > NOW()")
	case query.BERLANGSUNG:
		where = append(where, "ju.waktu_mulai <= NOW() AND ju.waktu_selesai >= NOW()")
	case query.SELESAI:
		where = append(where, "ju.waktu_selesai < NOW()")
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

func (r *SiswaUjianRepo) ListUjianSiswa(ctx context.Context, idSiswa int, filter query.ListUjianFilter) ([]ujian.ListUjian, error) {
	queryText, args := r.buildListUjianSiswaQuery(idSiswa, filter)

	rows, err := r.q.Query(ctx, queryText, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get list ujian siswa", "layer", "repo.db", "op", "siswa_ujian.list", "err", err)
		return nil, err
	}
	defer rows.Close()

	var results []ujian.ListUjian
	for rows.Next() {
		var (
			item           ujian.ListUjian
			idNamaKelas    sql.NullInt64
			namaKelas      sql.NullString
			deskripsiUjian sql.NullString
		)

		if err := rows.Scan(
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
			&item.StatusUjian,
			&item.IdPengawas,
			&item.NamaPengawas,
			&item.PengawasUsername,
			&item.IdSesi,
			&item.NamaSesi,
			&item.IdRuangan,
			&item.NamaRuangan,
			&deskripsiUjian,
			&item.Token,
			&item.AcakSoal,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning list ujian siswa", "layer", "repo.db", "op", "siswa_ujian.list.scan", "err", err)
			return nil, err
		}

		item.IdNamaKelas = nullInt64ToUjianIDPtr(idNamaKelas)
		item.NamaKelas = nullStringToPtr(namaKelas)
		item.DeskripsiUjian = nullStringToPtr(deskripsiUjian)

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating list ujian siswa", "layer", "repo.db", "op", "siswa_ujian.list.iter", "err", err)
		return nil, err
	}

	return results, nil
}


