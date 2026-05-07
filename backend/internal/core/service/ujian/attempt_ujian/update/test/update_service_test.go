package attemptujian_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	attemptujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeUpdateAttemptChecker struct {
	ownedRet     bool
	ownedErr     error
	checkCalled  bool
	gotSiswaID   int
	gotAttemptID ujian.ID
}

func (f *fakeUpdateAttemptChecker) CheckAttemptOwnershipBySiswa(_ context.Context, idSiswa int, idAttempt ujian.ID) (bool, error) {
	f.checkCalled = true
	f.gotSiswaID = idSiswa
	f.gotAttemptID = idAttempt
	if f.ownedErr != nil {
		return false, f.ownedErr
	}
	return f.ownedRet, nil
}

func (*fakeUpdateAttemptChecker) CheckValidSiswaInPesertaUjianById(context.Context, int, int) (bool, int, error) {
	return false, 0, nil
}

func (*fakeUpdateAttemptChecker) CheckTokenUjian(context.Context, string, int) (bool, error) {
	return false, nil
}

func (*fakeUpdateAttemptChecker) GetDeadlineUjian(context.Context, int) (time.Time, error) {
	return time.Time{}, nil
}

func statusPtr(v ujian.StatusAttempt) *ujian.StatusAttempt {
	return &v
}

func timePatchPtr(v time.Time) *time.Time {
	return &v
}

func idPatchPtr(v ujian.ID) *ujian.ID {
	return &v
}

func TestUpdateAttemptUjianService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	mulai := time.Date(2026, 4, 12, 8, 0, 0, 0, time.UTC)
	submit := time.Date(2026, 4, 12, 9, 0, 0, 0, time.UTC)
	deadline := time.Date(2026, 4, 12, 10, 0, 0, 0, time.UTC)
	status := ujian.StatusAttempt(" submitted ")

	tests := []struct {
		name       string
		idAttempt  ujian.ID
		payload    updatepatch.UpdateAttemptUjianPatch
		repo       *FakeAttemptUjianRepo
		wantErr    error
		wantUpdate bool
		assertion  func(*testing.T, *FakeAttemptUjianRepo)
	}{
		{
			name:      "Path 1 -> status attempt kosong setelah trim",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: statusPtr(ujian.StatusAttempt("   ")),
			},
			repo:    &FakeAttemptUjianRepo{},
			wantErr: coreerror.ErrMissingField,
		},
		{
			name:      "Path 2 -> id attempt tidak valid",
			idAttempt: 0,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: statusPtr(status),
			},
			repo:    &FakeAttemptUjianRepo{},
			wantErr: coreerror.ErrMissingId,
		},
		{
			name:      "Path 3 -> tidak ada field untuk diupdate",
			idAttempt: 10,
			payload:   updatepatch.UpdateAttemptUjianPatch{},
			repo:      &FakeAttemptUjianRepo{},
			wantErr:   coreerror.ErrNoFieldToUpdate,
		},
		{
			name:      "Path 4 -> status attempt tidak valid",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: statusPtr(ujian.StatusAttempt("unknown")),
			},
			repo:    &FakeAttemptUjianRepo{},
			wantErr: coreerror.ErrInvalidInput,
		},
		{
			name:      "Path 5 -> id peserta ujian tidak valid",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				IdPesertaUjian: idPatchPtr(0),
			},
			repo:    &FakeAttemptUjianRepo{},
			wantErr: coreerror.ErrMissingId,
		},
		{
			name:      "Path 6 -> waktu submit sebelum waktu mulai",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				WaktuMulai:  timePatchPtr(submit),
				WaktuSubmit: timePatchPtr(mulai),
			},
			repo:    &FakeAttemptUjianRepo{},
			wantErr: coreerror.ErrInvalidInput,
		},
		{
			name:      "Path 7 -> deadline tidak setelah waktu mulai",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				WaktuMulai: timePatchPtr(deadline),
				DeadlineAt: timePatchPtr(mulai),
			},
			repo:    &FakeAttemptUjianRepo{},
			wantErr: coreerror.ErrInvalidInput,
		},
		{
			name:      "Path 8 -> repo update attempt gagal",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: statusPtr(status),
				WaktuMulai:    timePatchPtr(mulai),
				WaktuSubmit:   timePatchPtr(submit),
				DeadlineAt:    timePatchPtr(deadline),
			},
			repo:       &FakeAttemptUjianRepo{UpdateAttemptUjianErr: repoErr},
			wantErr:    repoErr,
			wantUpdate: true,
		},
		{
			name:      "Path 9 -> berhasil update attempt",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: statusPtr(status),
				WaktuMulai:    timePatchPtr(mulai),
				WaktuSubmit:   timePatchPtr(submit),
				DeadlineAt:    timePatchPtr(deadline),
			},
			repo:       &FakeAttemptUjianRepo{},
			wantUpdate: true,
			assertion: func(t *testing.T, repo *FakeAttemptUjianRepo) {
				t.Helper()
				require.NotNil(t, repo.GotUpdatePatch.StatusAttempt)
				assert.Equal(t, ujian.ATTEMPT_SUBMITTED, *repo.GotUpdatePatch.StatusAttempt)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := attemptujian_service.NewUpdateAttemptUjianService(tc.repo)
			err := svc.UpdateAttemptUjian(ctx, tc.idAttempt, tc.payload)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantUpdate, tc.repo.UpdateAttemptUjianCalled)
			if tc.assertion != nil {
				tc.assertion(t, tc.repo)
			}
		})
	}
}

