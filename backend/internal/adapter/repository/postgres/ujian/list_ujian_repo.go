package ujianrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

type ListUjianRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewListUjianRepo(q pg.Executor, logger corelog.Logger) *ListUjianRepo {
	return &ListUjianRepo{q: q, logger: logger}
}

func (r *ListUjianRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *ListUjianRepo) GetAllUjian(ctx context.Context, filter query.ListUjianFilter) ([]ujian.ListUjian, error) {
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

	rows, err := r.q.Query(ctx, baseQuery, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get list ujian", "layer", "repo.db", "op", "ujian.get_all", "err", err)
		return nil, err
	}
	defer rows.Close()

	var results []ujian.ListUjian
	for rows.Next() {
		var (
			item        ujian.ListUjian
			idNamaKls   sql.NullInt64
			namaKelas   sql.NullString
			statusUjian sql.NullString
		)

		if err := rows.Scan(
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
			r.loggerFor(ctx).Error(ctx, "failed scanning list ujian", "layer", "repo.db", "op", "ujian.get_all.scan", "err", err)
			return nil, err
		}

		item.IdNamaKelas = nullInt64ToUjianIDPtr(idNamaKls)
		item.NamaKelas = nullStringToPtr(namaKelas)
		item.StatusUjian = nullStringToStatusUjianPtr(statusUjian)

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating list ujian", "layer", "repo.db", "op", "ujian.get_all.iter", "err", err)
		return nil, err
	}

	return results, nil
}

func (r *ListUjianRepo) GetUjianById(ctx context.Context, id ujian.ID) (ujian.ListUjian, error) {
	query := `
		SELECT
			u.id_ujian,
			u.id_bank_soal,
			u.id_guru,
			u.nama_ujian,
			u.deskripsi_ujian,
			pg.username AS pembuat_username,
			u.id_kelas,
			u.id_nama_kelas,
			u.acak_soal,
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
			ju.token,
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
		WHERE u.id_ujian = $1
	`
	var (
		item           ujian.ListUjian
		idNamaKls      sql.NullInt64
		namaKelas      sql.NullString
		deskripsiUjian sql.NullString
		statusUjian    sql.NullString
	)

	err := r.q.QueryRow(ctx, query, id).Scan(
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
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return ujian.ListUjian{}, coreerror.ErrNotFound
	}
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed scanning list ujian", "layer", "repo.db", "op", "ujian.get_by_id.scan", "err", err)
		return ujian.ListUjian{}, err
	}

	item.IdNamaKelas = nullInt64ToUjianIDPtr(idNamaKls)
	item.NamaKelas = nullStringToPtr(namaKelas)
	item.DeskripsiUjian = nullStringToPtr(deskripsiUjian)
	item.StatusUjian = nullStringToStatusUjianPtr(statusUjian)

	return item, nil
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
