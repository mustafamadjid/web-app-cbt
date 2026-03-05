package ujianrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

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
			item      ujian.ListUjian
			idNamaKls sql.NullInt64
			namaKelas sql.NullString
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
			&item.StatusUjian,
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

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating list ujian", "layer", "repo.db", "op", "ujian.get_all.iter", "err", err)
		return nil, err
	}

	return results, nil
}

func (r *ListUjianRepo) GetUjianById(ctx context.Context, id ujian.ID) (ujian.Ujian, error) {
	const queryText = `
		SELECT
			id_ujian,
			id_bank_soal,
			id_kelas,
			id_nama_kelas,
			id_guru,
			nama_ujian,
			deskripsi_ujian,
			acak_soal,
			created_at,
			updated_at
		FROM ujian
		WHERE id_ujian = $1
	`

	var (
		item        ujian.Ujian
		idNamaKelas sql.NullInt64
		deskripsi   sql.NullString
		updatedAt   sql.NullTime
	)

	if err := r.q.QueryRow(ctx, queryText, id).Scan(
		&item.IdUjian,
		&item.IdBankSoal,
		&item.IdKelas,
		&idNamaKelas,
		&item.IdGuru,
		&item.NamaUjian,
		&deskripsi,
		&item.AcakSoal,
		&item.CreatedAt,
		&updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ujian.Ujian{}, coreerror.ErrNotFound
		}
		r.loggerFor(ctx).Error(ctx, "failed get ujian by id", "layer", "repo.db", "op", "ujian.get_by_id", "err", err)
		return ujian.Ujian{}, err
	}

	item.IdNamaKelas = nullInt64ToUjianIDPtr(idNamaKelas)
	item.DeskripsiUjian = nullStringToPtr(deskripsi)
	item.UpdatedAt = nullTimeToPtr(updatedAt)

	return item, nil
}

func (r *ListUjianRepo) GetJadwalUjianById(ctx context.Context, id ujian.ID) (ujian.JadwalUjian, error) {
	const queryText = `
		SELECT
			id_jadwal_ujian,
			id_ujian,
			id_sesi,
			id_ruangan,
			id_pengawas,
			tanggal_ujian,
			waktu_mulai,
			waktu_selesai,
			token,
			status_ujian,
			created_at,
			updated_at
		FROM jadwal_ujian
		WHERE id_jadwal_ujian = $1
	`

	var (
		item       ujian.JadwalUjian
		idPengawas sql.NullInt64
		updatedAt  sql.NullTime
	)

	if err := r.q.QueryRow(ctx, queryText, id).Scan(
		&item.IdJadwalUjian,
		&item.IdUjian,
		&item.IdSesi,
		&item.IdRuangan,
		&idPengawas,
		&item.TanggalUjian,
		&item.WaktuMulai,
		&item.WaktuSelesai,
		&item.Token,
		&item.StatusUjian,
		&item.CreatedAt,
		&updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ujian.JadwalUjian{}, coreerror.ErrNotFound
		}
		r.loggerFor(ctx).Error(ctx, "failed get jadwal ujian by id", "layer", "repo.db", "op", "ujian.get_jadwal_by_id", "err", err)
		return ujian.JadwalUjian{}, err
	}

	if idPengawas.Valid {
		item.IdPengawas = ujian.ID(idPengawas.Int64)
	}
	item.UpdatedAt = nullTimeToPtr(updatedAt)

	return item, nil
}

