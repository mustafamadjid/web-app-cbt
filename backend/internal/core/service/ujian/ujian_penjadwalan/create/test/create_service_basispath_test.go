package ujian_service_test

import (
	"context"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/create"
	"github.com/stretchr/testify/assert"
)

func TestCreateUjianService_BasisPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		repo       *fakeCreateUjianRepo
		wantErr    error
		wantCreate bool
	}{
		{
			name:       "path 1 -> repo conflict diteruskan",
			repo:       &fakeCreateUjianRepo{createErr: coreerror.ErrConflict},
			wantErr:    coreerror.ErrConflict,
			wantCreate: true,
		},
		{
			name:       "path 2 -> create ujian berhasil",
			repo:       &fakeCreateUjianRepo{},
			wantCreate: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := ujian_service.NewCreateUjianService(tc.repo)
			err := svc.CreateUjianService(context.Background(), baseCreateUjianData())

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantCreate, tc.repo.createCalled)
		})
	}
}
