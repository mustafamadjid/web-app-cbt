package ujianlistrepo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

type listUjianScanner interface {
	Scan(dest ...any) error
}

func (r *ListUjianRepo) buildListUjianQuery(filter query.ListUjianFilter) (string, []any) {
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

	where := make([]string, 0, 5)
	args := make([]any, 0, 7)

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("(u.nama_ujian ILIKE $%d OR COALESCE(nk.nama_kelas, '') ILIKE $%d OR p.nama_lengkap ILIKE $%d OR s.nama_sesi ILIKE $%d OR ru.nama_ruangan ILIKE $%d)", idx, idx, idx, idx, idx))
	}

	if filter.TanggalUjian != nil {
		args = append(args, *filter.TanggalUjian)
		where = append(where, fmt.Sprintf("ju.tanggal_ujian = $%d::date", len(args)))
	}

	if filter.Tahun != nil {
		args = append(args, *filter.Tahun)
		where = append(where, fmt.Sprintf("EXTRACT(YEAR FROM ju.tanggal_ujian)::text = $%d", len(args)))
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

	switch filter.KategoriUjian {
	case query.MENDATANG:
		where = append(where, "ju.waktu_mulai > NOW()")
	case query.BERLANGSUNG:
		where = append(where, "ju.waktu_mulai <= NOW() AND ju.waktu_selesai >= NOW()")
	case query.SELESAI:
		where = append(where, "ju.waktu_selesai < NOW()")
	case query.DIBATALKAN:
		where = append(where, "ju.status_ujian = 'DIBATALKAN'")
	}

	if len(where) > 0 {
		baseQuery = fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(where, " AND "))
	}

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

func scanListUjianSummaryRow(row listUjianScanner) (ujian.ListUjian, error) {
	var (
		item        ujian.ListUjian
		idNamaKls   sql.NullInt64
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
		&idNamaKls,
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

	item.IdNamaKelas = nullInt64ToUjianIDPtr(idNamaKls)
	item.NamaKelas = nullStringToPtr(namaKelas)
	item.StatusUjian = nullStringToStatusUjianPtr(statusUjian)

	return item, nil
}

func (r *ListUjianRepo) scanListUjianSummaryRows(ctx context.Context, op string, rows pgx.Rows) ([]ujian.ListUjian, error) {
	var results []ujian.ListUjian
	for rows.Next() {
		item, err := scanListUjianSummaryRow(rows)
		if err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning list ujian", "layer", "repo.db", "op", op, "err", err)
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating list ujian", "layer", "repo.db", "op", op, "err", err)
		return nil, err
	}

	return results, nil
}

func scanListUjianDetailRow(row listUjianScanner) (ujian.ListUjian, error) {
	var (
		item           ujian.ListUjian
		idNamaKls      sql.NullInt64
		namaKelas      sql.NullString
		deskripsiUjian sql.NullString
		statusUjian    sql.NullString
	)

	if err := row.Scan(
		&item.IdUjian,
		&item.IdBankSoal,
		&item.IdGuru,
		&item.NamaUjian,
		&deskripsiUjian,
		&item.PembuatUsername,
		&item.IdKelas,
		&idNamaKls,
		&item.AcakSoal,
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
		&item.Token,
		&item.NamaSesi,
		&item.IdRuangan,
		&item.NamaRuangan,
	); err != nil {
		return ujian.ListUjian{}, err
	}

	item.IdNamaKelas = nullInt64ToUjianIDPtr(idNamaKls)
	item.NamaKelas = nullStringToPtr(namaKelas)
	item.DeskripsiUjian = nullStringToPtr(deskripsiUjian)
	item.StatusUjian = nullStringToStatusUjianPtr(statusUjian)

	return item, nil
}

func (r *ListSoalUjianRepo) scanSoalUjianRows(ctx context.Context, op string, rows pgx.Rows) ([]ujian.SoalUjianSiswa, error) {
	var (
		itemSoalResult []ujian.SoalUjianSiswa
		orderedSoalIDs []ujian.ID
		itemsBySoalID  = make(map[ujian.ID]*ujian.SoalUjianSiswa)
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
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning soal ujian", "op", op, "err", err)
			return nil, err
		}

		item, exists := itemsBySoalID[idSoal]
		if !exists {
			item = &ujian.SoalUjianSiswa{
				IdSoal:            idSoal,
				IdBankSoalVersion: idBankSoalVersion,
				TipeSoal:          tipeSoal,
				Pertanyaan:        pertanyaan,
				BobotSoal:         bobotSoal,
				NoUrutSoal:        noUrutSoal,
			}
			if gambar.Valid {
				item.Gambar = gambar.String
			}
			itemsBySoalID[idSoal] = item
			orderedSoalIDs = append(orderedSoalIDs, idSoal)
		}

		if idPilgan.Valid {
			opsi := ujian.OpsiPilganUjian{
				IdPilihanGanda: ujian.ID(idPilgan.Int64),
				IdSoal:         item.IdSoal,
			}
			if isiPilgan.Valid {
				opsi.IsiPilihan = isiPilgan.String
			}
			if isBenar.Valid {
				opsi.IsBenar = isBenar.Bool
			}
			item.OpsiJawaban = append(item.OpsiJawaban, opsi)
		}
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating soal ujian rows", "op", op, "err", err)
		return nil, err
	}

	for _, idSoal := range orderedSoalIDs {
		itemSoalResult = append(itemSoalResult, *itemsBySoalID[idSoal])
	}

	return itemSoalResult, nil
}

