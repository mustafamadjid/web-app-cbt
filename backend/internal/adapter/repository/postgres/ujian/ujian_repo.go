package ujianrepo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type UjianRepo struct {
	q      pg.Executor
	logger corelog.Logger
	pool   *pgxpool.Pool
}

func NewUjianRepo(q pg.Executor, logger corelog.Logger, pool *pgxpool.Pool) *UjianRepo {
	return &UjianRepo{q: q, logger: logger, pool: pool}
}

func (r *UjianRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *UjianRepo) GetIdUjianByAttempt(ctx context.Context, idAttempt ujian.ID) (ujian.ID, error) {
	query := `
		SELECT 
			ju.id_ujian
		FROM attempt_ujian au
		JOIN peserta_ujian pu 
			ON pu.id_peserta_ujian = au.id_peserta_ujian
		JOIN jadwal_ujian ju 
			ON ju.id_jadwal_ujian = pu.id_jadwal_ujian
		WHERE au.id_attempt = $1;
	`

	var idUjian ujian.ID

	err := r.q.QueryRow(ctx, query, idAttempt).Scan(&idUjian)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, coreerror.ErrNotFound
		}
		r.loggerFor(ctx).Error(ctx, "failed get ujian id by attempt id", "layer", "repo.db", "op", "ujian.get_id_by_attempt_id", "attempt_id", idAttempt, "err", err)
		return 0, err
	}
	return idUjian, nil

}

