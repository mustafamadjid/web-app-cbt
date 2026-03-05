package ujianrepo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type PesertaUjianRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewPesertaUjianRepo(q pg.Executor, logger corelog.Logger) *PesertaUjianRepo {
	return &PesertaUjianRepo{q: q, logger: logger}
}

func (r *PesertaUjianRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *PesertaUjianRepo) CreatePesertaUjian(ctx context.Context, peserta ujian.PesertaUjian) (ujian.ID, error) {
	query := `
		INSERT INTO peserta_ujian (
			id_jadwal_ujian,
			id_siswa,
			waktu_mulai,
			waktu_submit,
			nilai_ujian
		)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id_peserta_ujian
	`

	var idPesertaUjian ujian.ID
	err := r.q.QueryRow(
		ctx,
		query,
		peserta.IdJadwalUjian,
		peserta.IdSiswa,
		peserta.WaktuMulai,
		peserta.WaktuSubmit,
		peserta.NilaiUjian,
	).Scan(&idPesertaUjian)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed creating peserta ujian", "layer", "repo.db", "op", "peserta_ujian.create", "err", err)
		return 0, err
	}

	return idPesertaUjian, nil
}

func (r *PesertaUjianRepo) UpdatePesertaUjian(ctx context.Context, id ujian.ID, payload updatepatch.UpdatePesertaUjianPatch) error {
	tag, err := r.q.Exec(
		ctx,
		`UPDATE peserta_ujian
		SET
			id_jadwal_ujian = COALESCE($1, id_jadwal_ujian),
			id_siswa = COALESCE($2, id_siswa),
			waktu_mulai = COALESCE($3, waktu_mulai),
			waktu_submit = COALESCE($4, waktu_submit),
			nilai_ujian = COALESCE($5, nilai_ujian),
			updated_at = NOW()
		WHERE id_peserta_ujian = $6`,
		payload.IdJadwalUjian,
		payload.IdSiswa,
		payload.WaktuMulai,
		payload.WaktuSubmit,
		payload.NilaiUjian,
		id,
	)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed updating peserta ujian", "layer", "repo.db", "op", "peserta_ujian.update", "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	return nil
}

func (r *PesertaUjianRepo) DeletePesertaUjian(ctx context.Context, id ujian.ID) error {
	tag, err := r.q.Exec(ctx, `DELETE FROM peserta_ujian WHERE id_peserta_ujian = $1`, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23001" || pgErr.Code == "23503" {
				return coreerror.ErrDeleteRestricted
			}
		}
		r.loggerFor(ctx).Error(ctx, "failed deleting peserta ujian", "layer", "repo.db", "op", "peserta_ujian.delete", "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	return nil
}
