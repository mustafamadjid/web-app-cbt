package user_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/get"
	"github.com/stretchr/testify/assert"
)

func TestGetGuruService_FindProfilGuruByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	wantGuru := user.DataGuru{IdPengguna: 10, Username: "guru"}
	repoErr := errors.New("repo error")

	tests := []struct {
		name       string
		id         user.ID
		repo       *FakeProfilGuruRepo
		wantErr    error
		wantGuru   user.DataGuru
		wantCalled bool
	}{
		{
			name:       "Path 1 -> profil guru ditemukan",
			id:         10,
			repo:       &FakeProfilGuruRepo{Result: wantGuru},
			wantGuru:   wantGuru,
			wantCalled: true,
		},
		{
			name:       "Path 2 -> repo find profil guru gagal",
			id:         20,
			repo:       &FakeProfilGuruRepo{Err: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var listRepo outuser.GetGuruListRepo
			service := user_service.NewGetListGuruService(listRepo, tc.repo)

			result, err := service.FindProfilGuruByID(ctx, tc.id)

			assert.Equal(t, tc.wantCalled, tc.repo.Called)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.wantGuru, result)
		})
	}
}
