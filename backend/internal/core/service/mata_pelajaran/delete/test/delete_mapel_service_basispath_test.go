package matapelajajaran_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	mapel_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/delete"
	"github.com/stretchr/testify/assert"
)

func TestDeleteMapelService_BranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	genericErr := errors.New("generic delete error")

	tests := []struct {
		name       string
		idMapel    int
		repo       *FakeMapelRepo
		wantErr    error
		wantDelete bool
	}{
		{
			name:       "Branch 1 -> id mapel tidak valid",
			idMapel:    0,
			repo:       &FakeMapelRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantDelete: false,
		},
		{
			name:       "Branch 2 -> delete ditolak constraint",
			idMapel:    1,
			repo:       &FakeMapelRepo{DeleteMapelErr: coreerror.ErrDeleteRestricted},
			wantErr:    coreerror.ErrDeleteRestricted,
			wantDelete: true,
		},
		{
			name:       "Branch 3 -> repo delete gagal",
			idMapel:    1,
			repo:       &FakeMapelRepo{DeleteMapelErr: genericErr},
			wantErr:    genericErr,
			wantDelete: true,
		},
		{
			name:       "Branch 4 -> delete mapel berhasil",
			idMapel:    1,
			repo:       &FakeMapelRepo{},
			wantErr:    nil,
			wantDelete: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mapel_service.NewDeleteMapelService(tt.repo)
			err := svc.DeleteMapelService(ctx, tt.idMapel)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantDelete, tt.repo.DeleteMapelCalled)
			if tt.repo.DeleteMapelCalled {
				assert.Equal(t, tt.idMapel, tt.repo.GotDeleteId)
			}
		})
	}
}
