package user_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/get"
	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/get/fake_test"
	"github.com/stretchr/testify/assert"
)

func TestGetGuruService_ListGuru(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	bidangInput := "  Matematika  "
	bidangExpected := "Matematika"
	invalidStatus := user.StatusAkun("BANNED")
	statusAktif := user.AKTIF
	items := []query.GuruListItem{{Username: "guru-1"}}
	repoErr := errors.New("repo error")

	tests := []struct {
		name       string
		filter     query.ListGuruFilter
		repo       *fake_test.FakeGetGuruRepo
		wantItems  []query.GuruListItem
		wantErr    error
		wantErrMsg string
		wantCalled bool
		wantFilter *query.ListGuruFilter
	}{
		{
			name: "Branch 1 -> semua patch berhasil",
			filter: query.ListGuruFilter{
				Search:   "  guru ",
				Limit:    0,
				Offset:   -10,
				SortBy:   "",
				Bidang:   &bidangInput,
				SortDesc: true,
			},
			repo:       &fake_test.FakeGetGuruRepo{Items: items},
			wantItems:  items,
			wantCalled: true,
			wantFilter: &query.ListGuruFilter{
				Search:   "guru",
				Limit:    20,
				Offset:   0,
				SortBy:   "created_at",
				SortDesc: true,
				Status:   &statusAktif,
				Bidang:   &bidangExpected,
			},
		},
		{
			name: "Branch 2 -> sortBy tidak valid",
			filter: query.ListGuruFilter{
				SortBy: "invalid",
			},
			repo:       &fake_test.FakeGetGuruRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name: "Branch 3 -> status tidak valid",
			filter: query.ListGuruFilter{
				SortBy: "created_at",
				Status: &invalidStatus,
			},
			repo:       &fake_test.FakeGetGuruRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name: "Branch 4 -> bidang kosong",
			filter: query.ListGuruFilter{
				SortBy: "created_at",
				Bidang: func() *string {
					v := "   "
					return &v
				}(),
			},
			repo:       &fake_test.FakeGetGuruRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name: "Branch 5 -> repo error",
			filter: query.ListGuruFilter{
				SortBy: "created_at",
			},
			repo:       &fake_test.FakeGetGuruRepo{Err: repoErr},
			wantErrMsg: "repo error",
			wantCalled: true,
			wantFilter: &query.ListGuruFilter{
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
			service := user_service.NewGetListGuruService(tc.repo, tc.repo)

			result, err := service.ListGuru(ctx, tc.filter)

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