func (r *ListUjianRepo) GetAllPesertaUjian(ctx context.Context, peserta ujian.PesertaUjian) ([]ujian.PesertaUjian, error) {
	baseQuery := `
		SELECT
			id_peserta_ujian,
			id_jadwal_ujian,
			id_siswa,
			waktu_mulai,
			waktu_submit,
			nilai_ujian,
			created_at,
			updated_at
		FROM peserta_ujian
	`

	where := make([]string, 0, 6)
	args := make([]any, 0, 6)

	if peserta.IdPesertaUjian > 0 {
		args = append(args, peserta.IdPesertaUjian)
		where = append(where, fmt.Sprintf("id_peserta_ujian = $%d", len(args)))
	}

	if peserta.IdJadwalUjian > 0 {
		args = append(args, peserta.IdJadwalUjian)
		where = append(where, fmt.Sprintf("id_jadwal_ujian = $%d", len(args)))
	}

	if peserta.IdSiswa > 0 {
		args = append(args, peserta.IdSiswa)
		where = append(where, fmt.Sprintf("id_siswa = $%d", len(args)))
	}

	if peserta.WaktuMulai != nil {
		args = append(args, *peserta.WaktuMulai)
		where = append(where, fmt.Sprintf("waktu_mulai = $%d", len(args)))
	}

	if peserta.WaktuSubmit != nil {
		args = append(args, *peserta.WaktuSubmit)
		where = append(where, fmt.Sprintf("waktu_submit = $%d", len(args)))
	}

	if peserta.NilaiUjian != nil {
		args = append(args, *peserta.NilaiUjian)
		where = append(where, fmt.Sprintf("nilai_ujian = $%d", len(args)))
	}

	if len(where) > 0 {
		baseQuery = fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(where, " AND "))
	}

	baseQuery = fmt.Sprintf("%s ORDER BY id_peserta_ujian ASC", baseQuery)

	rows, err := r.q.Query(ctx, baseQuery, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get peserta ujian", "layer", "repo.db", "op", "ujian.get_peserta", "err", err)
		return nil, err
	}
	defer rows.Close()

	var results []ujian.PesertaUjian
	for rows.Next() {
		var (
			item        ujian.PesertaUjian
			waktuMulai  sql.NullTime
			waktuSubmit sql.NullTime
			nilaiUjian  sql.NullFloat64
			updatedAt   sql.NullTime
		)

		if err := rows.Scan(
			&item.IdPesertaUjian,
			&item.IdJadwalUjian,
			&item.IdSiswa,
			&waktuMulai,
			&waktuSubmit,
			&nilaiUjian,
			&item.CreatedAt,
			&updatedAt,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning peserta ujian", "layer", "repo.db", "op", "ujian.get_peserta.scan", "err", err)
			return nil, err
		}

		item.WaktuMulai = nullTimeToPtr(waktuMulai)
		item.WaktuSubmit = nullTimeToPtr(waktuSubmit)
		item.NilaiUjian = nullFloat64ToPtr(nilaiUjian)
		item.UpdatedAt = nullTimeToPtr(updatedAt)

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating peserta ujian", "layer", "repo.db", "op", "ujian.get_peserta.iter", "err", err)
		return nil, err
	}

	return results, nil
}

func (r *ListUjianRepo) GetPesertaUjianBySiswa(ctx context.Context, idSiswa ujian.ID, peserta ujian.PesertaUjian) (ujian.PesertaUjian, error) {
	baseQuery := `
		SELECT
			id_peserta_ujian,
			id_jadwal_ujian,
			id_siswa,
			waktu_mulai,
			waktu_submit,
			nilai_ujian,
			created_at,
			updated_at
		FROM peserta_ujian
		WHERE id_siswa = $1
	`

	args := make([]any, 0, 6)
	args = append(args, idSiswa)

	if peserta.IdPesertaUjian > 0 {
		args = append(args, peserta.IdPesertaUjian)
		baseQuery = fmt.Sprintf("%s AND id_peserta_ujian = $%d", baseQuery, len(args))
	}

	if peserta.IdJadwalUjian > 0 {
		args = append(args, peserta.IdJadwalUjian)
		baseQuery = fmt.Sprintf("%s AND id_jadwal_ujian = $%d", baseQuery, len(args))
	}

	if peserta.WaktuMulai != nil {
		args = append(args, *peserta.WaktuMulai)
		baseQuery = fmt.Sprintf("%s AND waktu_mulai = $%d", baseQuery, len(args))
	}

	if peserta.WaktuSubmit != nil {
		args = append(args, *peserta.WaktuSubmit)
		baseQuery = fmt.Sprintf("%s AND waktu_submit = $%d", baseQuery, len(args))
	}

	if peserta.NilaiUjian != nil {
		args = append(args, *peserta.NilaiUjian)
		baseQuery = fmt.Sprintf("%s AND nilai_ujian = $%d", baseQuery, len(args))
	}

	baseQuery = fmt.Sprintf("%s ORDER BY id_peserta_ujian ASC LIMIT 1", baseQuery)

	var (
		item        ujian.PesertaUjian
		waktuMulai  sql.NullTime
		waktuSubmit sql.NullTime
		nilaiUjian  sql.NullFloat64
		updatedAt   sql.NullTime
	)

	if err := r.q.QueryRow(ctx, baseQuery, args...).Scan(
		&item.IdPesertaUjian,
		&item.IdJadwalUjian,
		&item.IdSiswa,
		&waktuMulai,
		&waktuSubmit,
		&nilaiUjian,
		&item.CreatedAt,
		&updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ujian.PesertaUjian{}, coreerror.ErrNotFound
		}
		r.loggerFor(ctx).Error(ctx, "failed get peserta ujian by siswa", "layer", "repo.db", "op", "ujian.get_peserta_by_siswa", "err", err)
		return ujian.PesertaUjian{}, err
	}

	item.WaktuMulai = nullTimeToPtr(waktuMulai)
	item.WaktuSubmit = nullTimeToPtr(waktuSubmit)
	item.NilaiUjian = nullFloat64ToPtr(nilaiUjian)
	item.UpdatedAt = nullTimeToPtr(updatedAt)

	return item, nil
}

