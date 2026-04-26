package worker_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	importsoal "github.com/mustafamadjid/web-app-cbt/internal/core/domain/import_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	importversion "github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/import_version"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/import_soal/worker"
	"github.com/stretchr/testify/assert"
)

// testLogger implements corelog.Logger for testing purposes.
type testLogger struct{}

func (l testLogger) With(_ ...any) corelog.Logger                { return l }
func (l testLogger) Info(_ context.Context, _ string, _ ...any)  {}
func (l testLogger) Error(_ context.Context, _ string, _ ...any) {}

func newWorkerForTest(jobRepo *FakeImportSoalJobRepo, batchRepo *FakeIsiSoalBatchRepo, imageDir string, interval time.Duration) *worker.Worker {
	importSvc := importversion.NewService(batchRepo)
	return worker.NewWorker(jobRepo, importSvc, imageDir, "/uploads/image", interval, testLogger{})
}

func TestWorkerProcessJobs_BasisPath(t *testing.T) {
	t.Parallel()

	t.Run("Path 1 -> GetPendingJobs error", func(t *testing.T) {
		t.Parallel()

		jobRepo := &FakeImportSoalJobRepo{
			GetPendingJobsFn: func(_ context.Context, _ int) ([]importsoal.ImportSoalJob, error) {
				return nil, errors.New("get pending jobs error")
			},
		}
		batchRepo := &FakeIsiSoalBatchRepo{}

		tmpDir := t.TempDir()
		w := newWorkerForTest(jobRepo, batchRepo, tmpDir, 1*time.Second)

		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel()
		w.Start(cancelCtx)
		// Worker should exit without panic when context is cancelled
	})

	t.Run("Path 2 -> no pending jobs", func(t *testing.T) {
		t.Parallel()

		jobRepo := &FakeImportSoalJobRepo{
			GetPendingJobsFn: func(_ context.Context, _ int) ([]importsoal.ImportSoalJob, error) {
				return nil, nil
			},
		}
		batchRepo := &FakeIsiSoalBatchRepo{}

		tmpDir := t.TempDir()
		w := newWorkerForTest(jobRepo, batchRepo, tmpDir, 1*time.Second)

		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel()
		w.Start(cancelCtx)
	})

	t.Run("Path 3 -> processOneJob UpdateJobStatus (mark processing) error", func(t *testing.T) {
		t.Parallel()

		updateErr := errors.New("update status error")
		var mu sync.Mutex
		updateCalled := false

		jobRepo := &FakeImportSoalJobRepo{
			GetPendingJobsFn: func(_ context.Context, _ int) ([]importsoal.ImportSoalJob, error) {
				mu.Lock()
				defer mu.Unlock()
				if !updateCalled {
					return []importsoal.ImportSoalJob{
						{IDJob: 1, IDBankSoal: 1, FilePath: "/tmp/test.docx"},
					}, nil
				}
				return nil, nil
			},
			UpdateJobStatusFn: func(_ context.Context, _ int64, _ importsoal.JobStatus, _, _ string, _ int) error {
				mu.Lock()
				updateCalled = true
				mu.Unlock()
				return updateErr
			},
		}
		batchRepo := &FakeIsiSoalBatchRepo{}

		tmpDir := t.TempDir()
		// Use very short poll interval so ticker fires before timeout
		w := newWorkerForTest(jobRepo, batchRepo, tmpDir, 10*time.Millisecond)

		cancelCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		w.Start(cancelCtx)
	})

	t.Run("Path 4 -> processOneJob ReadFile error", func(t *testing.T) {
		t.Parallel()

		var mu sync.Mutex
		statusUpdates := make(map[importsoal.JobStatus]bool)
		jobReturned := false

		jobRepo := &FakeImportSoalJobRepo{
			GetPendingJobsFn: func(_ context.Context, _ int) ([]importsoal.ImportSoalJob, error) {
				mu.Lock()
				defer mu.Unlock()
				if !jobReturned {
					jobReturned = true
					return []importsoal.ImportSoalJob{
						{IDJob: 1, IDBankSoal: 1, FilePath: "/nonexistent/path/test.docx"},
					}, nil
				}
				return nil, nil
			},
			UpdateJobStatusFn: func(_ context.Context, _ int64, status importsoal.JobStatus, _, _ string, _ int) error {
				mu.Lock()
				statusUpdates[status] = true
				mu.Unlock()
				return nil
			},
		}
		batchRepo := &FakeIsiSoalBatchRepo{}

		tmpDir := t.TempDir()
		w := newWorkerForTest(jobRepo, batchRepo, tmpDir, 10*time.Millisecond)

		cancelCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		w.Start(cancelCtx)

		mu.Lock()
		defer mu.Unlock()
		assert.True(t, statusUpdates[importsoal.StatusProcessing])
		assert.True(t, statusUpdates[importsoal.StatusFailed])
	})

	t.Run("Path 5 -> processOneJob ExtractParagraphs error (invalid docx)", func(t *testing.T) {
		t.Parallel()

		var mu sync.Mutex
		statusUpdates := make(map[importsoal.JobStatus]bool)
		jobReturned := false

		tmpDir := t.TempDir()
		invalidFile := filepath.Join(tmpDir, "invalid.docx")
		err := os.WriteFile(invalidFile, []byte("bukan docx"), 0o644)
		assert.NoError(t, err)

		jobRepo := &FakeImportSoalJobRepo{
			GetPendingJobsFn: func(_ context.Context, _ int) ([]importsoal.ImportSoalJob, error) {
				mu.Lock()
				defer mu.Unlock()
				if !jobReturned {
					jobReturned = true
					return []importsoal.ImportSoalJob{
						{IDJob: 1, IDBankSoal: 1, FilePath: invalidFile},
					}, nil
				}
				return nil, nil
			},
			UpdateJobStatusFn: func(_ context.Context, _ int64, status importsoal.JobStatus, _, _ string, _ int) error {
				mu.Lock()
				statusUpdates[status] = true
				mu.Unlock()
				return nil
			},
		}
		batchRepo := &FakeIsiSoalBatchRepo{}

		w := newWorkerForTest(jobRepo, batchRepo, tmpDir, 10*time.Millisecond)

		cancelCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		w.Start(cancelCtx)

		mu.Lock()
		defer mu.Unlock()
		assert.True(t, statusUpdates[importsoal.StatusProcessing])
		assert.True(t, statusUpdates[importsoal.StatusFailed])
	})
}
