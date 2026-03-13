package attemptujian_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	attemptujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/delete"
	fake_repo "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/delete/fake_repo"
	"github.com/stretchr/testify/assert"
)

func TestDeleteAttemptUjianService_BranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")

	tests := []struct {
		name       string
		idAttempt  ujian.ID
		repo       *fake_repo.FakeAttemptUjianRepo
		wantErr    error
		wantDelete bool
	}{
		{
			name:       "Branch 1 -> idAttempt <= 0",
			idAttempt:  0,
			repo:       &fake_repo.FakeAttemptUjianRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantDelete: false,
		},
		{
			name:       "Branch 2 -> repo DeleteAttemptUjian error",
			idAttempt:  10,
			repo:       &fake_repo.FakeAttemptUjianRepo{DeleteAttemptUjianErr: repoErr},
			wantErr:    repoErr,
			wantDelete: true,
		},
		{
			name:       "Branch 3 -> berhasil delete attempt ujian",
			idAttempt:  10,
			repo:       &fake_repo.FakeAttemptUjianRepo{},
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
			if tc.repo.DeleteAttemptUjianCalled {
				assert.Equal(t, tc.idAttempt, tc.repo.GotDeleteAttemptID)
			}
		})
	}
}
