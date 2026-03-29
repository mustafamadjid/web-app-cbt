package kelasrepo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
)

type KelasRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewKelasRepo(q pg.Executor, logger corelog.Logger) *KelasRepo {
	return &KelasRepo{q: q, logger: logger}
}

func (r *KelasRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *KelasRepo) GetKelas(ctx context.Context, filter query.ListKelasFilter) ([]kelas.FullKelasData, error) {
	queryText, args := r.buildListKelasQuery(filter)

	rows, err := r.q.Query(ctx, queryText, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed listing kelas", "op", "kelas_repo.list", "err", err)
		return nil, err
	}
	defer rows.Close()

	return r.scanKelasRows(ctx, "kelas_repo.list", rows)
}

func (r *KelasRepo) GetKelasById(ctx context.Context, idTingkatKelas int, idNamaKelas int) (kelas.KelasData, error) {
	query := `
		SELECT 
		tk.id_kelas,
		tk.tingkat_kelas,
		k.id_nama_kelas,
		k.nama_kelas
		FROM kelas tk
		INNER JOIN nama_kelas k ON tk.id_kelas = k.id_kelas
		WHERE tk.id_kelas = $1 AND k.id_nama_kelas = $2
	`

	item, err := scanKelasRow(r.q.QueryRow(ctx, query, idTingkatKelas, idNamaKelas))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return kelas.KelasData{}, coreerror.ErrNotFound
		}

		r.loggerFor(ctx).Error(ctx, "failed scanning kelas", "op", "kelas_repo.scan", "err", err)
		return kelas.KelasData{}, err
	}

	return item, nil
}

func (r *KelasRepo) CreateTingkatKelas(ctx context.Context, tingkatKelas int) error {
	const query = `
		INSERT INTO kelas (tingkat_kelas) VALUES ($1)
		RETURNING id_kelas
	`
	var kelasId kelas.ID
	err := r.q.QueryRow(
		ctx,
		query,
		tingkatKelas,
	).Scan(&kelasId)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed creating tingkat kelas", "op", "kelas_repo.create", "err", err)
		return err
	}

	return nil
}

func (r *KelasRepo) CreateNamaKelas(ctx context.Context, namaKelas kelas.NamaKelas) error {
	const query = `
		INSERT INTO nama_kelas (id_kelas, nama_kelas) VALUES ($1, $2)
		RETURNING id_nama_kelas
	`
	var kelasId kelas.ID
	err := r.q.QueryRow(
		ctx,
		query,
		namaKelas.IdTingkatKelas,
		namaKelas.NamaKelas,
	).Scan(&kelasId)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed creating nama kelas", "op", "kelas_repo.create", "err", err)
		return err
	}

	return nil
}

func (r *KelasRepo) ExistTingkatKelas(ctx context.Context, tingkatKelas int) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM kelas
			WHERE tingkat_kelas = $1
		)
	`

	var exists bool
	err := r.q.QueryRow(
		ctx,
		query,
		tingkatKelas,
	).Scan(&exists)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed checking tingkat kelas", "op", "kelas_repo.check_exist_tingkat_kelas", "err", err)
		return false, err
	}

	return exists, nil
}
func (r *KelasRepo) ExistNamaKelas(ctx context.Context, namaKelas string) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM nama_kelas
			WHERE nama_kelas = $1
		)
	`

	var exist bool
	err := r.q.QueryRow(
		ctx,
		query,
		namaKelas,
	).Scan(&exist)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed checking nama kelas", "op", "kelas_repo.check_exist_nama_kelas", "err", err)
		return false, err
	}

	return exist, nil
}

func (r *KelasRepo) UpdateNamaKelas(ctx context.Context, idNamaKelas int, dataUpdate updatepatch.NamaKelasPatch) error {
	const query = `
		UPDATE nama_kelas
		SET
			id_kelas = COALESCE($1, id_kelas),
			nama_kelas = COALESCE($2, nama_kelas),
			updated_at = now()
		WHERE id_nama_kelas = $3
	`

	tag, err := r.q.Exec(
		ctx,
		query,
		dataUpdate.IdTingkatKelas,
		dataUpdate.NamaKelas,
		idNamaKelas,
	)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed updating nama kelas", "op", "kelas_repo.update", "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	return nil
}

func (r *KelasRepo) DeleteNamaKelas(ctx context.Context, idNamaKelas int) error {
	const query = `
		DELETE FROM nama_kelas
		WHERE id_nama_kelas = $1
	`

	tag, err := r.q.Exec(
		ctx,
		query,
		idNamaKelas,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23001" {
				return coreerror.ErrDeleteRestricted
			}
		}

		r.loggerFor(ctx).Error(ctx, "failed deleting nama kelas", "op", "kelas_repo.delete", "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	return nil
}