func (r *ListUjianRepo) GetAllJawabanUjianSiswa(ctx context.Context, jawaban ujian.JawabanUjianSiswa) ([]ujian.JawabanUjianSiswa, error) {
	baseQuery := `
		SELECT
			id_jawaban,
			id_peserta_ujian,
			id_soal,
			id_pilihan,
			jawaban_essay,
			is_benar,
			waktu_jawab
		FROM jawaban_ujian_siswa
	`

	where := make([]string, 0, 7)
	args := make([]any, 0, 7)

	if jawaban.IdJawaban > 0 {
		args = append(args, jawaban.IdJawaban)
		where = append(where, fmt.Sprintf("id_jawaban = $%d", len(args)))
	}

	if jawaban.IdPesertaUjian > 0 {
		args = append(args, jawaban.IdPesertaUjian)
		where = append(where, fmt.Sprintf("id_peserta_ujian = $%d", len(args)))
	}

	if jawaban.IdSoal > 0 {
		args = append(args, jawaban.IdSoal)
		where = append(where, fmt.Sprintf("id_soal = $%d", len(args)))
	}

	if jawaban.IdPilihan != nil {
		args = append(args, *jawaban.IdPilihan)
		where = append(where, fmt.Sprintf("id_pilihan = $%d", len(args)))
	}

	if jawaban.JawabanEssay != nil {
		args = append(args, *jawaban.JawabanEssay)
		where = append(where, fmt.Sprintf("jawaban_essay = $%d", len(args)))
	}

	if jawaban.IsBenar != nil {
		args = append(args, *jawaban.IsBenar)
		where = append(where, fmt.Sprintf("is_benar = $%d", len(args)))
	}

	if jawaban.WaktuJawab != nil {
		args = append(args, *jawaban.WaktuJawab)
		where = append(where, fmt.Sprintf("waktu_jawab = $%d", len(args)))
	}

	if len(where) > 0 {
		baseQuery = fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(where, " AND "))
	}

	baseQuery = fmt.Sprintf("%s ORDER BY id_jawaban ASC", baseQuery)

	rows, err := r.q.Query(ctx, baseQuery, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get jawaban ujian siswa", "layer", "repo.db", "op", "ujian.get_jawaban", "err", err)
		return nil, err
	}
	defer rows.Close()

	var results []ujian.JawabanUjianSiswa
	for rows.Next() {
		var (
			item         ujian.JawabanUjianSiswa
			idPilihan    sql.NullInt64
			jawabanEssay sql.NullString
			isBenar      sql.NullBool
			waktuJawab   sql.NullTime
		)

		if err := rows.Scan(
			&item.IdJawaban,
			&item.IdPesertaUjian,
			&item.IdSoal,
			&idPilihan,
			&jawabanEssay,
			&isBenar,
			&waktuJawab,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning jawaban ujian siswa", "layer", "repo.db", "op", "ujian.get_jawaban.scan", "err", err)
			return nil, err
		}

		item.IdPilihan = nullInt64ToUjianIDPtr(idPilihan)
		item.JawabanEssay = nullStringToPtr(jawabanEssay)
		item.IsBenar = nullBoolToPtr(isBenar)
		item.WaktuJawab = nullTimeToPtr(waktuJawab)

		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating jawaban ujian siswa", "layer", "repo.db", "op", "ujian.get_jawaban.iter", "err", err)
		return nil, err
	}

	return results, nil
}

