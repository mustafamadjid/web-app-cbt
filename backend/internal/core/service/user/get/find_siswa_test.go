package user_service

import (
	"context"
	"errors"
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"

	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/get/fake_test"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

func TestGetSiswaService_FindProfilSiswaByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	wantSiswa := user.DataSiswa{IdPengguna: 11, Username: "siswa"}
	repoErr := errors.New("repo error")

	tests := []struct {
		name       string
		id         user.ID
		repo       *fake_test.FakeProfilSiswaRepo
		wantErr    error
		wantSiswa  user.DataSiswa
		wantCalled bool
	}{
		{
			name:       "success",
			id:         11,
			repo:       &fake_test.FakeProfilSiswaRepo{Result: wantSiswa},
			wantSiswa:  wantSiswa,
			wantCalled: true,
		},
		{
			name:       "error",
			id:         12,
			repo:       &fake_test.FakeProfilSiswaRepo{Err: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var listRepo outuser.GetListSiswaRepo
			service := NewGetListSiswaService(listRepo, tc.repo)

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
