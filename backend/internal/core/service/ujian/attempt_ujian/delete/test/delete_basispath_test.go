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

func TestDeleteAttemptUjianService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")

	tests := []struct {
		name       string
		idAttempt  ujian.ID
		repo       *FakeAttemptUjianRepo
		wantErr    error
		wantDelete bool
	}{
		{
			name:      "path 1 -> id attempt tidak valid",
			idAttempt: 0,
			repo:      &FakeAttemptUjianRepo{},
			wantErr:   coreerror.ErrMissingId,
		},
		{
			name:       "path 2 -> repo delete attempt gagal",
			idAttempt:  12,
			repo:       &FakeAttemptUjianRepo{DeleteAttemptUjianErr: repoErr},
			wantErr:    repoErr,
			wantDelete: true,
		},
		{
			name:       "path 3 -> berhasil delete attempt",
			idAttempt:  12,
			repo:       &FakeAttemptUjianRepo{},
			wantDelete: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := attemptujian_service.NewDeleteAttemptUjianService(tc.repo)
			err := svc.DeleteAttemptUjian(ctx, tc.idAttempt)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantDelete, tc.repo.DeleteAttemptUjianCalled)
		})
	}
}
