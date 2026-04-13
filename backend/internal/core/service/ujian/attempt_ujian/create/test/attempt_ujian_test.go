package attemptujian_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	attemptujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/create"
	"github.com/stretchr/testify/assert"
)

func TestAttemptUjianService_TreatCreateRaceAsResume(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
	deadline := now.Add(30 * time.Minute)
	checker := &fakeSiswaUjianChecker{
		checkValidRet:   true,
		checkValidIDRet: 77,
		checkTokenRet:   true,
		getDeadlineRet:  deadline,
	}
	attemptRepo := &fakeAttemptUjianRepo{
		createErr: coreerror.ErrSiswaHasActiveAttempt,
	}

	svc := attemptujian_service.NewAttemptUjianService(checker, attemptRepo)

	err := svc.AttemptUjian(context.Background(), 10, 20, "TOKEN", now)

	assert.NoError(t, err)
	assert.True(t, attemptRepo.createCalled)
}

func TestAttemptUjianService_ReturnsCreateErrorWhenNotResume(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
	deadline := now.Add(30 * time.Minute)
	expectedErr := errors.New("create failed")
	checker := &fakeSiswaUjianChecker{
		checkValidRet:   true,
		checkValidIDRet: 77,
		checkTokenRet:   true,
		getDeadlineRet:  deadline,
	}
	attemptRepo := &fakeAttemptUjianRepo{
		createErr: expectedErr,
	}

	svc := attemptujian_service.NewAttemptUjianService(checker, attemptRepo)

	err := svc.AttemptUjian(context.Background(), 10, 20, "TOKEN", now)

	assert.ErrorIs(t, err, expectedErr)
	assert.True(t, attemptRepo.createCalled)
}

type fakeSiswaUjianChecker struct {
	checkValidRet   bool
	checkValidIDRet int
	checkValidErr   error
	checkTokenRet   bool
	checkTokenErr   error
	getDeadlineRet  time.Time
	getDeadlineErr  error
}

func (f *fakeSiswaUjianChecker) CheckAttemptOwnershipBySiswa(context.Context, int, ujian.ID) (bool, error) {
	return false, nil
}

func (f *fakeSiswaUjianChecker) CheckValidSiswaInPesertaUjianById(context.Context, int, int) (bool, int, error) {
	return f.checkValidRet, f.checkValidIDRet, f.checkValidErr
}

func (f *fakeSiswaUjianChecker) CheckTokenUjian(context.Context, string, int) (bool, error) {
	return f.checkTokenRet, f.checkTokenErr
}

func (f *fakeSiswaUjianChecker) GetDeadlineUjian(context.Context, int) (time.Time, error) {
	return f.getDeadlineRet, f.getDeadlineErr
}

type fakeAttemptUjianRepo struct {
	createErr    error
	createCalled bool
}

func (f *fakeAttemptUjianRepo) GetAttemptUjianById(context.Context, ujian.ID) (ujian.AttemptUjian, error) {
	return ujian.AttemptUjian{}, nil
}

func (f *fakeAttemptUjianRepo) CreateAttemptUjian(context.Context, ujian.AttemptUjian) error {
	f.createCalled = true
	return f.createErr
}

func (f *fakeAttemptUjianRepo) UpdateAttemptUjian(context.Context, ujian.ID, updatepatch.UpdateAttemptUjianPatch) error {
	return nil
}

func (f *fakeAttemptUjianRepo) DeleteAttemptUjian(context.Context, ujian.ID) error {
	return nil
}

func (f *fakeAttemptUjianRepo) SubmitAttemptUjian(context.Context, ujian.ID) error {
	return nil
}

func (f *fakeAttemptUjianRepo) ListPesertaUjianAttemptSubmittedByIdJadwalUjian(context.Context, ujian.ID) ([]ujian.PesertaUjianSubmitted, error) {
	return nil, nil
}