func (r *UjianRepo) CreateUjian(ctx context.Context, data ujian.PenjadwalanUjian) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed begin tx create ujian", "layer", "repo.db", "op", "ujian.create.begin_tx", "err", err)
		return err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	var idUjian int64
	var idJadwalUjian int64

	err = tx.QueryRow(
		ctx,
		`INSERT INTO ujian (
			id_bank_soal,
			id_kelas,
			id_nama_kelas,
			id_guru,
			nama_ujian,
			deskripsi_ujian,
			acak_soal
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id_ujian`,
		data.Ujian.IdBankSoal,
		data.Ujian.IdKelas,
		data.Ujian.IdNamaKelas,
		data.Ujian.IdGuru,
		data.Ujian.NamaUjian,
		data.Ujian.DeskripsiUjian,
		data.Ujian.AcakSoal,
	).Scan(&idUjian)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed to insert ujian", "layer", "repo.db", "op", "ujian.create.insert_ujian", "err", err)
		return err
	}

	err = tx.QueryRow(
		ctx,
		`INSERT INTO jadwal_ujian (
			id_ujian,
			id_sesi,
			id_ruangan,
			id_pengawas,
			tanggal_ujian,
			waktu_mulai,
			waktu_selesai,
			token,
			status_ujian
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id_jadwal_ujian`,
		idUjian,
		data.JadwalUjian.IdSesi,
		data.JadwalUjian.IdRuangan,
		data.JadwalUjian.IdPengawas,
		data.JadwalUjian.TanggalUjian,
		data.JadwalUjian.WaktuMulai,
		data.JadwalUjian.WaktuSelesai,
		data.JadwalUjian.Token,
		data.JadwalUjian.StatusUjian,
	).Scan(&idJadwalUjian)
	if err != nil {
		if mappedErr := mapJadwalUjianConflictError(err); mappedErr != nil {
			return mappedErr
		}

		r.loggerFor(ctx).Error(ctx, "failed to insert jadwal ujian", "layer", "repo.db", "op", "ujian.create.insert_jadwal", "err", err)
		return err
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO peserta_ujian (id_jadwal_ujian, id_siswa)
		SELECT $1, ps.id_pengguna
		FROM profil_siswa ps
		JOIN nama_kelas nk
			ON nk.id_nama_kelas = ps.id_nama_kelas
		WHERE nk.id_kelas = $2
			AND ($3::bigint IS NULL OR ps.id_nama_kelas = $3)`,
		idJadwalUjian,
		data.Ujian.IdKelas,
		data.Ujian.IdNamaKelas,
	)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed to insert peserta ujian", "layer", "repo.db", "op", "ujian.create.insert_peserta", "err", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed commit create ujian", "layer", "repo.db", "op", "ujian.create.commit", "err", err)
		return err
	}

	committed = true
	return nil
}

func (r *UjianRepo) UpdateUjian(ctx context.Context, id ujian.ID, payload updatepatch.UpdatePenjadwalanUjian) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed begin tx update ujian", "layer", "repo.db", "op", "ujian.update.begin_tx", "err", err)
		return err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	if hasUpdateUjianPatch(payload.Ujian) {
		tag, err := tx.Exec(
			ctx,
			`UPDATE ujian
			SET
				id_bank_soal = COALESCE($1, id_bank_soal),
				id_kelas = COALESCE($2, id_kelas),
				id_nama_kelas = CASE
					WHEN $3::bigint IS NULL THEN id_nama_kelas
					WHEN $3::bigint = 0 THEN NULL
					ELSE $3::bigint
				END,
				id_guru = COALESCE($4, id_guru),
				nama_ujian = COALESCE($5, nama_ujian),
				deskripsi_ujian = COALESCE($6, deskripsi_ujian),
				acak_soal = COALESCE($7, acak_soal),
				updated_at = NOW()
			WHERE id_ujian = $8`,
			payload.Ujian.IdBankSoal,
			payload.Ujian.IdKelas,
			payload.Ujian.IdNamaKelas,
			payload.Ujian.IdGuru,
			payload.Ujian.NamaUjian,
			payload.Ujian.DeskripsiUjian,
			payload.Ujian.AcakSoal,
			id,
		)
		if err != nil {
			r.loggerFor(ctx).Error(ctx, "failed update ujian", "layer", "repo.db", "op", "ujian.update.ujian", "err", err)
			return err
		}
		if tag.RowsAffected() == 0 {
			return coreerror.ErrNotFound
		}
	}

	if hasUpdateJadwalUjianPatch(payload.JadwalUjian) {
		tag, err := tx.Exec(
			ctx,
			`UPDATE jadwal_ujian
			SET
				id_ujian = COALESCE($1, id_ujian),
				id_sesi = COALESCE($2, id_sesi),
				id_ruangan = COALESCE($3, id_ruangan),
				id_pengawas = COALESCE($4, id_pengawas),
				tanggal_ujian = COALESCE($5, tanggal_ujian),
				token = COALESCE($6, token),
				waktu_mulai = COALESCE($7, waktu_mulai),
				waktu_selesai = COALESCE($8, waktu_selesai),
				status_ujian = COALESCE($9, status_ujian),
				updated_at = NOW()
			WHERE id_ujian = $10`,
			payload.JadwalUjian.IdUjian,
			payload.JadwalUjian.IdSesi,
			payload.JadwalUjian.IdRuangan,
			payload.JadwalUjian.IdPengawas,
			payload.JadwalUjian.TanggalUjian,
			payload.JadwalUjian.Token,
			payload.JadwalUjian.WaktuMulai,
			payload.JadwalUjian.WaktuSelesai,
			payload.JadwalUjian.StatusUjian,
			id,
		)
		if err != nil {
			if mappedErr := mapJadwalUjianConflictError(err); mappedErr != nil {
				return mappedErr
			}

			r.loggerFor(ctx).Error(ctx, "failed update jadwal ujian", "layer", "repo.db", "op", "ujian.update.jadwal", "err", err)
			return err
		}
		if tag.RowsAffected() == 0 {
			return coreerror.ErrNotFound
		}
	}

	if err := tx.Commit(ctx); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed commit update ujian", "layer", "repo.db", "op", "ujian.update.commit", "err", err)
		return err
	}

	committed = true
	return nil
}

func (r *UjianRepo) DeleteUjian(ctx context.Context, id ujian.ID) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed begin tx delete ujian", "layer", "repo.db", "op", "ujian.delete.begin_tx", "err", err)
		return err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	if _, err := tx.Exec(ctx, `DELETE FROM jadwal_ujian WHERE id_ujian = $1`, id); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed delete jadwal ujian by id ujian", "layer", "repo.db", "op", "ujian.delete.jadwal", "err", err)
		return err
	}

	tag, err := tx.Exec(ctx, `DELETE FROM ujian WHERE id_ujian = $1`, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23001" || pgErr.Code == "23503" {
				return coreerror.ErrDeleteRestricted
			}
		}
		r.loggerFor(ctx).Error(ctx, "failed delete ujian", "layer", "repo.db", "op", "ujian.delete.ujian", "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed commit delete ujian", "layer", "repo.db", "op", "ujian.delete.commit", "err", err)
		return err
	}

	committed = true
	return nil
}

func hasUpdateUjianPatch(payload updatepatch.UpdateUjianPatch) bool {
	return payload.IdBankSoal != nil ||
		payload.IdKelas != nil ||
		payload.IdNamaKelas != nil ||
		payload.IdGuru != nil ||
		payload.NamaUjian != nil ||
		payload.DeskripsiUjian != nil ||
		payload.AcakSoal != nil
}

func hasUpdateJadwalUjianPatch(payload updatepatch.UpdateJadwalUjianPatch) bool {
	return payload.IdUjian != nil ||
		payload.IdSesi != nil ||
		payload.IdRuangan != nil ||
		payload.IdPengawas != nil ||
		payload.TanggalUjian != nil ||
		payload.Token != nil ||
		payload.WaktuMulai != nil ||
		payload.WaktuSelesai != nil ||
		payload.StatusUjian != nil
}

func mapJadwalUjianConflictError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23P01" {
		return nil
	}

	if pgErr.ConstraintName == "excl_jadwal_ujian_ruangan_sesi_waktu_active" {
		return coreerror.ErrConflict
	}

	return nil
}
