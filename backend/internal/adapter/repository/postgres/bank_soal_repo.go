package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
)

type BankSoalRepo struct {
	q      Executor
	logger corelog.Logger
}

func NewBankSoalRepo(q Executor, logger corelog.Logger) *BankSoalRepo {
	return &BankSoalRepo{q: q, logger: logger}
}

func (r *BankSoalRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

const bankSoalSelectColumns = `
	SELECT id_bank_soal,id_mapel,id_kelas,id_pengguna,nama_bank_soal,deskripsi,materi,created_at,(id_bank_soal_version_aktif IS NOT NULL) AS soal_uploaded
	FROM bank_soal
`

func (r *BankSoalRepo) buildListBankSoalQuery(filter query.BankSoalFilter, uploadedOnly bool) (string, []any) {
	baseQuery := bankSoalSelectColumns
	where := make([]string, 0, 4)
	args := make([]any, 0, 4)

	if uploadedOnly {
		where = append(where, "id_bank_soal_version_aktif IS NOT NULL")
	}

	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		idx := len(args)
		where = append(where, fmt.Sprintf("(nama_bank_soal ILIKE $%d OR deskripsi ILIKE $%d OR materi ILIKE $%d)", idx, idx, idx))
	}

	if filter.TingkatKelas != nil {
		args = append(args, *filter.TingkatKelas)
		where = append(where, fmt.Sprintf("id_kelas = $%d", len(args)))
	}

	if filter.Mapel != nil {
		args = append(args, *filter.Mapel)
		where = append(where, fmt.Sprintf("id_mapel = $%d", len(args)))
	}

	if len(where) > 0 {
		baseQuery = fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(where, " AND "))
	}

	baseQuery = fmt.Sprintf("%s ORDER BY created_at ASC", baseQuery)

	if filter.Limit > 0 {
		args = append(args, filter.Limit)
		limitIndex := len(args)
		args = append(args, filter.Offset)
		offsetIndex := len(args)
		baseQuery = fmt.Sprintf("%s LIMIT $%d OFFSET $%d", baseQuery, limitIndex, offsetIndex)
	}

	return baseQuery, args
}

func (r *BankSoalRepo) scanBankSoalRows(ctx context.Context, op string, rows pgx.Rows) ([]bank_soal.BankSoal, error) {
	var results []bank_soal.BankSoal
	for rows.Next() {
		var item bank_soal.BankSoal
		if err := rows.Scan(
			&item.IdBankSoal,
			&item.IdMapel,
			&item.IdKelas,
			&item.IdPengguna,
			&item.NamaBankSoal,
			&item.Deskripsi,
			&item.Materi,
			&item.CreatedAt,
			&item.SoalUploaded,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scan bank soal", "layer", "repo.db", "op", op, "err", err)
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		r.loggerFor(ctx).Error(ctx, "failed iterating bank soal rows", "layer", "repo.db", "op", op, "err", err)
		return nil, err
	}

	return results, nil
}

func (r *BankSoalRepo) GetBankSoal(ctx context.Context, filter query.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	queryText, args := r.buildListBankSoalQuery(filter, false)

	rows, err := r.q.Query(ctx, queryText, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get bank soal", "layer", "repo.db", "op", "bank_soal.get", "err", err)
		return nil, err
	}
	defer rows.Close()

	results, err := r.scanBankSoalRows(ctx, "bank_soal.get", rows)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *BankSoalRepo) GetBankSoalUploaded(ctx context.Context, filter query.BankSoalFilter) ([]bank_soal.BankSoal, error) {
	queryText, args := r.buildListBankSoalQuery(filter, true)

	rows, err := r.q.Query(ctx, queryText, args...)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get uploaded bank soal", "layer", "repo.db", "op", "bank_soal.get_uploaded", "err", err)
		return nil, err
	}
	defer rows.Close()

	results, err := r.scanBankSoalRows(ctx, "bank_soal.get_uploaded", rows)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *BankSoalRepo) GetBankSoalByGuru(ctx context.Context, idPengguna bank_soal.ID) ([]bank_soal.BankSoal, error) {
	const queryText = `
		SELECT id_bank_soal,id_mapel,id_kelas,id_pengguna,nama_bank_soal,deskripsi,materi,created_at,(id_bank_soal_version_aktif IS NOT NULL) AS soal_uploaded
		FROM bank_soal
		WHERE id_pengguna = $1
		ORDER BY created_at ASC
	`

	rows, err := r.q.Query(ctx, queryText, idPengguna)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed get bank soal by guru", "layer", "repo.db", "op", "bank_soal.get_by_guru", "err", err)
		return nil, err
	}
	defer rows.Close()

	results, err := r.scanBankSoalRows(ctx, "bank_soal.get_by_guru", rows)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *BankSoalRepo) GetBankSoalById(ctx context.Context, idBankSoal bank_soal.ID) (bank_soal.BankSoal, error) {
	const queryText = `
		SELECT id_bank_soal,id_mapel,id_kelas,id_pengguna,nama_bank_soal,deskripsi,materi,created_at,(id_bank_soal_version_aktif IS NOT NULL) AS soal_uploaded
		FROM bank_soal
		WHERE id_bank_soal = $1
	`

	var item bank_soal.BankSoal
	if err := r.q.QueryRow(ctx, queryText, idBankSoal).Scan(
		&item.IdBankSoal,
		&item.IdMapel,
		&item.IdKelas,
		&item.IdPengguna,
		&item.NamaBankSoal,
		&item.Deskripsi,
		&item.Materi,
		&item.CreatedAt,
		&item.SoalUploaded,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return bank_soal.BankSoal{}, coreerror.ErrNotFound
		}
		r.loggerFor(ctx).Error(ctx, "failed get bank soal by id", "layer", "repo.db", "op", "bank_soal.get_by_id", "err", err)
		return bank_soal.BankSoal{}, err
	}

	return item, nil
}

