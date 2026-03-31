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

func TestGetSiswaService_FindProfilSiswaByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	wantSiswa := user.DataSiswa{IdPengguna: 11, Username: "siswa"}
	repoErr := errors.New("repo error")

	tests := []struct {
		name       string
		id         user.ID
		repo       *FakeProfilSiswaRepo
		wantErr    error
		wantSiswa  user.DataSiswa
		wantCalled bool
	}{
		{
			name:       "success",
			id:         11,
			repo:       &FakeProfilSiswaRepo{Result: wantSiswa},
			wantSiswa:  wantSiswa,
			wantCalled: true,
		},
		{
			name:       "error",
			id:         12,
			repo:       &FakeProfilSiswaRepo{Err: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var listRepo outuser.GetListSiswaRepo
			service := user_service.NewGetListSiswaService(listRepo, tc.repo)

			result, err := service.FindProfilSiswaByID(ctx, tc.id)

			assert.Equal(t, tc.wantCalled, tc.repo.Called)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.wantSiswa, result)
		})
	}
}
