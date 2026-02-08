package user_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/get"
	"github.com/stretchr/testify/assert"
)

type fakeGetSiswaRepo struct {
	items     []query.SiswaListItem
	err       error
	called    bool
	gotFilter query.ListSiswaFilter
}

func (f *fakeGetSiswaRepo) GetListSiswa(ctx context.Context, filter query.ListSiswaFilter) ([]query.SiswaListItem, error) {
	f.called = true
	f.gotFilter = filter
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

func (f *fakeGetSiswaRepo) FindProfilSiswaByID(ctx context.Context, id user.ID) (user.DataSiswa, error) {
	return user.DataSiswa{}, nil
}

func (f *fakeGetSiswaRepo) ExistByNISN(ctx context.Context, nisn string) (bool, error) {
	return false, nil
}

func (f *fakeGetSiswaRepo) CreateProfilSiswa(ctx context.Context, g user.ProfilSiswa) (user.ID, error) {
	return 0, nil
}

func (f *fakeGetSiswaRepo) UpdateProfilSiswa(ctx context.Context, idPengguna user.ID, profilSiswa outuser.UpdateProfilSiswaPatch) error {
	return nil
}

func TestGetSiswaService_ListSiswa(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	statusAktif := user.AKTIF
	invalidStatus := user.StatusAkun("BANNED")
	currentYear := time.Now().Year()
	validAngkatan := currentYear
	invalidAngkatan := 2010
	validTingkat := 1
	invalidTingkat := -2
	validJenisKelamin := 2
	invalidJenisKelamin := 3
	items := []query.SiswaListItem{{Username: "siswa-1"}}
	repoErr := errors.New("repo error")

	tests := []struct {
		name       string
		filter     query.ListSiswaFilter
		repo       *fakeGetSiswaRepo
		wantItems  []query.SiswaListItem
		wantErr    error
		wantErrMsg string
		wantCalled bool
		wantFilter *query.ListSiswaFilter
	}{
		{
			name: "Branch 1 -> semua patch berhasil",
			filter: query.ListSiswaFilter{
				Search:       "  siswa ",
				Limit:        0,
				Offset:       -1,
				SortBy:       "",
				SortDesc:     true,
				Angkatan:     &validAngkatan,
				TingkatKelas: &validTingkat,
				JenisKelamin: &validJenisKelamin,
			},
			repo:       &fakeGetSiswaRepo{items: items},
			wantItems:  items,
			wantCalled: true,
			wantFilter: &query.ListSiswaFilter{
				Search:       "siswa",
				Limit:        20,
				Offset:       0,
				SortBy:       "created_at",
				SortDesc:     true,
				Status:       &statusAktif,
				Angkatan:     &validAngkatan,
				TingkatKelas: &validTingkat,
				JenisKelamin: &validJenisKelamin,
			},
		},
		{
			name: "Branch 2 -> sortBy tidak valid",
			filter: query.ListSiswaFilter{
				SortBy: "invalid",
			},
			repo:       &fakeGetSiswaRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name: "Branch 3 -> status tidak valid",
			filter: query.ListSiswaFilter{
				SortBy: "created_at",
				Status: &invalidStatus,
			},
			repo:       &fakeGetSiswaRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name: "Branch 4 -> angkatan tidak valid",
			filter: query.ListSiswaFilter{
				SortBy:   "created_at",
				Angkatan: &invalidAngkatan,
			},
			repo:       &fakeGetSiswaRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name: "Branch 5 -> tingkat kelas tidak valid",
			filter: query.ListSiswaFilter{
				SortBy:       "created_at",
				TingkatKelas: &invalidTingkat,
			},
			repo:       &fakeGetSiswaRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name: "Branch 6 -> jenis kelamin tidak valid",
			filter: query.ListSiswaFilter{
				SortBy:       "created_at",
				JenisKelamin: &invalidJenisKelamin,
			},
			repo:       &fakeGetSiswaRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name: "Branch 7 -> repo error",
			filter: query.ListSiswaFilter{
				SortBy: "created_at",
			},
			repo:       &fakeGetSiswaRepo{err: repoErr},
			wantErrMsg: "repo error",
			wantCalled: true,
			wantFilter: &query.ListSiswaFilter{
				Limit:  20,
				Offset: 0,
				SortBy: "created_at",
				Status: &statusAktif,
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			service := user_service.NewGetListSiswaService(tc.repo, tc.repo)

			result, err := service.ListSiswa(ctx, tc.filter)

			assert.Equal(t, tc.wantCalled, tc.repo.called)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			if tc.wantErrMsg != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.wantErrMsg)
				if tc.wantFilter != nil {
					assert.Equal(t, *tc.wantFilter, tc.repo.gotFilter)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantItems, result)
			if tc.wantFilter != nil {
				assert.Equal(t, *tc.wantFilter, tc.repo.gotFilter)
			}
		})
	}
}
