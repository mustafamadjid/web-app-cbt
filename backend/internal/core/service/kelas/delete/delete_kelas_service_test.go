package kelas_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	kelas_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/delete"
	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/delete/fake_test"
)

func TestDeleteKelasService_DeleteNamaKelas(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("delete error")

	tests := []struct {
		name        string
		idNamaKelas int
		repo        *fake_test.FakeKelasRepo
		wantErr     error
	}{
		{
			name:        "Branch 1 -> gagal delete",
			idNamaKelas: 10,
			repo:        &fake_test.FakeKelasRepo{DeleteErr: repoErr},
			wantErr:     repoErr,
		},
		{
			name:        "Branch 2 -> sukses delete",
			idNamaKelas: 11,
			repo:        &fake_test.FakeKelasRepo{},
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := kelas_service.NewDeleteKelasService(tt.repo)
			err := svc.DeleteNamaKelas(ctx, tt.idNamaKelas)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.True(t, tt.repo.DeleteCalled)
			assert.Equal(t, tt.idNamaKelas, tt.repo.GotIDNamaKelas)
		})
	}
}
