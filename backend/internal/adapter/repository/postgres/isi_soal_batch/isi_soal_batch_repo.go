package isisoalbatchrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	importsoalrepo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/import_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type IsiSoalBatchRepo struct {
	pool   *pgxpool.Pool
	logger corelog.Logger
}

func NewIsiSoalBatchRepo(pool *pgxpool.Pool, logger corelog.Logger) *IsiSoalBatchRepo {
	return &IsiSoalBatchRepo{pool: pool, logger: logger}
}

func (r *IsiSoalBatchRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

const insertBatchSize = 500

func (r *IsiSoalBatchRepo) ImportBankSoalVersion(ctx context.Context, bankID, userID int64, payload importsoalrepo.ImportBankSoalVersionPayload) (int64, error) {
	logger := r.loggerFor(ctx)

	if bankID <= 0 {
		return 0, fmt.Errorf("%w: bank id must be greater than 0", coreerror.ErrInvalidInput)
	}
	if userID <= 0 {
		return 0, fmt.Errorf("%w: user id must be greater than 0", coreerror.ErrInvalidInput)
	}
	if len(payload.SoalList) == 0 {
		return 0, fmt.Errorf("%w: payload is empty", coreerror.ErrInvalidInput)
	}
	if err := validateSingleCorrectOption(payload.SoalList); err != nil {
		return 0, err
	}

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		logger.Error(ctx, "failed begin tx import bank soal version", "layer", "repo.db", "op", "import_soal_version.begin_tx", "err", err)
		return 0, err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	var activeVersionID sql.NullInt64
	err = tx.QueryRow(ctx, `
		SELECT id_bank_soal_version_aktif
		FROM bank_soal
		WHERE id_bank_soal = $1
		FOR UPDATE
	`, bankID).Scan(&activeVersionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, coreerror.ErrBankSoalNotFound
		}
		logger.Error(ctx, "failed lock bank soal row", "layer", "repo.db", "op", "import_soal_version.lock_bank_soal", "bank_id", bankID, "err", err)
		return 0, err
	}

	var newVersionID int64
	err = tx.QueryRow(ctx, `
		INSERT INTO bank_soal_version (id_bank_soal, version_no, status, created_by)
		SELECT
			$1,
			COALESCE(MAX(version_no), 0) + 1,
			'draft',
			$2
		FROM bank_soal_version
		WHERE id_bank_soal = $1
		RETURNING id_bank_soal_version
	`, bankID, userID).Scan(&newVersionID)
	if err != nil {
		logger.Error(ctx, "failed create draft bank soal version", "layer", "repo.db", "op", "import_soal_version.create_draft", "bank_id", bankID, "err", err)
		return 0, normalizeImportVersionErr(err)
	}

	soalIDs, err := r.insertIsiSoalInBatches(ctx, tx, newVersionID, payload.SoalList)
	if err != nil {
		return 0, normalizeImportVersionErr(err)
	}

	if err := r.insertOpsiInBatches(ctx, tx, soalIDs, payload.SoalList); err != nil {
		return 0, normalizeImportVersionErr(err)
	}

	if activeVersionID.Valid {
		_, err = tx.Exec(ctx, `
			UPDATE bank_soal_version
			SET status = 'archived'
			WHERE id_bank_soal_version = $1
			  AND id_bank_soal = $2
		`, activeVersionID.Int64, bankID)
		if err != nil {
			logger.Error(ctx, "failed archive previous active version", "layer", "repo.db", "op", "import_soal_version.archive_old", "bank_id", bankID, "old_version_id", activeVersionID.Int64, "err", err)
			return 0, normalizeImportVersionErr(err)
		}
	}

	_, err = tx.Exec(ctx, `
		UPDATE bank_soal_version
		SET status = 'published'
		WHERE id_bank_soal_version = $1
		  AND id_bank_soal = $2
	`, newVersionID, bankID)
	if err != nil {
		logger.Error(ctx, "failed publish new bank soal version", "layer", "repo.db", "op", "import_soal_version.publish", "bank_id", bankID, "new_version_id", newVersionID, "err", err)
		return 0, normalizeImportVersionErr(err)
	}

	tag, err := tx.Exec(ctx, `
		UPDATE bank_soal
		SET id_bank_soal_version_aktif = $1, updated_at = now()
		WHERE id_bank_soal = $2
	`, newVersionID, bankID)
	if err != nil {
		logger.Error(ctx, "failed set active bank soal version", "layer", "repo.db", "op", "import_soal_version.update_active", "bank_id", bankID, "new_version_id", newVersionID, "err", err)
		return 0, normalizeImportVersionErr(err)
	}
	if tag.RowsAffected() == 0 {
		return 0, coreerror.ErrBankSoalNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error(ctx, "failed commit import bank soal version", "layer", "repo.db", "op", "import_soal_version.commit", "bank_id", bankID, "new_version_id", newVersionID, "err", err)
		return 0, normalizeImportVersionErr(err)
	}

	committed = true
	return newVersionID, nil
}

