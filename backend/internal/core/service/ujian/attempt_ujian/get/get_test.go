package attemptujian_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	attemptujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/get"
	"github.com/stretchr/testify/assert"
)

type fakeGetAttemptRepo struct {
	getFn     func(ctx context.Context, idAttempt ujian.ID) (ujian.AttemptUjian, error)
	getCalled bool
}

func (f *fakeGetAttemptRepo) GetAttemptUjianById(ctx context.Context, idAttempt ujian.ID) (ujian.AttemptUjian, error) {
	f.getCalled = true
	if f.getFn != nil {
		return f.getFn(ctx, idAttempt)
	}
	return ujian.AttemptUjian{}, nil
}

func TestGetAttemptUjianService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	expected := ujian.AttemptUjian{IdAttempt: 10, IdPesertaUjian: 5, StatusAttempt: ujian.ATTEMPT_IN_PROGRESS}

	tests := []struct {
		name       string
		idAttempt  ujian.ID
		repo       *fakeGetAttemptRepo
		wantErr    error
		wantGet    bool
		wantResult ujian.AttemptUjian
	}{
		{
			name:      "branch 1 -> id attempt tidak valid",
			idAttempt: 0,
			repo:      &fakeGetAttemptRepo{},
			wantErr:   coreerror.ErrMissingId,
			wantGet:   false,
		},
		{
			name:      "branch 2 -> repo error",
			idAttempt: 10,
			repo: &fakeGetAttemptRepo{
				getFn: func(context.Context, ujian.ID) (ujian.AttemptUjian, error) {
					return ujian.AttemptUjian{}, repoErr
				},
			},
			wantErr: repoErr,
			wantGet: true,
		},
		{
			name:      "happy path -> berhasil get attempt",
			idAttempt: 10,
			repo: &fakeGetAttemptRepo{
				getFn: func(_ context.Context, id ujian.ID) (ujian.AttemptUjian, error) {
					assert.Equal(t, ujian.ID(10), id)
					return expected, nil
				},
			},
			wantGet:    true,
			wantResult: expected,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := attemptujian_service.NewGetAttemptUjianService(tc.repo)
			result, err := svc.GetAttemptUjianById(ctx, tc.idAttempt)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.wantGet, tc.repo.getCalled)
			assert.Equal(t, tc.wantResult, result)
		})
	}
}
