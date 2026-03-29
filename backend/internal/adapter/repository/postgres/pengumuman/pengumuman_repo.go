package pengumumanrepo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type PengumumanRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewPengumumanRepo(q pg.Executor, logger corelog.Logger) *PengumumanRepo {
	return &PengumumanRepo{q: q, logger: logger}
}

func (r *PengumumanRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *PengumumanRepo) GetPengumumanActive(ctx context.Context) ([]pengumuman.Pengumuman, error) {
	const query = `
		SELECT
			id_pengumuman,
			id_pengguna,
			judul_pengumuman,
			isi_pengumuman,
			tanggal_rilis_pengumuman,
			tanggal_selesai_pengumuman,
			dokumen_pengumuman
		FROM pengumuman
		WHERE tanggal_rilis_pengumuman <= CURRENT_DATE
			AND tanggal_selesai_pengumuman >= CURRENT_DATE
		ORDER BY tanggal_rilis_pengumuman DESC, id_pengumuman DESC
	`

	rows, err := r.q.Query(ctx, query)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed getting active pengumuman", "layer", "repo.db", "op", "pengumuman_repo.get_active", "err", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanPengumumanRows(ctx, "pengumuman_repo.get_active", rows)
}

func (r *PengumumanRepo) GetPengumumanNonActive(ctx context.Context) ([]pengumuman.Pengumuman, error) {
	const query = `
		SELECT
			id_pengumuman,
			id_pengguna,
			judul_pengumuman,
			isi_pengumuman,
			tanggal_rilis_pengumuman,
			tanggal_selesai_pengumuman,
			dokumen_pengumuman
		FROM pengumuman
		WHERE tanggal_selesai_pengumuman < CURRENT_DATE
		ORDER BY tanggal_selesai_pengumuman DESC, id_pengumuman DESC
	`

	rows, err := r.q.Query(ctx, query)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed getting non active pengumuman", "layer", "repo.db", "op", "pengumuman_repo.get_non_active", "err", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanPengumumanRows(ctx, "pengumuman_repo.get_non_active", rows)
}

func (r *PengumumanRepo) GetPengumumanIncoming(ctx context.Context) ([]pengumuman.Pengumuman, error) {
	const query = `
		SELECT
			id_pengumuman,
			id_pengguna,
			judul_pengumuman,
			isi_pengumuman,
			tanggal_rilis_pengumuman,
			tanggal_selesai_pengumuman,
			dokumen_pengumuman
		FROM pengumuman
		WHERE tanggal_rilis_pengumuman > CURRENT_DATE
		ORDER BY tanggal_rilis_pengumuman ASC, id_pengumuman ASC
	`

	rows, err := r.q.Query(ctx, query)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed getting incoming pengumuman", "layer", "repo.db", "op", "pengumuman_repo.get_incoming", "err", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanPengumumanRows(ctx, "pengumuman_repo.get_incoming", rows)
}

func (r *PengumumanRepo) GetPengumumanById(ctx context.Context, idPengumuman pengumuman.ID) (pengumuman.Pengumuman, error) {
	const query = `
		SELECT
			id_pengumuman,
			id_pengguna,
			judul_pengumuman,
			isi_pengumuman,
			tanggal_rilis_pengumuman,
			tanggal_selesai_pengumuman,
			dokumen_pengumuman
		FROM pengumuman
		WHERE id_pengumuman = $1
	`

	item, err := scanPengumumanRow(r.q.QueryRow(ctx, query, idPengumuman))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pengumuman.Pengumuman{}, coreerror.ErrNotFound
		}
		r.loggerFor(ctx).Error(ctx, "failed getting pengumuman by id", "layer", "repo.db", "op", "pengumuman_repo.get_by_id", "pengumuman_id", idPengumuman, "err", err)
		return pengumuman.Pengumuman{}, err
	}

	return item, nil
}

func (r *PengumumanRepo) CreatePengumuman(ctx context.Context, pengumuman pengumuman.Pengumuman) error {
	const query = `
		INSERT INTO pengumuman (
			id_pengguna,
			judul_pengumuman,
			isi_pengumuman,
			tanggal_rilis_pengumuman,
			tanggal_selesai_pengumuman,
			dokumen_pengumuman
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.q.Exec(
		ctx,
		query,
		pengumuman.IdPengguna,
		pengumuman.JudulPengumuman,
		pengumuman.IsiPengumuman,
		pengumuman.TanggalRilisPengumuman,
		pengumuman.TanggalSelesaiPengumuman,
		pengumuman.DokumenPengumuman,
	)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed creating pengumuman", "layer", "repo.db", "op", "pengumuman_repo.create", "err", err)
		return err
	}

	return nil
}

func (r *PengumumanRepo) UpdatePengumuman(ctx context.Context, idPengumuman pengumuman.ID, dataUpdate updatepatch.PengumumanUpdatePatch) error {
	const query = `
		UPDATE pengumuman
		SET
			id_pengguna = COALESCE($1, id_pengguna),
			judul_pengumuman = COALESCE($2, judul_pengumuman),
			isi_pengumuman = COALESCE($3, isi_pengumuman),
			tanggal_rilis_pengumuman = COALESCE($4, tanggal_rilis_pengumuman),
			tanggal_selesai_pengumuman = COALESCE($5, tanggal_selesai_pengumuman),
			dokumen_pengumuman = COALESCE($6, dokumen_pengumuman),
			updated_at = now()
		WHERE id_pengumuman = $7
	`

	tag, err := r.q.Exec(
		ctx,
		query,
		dataUpdate.IdPengguna,
		dataUpdate.JudulPengumuman,
		dataUpdate.IsiPengumuman,
		dataUpdate.TanggalRilisPengumuman,
		dataUpdate.TanggalSelesaiPengumuman,
		dataUpdate.DokumenPengumuman,
		idPengumuman,
	)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed updating pengumuman", "layer", "repo.db", "op", "pengumuman_repo.update", "pengumuman_id", idPengumuman, "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	return nil
}

func (r *PengumumanRepo) DeletePengumuman(ctx context.Context, idPengumuman pengumuman.ID) error {
	const query = `
		DELETE FROM pengumuman
		WHERE id_pengumuman = $1
	`

	tag, err := r.q.Exec(ctx, query, idPengumuman)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23001" {
				return coreerror.ErrDeleteRestricted
			}
		}
		r.loggerFor(ctx).Error(ctx, "failed deleting pengumuman", "layer", "repo.db", "op", "pengumuman_repo.delete", "pengumuman_id", idPengumuman, "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	return nil
}
