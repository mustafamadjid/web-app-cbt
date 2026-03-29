package sesirepo

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
)

type SesiRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewSesirepo(q pg.Executor, logger corelog.Logger) *SesiRepo {
	return &SesiRepo{q: q, logger: logger}
}

func (r *SesiRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *SesiRepo) GetSesi(ctx context.Context, filter query.ListSesiFilter) ([]sesi.Sesi, error) {
	queryText, args := r.buildListSesiQuery(filter)

	rows, err := r.q.Query(ctx, queryText, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get sesi", "layer", "repo.db", "op", "sesi.get", "err", err)
		return nil, err
	}

	defer rows.Close()

	return r.scanSesiRows(ctx, "sesi.get", rows)
}

func (r *SesiRepo) GetSesiById(ctx context.Context, idSesi int) (sesi.Sesi, error) {
	query := `
		SELECT
			id_sesi,
			kode_sesi,
			nama_sesi,
		FROM sesi_ujian	
		WHERE id_sesi = $1
	`

	item, err := scanSesiRow(r.q.QueryRow(ctx, query, idSesi))
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed scanning sesi", "layer", "repo.db", "op", "sesi.get_by_id", "err", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return sesi.Sesi{}, coreerror.ErrNotFound
		}
		return sesi.Sesi{}, err
	}

	return item, nil
}

func (r *SesiRepo) GetSesiByKode(ctx context.Context, kodeSesi string) (sesi.Sesi, error) {
	query := `
		SELECT
			id_sesi,
			kode_sesi,
			nama_sesi,
		FROM sesi_ujian		
		WHERE kode_sesi = $1
	`

	item, err := scanSesiRow(r.q.QueryRow(ctx, query, kodeSesi))
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed scanning sesi", "layer", "repo.db", "op", "sesi.get_by_kode", "err", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return sesi.Sesi{}, coreerror.ErrNotFound
		}
		return sesi.Sesi{}, err
	}

	return item, nil
}

func (r *SesiRepo) ExistByKodeSesi(ctx context.Context, kodeSesi string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM sesi_ujian	
			WHERE kode_sesi = $1
		)
	`

	kodeSesi = strings.TrimSpace(kodeSesi)
	kodeSesi = strings.ToUpper(kodeSesi)

	var exist bool
	if err := r.q.QueryRow(ctx, query, kodeSesi).Scan(&exist); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed check exist sesi", "layer", "repo.db", "op", "sesi.exist_by_kode_sesi", "err", err)
		return false, err
	}
	return exist, nil
}

func (r *SesiRepo) CreateSesi(ctx context.Context, sesi sesi.Sesi) error {
	query := `
		INSERT INTO sesi_ujian (kode_sesi,nama_sesi)
		VALUES ($1,$2)
	`
	_, err := r.q.Exec(ctx, query,
		sesi.KodeSesi,
		sesi.NamaSesi,
	)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed insert sesi", "layer", "repo.db", "op", "sesi.create", "err", err)
		return err
	}
	return nil
}

func (r *SesiRepo) UpdateSesi(ctx context.Context, idSesi int, sesi updatepatch.UpdateSesiPatch) error {
	query := `
		UPDATE sesi_ujian	
		SET
			kode_sesi = COALESCE($1,kode_sesi),
			nama_sesi = COALESCE($2,nama_sesi)
			updated_at = now()
		WHERE id_sesi = $3s
	`
	tag, err := r.q.Exec(
		ctx,
		query,
		idSesi,
		sesi.KodeSesi,
		sesi.NamaSesi,
	)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed update sesi", "layer", "repo.db", "op", "sesi.update", "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		r.loggerFor(ctx).Error(ctx, "failed update sesi", "layer", "repo.db", "op", "sesi.update", "err", coreerror.ErrNotFound)
		return coreerror.ErrNotFound
	}

	return nil
}

func (r *SesiRepo) DeleteSesi(ctx context.Context, idSesi int) error {
	query := `
		DELETE FROM sesi_ujian	
		WHERE id_sesi = $1
	`
	tag, err := r.q.Exec(ctx, query, idSesi)
	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23001" {
				return coreerror.ErrDeleteRestricted
			}
		}
		r.loggerFor(ctx).Error(ctx, "failed delete sesi", "layer", "repo.db", "op", "sesi.delete", "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	return nil
}
