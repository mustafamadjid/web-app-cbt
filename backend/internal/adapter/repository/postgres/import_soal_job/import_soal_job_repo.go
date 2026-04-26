package importsoaljobrepo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	pg "github.com/mustafamadjid/web-app-cbt/internal/adapter/repository/postgres/contract"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type ImportSoalJobRepo struct {
	q      pg.Executor
	logger corelog.Logger
}

func NewImportSoalJobRepo(q pg.Executor, logger corelog.Logger) *ImportSoalJobRepo {
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" && pgErr.ConstraintName == "fk_job_bank_soal" {
			return 0, coreerror.ErrBankSoalNotFound
		}
		r.loggerFor(ctx).Error(ctx, "failed creating import job", "layer", "repo.db", "op", "import_soal_job.create", "err", err)
		return 0, err
	}

	return jobID, nil
}

func (r *ImportSoalJobRepo) GetPendingJobs(ctx context.Context, limit int) ([]importsoal.ImportSoalJob, error) {
	query := `
		SELECT id_job, id_bank_soal, id_pengguna, status, file_path, 
		       COALESCE(error_msg, ''), COALESCE(warning_msg, ''), COALESCE(total_soal, 0), created_at, updated_at
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

	return r.scanImportSoalJobRows(ctx, "import_soal_job.get_pending", rows)
}

func (r *ImportSoalJobRepo) UpdateJobStatus(ctx context.Context, jobID int64, status importsoal.JobStatus, errorMsg, warningMsg string, totalSoal int) error {
	query := `
		UPDATE import_soal_job
		SET status = $1, error_msg = $2, warning_msg = $3, total_soal = $4, updated_at = now()
		WHERE id_job = $5
	`

	tag, err := r.q.Exec(ctx, query, string(status), errorMsg, warningMsg, totalSoal, jobID)
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
		       COALESCE(error_msg, ''), COALESCE(warning_msg, ''), COALESCE(total_soal, 0), created_at, updated_at
		FROM import_soal_job
		WHERE id_job = $1
	`

	j, err := scanImportSoalJobRow(r.q.QueryRow(ctx, query, jobID))
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
		       COALESCE(error_msg, ''), COALESCE(warning_msg, ''), COALESCE(total_soal, 0), created_at, updated_at
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

	return r.scanImportSoalJobRows(ctx, "import_soal_job.get_by_bank_soal", rows)
}
