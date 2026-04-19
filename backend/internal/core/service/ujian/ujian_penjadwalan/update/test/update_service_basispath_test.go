package ujian_service_test

import (
	"context"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/update"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUjianService_BasisPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		repo       *fakeUpdateUjianRepo
		wantErr    error
		wantUpdate bool
	}{
		{
			name:       "path 1 -> repo conflict diteruskan",
			repo:       &fakeUpdateUjianRepo{updateErr: coreerror.ErrConflict},
			wantErr:    coreerror.ErrConflict,
			wantUpdate: true,
		},
		{
			name:       "path 2 -> update ujian berhasil",
			repo:       &fakeUpdateUjianRepo{},
			wantUpdate: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := ujian_service.NewUpdateUjianService(tc.repo)
			err := svc.UpdateUjianService(context.Background(), 10, validUpdatePayload())

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantUpdate, tc.repo.updateCalled)
		})
	}
}
