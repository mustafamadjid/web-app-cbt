package matapelajajaran_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	mapel_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/delete"
	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/delete/fake_test"
	"github.com/stretchr/testify/assert"
)

func TestDeleteMapelService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	genericErr := errors.New("generic delete error")

	tests := []struct {
		name       string
		idMapel    int
		repo       *fake_test.FakeMapelRepo
		wantErr    error
		wantDelete bool
	}{
		{
			name:       "Path 1 -> idMapel <= 0",
			idMapel:    0,
			repo:       &fake_test.FakeMapelRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantDelete: false,
		},
		{
			name:       "Path 2 -> DeleteMapel error ErrDeleteRestricted",
			idMapel:    1,
			repo:       &fake_test.FakeMapelRepo{DeleteMapelErr: coreerror.ErrDeleteRestricted},
			wantErr:    coreerror.ErrDeleteRestricted,
			wantDelete: true,
		},
		{
			name:       "Path 3 -> DeleteMapel generic error",
			idMapel:    1,
			repo:       &fake_test.FakeMapelRepo{DeleteMapelErr: genericErr},
			wantErr:    genericErr,
			wantDelete: true,
		},
		{
			name:       "Path 4 -> berhasil delete mapel",
			idMapel:    1,
			repo:       &fake_test.FakeMapelRepo{},
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