func (r *ListUjianRepo) GetJawabanBySiswa(ctx context.Context, idSiswa ujian.ID, jawaban ujian.JawabanUjianSiswa) (ujian.JawabanUjianSiswa, error) {
	baseQuery := `
		SELECT
			jus.id_jawaban,
			jus.id_peserta_ujian,
			jus.id_soal,
			jus.id_pilihan,
			jus.jawaban_essay,
			jus.is_benar,
			jus.waktu_jawab
		FROM jawaban_ujian_siswa jus
		JOIN peserta_ujian pu
			ON pu.id_peserta_ujian = jus.id_peserta_ujian
		WHERE pu.id_siswa = $1
	`

	args := make([]any, 0, 7)
	args = append(args, idSiswa)

	if jawaban.IdJawaban > 0 {
		args = append(args, jawaban.IdJawaban)
		baseQuery = fmt.Sprintf("%s AND jus.id_jawaban = $%d", baseQuery, len(args))
	}

	if jawaban.IdPesertaUjian > 0 {
		args = append(args, jawaban.IdPesertaUjian)
		baseQuery = fmt.Sprintf("%s AND jus.id_peserta_ujian = $%d", baseQuery, len(args))
	}

	if jawaban.IdSoal > 0 {
		args = append(args, jawaban.IdSoal)
		baseQuery = fmt.Sprintf("%s AND jus.id_soal = $%d", baseQuery, len(args))
	}

	if jawaban.IdPilihan != nil {
		args = append(args, *jawaban.IdPilihan)
		baseQuery = fmt.Sprintf("%s AND jus.id_pilihan = $%d", baseQuery, len(args))
	}

	if jawaban.JawabanEssay != nil {
		args = append(args, *jawaban.JawabanEssay)
		baseQuery = fmt.Sprintf("%s AND jus.jawaban_essay = $%d", baseQuery, len(args))
	}

	if jawaban.IsBenar != nil {
		args = append(args, *jawaban.IsBenar)
		baseQuery = fmt.Sprintf("%s AND jus.is_benar = $%d", baseQuery, len(args))
	}

	if jawaban.WaktuJawab != nil {
		args = append(args, *jawaban.WaktuJawab)
		baseQuery = fmt.Sprintf("%s AND jus.waktu_jawab = $%d", baseQuery, len(args))
	}

	baseQuery = fmt.Sprintf("%s ORDER BY jus.id_jawaban ASC LIMIT 1", baseQuery)

	var (
		item         ujian.JawabanUjianSiswa
		idPilihan    sql.NullInt64
		jawabanEssay sql.NullString
		isBenar      sql.NullBool
		waktuJawab   sql.NullTime
	)

	if err := r.q.QueryRow(ctx, baseQuery, args...).Scan(
		&item.IdJawaban,
		&item.IdPesertaUjian,
		&item.IdSoal,
		&idPilihan,
		&jawabanEssay,
		&isBenar,
		&waktuJawab,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ujian.JawabanUjianSiswa{}, coreerror.ErrNotFound
		}
		r.loggerFor(ctx).Error(ctx, "failed get jawaban by siswa", "layer", "repo.db", "op", "ujian.get_jawaban_by_siswa", "err", err)
		return ujian.JawabanUjianSiswa{}, err
	}

	item.IdPilihan = nullInt64ToUjianIDPtr(idPilihan)
	item.JawabanEssay = nullStringToPtr(jawabanEssay)
	item.IsBenar = nullBoolToPtr(isBenar)
	item.WaktuJawab = nullTimeToPtr(waktuJawab)

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

func nullTimeToPtr(v sql.NullTime) *time.Time {
	if !v.Valid {
		return nil
	}
	t := v.Time
	return &t
}

func nullFloat64ToPtr(v sql.NullFloat64) *float64 {
	if !v.Valid {
		return nil
	}
	f := v.Float64
	return &f
}

func nullBoolToPtr(v sql.NullBool) *bool {
	if !v.Valid {
		return nil
	}
	b := v.Bool
	return &b
}
