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

func TestGetAttemptUjianService_BranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	expected := ujian.AttemptUjian{IdAttempt: 10, IdPesertaUjian: 5, StatusAttempt: ujian.ATTEMPT_IN_PROGRESS}

	tests := []struct {
		name       string
		idAttempt  ujian.ID
		repo       *FakeAttemptUjianRepo
		wantErr    error
		wantGet    bool
		wantResult ujian.AttemptUjian
	}{
		{
			name:      "Branch 1 -> idAttempt <= 0",
			idAttempt: 0,
			repo:      &FakeAttemptUjianRepo{},
			wantErr:   coreerror.ErrMissingId,
			wantGet:   false,
		},
		{
			name:      "Branch 2 -> repo GetAttemptUjianById error",
			idAttempt: 10,
			repo:      &FakeAttemptUjianRepo{GetAttemptUjianByIdErr: repoErr},
			wantErr:   repoErr,
			wantGet:   true,
		},
		{
			name:       "Branch 3 -> berhasil get attempt ujian",
			idAttempt:  10,
			repo:       &FakeAttemptUjianRepo{GetAttemptUjianByIdRet: expected},
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

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantGet, tc.repo.GetAttemptUjianByIdCalled)
			if tc.repo.GetAttemptUjianByIdCalled {
				assert.Equal(t, tc.idAttempt, tc.repo.GotGetAttemptID)
			}
			assert.Equal(t, tc.wantResult, result)
		})
	}
}
