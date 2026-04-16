package ujianlistrepo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
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
	queryText, args := r.buildListUjianQuery(filter)

	rows, err := r.q.Query(ctx, queryText, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get list ujian", "layer", "repo.db", "op", "ujian.get_all", "err", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanListUjianSummaryRows(ctx, "ujian.get_all", rows)
}

func (r *ListUjianRepo) GetAllUjianSubmittedByIdSiswa(ctx context.Context, idSiswa int) ([]ujian.ListUjian, error) {
	queryText, args := r.buildListUjianSubmittedByIdSiswaQuery(idSiswa)

	rows, err := r.q.Query(ctx, queryText, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get submitted ujian list by siswa", "layer", "repo.db", "op", "ujian.get_all_submitted_by_id_siswa", "id_siswa", idSiswa, "err", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanListUjianSubmittedRows(ctx, "ujian.get_all_submitted_by_id_siswa", rows)
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
	item, err := scanListUjianDetailRow(r.q.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return ujian.ListUjian{}, coreerror.ErrNotFound
	}
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed scanning list ujian", "layer", "repo.db", "op", "ujian.get_by_id.scan", "err", err)
		return ujian.ListUjian{}, err
	}

	return item, nil
}