func (r *BankSoalRepo) CreateBankSoal(ctx context.Context, bankSoal bank_soal.BankSoal) error {
	const query = `
		INSERT INTO bank_soal (id_mapel,id_kelas,id_pengguna,nama_bank_soal,deskripsi,materi) 
		VALUES ($1,$2,$3,$4,$5,$6)
	`

	_, err := r.q.Exec(ctx, query,
		bankSoal.IdMapel,
		bankSoal.IdKelas,
		bankSoal.IdPengguna,
		bankSoal.NamaBankSoal,
		bankSoal.Deskripsi,
		bankSoal.Materi,
	)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed insert bank soal", "layer", "repo.db", "op", "bank_soal.create", "err", err)
		return err
	}
	return nil
}

func (r *BankSoalRepo) UpdateBankSoal(ctx context.Context, idBankSoal bank_soal.ID, bankSoal updatepatch.UpdateBankSoalPatch) error {
	const query = `
		UPDATE bank_soal
		SET
			id_mapel = COALESCE($1, id_mapel),
			id_kelas = COALESCE($2, id_kelas),
			id_pengguna = COALESCE($3, id_pengguna),
			nama_bank_soal = COALESCE($4, nama_bank_soal),
			deskripsi = COALESCE($5, deskripsi),
			materi = COALESCE($6, materi),
			updated_at = now()
		WHERE id_bank_soal = $7
	`

	tag, err := r.q.Exec(
		ctx,
		query,
		bankSoal.IdMapel,
		bankSoal.IdKelas,
		bankSoal.IdPengguna,
		bankSoal.NamaBankSoal,
		bankSoal.Deskripsi,
		bankSoal.Materi,
		idBankSoal,
	)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed update bank soal", "layer", "repo.db", "op", "bank_soal.update", "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	return nil
}

func (r *BankSoalRepo) DeleteBankSoal(ctx context.Context, idBankSoal bank_soal.ID) error {
	const query = `
		DELETE FROM bank_soal
		WHERE id_bank_soal = $1
	`

	tag, err := r.q.Exec(ctx, query, idBankSoal)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23001" {
				return coreerror.ErrDeleteRestricted
			}
		}
		r.loggerFor(ctx).Error(ctx, "failed delete bank soal", "layer", "repo.db", "op", "bank_soal.delete", "err", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return coreerror.ErrNotFound
	}

	return nil
}