func (r *ListSoalUjianRepo) scanSoalUjianSiswaRows(ctx context.Context, op string, rows pgx.Rows) ([]ujian.SoalUjianSiswa, error) {
	var (
		itemSoalResult []ujian.SoalUjianSiswa
		orderedSoalIDs []ujian.ID
		itemsBySoalID  = make(map[ujian.ID]*ujian.SoalUjianSiswa)
	)

	for rows.Next() {
		var (
			idSoal     ujian.ID
			tipeSoal   string
			pertanyaan string
			gambar     sql.NullString
			bobotSoal  float64
			noUrutSoal int
			idPilgan   sql.NullInt64
			isiPilgan  sql.NullString
		)

		if err := rows.Scan(
			&idSoal,
			&tipeSoal,
			&pertanyaan,
			&gambar,
			&bobotSoal,
			&noUrutSoal,
			&idPilgan,
			&isiPilgan,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning soal ujian", "op", op, "err", err)
			return nil, err
		}

		item, exists := itemsBySoalID[idSoal]
		if !exists {
			item = &ujian.SoalUjianSiswa{
				IdSoal:     idSoal,
				TipeSoal:   tipeSoal,
				Pertanyaan: pertanyaan,
				BobotSoal:  bobotSoal,
				NoUrutSoal: noUrutSoal,
			}
			if gambar.Valid {
				item.Gambar = gambar.String
			}
			itemsBySoalID[idSoal] = item
			orderedSoalIDs = append(orderedSoalIDs, idSoal)
		}

		if idPilgan.Valid {
			opsi := ujian.OpsiPilganUjian{
				IdPilihanGanda: ujian.ID(idPilgan.Int64),
			}
			if isiPilgan.Valid {
				opsi.IsiPilihan = isiPilgan.String
			}
			item.OpsiJawaban = append(item.OpsiJawaban, opsi)
		}
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating soal ujian rows", "op", op, "err", err)
		return nil, err
	}

	for _, idSoal := range orderedSoalIDs {
		itemSoalResult = append(itemSoalResult, *itemsBySoalID[idSoal])
	}

	return itemSoalResult, nil
}

func (r *ListSoalUjianRepo) scanOpsiPilganRows(ctx context.Context, op string, rows pgx.Rows) ([]ujian.OpsiPilganUjian, error) {
	var items []ujian.OpsiPilganUjian
	for rows.Next() {
		var item ujian.OpsiPilganUjian
		if err := rows.Scan(
			&item.IdPilihanGanda,
			&item.IdSoal,
			&item.IsiPilihan,
			&item.IsBenar,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning opsi pilihan ganda by bank soal", "layer", "repo.db", "op", op, "err", err)
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating opsi pilihan ganda by bank soal", "layer", "repo.db", "op", op, "err", err)
		return nil, err
	}

	return items, nil
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
