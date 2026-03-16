package siswaujian_service

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"github.com/stretchr/testify/assert"
)

type fakeSiswaUjianChecker struct {
	checkFn        func(ctx context.Context, idSiswa int, idJadwalUjian int) (bool, int, error)
	tokenFn        func(ctx context.Context, token string, idJadwalUjian int) (bool, error)
	deadlineFn     func(ctx context.Context, idJadwalUjian int) (time.Time, error)
	gotCheckSiswa  int
	gotCheckJadwal int
	gotToken       string
	gotTokenID     int
	gotDeadlineID  int
}

func (f *fakeSiswaUjianChecker) CheckValidSiswaInPesertaUjianById(ctx context.Context, idSiswa int, idJadwalUjian int) (bool, int, error) {
	f.gotCheckSiswa = idSiswa
	f.gotCheckJadwal = idJadwalUjian
	if f.checkFn != nil {
		return f.checkFn(ctx, idSiswa, idJadwalUjian)
	}
	return false, 0, nil
}

func (f *fakeSiswaUjianChecker) CheckTokenUjian(ctx context.Context, token string, idJadwalUjian int) (bool, error) {
	f.gotToken = token
	f.gotTokenID = idJadwalUjian
	if f.tokenFn != nil {
		return f.tokenFn(ctx, token, idJadwalUjian)
	}
	return false, nil
}

func (f *fakeSiswaUjianChecker) GetDeadlineUjian(ctx context.Context, idJadwalUjian int) (time.Time, error) {
	f.gotDeadlineID = idJadwalUjian
	if f.deadlineFn != nil {
		return f.deadlineFn(ctx, idJadwalUjian)
	}
	return time.Time{}, nil
}

type fakeAttemptUjianRepo struct {
	createFn     func(ctx context.Context, data ujian.AttemptUjian) error
	createCalled bool
	gotData      ujian.AttemptUjian
}

func (f *fakeAttemptUjianRepo) GetAttemptUjianById(context.Context, ujian.ID) (ujian.AttemptUjian, error) {
	return ujian.AttemptUjian{}, nil
}

func (f *fakeAttemptUjianRepo) CreateAttemptUjian(ctx context.Context, data ujian.AttemptUjian) error {
	f.createCalled = true
	f.gotData = data
	if f.createFn != nil {
		return f.createFn(ctx, data)
	}
	return nil
}

func (f *fakeAttemptUjianRepo) UpdateAttemptUjian(context.Context, ujian.ID, updatepatch.UpdateAttemptUjianPatch) error {
	return nil
}

func (f *fakeAttemptUjianRepo) DeleteAttemptUjian(context.Context, ujian.ID) error {
	return nil
}

func TestAttemptUjianService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	deadline := time.Date(2026, time.March, 11, 12, 0, 0, 0, time.UTC)
	waktuAttempt := deadline.Add(-30 * time.Minute)
	repoErr := errors.New("repo error")

	tests := []struct {
		name         string
		checker      *fakeSiswaUjianChecker
		attemptRepo  *fakeAttemptUjianRepo
		wantErr      error
		wantCreate   bool
		assertResult func(t *testing.T, checker *fakeSiswaUjianChecker, attemptRepo *fakeAttemptUjianRepo)
	}{
		{
			name: "pair siswa dan jadwal tidak ada",
			checker: &fakeSiswaUjianChecker{
				checkFn: func(context.Context, int, int) (bool, int, error) {
					return false, 0, nil
				},
			},
			attemptRepo: &fakeAttemptUjianRepo{},
			wantErr:     coreerror.ErrPesertaNotAllowedToAttemptJadwal,
			wantCreate:  false,
			assertResult: func(t *testing.T, checker *fakeSiswaUjianChecker, attemptRepo *fakeAttemptUjianRepo) {
				t.Helper()
				assert.Equal(t, 9, checker.gotCheckSiswa)
				assert.Equal(t, 21, checker.gotCheckJadwal)
			},
		},
		{
			name: "repo checker error",
			checker: &fakeSiswaUjianChecker{
				checkFn: func(context.Context, int, int) (bool, int, error) {
					return false, 0, repoErr
				},
			},
			attemptRepo: &fakeAttemptUjianRepo{},
			wantErr:     repoErr,
			wantCreate:  false,
		},
		{
			name: "happy path",
			checker: &fakeSiswaUjianChecker{
				checkFn: func(context.Context, int, int) (bool, int, error) {
					return true, 17, nil
				},
				deadlineFn: func(context.Context, int) (time.Time, error) {
					return deadline, nil
				},
				tokenFn: func(context.Context, string, int) (bool, error) {
					return true, nil
				},
			},
			attemptRepo: &fakeAttemptUjianRepo{},
			wantCreate:  true,
			assertResult: func(t *testing.T, checker *fakeSiswaUjianChecker, attemptRepo *fakeAttemptUjianRepo) {
				t.Helper()
				assert.Equal(t, "TOKEN-123", checker.gotToken)
				assert.Equal(t, 21, checker.gotTokenID)
				assert.Equal(t, 21, checker.gotDeadlineID)
				assert.Equal(t, ujian.ID(17), attemptRepo.gotData.IdPesertaUjian)
				assert.Equal(t, waktuAttempt, *attemptRepo.gotData.WaktuMulai)
				assert.Equal(t, deadline, *attemptRepo.gotData.DeadlineAt)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := NewAttemptUjianService(tc.checker, tc.attemptRepo)
			err := svc.AttemptUjian(ctx, 9, 21, " token-123 ", waktuAttempt)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.wantCreate, tc.attemptRepo.createCalled)

			if tc.assertResult != nil {
				tc.assertResult(t, tc.checker, tc.attemptRepo)
			}
		})
	}
}
