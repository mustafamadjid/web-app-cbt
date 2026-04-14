package gradingujian_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	gradingujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/list/list_ujian_essay_ungraded"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeListEssayUngradedRepo struct {
	listRet    []ujian.ListUjian
	listErr    error
	listCalled bool
	gotFilter  query.ListUjianEssayUngradedFilter
}

func (f *fakeListEssayUngradedRepo) ListUjianEssayUngraded(_ context.Context, filter query.ListUjianEssayUngradedFilter) ([]ujian.ListUjian, error) {
	f.listCalled = true
	f.gotFilter = filter
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.listRet, nil
}

func stringPtr(v string) *string {
	return &v
}

func runListEssayUngradedCases(t *testing.T, prefix string) {
	t.Helper()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	tanggal := "2026-06-15"
	expected := []ujian.ListUjian{{IdUjian: 1, NamaUjian: "Essay UAS", TanggalUjian: time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)}}

	tests := []struct {
		name         string
		filter       query.ListUjianEssayUngradedFilter
		repo         *fakeListEssayUngradedRepo
		wantErr      error
		wantCalled   bool
		wantItems    []ujian.ListUjian
		assertFilter func(t *testing.T, got query.ListUjianEssayUngradedFilter)
	}{
		{
			name:       prefix + "1 -> filter list essay ungraded tidak valid",
			filter:     query.ListUjianEssayUngradedFilter{TanggalUjian: stringPtr(" ")},
			repo:       &fakeListEssayUngradedRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name:       prefix + "2 -> repo list essay ungraded gagal",
			filter:     query.ListUjianEssayUngradedFilter{TanggalUjian: &tanggal},
			repo:       &fakeListEssayUngradedRepo{listErr: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
		},
		{
			name:       prefix + "3 -> berhasil list essay ungraded",
			filter:     query.ListUjianEssayUngradedFilter{Search: "  essay  ", Limit: 70, Offset: -2},
			repo:       &fakeListEssayUngradedRepo{listRet: expected},
			wantCalled: true,
			wantItems:  expected,
			assertFilter: func(t *testing.T, got query.ListUjianEssayUngradedFilter) {
				t.Helper()
				assert.Equal(t, "essay", got.Search)
				assert.Equal(t, 50, got.Limit)
				assert.Equal(t, 0, got.Offset)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := gradingujian_service.NewListUjianEssayUngradedService(tc.repo)
			got, err := svc.ListUjianEssayUngraded(ctx, tc.filter)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantCalled, tc.repo.listCalled)
			assert.Equal(t, tc.wantItems, got)
			if tc.wantCalled && tc.assertFilter != nil {
				tc.assertFilter(t, tc.repo.gotFilter)
			}
		})
	}
}

func TestListUjianEssayUngradedService_BranchCoverage(t *testing.T) {
	t.Parallel()
	runListEssayUngradedCases(t, "branch ")
}

func TestListUjianEssayUngradedService_BasisPath(t *testing.T) {
	t.Parallel()
	runListEssayUngradedCases(t, "path ")
}

func TestListUjianEssayUngradedService_FilterSanitizer(t *testing.T) {
	t.Parallel()

	repo := &fakeListEssayUngradedRepo{listRet: []ujian.ListUjian{}}
	svc := gradingujian_service.NewListUjianEssayUngradedService(repo)
	tanggal := " 2026-06-15 "
	tahun := " 2026 "
	bulan := " 6 "

	_, err := svc.ListUjianEssayUngraded(context.Background(), query.ListUjianEssayUngradedFilter{
		Search:       "  kelas  ",
		Limit:        0,
		Offset:       -1,
		TanggalUjian: &tanggal,
		Tahun:        &tahun,
		Bulan:        &bulan,
	})

	require.NoError(t, err)
	assert.Equal(t, "kelas", repo.gotFilter.Search)
	assert.Equal(t, 20, repo.gotFilter.Limit)
	assert.Equal(t, 0, repo.gotFilter.Offset)
	require.NotNil(t, repo.gotFilter.TanggalUjian)
	assert.Equal(t, "2026-06-15", *repo.gotFilter.TanggalUjian)
	require.NotNil(t, repo.gotFilter.Tahun)
	assert.Equal(t, "2026", *repo.gotFilter.Tahun)
	require.NotNil(t, repo.gotFilter.Bulan)
	assert.Equal(t, "6", *repo.gotFilter.Bulan)
}
