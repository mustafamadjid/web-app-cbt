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

func TestUpdateAttemptUjianService_BranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	validPesertaID := ujian.ID(5)
	start := time.Date(2026, 3, 10, 9, 0, 0, 0, time.UTC)
	submit := start.Add(20 * time.Minute)
	beforeStart := start.Add(-5 * time.Minute)
	deadline := start.Add(30 * time.Minute)
	sameAsStart := start
	zeroTime := time.Time{}

	tests := []struct {
		name          string
		idAttempt     ujian.ID
		payload       updatepatch.UpdateAttemptUjianPatch
		repo          *fake_repo.FakeAttemptUjianRepo
		wantErr       error
		wantUpdate    bool
		assertPatched func(t *testing.T, patch updatepatch.UpdateAttemptUjianPatch)
	}{
		{
			name:       "Branch 1 -> idAttempt <= 0",
			idAttempt:  0,
			payload:    updatepatch.UpdateAttemptUjianPatch{},
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantUpdate: false,
		},
		{
			name:       "Branch 2 -> tidak ada field yang diupdate",
			idAttempt:  10,
			payload:    updatepatch.UpdateAttemptUjianPatch{},
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrNoFieldToUpdate,
			wantUpdate: false,
		},
		{
			name:      "Branch 3 -> idPesertaUjian <= 0",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				IdPesertaUjian: idPtr(0),
			},
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantUpdate: false,
		},
		{
			name:      "Branch 4 -> StatusAttempt kosong setelah trim",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: statusPtr("   "),
			},
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantUpdate: false,
		},
		{
			name:      "Branch 5 -> StatusAttempt tidak valid",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: statusPtr("finished"),
			},
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantUpdate: false,
		},
		{
			name:      "Branch 6 -> WaktuMulai zero value",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				WaktuMulai: &zeroTime,
			},
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantUpdate: false,
		},
		{
			name:      "Branch 7 -> WaktuSubmit zero value",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				WaktuSubmit: &zeroTime,
			},
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantUpdate: false,
		},
		{
			name:      "Branch 8 -> DeadlineAt zero value",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				DeadlineAt: &zeroTime,
			},
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantUpdate: false,
		},
		{
			name:      "Branch 9 -> WaktuSubmit sebelum WaktuMulai",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				WaktuMulai:  &start,
				WaktuSubmit: &beforeStart,
			},
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantUpdate: false,
		},
		{
			name:      "Branch 10 -> DeadlineAt tidak lebih besar dari WaktuMulai",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				WaktuMulai: &start,
				DeadlineAt: &sameAsStart,
			},
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantUpdate: false,
		},
		{
			name:      "Branch 11 -> repo UpdateAttemptUjian error",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: statusPtr(" submitted "),
				WaktuSubmit:   &submit,
			},
			repo:       &fake_repo.FakeAttemptUjianRepo{UpdateAttemptUjianErr: repoErr},
			wantErr:    repoErr,
			wantUpdate: true,
			assertPatched: func(t *testing.T, patch updatepatch.UpdateAttemptUjianPatch) {
				t.Helper()
				if assert.NotNil(t, patch.StatusAttempt) {
					assert.Equal(t, ujian.ATTEMPT_SUBMITTED, *patch.StatusAttempt)
				}
				assert.Equal(t, &submit, patch.WaktuSubmit)
			},
		},
		{
			name:      "Branch 12 -> happy path status dinormalisasi dan patch valid",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				IdPesertaUjian: &validPesertaID,
				StatusAttempt:  statusPtr(" Expired "),
				WaktuMulai:     &start,
				DeadlineAt:     &deadline,
			},
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantUpdate: true,
			assertPatched: func(t *testing.T, patch updatepatch.UpdateAttemptUjianPatch) {
				t.Helper()
				if assert.NotNil(t, patch.IdPesertaUjian) {
					assert.Equal(t, validPesertaID, *patch.IdPesertaUjian)
				}
				if assert.NotNil(t, patch.StatusAttempt) {
					assert.Equal(t, ujian.ATTEMPT_EXPIRED, *patch.StatusAttempt)
				}
				assert.Equal(t, &start, patch.WaktuMulai)
				assert.Equal(t, &deadline, patch.DeadlineAt)
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
			if tc.wantUpdate {
				assert.Equal(t, tc.idAttempt, tc.repo.GotUpdateAttemptID)
			}
			if tc.assertPatched != nil {
				tc.assertPatched(t, tc.repo.GotUpdatePatch)
			}
		})
	}
}

func idPtr(v ujian.ID) *ujian.ID {
	return &v
}

func statusPtr(v string) *ujian.StatusAttempt {
	status := ujian.StatusAttempt(v)
	return &status
}