func TestSiswaUpdateAttemptUjianService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	checkErr := errors.New("checker error")
	repoErr := errors.New("repo error")
	submit := time.Date(2026, 4, 12, 9, 0, 0, 0, time.UTC)
	statusSubmitted := ujian.ATTEMPT_SUBMITTED
	statusExpired := ujian.ATTEMPT_EXPIRED

	tests := []struct {
		name       string
		idSiswa    int
		idAttempt  ujian.ID
		payload    updatepatch.UpdateAttemptUjianPatch
		checker    *fakeUpdateAttemptChecker
		repo       *FakeAttemptUjianRepo
		wantErr    error
		wantCheck  bool
		wantUpdate bool
	}{
		{
			name:      "Path 1 -> id siswa tidak valid",
			idSiswa:   0,
			idAttempt: 11,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: &statusSubmitted,
				WaktuSubmit:   timePatchPtr(submit),
			},
			checker: &fakeUpdateAttemptChecker{},
			repo:    &FakeAttemptUjianRepo{},
			wantErr: coreerror.ErrMissingId,
		},
		{
			name:      "Path 2 -> payload siswa tidak valid",
			idSiswa:   9,
			idAttempt: 11,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: &statusSubmitted,
			},
			checker: &fakeUpdateAttemptChecker{},
			repo:    &FakeAttemptUjianRepo{},
			wantErr: coreerror.ErrMissingField,
		},
		{
			name:      "Path 3 -> checker ownership gagal",
			idSiswa:   9,
			idAttempt: 11,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: &statusSubmitted,
				WaktuSubmit:   timePatchPtr(submit),
			},
			checker:   &fakeUpdateAttemptChecker{ownedErr: checkErr},
			repo:      &FakeAttemptUjianRepo{},
			wantErr:   checkErr,
			wantCheck: true,
		},
		{
			name:      "Path 4 -> attempt bukan milik siswa",
			idSiswa:   9,
			idAttempt: 11,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: &statusSubmitted,
				WaktuSubmit:   timePatchPtr(submit),
			},
			checker:   &fakeUpdateAttemptChecker{ownedRet: false},
			repo:      &FakeAttemptUjianRepo{},
			wantErr:   coreerror.ErrNotFound,
			wantCheck: true,
		},
		{
			name:      "Path 5 -> update attempt siswa gagal",
			idSiswa:   9,
			idAttempt: 11,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: &statusExpired,
			},
			checker:    &fakeUpdateAttemptChecker{ownedRet: true},
			repo:       &FakeAttemptUjianRepo{UpdateAttemptUjianErr: repoErr},
			wantErr:    repoErr,
			wantCheck:  true,
			wantUpdate: true,
		},
		{
			name:      "Path 6 -> update attempt siswa berhasil submitted",
			idSiswa:   9,
			idAttempt: 11,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: &statusSubmitted,
				WaktuSubmit:   timePatchPtr(submit),
			},
			checker:    &fakeUpdateAttemptChecker{ownedRet: true},
			repo:       &FakeAttemptUjianRepo{},
			wantCheck:  true,
			wantUpdate: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			updater := attemptujian_service.NewUpdateAttemptUjianService(tc.repo)
			svc := attemptujian_service.NewSiswaUpdateAttemptUjianService(tc.checker, updater)
			err := svc.UpdateAttemptUjian(ctx, tc.idSiswa, tc.idAttempt, tc.payload)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantCheck, tc.checker.checkCalled)
			assert.Equal(t, tc.wantUpdate, tc.repo.UpdateAttemptUjianCalled)
		})
	}
}
