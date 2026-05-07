package siswaujian_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	siswaujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/list"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeListUjianSiswaRepo struct {
	listRet    []ujian.ListUjian
	listErr    error
	listCalled bool
	gotSiswaID int
	gotFilter  query.ListUjianFilter
}

func (f *fakeListUjianSiswaRepo) ListUjianSiswa(_ context.Context, idSiswa int, filter query.ListUjianFilter) ([]ujian.ListUjian, error) {
	f.listCalled = true
	f.gotSiswaID = idSiswa
	f.gotFilter = filter
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.listRet, nil
}

func (*fakeListUjianSiswaRepo) GetWaktuSelesaiUjian(context.Context, int) (time.Time, error) {
	return time.Time{}, nil
}

func (*fakeListUjianSiswaRepo) GetActiveUjianAttemptBySiswa(context.Context, int, int) (ujian.AttemptUjian, error) {
	return ujian.AttemptUjian{}, nil
}

func strPtr(s string) *string {
	return &s
}

func TestListUjianSiswaService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	filterDate := "2026-05-15"
	expected := []ujian.ListUjian{{IdUjian: 5, NamaUjian: "UTS"}}

	tests := []struct {
		name         string
		idSiswa      int
		filter       query.ListUjianFilter
		repo         *fakeListUjianSiswaRepo
		wantErr      error
		wantCalled   bool
		wantItems    []ujian.ListUjian
		assertFilter func(t *testing.T, got query.ListUjianFilter)
	}{
		{
			name:       "Path 1 -> id siswa tidak valid",
			idSiswa:    0,
			filter:     query.ListUjianFilter{},
			repo:       &fakeListUjianSiswaRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantItems:  nil,
			wantCalled: false,
		},
		{
			name:      "Path 2 -> filter tidak valid",
			idSiswa:   7,
			filter:    query.ListUjianFilter{TanggalUjian: strPtr("  ")},
			repo:      &fakeListUjianSiswaRepo{},
			wantErr:   coreerror.ErrInvalidInput,
			wantItems: nil,
		},
		{
			name:       "Path 3 -> repo list ujian siswa gagal",
			idSiswa:    7,
			filter:     query.ListUjianFilter{TanggalUjian: &filterDate},
			repo:       &fakeListUjianSiswaRepo{listErr: repoErr},
			wantErr:    repoErr,
			wantItems:  nil,
			wantCalled: true,
		},
		{
			name:       "Path 4 -> berhasil list ujian siswa",
			idSiswa:    7,
			filter:     query.ListUjianFilter{Search: "  uts  ", Limit: 0, Offset: -1, KategoriUjian: query.SELESAI},
			repo:       &fakeListUjianSiswaRepo{listRet: expected},
			wantItems:  expected,
			wantCalled: true,
			assertFilter: func(t *testing.T, got query.ListUjianFilter) {
				t.Helper()
				assert.Equal(t, "uts", got.Search)
				assert.Equal(t, 20, got.Limit)
				assert.Equal(t, 0, got.Offset)
				assert.Equal(t, query.SELESAI, got.KategoriUjian)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := siswaujian_service.NewListUjianSiswaService(tc.repo)
			got, err := svc.ListUjianSiswa(ctx, tc.idSiswa, tc.filter)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantCalled, tc.repo.listCalled)
			assert.Equal(t, tc.wantItems, got)
			if tc.wantCalled {
				assert.Equal(t, tc.idSiswa, tc.repo.gotSiswaID)
				if tc.assertFilter != nil {
					tc.assertFilter(t, tc.repo.gotFilter)
				}
			}
		})
	}
}

func TestListUjianSiswaService_SanitizeKategori(t *testing.T) {
	t.Parallel()

	repo := &fakeListUjianSiswaRepo{listRet: []ujian.ListUjian{}}
	svc := siswaujian_service.NewListUjianSiswaService(repo)

	_, err := svc.ListUjianSiswa(context.Background(), 8, query.ListUjianFilter{
		Search:        "  mapel  ",
		Limit:         99,
		Offset:        -4,
		KategoriUjian: query.BERLANGSUNG,
	})

	require.NoError(t, err)
	assert.Equal(t, "mapel", repo.gotFilter.Search)
	assert.Equal(t, 50, repo.gotFilter.Limit)
	assert.Equal(t, 0, repo.gotFilter.Offset)
	assert.Equal(t, query.BERLANGSUNG, repo.gotFilter.KategoriUjian)
}
