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
)

type fakeUpdateAttemptRepo struct {
	updateFn     func(ctx context.Context, idAttempt ujian.ID, data updatepatch.UpdateAttemptUjianPatch) error
	updateCalled bool
	gotID        ujian.ID
	gotPatch     updatepatch.UpdateAttemptUjianPatch
}

func (f *fakeUpdateAttemptRepo) UpdateAttemptUjian(ctx context.Context, idAttempt ujian.ID, data updatepatch.UpdateAttemptUjianPatch) error {
	f.updateCalled = true
	f.gotID = idAttempt
	f.gotPatch = data

	if f.updateFn != nil {
		return f.updateFn(ctx, idAttempt, data)
	}

	return nil
}

func TestUpdateAttemptUjianService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	start := time.Date(2026, 3, 10, 9, 0, 0, 0, time.UTC)
	submit := start.Add(20 * time.Minute)
	beforeStart := start.Add(-5 * time.Minute)

	tests := []struct {
		name          string
		idAttempt     ujian.ID
		payload       updatepatch.UpdateAttemptUjianPatch
		repo          *fakeUpdateAttemptRepo
		wantErr       error
		wantUpdate    bool
		assertPatched func(t *testing.T, patch updatepatch.UpdateAttemptUjianPatch)
	}{
		{
			name:       "branch 1 -> id attempt tidak valid",
			idAttempt:  0,
			payload:    updatepatch.UpdateAttemptUjianPatch{},
			repo:       &fakeUpdateAttemptRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantUpdate: false,
		},
		{
			name:       "branch 2 -> tidak ada field yang diupdate",
			idAttempt:  10,
			payload:    updatepatch.UpdateAttemptUjianPatch{},
			repo:       &fakeUpdateAttemptRepo{},
			wantErr:    coreerror.ErrNoFieldToUpdate,
			wantUpdate: false,
		},
		{
			name:      "branch 3 -> id peserta ujian tidak valid",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				IdPesertaUjian: idPtr(0),
			},
			repo:       &fakeUpdateAttemptRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantUpdate: false,
		},
		{
			name:      "branch 4 -> status kosong setelah trim",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: statusPtr("   "),
			},
			repo:       &fakeUpdateAttemptRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantUpdate: false,
		},
		{
			name:      "branch 5 -> status tidak valid",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: statusPtr("finished"),
			},
			repo:       &fakeUpdateAttemptRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantUpdate: false,
		},
		{
			name:      "branch 6 -> waktu submit sebelum waktu mulai",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				WaktuMulai:  &start,
				WaktuSubmit: &beforeStart,
			},
			repo:       &fakeUpdateAttemptRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantUpdate: false,
		},
		{
			name:      "branch 7 -> repo error",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: statusPtr(" submitted "),
				WaktuSubmit:   &submit,
			},
			repo: &fakeUpdateAttemptRepo{
				updateFn: func(context.Context, ujian.ID, updatepatch.UpdateAttemptUjianPatch) error {
					return repoErr
				},
			},
			wantErr:    repoErr,
			wantUpdate: true,
		},
		{
			name:      "happy path -> status dinormalisasi",
			idAttempt: 10,
			payload: updatepatch.UpdateAttemptUjianPatch{
				StatusAttempt: statusPtr(" Expired "),
			},
			repo:       &fakeUpdateAttemptRepo{},
			wantUpdate: true,
			assertPatched: func(t *testing.T, patch updatepatch.UpdateAttemptUjianPatch) {
				t.Helper()
				if assert.NotNil(t, patch.StatusAttempt) {
					assert.Equal(t, ujian.ATTEMPT_EXPIRED, *patch.StatusAttempt)
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := attemptujian_service.NewUpdateAttemptUjianService(tc.repo)
			err := svc.UpdateAttemptUjian(ctx, tc.idAttempt, tc.payload)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.wantUpdate, tc.repo.updateCalled)
			if tc.wantUpdate {
				assert.Equal(t, tc.idAttempt, tc.repo.gotID)
			}
			if tc.assertPatched != nil {
				tc.assertPatched(t, tc.repo.gotPatch)
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
