package attemptujian_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	attemptujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/delete"
	"github.com/stretchr/testify/assert"
)

type fakeDeleteAttemptRepo struct {
	deleteFn     func(ctx context.Context, idAttempt ujian.ID) error
	deleteCalled bool
}

func (f *fakeDeleteAttemptRepo) DeleteAttemptUjian(ctx context.Context, idAttempt ujian.ID) error {
	f.deleteCalled = true
	if f.deleteFn != nil {
		return f.deleteFn(ctx, idAttempt)
	}
	return nil
}

func TestDeleteAttemptUjianService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")

	tests := []struct {
		name       string
		idAttempt  ujian.ID
		repo       *fakeDeleteAttemptRepo
		wantErr    error
		wantDelete bool
	}{
		{
			name:       "branch 1 -> id attempt tidak valid",
			idAttempt:  0,
			repo:       &fakeDeleteAttemptRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantDelete: false,
		},
		{
			name:      "branch 2 -> repo error",
			idAttempt: 10,
			repo: &fakeDeleteAttemptRepo{
				deleteFn: func(context.Context, ujian.ID) error {
					return repoErr
				},
			},
			wantErr:    repoErr,
			wantDelete: true,
		},
		{
			name:       "happy path -> berhasil delete",
			idAttempt:  10,
			repo:       &fakeDeleteAttemptRepo{},
			wantDelete: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := attemptujian_service.NewDeleteAttemptUjianService(tc.repo)
			err := svc.DeleteAttemptUjian(ctx, tc.idAttempt)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.wantDelete, tc.repo.deleteCalled)
		})
	}
}
