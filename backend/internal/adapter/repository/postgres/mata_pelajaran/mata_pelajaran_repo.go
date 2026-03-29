package matapelajaranrepo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
)

type MapelRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewMapelRepo(q pg.Executor, logger corelog.Logger) *MapelRepo {
	return &MapelRepo{q: q, logger: logger}
}

func (r *MapelRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *MapelRepo) GetMapel(ctx context.Context, filter query.ListMapelFilter) ([]matapelajaran.MataPelajaran, error) {
	queryText, args := r.buildListMapelQuery(filter)

	rows, err := r.q.Query(ctx, queryText, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get mapel", "layer", "core.service", "op", "matapelajaran.get", "err", err)
		return nil, err
	}

	defer rows.Close()

	return r.scanMapelRows(ctx, "matapelajaran.get", rows)
}

func (r *MapelRepo) GetMapelById(ctx context.Context, idMapel int) (matapelajaran.MataPelajaran, error) {
	query := `
		SELECT
			id_mapel,
			id_kelas,
			kode_mapel,
			nama_mapel,
			deskripsi
		FROM mata_pelajaran
		WHERE id_mapel = $1
	`

	item, err := scanMapelRow(r.q.QueryRow(ctx, query, idMapel))
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed scanning mapel", "layer", "core.service", "op", "matapelajaran.get_by_id", "err", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return matapelajaran.MataPelajaran{}, coreerror.ErrNotFound
		}
		return matapelajaran.MataPelajaran{}, err
	}
	return item, nil
}

func (r *MapelRepo) CreateMapel(ctx context.Context, mapel matapelajaran.MataPelajaran) error {
	query := `
		INSERT INTO mata_pelajaran (id_kelas, kode_mapel, nama_mapel, deskripsi) VALUES ($1, $2, $3, $4)
	`

	_, err := r.q.Exec(ctx, query, mapel.IdKelas, mapel.KodeMapel, mapel.NamaMapel, mapel.Deskripsi)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed creating mapel", "layer", "core.service", "op", "matapelajaran.create_mapel", "err", err)
		return err
	}
	return nil
}

func (r *MapelRepo) UpdateMapel(ctx context.Context, idMapel int, mapel updatepatch.UpdateMapelPatch) error {
	query := `
		UPDATE mata_pelajaran
		SET
			id_kelas = COALESCE($1, id_kelas),
			kode_mapel = COALESCE($2, kode_mapel),
			nama_mapel = COALESCE($3, nama_mapel),
			deskripsi = COALESCE($4, deskripsi),
			updated_at = now()
		WHERE id_mapel = $5
	`

	tag, err := r.q.Exec(
		ctx,
		query,
		mapel.IdKelas,
		mapel.KodeMapel,
		mapel.NamaMapel,
		mapel.Deskripsi,
		idMapel,
	)

	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed updating mapel", "layer", "repo.db", "op", "matapelajaran.update_mapel", "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}
	return nil
}

func (r *MapelRepo) DeleteMapel(ctx context.Context, idMapel int) error {
	query := `
		DELETE FROM mata_pelajaran
		WHERE id_mapel = $1
	`

	tag, err := r.q.Exec(ctx, query, idMapel)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23001" {
				return coreerror.ErrDeleteRestricted
			}
		}
		r.loggerFor(ctx).Error(ctx, "failed deleting mapel", "layer", "core.service", "op", "matapelajaran.delete_mapel", "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}
	return nil
}

func (r *MapelRepo) ExistKodeMapel(ctx context.Context, kodeMapel string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM mata_pelajaran WHERE kode_mapel = $1)`

	var exist bool
	if err := r.q.QueryRow(ctx, query, kodeMapel).Scan(&exist); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed scanning mapel", "layer", "core.service", "op", "matapelajaran.existKodeMapel", "err", err)
		return false, err
	}
	return exist, nil
}
