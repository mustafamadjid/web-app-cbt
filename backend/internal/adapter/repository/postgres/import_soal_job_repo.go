package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type ImportSoalJobRepo struct {
	q      Executor
	logger corelog.Logger
}

func NewImportSoalJobRepo(q Executor, logger corelog.Logger) *ImportSoalJobRepo {
	return &ImportSoalJobRepo{q: q, logger: logger}
}

func (r *ImportSoalJobRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *ImportSoalJobRepo) CreateJob(ctx context.Context, job importsoal.ImportSoalJob) (int64, error) {
	query := `
		INSERT INTO import_soal_job (id_bank_soal, id_pengguna, status, file_path)
		VALUES ($1, $2, $3, $4)
		RETURNING id_job
	`

	var jobID int64
	err := r.q.QueryRow(ctx, query, job.IDBankSoal, job.IDPengguna, string(job.Status), job.FilePath).Scan(&jobID)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed creating import job", "layer", "repo.db", "op", "import_soal_job.create", "err", err)
		return 0, err
	}

	return jobID, nil
}

func (r *ImportSoalJobRepo) GetPendingJobs(ctx context.Context, limit int) ([]importsoal.ImportSoalJob, error) {
	query := `
		SELECT id_job, id_bank_soal, id_pengguna, status, file_path, 
		       COALESCE(error_msg, ''), COALESCE(total_soal, 0), created_at, updated_at
		FROM import_soal_job
		WHERE status = 'pending'
		ORDER BY created_at ASC
		LIMIT $1
	`

	rows, err := r.q.Query(ctx, query, limit)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed fetching pending jobs", "layer", "repo.db", "op", "import_soal_job.get_pending", "err", err)
		return nil, err
	}
	defer rows.Close()

	var jobs []importsoal.ImportSoalJob
	for rows.Next() {
		var j importsoal.ImportSoalJob
		if err := rows.Scan(
			&j.IDJob, &j.IDBankSoal, &j.IDPengguna, &j.Status, &j.FilePath,
			&j.ErrorMsg, &j.TotalSoal, &j.CreatedAt, &j.UpdatedAt,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning pending job", "layer", "repo.db", "op", "import_soal_job.get_pending", "err", err)
			return nil, err
		}
		jobs = append(jobs, j)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *ImportSoalJobRepo) UpdateJobStatus(ctx context.Context, jobID int64, status importsoal.JobStatus, errorMsg string, totalSoal int) error {
	query := `
		UPDATE import_soal_job
		SET status = $1, error_msg = $2, total_soal = $3, updated_at = now()
		WHERE id_job = $4
	`

	tag, err := r.q.Exec(ctx, query, string(status), errorMsg, totalSoal, jobID)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed updating job status", "layer", "repo.db", "op", "import_soal_job.update_status", "err", err)
		return err
	}
	if tag.RowsAffected() == 0 {
		return coreerror.ErrImportJobNotFound
	}

	return nil
}

func (r *ImportSoalJobRepo) GetJobByID(ctx context.Context, jobID int64) (importsoal.ImportSoalJob, error) {
	query := `
		SELECT id_job, id_bank_soal, id_pengguna, status, file_path,
		       COALESCE(error_msg, ''), COALESCE(total_soal, 0), created_at, updated_at
		FROM import_soal_job
		WHERE id_job = $1
	`

	var j importsoal.ImportSoalJob
	err := r.q.QueryRow(ctx, query, jobID).Scan(
		&j.IDJob, &j.IDBankSoal, &j.IDPengguna, &j.Status, &j.FilePath,
		&j.ErrorMsg, &j.TotalSoal, &j.CreatedAt, &j.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return importsoal.ImportSoalJob{}, coreerror.ErrImportJobNotFound
		}
		r.loggerFor(ctx).Error(ctx, "failed getting job by id", "layer", "repo.db", "op", "import_soal_job.get_by_id", "err", err)
		return importsoal.ImportSoalJob{}, err
	}

	return j, nil
}

func (r *ImportSoalJobRepo) GetJobsByBankSoal(ctx context.Context, bankSoalID int64) ([]importsoal.ImportSoalJob, error) {
	query := `
		SELECT id_job, id_bank_soal, id_pengguna, status, file_path,
		       COALESCE(error_msg, ''), COALESCE(total_soal, 0), created_at, updated_at
		FROM import_soal_job
		WHERE id_bank_soal = $1
		ORDER BY created_at DESC
	`

	rows, err := r.q.Query(ctx, query, bankSoalID)
	if err != nil {
		r.loggerFor(ctx).Error(ctx, "failed getting jobs by bank soal", "layer", "repo.db", "op", "import_soal_job.get_by_bank_soal", "err", err)
		return nil, err
	}
	defer rows.Close()

	var jobs []importsoal.ImportSoalJob
	for rows.Next() {
		var j importsoal.ImportSoalJob
		if err := rows.Scan(
			&j.IDJob, &j.IDBankSoal, &j.IDPengguna, &j.Status, &j.FilePath,
			&j.ErrorMsg, &j.TotalSoal, &j.CreatedAt, &j.UpdatedAt,
		); err != nil {
			r.loggerFor(ctx).Error(ctx, "failed scanning job", "layer", "repo.db", "op", "import_soal_job.get_by_bank_soal", "err", err)
			return nil, err
		}
		jobs = append(jobs, j)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

// --- IsiSoalBatchRepo ---

type IsiSoalBatchRepo struct {
	q      Executor
	logger corelog.Logger
}

func NewIsiSoalBatchRepo(q Executor, logger corelog.Logger) *IsiSoalBatchRepo {
	return &IsiSoalBatchRepo{q: q, logger: logger}
}

func (r *IsiSoalBatchRepo) loggerFor(ctx context.Context) corelog.Logger {
	return corelog.FromContextOr(ctx, r.logger)
}

func (r *IsiSoalBatchRepo) InsertSoalBatch(ctx context.Context, bankSoalID int64, soalList []importsoal.ParsedSoal) error {
	for i, soal := range soalList {
		// Insert isi_soal
		insertSoal := `
			INSERT INTO isi_soal (id_bank_soal, tipe_soal, pertanyaan, gambar, bobot_soal)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id_soal
		`

		var gambar *string
		if soal.Gambar != "" {
			gambar = &soal.Gambar
		}

		var idSoal int64
		err := r.q.QueryRow(ctx, insertSoal, bankSoalID, soal.TipeSoal, soal.Pertanyaan, gambar, soal.BobotSoal).Scan(&idSoal)
		if err != nil {
			r.loggerFor(ctx).Error(ctx, "failed inserting soal", "layer", "repo.db", "op", "isi_soal_batch.insert", "soal_index", i, "err", err)
			return fmt.Errorf("insert soal ke-%d: %w", i+1, err)
		}

		// Insert opsi_pilihan_ganda (only for PG)
		if soal.TipeSoal == "pilihan_ganda" {
			for _, opsi := range soal.Opsi {
				insertOpsi := `
					INSERT INTO opsi_pilihan_ganda (id_soal, isi_pilihan, is_benar)
					VALUES ($1, $2, $3)
				`
				_, err := r.q.Exec(ctx, insertOpsi, idSoal, opsi.Isi, opsi.IsBenar)
				if err != nil {
					r.loggerFor(ctx).Error(ctx, "failed inserting opsi", "layer", "repo.db", "op", "isi_soal_batch.insert_opsi", "soal_index", i, "opsi", opsi.Label, "err", err)
					return fmt.Errorf("insert opsi %s soal ke-%d: %w", opsi.Label, i+1, err)
				}
			}
		}
	}

	return nil
}
