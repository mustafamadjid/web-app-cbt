package user_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/get"
	"github.com/stretchr/testify/assert"
)

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
		repo       *FakeGetSiswaRepo
		wantItems  []query.SiswaListItem
		wantErr    error
		wantErrMsg string
		wantCalled bool
		wantFilter *query.ListSiswaFilter
	}{
		{
			name: "Path 1 -> filter valid berhasil dilist",
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
			repo:       &FakeGetSiswaRepo{Items: items},
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
			name: "Path 2 -> sortBy tidak valid",
			filter: query.ListSiswaFilter{
				SortBy: "invalid",
			},
			repo:       &FakeGetSiswaRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name: "Path 3 -> status tidak valid",
			filter: query.ListSiswaFilter{
				SortBy: "created_at",
				Status: &invalidStatus,
			},
			repo:       &FakeGetSiswaRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name: "Path 4 -> angkatan tidak valid",
			filter: query.ListSiswaFilter{
				SortBy:   "created_at",
				Angkatan: &invalidAngkatan,
			},
			repo:       &FakeGetSiswaRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name: "Path 5 -> tingkat kelas tidak valid",
			filter: query.ListSiswaFilter{
				SortBy:       "created_at",
				TingkatKelas: &invalidTingkat,
			},
			repo:       &FakeGetSiswaRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name: "Path 6 -> jenis kelamin tidak valid",
			filter: query.ListSiswaFilter{
				SortBy:       "created_at",
				JenisKelamin: &invalidJenisKelamin,
			},
			repo:       &FakeGetSiswaRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name: "Path 7 -> repo error",
			filter: query.ListSiswaFilter{
				SortBy: "created_at",
			},
			repo:       &FakeGetSiswaRepo{Err: repoErr},
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

			assert.Equal(t, tc.wantCalled, tc.repo.Called)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			if tc.wantErrMsg != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.wantErrMsg)
				if tc.wantFilter != nil {
					assert.Equal(t, *tc.wantFilter, tc.repo.GotFilter)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantItems, result)
			if tc.wantFilter != nil {
				assert.Equal(t, *tc.wantFilter, tc.repo.GotFilter)
			}
		})
	}
}
