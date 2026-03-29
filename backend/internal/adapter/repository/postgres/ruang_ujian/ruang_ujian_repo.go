package ruangujianrepo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ruang_ujian"
)

type RuangUjianRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewRuangUjianRepo(q pg.Executor, logger corelog.Logger) *RuangUjianRepo {
	return &RuangUjianRepo{q: q, logger: logger}
}

func (r *RuangUjianRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *RuangUjianRepo) GetRuangUjian(ctx context.Context, filter query.ListRuangUjianFilter) ([]ruangujian.RuangUjian, error) {
	queryText, args := r.buildListRuangUjianQuery(filter)

	rows, err := r.q.Query(ctx, queryText, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get ruang ujian", "layer", "repo.db", "op", "ruang_ujian.get", "err", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanRuangUjianRows(ctx, "ruang_ujian.get", rows)
}

func (r *RuangUjianRepo) GetRuangUjianById(ctx context.Context, idRuangan int) (ruangujian.RuangUjian, error) {
	query := `
		SELECT
			id_ruangan,
			nama_ruangan,
			kode_ruang
		FROM ruang_ujian
		WHERE id_ruangan = $1
	`

	item, err := scanRuangUjianRow(r.q.QueryRow(ctx, query, idRuangan))
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed scanning ruang ujian", "layer", "repo.db", "op", "ruang_ujian.get_by_id", "err", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return ruangujian.RuangUjian{}, coreerror.ErrNotFound
		}
		return ruangujian.RuangUjian{}, err
	}

	return item, nil
}

func (r *RuangUjianRepo) GetRuangUjianByKode(ctx context.Context, kodeRuang string) (ruangujian.RuangUjian, error) {
	query := `
		SELECT
			id_ruangan,
			nama_ruangan,
			kode_ruang
		FROM ruang_ujian
		WHERE kode_ruang = $1
	`

	item, err := scanRuangUjianRow(r.q.QueryRow(ctx, query, kodeRuang))
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed scanning ruang ujian", "layer", "repo.db", "op", "ruang_ujian.get_by_kode", "err", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return ruangujian.RuangUjian{}, coreerror.ErrNotFound
		}
		return ruangujian.RuangUjian{}, err
	}

	return item, nil
}

func (r *RuangUjianRepo) ExistByKodeRuang(ctx context.Context, kodeRuang string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM ruang_ujian WHERE kode_ruang = $1)`

	var exists bool
	if err := r.q.QueryRow(ctx, query, kodeRuang).Scan(&exists); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed checking kode ruang", "layer", "repo.db", "op", "ruang_ujian.exist_by_kode", "err", err)
		return false, err
	}

	return exists, nil
}

func (r *RuangUjianRepo) CreateRuangUjian(ctx context.Context, ruangUjian ruangujian.RuangUjian) error {
	query := `
		INSERT INTO ruang_ujian (nama_ruangan, kode_ruang)
		VALUES ($1, $2)
	`

	_, err := r.q.Exec(ctx, query, ruangUjian.NamaRuangan, ruangUjian.KodeRuang)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed creating ruang ujian", "layer", "repo.db", "op", "ruang_ujian.create", "err", err)
		return err
	}

	return nil
}

func (r *RuangUjianRepo) UpdateRuangUjian(ctx context.Context, idRuangan int, ruangUjian updatepatch.UpdateRuangUjianPatch) error {
	query := `
		UPDATE ruang_ujian
		SET
			kode_ruang = COALESCE($1, kode_ruang),
			nama_ruangan = COALESCE($2, nama_ruangan),
			updated_at = now()
		WHERE id_ruangan = $3
	`

	tag, err := r.q.Exec(
		ctx,
		query,
		ruangUjian.KodeRuang,
		ruangUjian.NamaRuang,
		idRuangan,
	)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed updating ruang ujian", "layer", "repo.db", "op", "ruang_ujian.update", "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	return nil
}

func (r *RuangUjianRepo) DeleteRuangUjian(ctx context.Context, idRuangan int) error {
	query := `
		DELETE FROM ruang_ujian
		WHERE id_ruangan = $1
	`

	tag, err := r.q.Exec(ctx, query, idRuangan)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23001" {
				return coreerror.ErrDeleteRestricted
			}
		}

		r.loggerFor(ctx).Error(ctx, "failed deleting ruang ujian", "layer", "repo.db", "op", "ruang_ujian.delete", "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}
	return nil
}
