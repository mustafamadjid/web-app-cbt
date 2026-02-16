package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
)

type SesiRepo struct {
	q Executor
	logger corelog.Logger
}	

func NewSesirepo(q Executor, logger corelog.Logger) *SesiRepo {
	return &SesiRepo{q: q, logger: logger}
}

func (r *SesiRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}


func (r *SesiRepo)GetSesi(ctx context.Context,filter query.ListSesiFilter) ([]sesi.Sesi, error) {
	query := `
		SELECT
			id_sesi,
			kode_sesi,
			nama_sesi,
		FROM sesi	
	`

	where := make([]string, 0, 1)
	args := make([]any, 0, 3)

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("(nama_sesi ILIKE $%d OR kode_sesi ILIKE $%d)", idx, idx))
	}

	if len(where) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(where, " AND "))
	}

	query = fmt.Sprintf("%s ORDER BY id_sesi ASC", query)

	if filter.Limit > 0 {
		args = append(args, filter.Limit)
		limitIndex := len(args)
		args = append(args, filter.Offset)
		offsetIndex := len(args)
		query = fmt.Sprintf("%s LIMIT $%d OFFSET $%d", query, limitIndex, offsetIndex)
	}

	rows, err := r.q.Query(ctx, query, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get sesi", "layer", "repo.db", "op", "sesi.get", "err", err)
		return nil, err
	}

	defer rows.Close()

	var results []sesi.Sesi
	for rows.Next() {
		var item sesi.Sesi
		if err := rows.Scan(
			&item.IdSesi,
			&item.KodeSesi,
			&item.NamaSesi,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning sesi", "layer", "repo.db", "op", "sesi.get", "err", err)
			return nil, err
		}
		results = append(results, item)
	}

	return results, nil
}

func(r *SesiRepo)GetSesiById(ctx context.Context,idSesi int) (sesi.Sesi, error) {
	query := `
		SELECT
			id_sesi,
			kode_sesi,
			nama_sesi,
		FROM sesi	
		WHERE id_sesi = $1
	`

	rows := r.q.QueryRow(ctx,query,idSesi)

	var item sesi.Sesi
	if err := rows.Scan(
		&item.IdSesi,
		&item.KodeSesi,
		&item.NamaSesi,
	); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed scanning sesi", "layer", "repo.db", "op", "sesi.get_by_id", "err", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return sesi.Sesi{}, coreerror.ErrNotFound
		}
		return sesi.Sesi{}, err
	}

	return item, nil
}

func(r *SesiRepo)GetSesiByKode(ctx context.Context,kodeSesi string) (sesi.Sesi, error) {
	query := `
		SELECT
			id_sesi,
			kode_sesi,
			nama_sesi,
		FROM sesi	
		WHERE kode_sesi = $1
	`

	rows := r.q.QueryRow(ctx,query,kodeSesi)

	var item sesi.Sesi
	if err := rows.Scan(
		&item.IdSesi,
		&item.KodeSesi,
		&item.NamaSesi,
	); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed scanning sesi", "layer", "repo.db", "op", "sesi.get_by_kode", "err", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return sesi.Sesi{}, coreerror.ErrNotFound
		}
		return sesi.Sesi{}, err
	}

	return item, nil
}

func(r *SesiRepo)ExistByKodeSesi(ctx context.Context,kodeSesi string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM sesi
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

func(r *SesiRepo)CreateSesi(ctx context.Context, sesi sesi.Sesi) error {
	query := `
		INSERT INTO sesi (id_sesi,kode_sesi,nama_sesi)
		VALUES ($1,$2,$3)
	`
	_, err := r.q.Exec(ctx, query,
		sesi.IdSesi,
		sesi.KodeSesi,
		sesi.NamaSesi,
	)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed insert sesi", "layer", "repo.db", "op", "sesi.create", "err", err)
		return err
	}
	return nil
}

func (r *SesiRepo)UpdateSesi(ctx context.Context, idSesi int, sesi updatepatch.UpdateSesiPatch) error {
	query := `
		UPDATE sesi
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

func (r *SesiRepo)DeleteSesi(ctx context.Context, idSesi int) error {
	query := `
		DELETE FROM sesi
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