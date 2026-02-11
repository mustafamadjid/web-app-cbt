package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
)

type KelasRepo struct {
	q      Executor
	logger corelog.Logger
}

func NewKelasRepo(q Executor, logger corelog.Logger) *KelasRepo {
	return &KelasRepo{q: q, logger: logger}
}

func (r *KelasRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *KelasRepo) GetKelas(ctx context.Context, filter query.ListKelasFilter) ([]kelas.FullKelasData, error) {
	baseQuery := `
	SELECT 
	tk.id_kelas,
	tk.tingkat_kelas,
	k.id_nama_kelas,
	k.nama_kelas
	FROM kelas tk
	LEFT JOIN nama_kelas k ON tk.id_kelas = k.id_kelas
	`

	where := make([]string, 0, 2)
	args := make([]any, 0, 4)

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("(k.nama_kelas ILIKE $%d OR tk.tingkat_kelas::text ILIKE $%d)", idx, idx))
	}

	if filter.TingkatKelas != nil {
		args = append(args, *filter.TingkatKelas)
		where = append(where, fmt.Sprintf("tk.tingkat_kelas = $%d", len(args)))
	}

	if len(where) > 0 {
		baseQuery = fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(where, " AND "))
	}

	baseQuery = fmt.Sprintf("%s ORDER BY tk.tingkat_kelas ASC, k.nama_kelas ASC", baseQuery)

	if filter.Limit > 0 {
		args = append(args, filter.Limit)
		limitIndex := len(args)
		args = append(args, filter.Offset)
		offsetIndex := len(args)
		baseQuery = fmt.Sprintf("%s LIMIT $%d OFFSET $%d", baseQuery, limitIndex, offsetIndex)
	}

	rows, err := r.q.Query(ctx, baseQuery, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed listing kelas", "op", "kelas_repo.list", "err", err)
		return nil, err
	}
	defer rows.Close()

	var (
		itemsTingkat []kelas.TingkatKelas
		itemsNama    []kelas.NamaKelas
	)

	seenTingkat := make(map[int]bool)
	for rows.Next() {
		var (
			idTingkat int
			tingkat   int
			idNama    sql.NullInt64
			namaKelas sql.NullString
		)

		if err := rows.Scan(&idTingkat, &tingkat, &idNama, &namaKelas); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning kelas", "op", "kelas_repo.scan", "err", err)
			return nil, err
		}

		if !seenTingkat[idTingkat] {
			seenTingkat[idTingkat] = true
			itemsTingkat = append(itemsTingkat, kelas.TingkatKelas{
				IdTingkatKelas: kelas.ID(idTingkat),
				TingkatKelas:   tingkat,
			})
		}

		if idNama.Valid && namaKelas.Valid {
			itemsNama = append(itemsNama, kelas.NamaKelas{
				IdNamaKelas:    kelas.ID(idNama.Int64),
				IdTingkatKelas: kelas.ID(idTingkat),
				NamaKelas:      namaKelas.String,
			})

		}
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating kelas", "op", "kelas_repo.iter", "err", err)
		return nil, err
	}

	return []kelas.FullKelasData{{
		ItemsTingkatKelas: itemsTingkat,
		ItemsNamaKelas:    itemsNama,
	}}, nil
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

	rows := r.q.QueryRow(ctx, query, idTingkatKelas, idNamaKelas)

	var (
		itemsTingkatKelas kelas.TingkatKelas
		itemsNamaKelas    kelas.NamaKelas
	)

	if err := rows.Scan(&itemsTingkatKelas.IdTingkatKelas, &itemsTingkatKelas.TingkatKelas, &itemsNamaKelas.IdNamaKelas, &itemsNamaKelas.NamaKelas); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return kelas.KelasData{}, coreerror.ErrNotFound
		}

		r.loggerFor(ctx).Error(ctx, "failed scanning kelas", "op", "kelas_repo.scan", "err", err)
		return kelas.KelasData{}, err
	}

	itemsNamaKelas.IdTingkatKelas = itemsTingkatKelas.IdTingkatKelas

	return kelas.KelasData{
		ItemsTingkatKelas: itemsTingkatKelas,
		ItemsNamaKelas:    itemsNamaKelas,
	}, nil
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
