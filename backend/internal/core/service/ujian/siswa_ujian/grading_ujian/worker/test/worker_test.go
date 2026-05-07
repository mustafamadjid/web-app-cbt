package gradingujian_service_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	gradingujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/worker"
	"github.com/stretchr/testify/assert"
)

type gradingTestLogger struct{}

func (gradingTestLogger) With(...any) corelog.Logger            { return gradingTestLogger{} }
func (gradingTestLogger) Info(context.Context, string, ...any)  {}
func (gradingTestLogger) Error(context.Context, string, ...any) {}

type fakeGradingWorkerRepo struct {
	claimErr   error
	claimJobs  []ujian.GradingJob
	claimCalls int

	updateErr     error
	updates       []ujian.JobStatus
	errorCodes    []string
	errorMessages []string
	updatedJobIDs []int
	mu            sync.Mutex
}

func (f *fakeGradingWorkerRepo) ClaimQueuedJobs(context.Context, int) ([]ujian.GradingJob, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.claimCalls++
	if f.claimErr != nil {
		return nil, f.claimErr
	}
	if f.claimCalls > 1 {
		return nil, nil
	}
	return f.claimJobs, nil
}

func (f *fakeGradingWorkerRepo) UpdateStatusJob(_ context.Context, jobID int, statusJob ujian.JobStatus, errorMsg string, errorCode string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.updatedJobIDs = append(f.updatedJobIDs, jobID)
	f.updates = append(f.updates, statusJob)
	f.errorMessages = append(f.errorMessages, errorMsg)
	f.errorCodes = append(f.errorCodes, errorCode)
	return f.updateErr
}

type fakeGradingExecutor struct {
	gradingErr      error
	statistikErr    error
	gradingCalled   bool
	statistikCalled bool
	gotAttemptID    int
}

func (f *fakeGradingExecutor) GradingUjianPilgan(_ context.Context, idAttempt int) error {
	f.gradingCalled = true
	f.gotAttemptID = idAttempt
	return f.gradingErr
}

func (f *fakeGradingExecutor) StatistikUjian(_ context.Context, idAttempt int) error {
	f.statistikCalled = true
	f.gotAttemptID = idAttempt
	return f.statistikErr
}

type fakeExecutorSvc struct {
	gradingErr      error
	statistikErr    error
	gradingCalled   bool
	statistikCalled bool
}

func (f *fakeExecutorSvc) GradingUjianPilgan(context.Context, int) error {
	f.gradingCalled = true
	return f.gradingErr
}

func (f *fakeExecutorSvc) StatistikUjian(context.Context, int) error {
	f.statistikCalled = true
	return f.statistikErr
}

func TestCompositeGradingUjianExecutor_BasisPath(t *testing.T) {
	t.Parallel()

	t.Run("Path 1 -> grading service belum terpasang", func(t *testing.T) {
		exec := gradingujian_service.NewCompositeGradingUjianExecutor(nil, &fakeExecutorSvc{})
		err := exec.GradingUjianPilgan(context.Background(), 7)
		assert.Error(t, err)
	})

	t.Run("Path 2 -> statistik service belum terpasang", func(t *testing.T) {
		exec := gradingujian_service.NewCompositeGradingUjianExecutor(&fakeExecutorSvc{}, nil)
		err := exec.StatistikUjian(context.Background(), 7)
		assert.Error(t, err)
	})

	t.Run("Path 3 -> composite executor meneruskan panggilan", func(t *testing.T) {
		gradingSvc := &fakeExecutorSvc{}
		statistikSvc := &fakeExecutorSvc{}
		exec := gradingujian_service.NewCompositeGradingUjianExecutor(gradingSvc, statistikSvc)

		errGrading := exec.GradingUjianPilgan(context.Background(), 7)
		errStat := exec.StatistikUjian(context.Background(), 7)

		assert.NoError(t, errGrading)
		assert.NoError(t, errStat)
		assert.True(t, gradingSvc.gradingCalled)
		assert.True(t, statistikSvc.statistikCalled)
	})
}

func TestGradingUjianWorkerService_BasisPath(t *testing.T) {
	t.Parallel()

	runWorker := func(repo *fakeGradingWorkerRepo, executor gradingujian_service.GradingUjianExecutor) {
		t.Helper()
		ctx := corelog.WithLogger(context.Background(), gradingTestLogger{})
		svc := gradingujian_service.NewGradingUjianWorkerService(repo, executor, 10*time.Millisecond)
		timeoutCtx, cancel := context.WithTimeout(ctx, 80*time.Millisecond)
		defer cancel()
		svc.Start(timeoutCtx)
	}

	t.Run("Path 1 -> claim queued jobs gagal", func(t *testing.T) {
		t.Parallel()

		repo := &fakeGradingWorkerRepo{claimErr: errors.New("claim error")}
		runWorker(repo, &fakeGradingExecutor{})

		assert.Empty(t, repo.updates)
	})

	t.Run("Path 2 -> grading service belum terpasang", func(t *testing.T) {
		t.Parallel()

		repo := &fakeGradingWorkerRepo{claimJobs: []ujian.GradingJob{{IDgradingJob: 1, IDAttempt: 10}}}
		runWorker(repo, nil)

		assert.Contains(t, repo.updates, ujian.StatusFailed)
		assert.Contains(t, repo.errorCodes, "GRADING_FAILED")
	})

	t.Run("Path 3 -> grading ujian gagal", func(t *testing.T) {
		t.Parallel()

		repo := &fakeGradingWorkerRepo{claimJobs: []ujian.GradingJob{{IDgradingJob: 2, IDAttempt: 11}}}
		executor := &fakeGradingExecutor{gradingErr: errors.New("grading error")}
		runWorker(repo, executor)

		assert.True(t, executor.gradingCalled)
		assert.Contains(t, repo.updates, ujian.StatusFailed)
	})

	t.Run("Path 4 -> statistik ujian gagal", func(t *testing.T) {
		t.Parallel()

		repo := &fakeGradingWorkerRepo{claimJobs: []ujian.GradingJob{{IDgradingJob: 3, IDAttempt: 12}}}
		executor := &fakeGradingExecutor{statistikErr: errors.New("statistik error")}
		runWorker(repo, executor)

		assert.True(t, executor.gradingCalled)
		assert.True(t, executor.statistikCalled)
		assert.Contains(t, repo.updates, ujian.StatusFailed)
	})

	t.Run("Path 5 -> update status done gagal", func(t *testing.T) {
		t.Parallel()

		repo := &fakeGradingWorkerRepo{
			claimJobs: []ujian.GradingJob{{IDgradingJob: 4, IDAttempt: 13}},
			updateErr: errors.New("update done error"),
		}
		executor := &fakeGradingExecutor{}
		runWorker(repo, executor)

		assert.True(t, executor.gradingCalled)
		assert.True(t, executor.statistikCalled)
		assert.Contains(t, repo.updates, ujian.StatusDone)
	})

	t.Run("Path 6 -> grading job selesai", func(t *testing.T) {
		t.Parallel()

		repo := &fakeGradingWorkerRepo{claimJobs: []ujian.GradingJob{{IDgradingJob: 5, IDAttempt: 14}}}
		executor := &fakeGradingExecutor{}
		runWorker(repo, executor)

		assert.True(t, executor.gradingCalled)
		assert.True(t, executor.statistikCalled)
		assert.Contains(t, repo.updates, ujian.StatusDone)
		assert.Equal(t, 14, executor.gotAttemptID)
	})
}
