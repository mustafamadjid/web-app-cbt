package banksoalrepo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/bank_soal"
)

type BankSoalRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewBankSoalRepo(q pg.Executor, logger corelog.Logger) *BankSoalRepo {
	return &BankSoalRepo{q: q, logger: logger}
}

func (r *BankSoalRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
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
		SELECT 
			b.id_bank_soal,
			b.id_mapel,
			b.id_kelas,
			b.id_pengguna,
			b.nama_bank_soal,
			b.deskripsi,
			b.materi,
			b.created_at,
			(b.id_bank_soal_version_aktif IS NOT NULL) AS soal_uploaded,
			k.tingkat_kelas,
			m.nama_mapel,
			p.nama_lengkap
		FROM bank_soal b
		JOIN kelas k ON b.id_kelas = k.id_kelas
		JOIN mata_pelajaran m ON b.id_mapel = m.id_mapel
		JOIN pengguna p ON b.id_pengguna = p.id_pengguna
		WHERE b.id_pengguna = $1
		ORDER BY b.created_at ASC
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
		SELECT 
			b.id_bank_soal,
			b.id_mapel,
			b.id_kelas,
			b.id_pengguna,
			b.nama_bank_soal,
			b.deskripsi,
			b.materi,
			b.created_at,
			(b.id_bank_soal_version_aktif IS NOT NULL) AS soal_uploaded,
			k.tingkat_kelas,
			m.nama_mapel,
			p.nama_lengkap
		FROM bank_soal b
		JOIN kelas k ON b.id_kelas = k.id_kelas
		JOIN mata_pelajaran m ON b.id_mapel = m.id_mapel
		JOIN pengguna p ON b.id_pengguna = p.id_pengguna
		WHERE b.id_bank_soal = $1
	`

	item, err := scanBankSoalRow(r.q.QueryRow(ctx, queryText, idBankSoal))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return bank_soal.BankSoal{}, coreerror.ErrNotFound
		}
		r.loggerFor(ctx).Error(ctx, "failed get bank soal by id", "layer", "repo.db", "op", "bank_soal.get_by_id", "err", err)
		return bank_soal.BankSoal{}, err
	}

	return item, nil
}

func (r *BankSoalRepo) GetIdBankSoalByAttemptId(ctx context.Context, idAttempt ujian.ID) (ujian.ID, error) {
	const queryText = `
		SELECT u.id_bank_soal
		FROM attempt_ujian au
		JOIN peserta_ujian pu ON pu.id_peserta_ujian = au.id_peserta_ujian
		JOIN jadwal_ujian ju ON ju.id_jadwal_ujian = pu.id_jadwal_ujian
		JOIN ujian u ON u.id_ujian = ju.id_ujian
		WHERE au.id_attempt = $1
	`

	var idBankSoal int
	if err := r.q.QueryRow(ctx, queryText, idAttempt).Scan(&idBankSoal); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			r.loggerFor(ctx).Error(ctx, "failed get bank soal id by attempt id", "layer", "repo.db", "op", "bank_soal.get_id_by_attempt_id", "attempt_id", idAttempt, "err", err)
		}
		return 0, err
	}

	return ujian.ID(idBankSoal), nil
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
