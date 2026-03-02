package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
)

type MapelRepo struct {
	q Executor
	logger corelog.Logger
}

func NewMapelRepo(q Executor, logger corelog.Logger) *MapelRepo {
	return &MapelRepo{q: q, logger: logger}
}

func (r *MapelRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *MapelRepo)GetMapel(ctx context.Context, filter query.ListMapelFilter) ([]matapelajaran.MataPelajaran, error){
	query := `
		SELECT 
			id_mapel,
			id_kelas,
			kode_mapel,
			nama_mapel,
			deskripsi
		FROM mata_pelajaran
	`

	where := make([]string, 0, 2)
	args := make([]any, 0, 4)

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("(nama_mapel ILIKE $%d OR deskripsi ILIKE $%d OR kode_mapel ILIKE $%d)",idx, idx, idx,))
	}

	if filter.NamaMapel != nil {
		args = append(args, *filter.NamaMapel)
		where = append(where, fmt.Sprintf("nama_mapel = $%d", len(args)))
	}

	if filter.TingkatKelas != nil {
		args = append(args, *filter.TingkatKelas)
		where = append(where, fmt.Sprintf("id_kelas = $%d", len(args)))
	}

	if len(where) > 0{
		query = fmt.Sprintf("%s WHERE %s", query,strings.Join(where, " AND "))
	}

	if filter.Limit > 0 {
		args = append(args, filter.Limit)
		limitIndex := len(args)
		args = append(args, filter.Offset)
		offsetIndex := len(args)
		query = fmt.Sprintf("%s ORDER BY created_at ASC LIMIT $%d OFFSET $%d", query, limitIndex, offsetIndex)
	}

	rows, err := r.q.Query(ctx, query, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get mapel", "layer", "core.service", "op", "matapelajaran.get", "err", err)
		return nil, err
	}

	defer rows.Close()

	var results []matapelajaran.MataPelajaran
	for rows.Next() {
		var item matapelajaran.MataPelajaran
		if err := rows.Scan(
			&item.IdMapel,
			&item.IdKelas,
			&item.KodeMapel,
			&item.NamaMapel,
			&item.Deskripsi,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning mapel", "layer", "core.service", "op", "matapelajaran.get", "err", err)
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}

func (r *MapelRepo)GetMapelById(ctx context.Context, idMapel int) (matapelajaran.MataPelajaran, error){
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

	rows := r.q.QueryRow(ctx, query, idMapel)

	var item matapelajaran.MataPelajaran
	if err := rows.Scan(
		&item.IdMapel,
		&item.IdKelas,
		&item.KodeMapel,
		&item.NamaMapel,
		&item.Deskripsi,
	); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed scanning mapel", "layer", "core.service", "op", "matapelajaran.get_by_id", "err", err)
		
		if errors.Is(err, pgx.ErrNoRows) {
			return matapelajaran.MataPelajaran{}, coreerror.ErrNotFound
		}
		return matapelajaran.MataPelajaran{}, err
	}
	return item, nil
}

func (r *MapelRepo)	CreateMapel(ctx context.Context, mapel matapelajaran.MataPelajaran) error {
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

func (r *MapelRepo)	UpdateMapel(ctx context.Context,idMapel int, mapel updatepatch.UpdateMapelPatch) error {
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

func (r *MapelRepo)	DeleteMapel(ctx context.Context, idMapel int) error {
	query := `
		DELETE FROM mata_pelajaran
		WHERE id_mapel = $1
	`

	tag, err := r.q.Exec(ctx, query, idMapel)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err,&pgErr){
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

func (r *MapelRepo)ExistKodeMapel(ctx context.Context, kodeMapel string)	(bool,error) {
	query := `SELECT EXISTS(SELECT 1 FROM mata_pelajaran WHERE kode_mapel = $1)`

	var exist bool
	if err := r.q.QueryRow(ctx, query, kodeMapel).Scan(&exist); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed scanning mapel", "layer", "core.service", "op", "matapelajaran.existKodeMapel", "err", err)
		return false, err
	}
	return exist, nil
}


