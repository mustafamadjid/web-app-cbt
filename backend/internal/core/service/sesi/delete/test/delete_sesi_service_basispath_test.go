package sesi_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	sesi_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/delete"
	"github.com/stretchr/testify/assert"
)

func TestDeleteSesiService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("delete sesi error")

	tests := []struct {
		name       string
		idSesi     int
		repo       *FakeSesiRepo
		wantErr    error
		wantDelete bool
	}{
		{
			name:       "Path 1 -> idSesi == 0",
			idSesi:     0,
			repo:       &FakeSesiRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantDelete: false,
		},
		{
			name:       "Path 2 -> DeleteSesi error",
			idSesi:     1,
			repo:       &FakeSesiRepo{DeleteSesiErr: repoErr},
			wantErr:    repoErr,
			wantDelete: true,
		},
		{
			name:       "Path 3 -> berhasil delete sesi",
			idSesi:     1,
			repo:       &FakeSesiRepo{},
			wantErr:    nil,
			wantDelete: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := sesi_service.NewDeleteSesiService(tt.repo)
			err := svc.DeleteSesiService(ctx, tt.idSesi)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantDelete, tt.repo.DeleteSesiCalled)
			if tt.repo.DeleteSesiCalled {
				assert.Equal(t, tt.idSesi, tt.repo.GotDeleteId)
			}
		})
	}
}
