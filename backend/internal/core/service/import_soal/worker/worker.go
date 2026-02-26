package worker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	importsoal_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/import_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/parser"
)

type Worker struct {
	jobRepo       importsoal_repo.ImportSoalJobRepo
	soalBatchRepo importsoal_repo.IsiSoalBatchRepo
	imageDir      string // directory to store extracted images (e.g. uploads/image/bank_soal)
	imageRoute    string // route prefix for serving images (e.g. /uploads/image/bank_soal)
	pollInterval  time.Duration
	logger        corelog.Logger
}

func NewWorker(
	jobRepo importsoal_repo.ImportSoalJobRepo,
	soalBatchRepo importsoal_repo.IsiSoalBatchRepo,
	imageDir string,
	imageRoute string,
	pollInterval time.Duration,
	logger corelog.Logger,
) *Worker {
	return &Worker{
		jobRepo:       jobRepo,
		soalBatchRepo: soalBatchRepo,
		imageDir:      imageDir,
		imageRoute:    imageRoute,
		pollInterval:  pollInterval,
		logger:        logger,
	}
}

// Start runs the worker loop in a blocking fashion. Call this in a goroutine.
func (w *Worker) Start(ctx context.Context) {
	w.logger.Info(ctx, "import soal worker started", "poll_interval", w.pollInterval.String())

	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info(ctx, "import soal worker stopped")
			return
		case <-ticker.C:
			w.processJobs(ctx)
		}
	}
}

func (w *Worker) processJobs(ctx context.Context) {
	jobs, err := w.jobRepo.GetPendingJobs(ctx, 5)
	if err != nil {
		w.logger.Error(ctx, "failed fetching pending jobs", "layer", "worker", "err", err)
		return
	}

	for _, job := range jobs {
		w.processOneJob(ctx, job)
	}
}

func (w *Worker) processOneJob(ctx context.Context, job importsoal.ImportSoalJob) {
	logger := w.logger.With("job_id", job.IDJob, "bank_soal_id", job.IDBankSoal)

	// 1. Mark as processing
	if err := w.jobRepo.UpdateJobStatus(ctx, job.IDJob, importsoal.StatusProcessing, "", 0); err != nil {
		logger.Error(ctx, "failed updating job to processing", "err", err)
		return
	}

	// 2. Read file from disk
	data, err := os.ReadFile(job.FilePath)
	if err != nil {
		errMsg := fmt.Sprintf("gagal membaca file: %v", err)
		logger.Error(ctx, errMsg, "file_path", job.FilePath)
		_ = w.jobRepo.UpdateJobStatus(ctx, job.IDJob, importsoal.StatusFailed, errMsg, 0)
		return
	}

	// 3. Extract paragraphs from DOCX
	paragraphs, err := parser.ExtractParagraphs(data)
	if err != nil {
		errMsg := fmt.Sprintf("gagal extract paragraf dari docx: %v", err)
		logger.Error(ctx, errMsg)
		_ = w.jobRepo.UpdateJobStatus(ctx, job.IDJob, importsoal.StatusFailed, errMsg, 0)
		return
	}

	// 4. Parse markers
	soalList, err := parser.ParseMarkers(paragraphs)
	if err != nil {
		errMsg := fmt.Sprintf("gagal parsing marker: %v", err)
		logger.Error(ctx, errMsg)
		_ = w.jobRepo.UpdateJobStatus(ctx, job.IDJob, importsoal.StatusFailed, errMsg, 0)
		return
	}

	// 5. Validate parsed soal
	if err := parser.ValidateParsedSoal(soalList); err != nil {
		errMsg := fmt.Sprintf("validasi soal gagal: %v", err)
		logger.Error(ctx, errMsg)
		_ = w.jobRepo.UpdateJobStatus(ctx, job.IDJob, importsoal.StatusFailed, errMsg, 0)
		return
	}

	// 6. Extract and save images if any soal references them
	if err := w.saveImages(data, soalList); err != nil {
		errMsg := fmt.Sprintf("gagal menyimpan gambar: %v", err)
		logger.Error(ctx, errMsg)
		_ = w.jobRepo.UpdateJobStatus(ctx, job.IDJob, importsoal.StatusFailed, errMsg, 0)
		return
	}

	// 7. Insert batch into DB
	if err := w.soalBatchRepo.InsertSoalBatch(ctx, job.IDBankSoal, soalList); err != nil {
		errMsg := fmt.Sprintf("gagal insert soal ke database: %v", err)
		logger.Error(ctx, errMsg)
		_ = w.jobRepo.UpdateJobStatus(ctx, job.IDJob, importsoal.StatusFailed, errMsg, 0)
		return
	}

	// 8. Mark as completed
	if err := w.jobRepo.UpdateJobStatus(ctx, job.IDJob, importsoal.StatusCompleted, "", len(soalList)); err != nil {
		logger.Error(ctx, "failed updating job to completed", "err", err)
		return
	}

	logger.Info(ctx, "import job completed", "total_soal", len(soalList))
}

// saveImages extracts images referenced by [IMG] markers from the DOCX and
// saves them to the bank_soal image directory.
func (w *Worker) saveImages(docxData []byte, soalList []importsoal.ParsedSoal) error {
	// Collect referenced image names
	needed := make(map[string]bool)
	for _, s := range soalList {
		if s.Gambar != "" {
			needed[s.Gambar] = true
		}
	}
	if len(needed) == 0 {
		return nil
	}

	// Extract all images from DOCX ZIP
	images, err := parser.ExtractImageFiles(docxData)
	if err != nil {
		return fmt.Errorf("extract images: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(w.imageDir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", w.imageDir, err)
	}

	// Save referenced images and update soal gambar paths
	for i := range soalList {
		if soalList[i].Gambar == "" {
			continue
		}
		imgName := soalList[i].Gambar
		imgData, ok := images[imgName]
		if !ok {
			return fmt.Errorf("gambar %q tidak ditemukan di dalam file docx", imgName)
		}

		dstPath := filepath.Join(w.imageDir, imgName)
		if err := os.WriteFile(dstPath, imgData, 0o644); err != nil {
			return fmt.Errorf("write image %s: %w", imgName, err)
		}

		// Update path to relative route
		soalList[i].Gambar = strings.TrimRight(w.imageRoute, "/") + "/" + imgName
	}

	return nil
}
