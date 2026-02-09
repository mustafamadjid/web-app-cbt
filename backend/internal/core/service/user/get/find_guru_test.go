package user_service

import (
	"context"
	"errors"
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/stretchr/testify/assert"

	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/get/fake_test"
)

func TestGetGuruService_FindProfilGuruByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	wantGuru := user.DataGuru{IdPengguna: 10, Username: "guru"}
	repoErr := errors.New("repo error")

	tests := []struct {
		name       string
		id         user.ID
		repo       *fake_test.FakeProfilGuruRepo
		wantErr    error
		wantGuru   user.DataGuru
		wantCalled bool
	}{
		{
			name:       "success",
			id:         10,
			repo:       &fake_test.FakeProfilGuruRepo{Result: wantGuru},
			wantGuru:   wantGuru,
			wantCalled: true,
		},
		{
			name:       "error",
			id:         20,
			repo:       &fake_test.FakeProfilGuruRepo{Err: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var listRepo outuser.GetGuruListRepo
			service := NewGetListGuruService(listRepo, tc.repo)

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
