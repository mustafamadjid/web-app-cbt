package user_service

import (
	"context"
	"errors"
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
	"github.com/stretchr/testify/assert"
)

type fakeProfilGuruRepo struct {
	result user.DataGuru
	err    error
	called bool
	gotID  user.ID
}

func (f *fakeProfilGuruRepo) FindProfilGuruByID(ctx context.Context, id user.ID) (user.DataGuru, error) {
	f.called = true
	f.gotID = id
	if f.err != nil {
		return user.DataGuru{}, f.err
	}
	return f.result, nil
}

func (f *fakeProfilGuruRepo) ExistByNIP(ctx context.Context, nip user.NIP) (bool, error) {
	return false, nil
}

func (f *fakeProfilGuruRepo) CreateProfilGuru(ctx context.Context, g user.ProfilGuru) (user.ID, error) {
	return 0, nil
}

func (f *fakeProfilGuruRepo) UpdateProfilGuru(ctx context.Context, idPengguna user.ID, profilGuru outuser.UpdateProfilGuruPatch) error {
	return nil
}

func TestGetGuruService_FindProfilGuruByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	wantGuru := user.DataGuru{IdPengguna: 10, Username: "guru"}
	repoErr := errors.New("repo error")

	tests := []struct {
		name       string
		id         user.ID
		repo       *fakeProfilGuruRepo
		wantErr    error
		wantGuru   user.DataGuru
		wantCalled bool
	}{
		{
			name:       "success",
			id:         10,
			repo:       &fakeProfilGuruRepo{result: wantGuru},
			wantGuru:   wantGuru,
			wantCalled: true,
		},
		{
			name:       "error",
			id:         20,
			repo:       &fakeProfilGuruRepo{err: repoErr},
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

			assert.Equal(t, tc.wantCalled, tc.repo.called)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.wantGuru, result)
		})
	}
}