func (r *IsiSoalBatchRepo) insertIsiSoalInBatches(ctx context.Context, tx pgx.Tx, versionID int64, soalList []importsoal.ParsedSoal) ([]int64, error) {
	soalIDs := make([]int64, len(soalList))

	for start := 0; start < len(soalList); start += insertBatchSize {
		end := minInt(start+insertBatchSize, len(soalList))

		batch := &pgx.Batch{}
		for i := start; i < end; i++ {
			soal := soalList[i]
			noUrut := soal.NoUrutSoal
			if noUrut <= 0 {
				noUrut = i + 1
			}

			batch.Queue(`
				INSERT INTO isi_soal (id_bank_soal_version, tipe_soal, pertanyaan, gambar, bobot_soal, no_urut_soal, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, now(), now())
				RETURNING id_soal
			`, versionID, soal.TipeSoal, soal.Pertanyaan, normalizeSoalImage(soal.Gambar), soal.BobotSoal, noUrut)
		}

		results := tx.SendBatch(ctx, batch)
		for i := start; i < end; i++ {
			if err := results.QueryRow().Scan(&soalIDs[i]); err != nil {
				_ = results.Close()
				return nil, fmt.Errorf("insert isi_soal batch failed at index %d: %w", i, err)
			}
		}
		if err := results.Close(); err != nil {
			return nil, fmt.Errorf("close isi_soal batch failed: %w", err)
		}
	}

	return soalIDs, nil
}

func (r *IsiSoalBatchRepo) insertOpsiInBatches(ctx context.Context, tx pgx.Tx, soalIDs []int64, soalList []importsoal.ParsedSoal) error {
	type opsiInsert struct {
		SoalID    int64
		Isi       string
		IsBenar   bool
		SoalIdx   int
		OpsiLabel string
	}

	rows := make([]opsiInsert, 0)
	for i, soal := range soalList {
		if soal.TipeSoal != "pilihan_ganda" {
			continue
		}
		for _, opsi := range soal.Opsi {
			rows = append(rows, opsiInsert{
				SoalID:    soalIDs[i],
				Isi:       opsi.Isi,
				IsBenar:   opsi.IsBenar,
				SoalIdx:   i,
				OpsiLabel: opsi.Label,
			})
		}
	}

	for start := 0; start < len(rows); start += insertBatchSize {
		end := minInt(start+insertBatchSize, len(rows))

		batch := &pgx.Batch{}
		for i := start; i < end; i++ {
			row := rows[i]
			batch.Queue(`
				INSERT INTO opsi_pilihan_ganda (id_soal, isi_pilihan, is_benar, created_at, updated_at)
				VALUES ($1, $2, $3, now(), now())
			`, row.SoalID, row.Isi, row.IsBenar)
		}

		results := tx.SendBatch(ctx, batch)
		for i := start; i < end; i++ {
			if _, err := results.Exec(); err != nil {
				_ = results.Close()
				row := rows[i]
				return fmt.Errorf("insert opsi gagal pada soal index %d opsi %q: %w", row.SoalIdx, row.OpsiLabel, err)
			}
		}
		if err := results.Close(); err != nil {
			return fmt.Errorf("close opsi batch failed: %w", err)
		}
	}

	return nil
}

func normalizeImportVersionErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return coreerror.ErrConflict
		}
	}
	return err
}
