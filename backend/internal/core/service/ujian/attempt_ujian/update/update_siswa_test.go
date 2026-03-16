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
	fake_repo "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/update/fake_repo"
	"github.com/stretchr/testify/assert"
)

type fakeOwnershipChecker struct {
	checkFn      func(ctx context.Context, idSiswa int, idAttempt ujian.ID) (bool, error)
	checkCalled  bool
	gotSiswaID   int
	gotAttemptID ujian.ID
}

func (f *fakeOwnershipChecker) CheckAttemptOwnershipBySiswa(ctx context.Context, idSiswa int, idAttempt ujian.ID) (bool, error) {
	f.checkCalled = true
	f.gotSiswaID = idSiswa
	f.gotAttemptID = idAttempt
	if f.checkFn != nil {
		return f.checkFn(ctx, idSiswa, idAttempt)
	}
	return false, nil
}

func TestSiswaUpdateAttemptUjianService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	submitTime := time.Date(2026, 3, 16, 11, 45, 0, 0, time.UTC)
	repoErr := errors.New("repo error")

	tests := []struct {
		name        string
		idSiswa     int
		idAttempt   ujian.ID
		payload     updatepatch.UpdateAttemptUjianPatch
		checker     *fakeOwnershipChecker
		repo        *fake_repo.FakeAttemptUjianRepo
		wantErr     error
		wantCheck   bool
		wantUpdate  bool
		assertPatch func(t *testing.T, patch updatepatch.UpdateAttemptUjianPatch)
	}{
		{
			name:       "invalid siswa id",
			idAttempt:  7,
			payload:    updatepatch.UpdateAttemptUjianPatch{StatusAttempt: statusPtr("submitted"), WaktuSubmit: &submitTime},
			checker:    &fakeOwnershipChecker{},
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCheck:  false,
			wantUpdate: false,
		},
		{
			name:       "missing status while submit time provided",
			idSiswa:    9,
			idAttempt:  7,
			payload:    updatepatch.UpdateAttemptUjianPatch{WaktuSubmit: &submitTime},
			checker:    &fakeOwnershipChecker{},
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantCheck:  false,
			wantUpdate: false,
		},
		{
			name:      "submitted requires submit time",
			idSiswa:   9,
			idAttempt: 7,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: statusPtr("submitted"),
			},
			checker:    &fakeOwnershipChecker{},
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantCheck:  false,
			wantUpdate: false,
		},
		{
			name:      "expired rejects submit time",
			idSiswa:   9,
			idAttempt: 7,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: statusPtr("expired"),
				WaktuSubmit:   &submitTime,
			},
			checker:    &fakeOwnershipChecker{},
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCheck:  false,
			wantUpdate: false,
		},
		{
			name:      "ownership check error",
			idSiswa:   9,
			idAttempt: 7,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: statusPtr("expired"),
			},
			checker: &fakeOwnershipChecker{
				checkFn: func(context.Context, int, ujian.ID) (bool, error) {
					return false, repoErr
				},
			},
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantErr:    repoErr,
			wantCheck:  true,
			wantUpdate: false,
		},
		{
			name:      "attempt not owned by siswa",
			idSiswa:   9,
			idAttempt: 7,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: statusPtr("expired"),
			},
			checker: &fakeOwnershipChecker{
				checkFn: func(context.Context, int, ujian.ID) (bool, error) {
					return false, nil
				},
			},
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrNotFound,
			wantCheck:  true,
			wantUpdate: false,
		},
		{
			name:      "generic updater error bubbles up",
			idSiswa:   9,
			idAttempt: 7,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: statusPtr("submitted"),
				WaktuSubmit:   &submitTime,
			},
			checker: &fakeOwnershipChecker{
				checkFn: func(context.Context, int, ujian.ID) (bool, error) {
					return true, nil
				},
			},
			repo:       &fake_repo.FakeAttemptUjianRepo{UpdateAttemptUjianErr: repoErr},
			wantErr:    repoErr,
			wantCheck:  true,
			wantUpdate: true,
			assertPatch: func(t *testing.T, patch updatepatch.UpdateAttemptUjianPatch) {
				t.Helper()
				assert.Equal(t, &submitTime, patch.WaktuSubmit)
				if assert.NotNil(t, patch.StatusAttempt) {
					assert.Equal(t, ujian.ATTEMPT_SUBMITTED, *patch.StatusAttempt)
				}
			},
		},
		{
			name:      "success expired",
			idSiswa:   9,
			idAttempt: 7,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: statusPtr(" expired "),
			},
			checker: &fakeOwnershipChecker{
				checkFn: func(context.Context, int, ujian.ID) (bool, error) {
					return true, nil
				},
			},
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantCheck:  true,
			wantUpdate: true,
			assertPatch: func(t *testing.T, patch updatepatch.UpdateAttemptUjianPatch) {
				t.Helper()
				if assert.NotNil(t, patch.StatusAttempt) {
					assert.Equal(t, ujian.ATTEMPT_EXPIRED, *patch.StatusAttempt)
				}
				assert.Nil(t, patch.WaktuSubmit)
			},
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
			if tc.wantCheck {
				assert.Equal(t, tc.idSiswa, tc.checker.gotSiswaID)
				assert.Equal(t, tc.idAttempt, tc.checker.gotAttemptID)
			}
			if tc.assertPatch != nil {
				tc.assertPatch(t, tc.repo.GotUpdatePatch)
			}
		})
	}
}
