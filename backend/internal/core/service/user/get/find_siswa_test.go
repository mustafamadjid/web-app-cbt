package user_service

import (
	"context"
	"errors"
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
	"github.com/stretchr/testify/assert"
)

type fakeProfilSiswaRepo struct {
	result user.DataSiswa
	err    error
	called bool
	gotID  user.ID
}

func (f *fakeProfilSiswaRepo) FindProfilSiswaByID(ctx context.Context, id user.ID) (user.DataSiswa, error) {
	f.called = true
	f.gotID = id
	if f.err != nil {
		return user.DataSiswa{}, f.err
	}
	return f.result, nil
}

func (f *fakeProfilSiswaRepo) ExistByNISN(ctx context.Context, nisn string) (bool, error) {
	return false, nil
}

func (f *fakeProfilSiswaRepo) CreateProfilSiswa(ctx context.Context, g user.ProfilSiswa) (user.ID, error) {
	return 0, nil
}

func (f *fakeProfilSiswaRepo) UpdateProfilSiswa(ctx context.Context, idPengguna user.ID, profilSiswa outuser.UpdateProfilSiswaPatch) error {
	return nil
}

func TestGetSiswaService_FindProfilSiswaByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	wantSiswa := user.DataSiswa{IdPengguna: 11, Username: "siswa"}
	repoErr := errors.New("repo error")

	tests := []struct {
		name       string
		id         user.ID
		repo       *fakeProfilSiswaRepo
		wantErr    error
		wantSiswa  user.DataSiswa
		wantCalled bool
	}{
		{
			name:       "success",
			id:         11,
			repo:       &fakeProfilSiswaRepo{result: wantSiswa},
			wantSiswa:  wantSiswa,
			wantCalled: true,
		},
		{
			name:       "error",
			id:         12,
			repo:       &fakeProfilSiswaRepo{err: repoErr},
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

			assert.Equal(t, tc.wantCalled, tc.repo.called)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.wantSiswa, result)
		})
	}
}
