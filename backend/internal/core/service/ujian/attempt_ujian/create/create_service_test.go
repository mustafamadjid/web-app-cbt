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

type fakeAttemptRepo struct {
	createFn         func(ctx context.Context, data ujian.AttemptUjian) error
	createCalled     bool
	gotCreateAttempt ujian.AttemptUjian
}

func (f *fakeAttemptRepo) CreateAttemptUjian(ctx context.Context, data ujian.AttemptUjian) error {
	f.createCalled = true
	f.gotCreateAttempt = data

	if f.createFn != nil {
		return f.createFn(ctx, data)
	}

	return nil
}

func (f *fakeAttemptRepo) GetAttemptUjianById(context.Context, ujian.ID) (ujian.AttemptUjian, error) {
	return ujian.AttemptUjian{}, nil
}

func (f *fakeAttemptRepo) UpdateAttemptUjian(context.Context, ujian.ID, updatepatch.UpdateAttemptUjianPatch) error {
	return nil
}

func (f *fakeAttemptRepo) DeleteAttemptUjian(context.Context, ujian.ID) error {
	return nil
}

func TestCreateAttemptUjianService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	start := time.Date(2026, 3, 10, 9, 0, 0, 0, time.UTC)
	submit := start.Add(30 * time.Minute)
	beforeStart := start.Add(-10 * time.Minute)
	deadline := start.Add(60 * time.Minute)

	tests := []struct {
		name          string
		repo          *fakeAttemptRepo
		payload       ujian.AttemptUjian
		expectErr     error
		expectCreate  bool
		assertCreated func(t *testing.T, got ujian.AttemptUjian)
	}{
		{
			name:         "branch 1 -> id peserta ujian kosong",
			repo:         &fakeAttemptRepo{},
			payload:      ujian.AttemptUjian{},
			expectErr:    coreerror.ErrMissingId,
			expectCreate: false,
		},
		{
			name: "branch 2 -> status attempt tidak valid",
			repo: &fakeAttemptRepo{},
			payload: ujian.AttemptUjian{
				IdPesertaUjian: 10,
				StatusAttempt:  "finished",
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectCreate: false,
		},
		{
			name: "branch 3 -> submitted tanpa waktu submit",
			repo: &fakeAttemptRepo{},
			payload: ujian.AttemptUjian{
				IdPesertaUjian: 10,
				StatusAttempt:  " submitted ",
				WaktuMulai:     &start,
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectCreate: false,
		},
		{
			name: "branch 4 -> waktu submit sebelum waktu mulai",
			repo: &fakeAttemptRepo{},
			payload: ujian.AttemptUjian{
				IdPesertaUjian: 10,
				StatusAttempt:  "submitted",
				WaktuMulai:     &start,
				WaktuSubmit:    &beforeStart,
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectCreate: false,
		},
		{
			name: "branch 5 -> deadline tidak valid",
			repo: &fakeAttemptRepo{},
			payload: ujian.AttemptUjian{
				IdPesertaUjian: 10,
				StatusAttempt:  "in_progress",
				WaktuMulai:     &start,
				DeadlineAt:     &start,
			},
			expectErr:    coreerror.ErrInvalidInput,
			expectCreate: false,
		},
		{
			name: "branch 6 -> gagal create repo",
			repo: &fakeAttemptRepo{
				createFn: func(context.Context, ujian.AttemptUjian) error {
					return repoErr
				},
			},
			payload: ujian.AttemptUjian{
				IdPesertaUjian: 10,
				StatusAttempt:  "  submitted  ",
				WaktuMulai:     &start,
				WaktuSubmit:    &submit,
				DeadlineAt:     &deadline,
			},
			expectErr:    repoErr,
			expectCreate: true,
		},
		{
			name: "happy path -> status dinormalisasi dan default terisi",
			repo: &fakeAttemptRepo{},
			payload: ujian.AttemptUjian{
				IdPesertaUjian: 10,
				StatusAttempt:  "   ",
				WaktuMulai:     &start,
				DeadlineAt:     &deadline,
			},
			expectErr:    nil,
			expectCreate: true,
			assertCreated: func(t *testing.T, got ujian.AttemptUjian) {
				t.Helper()
				assert.Equal(t, ujian.ID(10), got.IdPesertaUjian)
				assert.Equal(t, ujian.ATTEMPT_IN_PROGRESS, got.StatusAttempt)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := attemptujian_service.NewCreateAttemptUjianService(tc.repo)
			err := svc.CreateAttemptUjian(ctx, tc.payload)

			assert.ErrorIs(t, err, tc.expectErr)
			assert.Equal(t, tc.expectCreate, tc.repo.createCalled)

			if tc.assertCreated != nil {
				tc.assertCreated(t, tc.repo.gotCreateAttempt)
			}
		})
	}
}
