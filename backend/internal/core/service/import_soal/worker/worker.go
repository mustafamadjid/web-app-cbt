package worker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	content "github.com/mustafamadjid/web-app-cbt/internal/core/domain/content"
	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	importsoal_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/import_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	importversion "github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/import_version"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/parser"
)

type Worker struct {
	jobRepo      importsoal_repo.ImportSoalJobRepo
	importSvc    *importversion.Service
	imageDir     string // directory to store extracted images (e.g. uploads/image/bank_soal)
	imageRoute   string // route prefix for serving images (e.g. /uploads/image/bank_soal)
	pollInterval time.Duration
	logger       corelog.Logger
}

func NewWorker(
	jobRepo importsoal_repo.ImportSoalJobRepo,
	importSvc *importversion.Service,
	imageDir string,
	imageRoute string,
	pollInterval time.Duration,
	logger corelog.Logger,
) *Worker {
	return &Worker{
		jobRepo:      jobRepo,
		importSvc:    importSvc,
		imageDir:     imageDir,
		imageRoute:   imageRoute,
		pollInterval: pollInterval,
		logger:       logger,
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
			w.logger.Info(ctx, "import soal worker stopped", "err", ctx.Err())
			return

		case <-ticker.C:
			tickCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			func() {
				defer cancel()
				defer func() {
					if r := recover(); r != nil {
						w.logger.Error(tickCtx, "panic in processJobs",
							"panic", r,
							"stack", string(debug.Stack()),
						)
					}
				}()
				w.processJobs(tickCtx)
			}()
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
	if err := w.jobRepo.UpdateJobStatus(ctx, job.IDJob, importsoal.StatusProcessing, "", "", 0); err != nil {
		logger.Error(ctx, "failed updating job to processing", "err", err)
		return
	}

	// 2. Read file from disk
	data, err := os.ReadFile(job.FilePath)
	if err != nil {
		errMsg := fmt.Sprintf("gagal membaca file: %v", err)
		logger.Error(ctx, errMsg, "file_path", job.FilePath)
		_ = w.jobRepo.UpdateJobStatus(ctx, job.IDJob, importsoal.StatusFailed, errMsg, "", 0)
		return
	}

	// 3. Extract rich paragraphs from DOCX
	paragraphs, warnings, err := parser.ExtractParagraphContents(data)
	if err != nil {
		errMsg := fmt.Sprintf("gagal extract paragraf dari docx: %v", err)
		logger.Error(ctx, errMsg)
		_ = w.jobRepo.UpdateJobStatus(ctx, job.IDJob, importsoal.StatusFailed, errMsg, "", 0)
		return
	}

	// 4. Parse markers
	soalList, parseWarnings, err := parser.ParseMarkersFromContent(paragraphs, data)
	if err != nil {
		errMsg := fmt.Sprintf("gagal parsing marker: %v", err)
		logger.Error(ctx, errMsg)
		_ = w.jobRepo.UpdateJobStatus(ctx, job.IDJob, importsoal.StatusFailed, errMsg, "", 0)
		return
	}
	warnings = append(warnings, parseWarnings...)

	// 5. Validate parsed soal
	if err := parser.ValidateParsedSoal(soalList); err != nil {
		errMsg := fmt.Sprintf("validasi soal gagal: %v", err)
		logger.Error(ctx, errMsg)
		_ = w.jobRepo.UpdateJobStatus(ctx, job.IDJob, importsoal.StatusFailed, errMsg, "", 0)
		return
	}

	// 6. Extract and save images if any soal references them
	if err := w.saveImages(data, soalList); err != nil {
		errMsg := fmt.Sprintf("gagal menyimpan gambar: %v", err)
		logger.Error(ctx, errMsg)
		_ = w.jobRepo.UpdateJobStatus(ctx, job.IDJob, importsoal.StatusFailed, errMsg, "", 0)
		return
	}

	if w.importSvc == nil {
		errMsg := "service import versi bank soal belum terpasang"
		logger.Error(ctx, errMsg)
		_ = w.jobRepo.UpdateJobStatus(ctx, job.IDJob, importsoal.StatusFailed, errMsg, "", 0)
		return
	}

	// 7. Persist parsed questions as a new bank_soal version.
	importResult, err := w.importSvc.Execute(ctx, importversion.Cmd{
		BankID:  job.IDBankSoal,
		UserID:  job.IDPengguna,
		Payload: soalList,
	})
	if err != nil {
		errMsg := fmt.Sprintf("gagal insert soal ke database: %v", err)
		if errors.Is(err, coreerror.ErrConflict) {
			errMsg = "konflik publish versi bank soal, silakan ulangi import"
		}
		logger.Error(ctx, errMsg)
		_ = w.jobRepo.UpdateJobStatus(ctx, job.IDJob, importsoal.StatusFailed, errMsg, "", 0)
		return
	}

	// 8. Mark as completed
	warningMsg := joinWarnings(warnings)
	if err := w.jobRepo.UpdateJobStatus(ctx, job.IDJob, importsoal.StatusCompleted, "", warningMsg, len(soalList)); err != nil {
		logger.Error(ctx, "failed updating job to completed", "err", err)
		return
	}

	logger.Info(ctx, "import job completed", "total_soal", len(soalList), "new_version_id", importResult.VersionID)
}

func joinWarnings(items []string) string {
	if len(items) == 0 {
		return ""
	}
	seen := make(map[string]struct{}, len(items))
	out := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, exists := seen[item]; exists {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return strings.Join(out, "; ")
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
		collectRichContentImages(s.PertanyaanContent, needed)
		for _, opsi := range s.Opsi {
			collectRichContentImages(opsi.IsiContent, needed)
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

	// Save referenced images once.
	for imgName := range needed {
		imgData, ok := images[imgName]
		if !ok {
			return fmt.Errorf("gambar %q tidak ditemukan di dalam file docx", imgName)
		}

		dstPath := filepath.Join(w.imageDir, imgName)
		if err := os.WriteFile(dstPath, imgData, 0o644); err != nil {
			return fmt.Errorf("write image %s: %w", imgName, err)
		}
	}

	routePrefix := strings.TrimRight(w.imageRoute, "/")
	for i := range soalList {
		if soalList[i].Gambar != "" {
			soalList[i].Gambar = routePrefix + "/" + soalList[i].Gambar
		}
		rewriteRichContentImages(&soalList[i].PertanyaanContent, routePrefix)
		for j := range soalList[i].Opsi {
			rewriteRichContentImages(&soalList[i].Opsi[j].IsiContent, routePrefix)
		}
	}

	return nil
}

func collectRichContentImages(value content.RichContent, needed map[string]bool) {
	for _, block := range value.Blocks {
		for _, child := range block.Children {
			if child.Type == "image" && strings.TrimSpace(child.Src) != "" {
				needed[child.Src] = true
			}
		}
	}
}

func rewriteRichContentImages(value *content.RichContent, routePrefix string) {
	if value == nil {
		return
	}
	for blockIdx := range value.Blocks {
		for childIdx := range value.Blocks[blockIdx].Children {
			child := &value.Blocks[blockIdx].Children[childIdx]
			if child.Type != "image" || strings.TrimSpace(child.Src) == "" {
				continue
			}
			child.Src = routePrefix + "/" + child.Src
		}
	}
}
