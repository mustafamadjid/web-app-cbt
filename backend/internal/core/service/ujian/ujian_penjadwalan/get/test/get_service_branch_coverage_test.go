package ujian_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/get"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeGetUjianRepo struct {
	getAllRet    []ujian.ListUjian
	getAllErr    error
	getAllCalled bool
	gotFilter    query.ListUjianFilter

	getByIDRet    ujian.ListUjian
	getByIDErr    error
	getByIDCalled bool
	gotID         ujian.ID
}

func strPtr(s string) *string {
	return &s
}

func (f *fakeGetUjianRepo) GetAllUjian(_ context.Context, filter query.ListUjianFilter) ([]ujian.ListUjian, error) {
	f.getAllCalled = true
	f.gotFilter = filter
	if f.getAllErr != nil {
		return nil, f.getAllErr
	}
	return f.getAllRet, nil
}

func (f *fakeGetUjianRepo) GetUjianById(_ context.Context, id ujian.ID) (ujian.ListUjian, error) {
	f.getByIDCalled = true
	f.gotID = id
	if f.getByIDErr != nil {
		return ujian.ListUjian{}, f.getByIDErr
	}
	return f.getByIDRet, nil
}

func buildValidListUjian() ujian.ListUjian {
	namaKelas := "  XII-A  "
	status := ujian.MULAI
	return ujian.ListUjian{
		IdUjian:          1,
		IdBankSoal:       2,
		IdGuru:           3,
		NamaUjian:        "  UAS  ",
		PembuatUsername:  "  guru1  ",
		IdKelas:          4,
		IdJadwalUjian:    5,
		IdPengawas:       6,
		IdSesi:           7,
		IdRuangan:        8,
		TingkatKelas:     12,
		NamaKelas:        &namaKelas,
		TanggalUjian:     time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC),
		WaktuMulai:       time.Date(2026, 6, 10, 8, 0, 0, 0, time.UTC),
		WaktuSelesai:     time.Date(2026, 6, 10, 10, 0, 0, 0, time.UTC),
		StatusUjian:      &status,
		NamaPengawas:     "  Pengawas  ",
		PengawasUsername: "  pengawas1  ",
		NamaSesi:         "  Pagi  ",
		NamaRuangan:      "  Lab 1  ",
	}
}

func TestGetUjianService_GetAllUjian_BranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	filterDate := "2026-06-10"
	expected := buildValidListUjian()

	tests := []struct {
		name         string
		filter       query.ListUjianFilter
		repo         *fakeGetUjianRepo
		wantErr      error
		wantCalled   bool
		wantItems    []ujian.ListUjian
		assertFilter func(t *testing.T, got query.ListUjianFilter)
	}{
		{
			name:       "branch 1 -> filter ujian tidak valid",
			filter:     query.ListUjianFilter{TanggalUjian: strPtr(" ")},
			repo:       &fakeGetUjianRepo{},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: false,
		},
		{
			name:       "branch 2 -> repo get all ujian gagal",
			filter:     query.ListUjianFilter{TanggalUjian: &filterDate},
			repo:       &fakeGetUjianRepo{getAllErr: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
		},
		{
			name:       "branch 3 -> item ujian dari repo tidak valid",
			filter:     query.ListUjianFilter{},
			repo:       &fakeGetUjianRepo{getAllRet: []ujian.ListUjian{{IdUjian: 0}}},
			wantErr:    coreerror.ErrInvalidInput,
			wantCalled: true,
		},
		{
			name:       "branch 4 -> berhasil get all ujian",
			filter:     query.ListUjianFilter{Search: "  uas  ", Limit: 70, Offset: -3, KategoriUjian: query.BERLANGSUNG},
			repo:       &fakeGetUjianRepo{getAllRet: []ujian.ListUjian{expected}},
			wantCalled: true,
			wantItems: []ujian.ListUjian{func() ujian.ListUjian {
				item := expected
				namaKelas := "XII-A"
				item.NamaKelas = &namaKelas
				item.NamaUjian = "UAS"
				item.PembuatUsername = "guru1"
				item.NamaPengawas = "Pengawas"
				item.PengawasUsername = "pengawas1"
				item.NamaSesi = "Pagi"
				item.NamaRuangan = "Lab 1"
				return item
			}()},
			assertFilter: func(t *testing.T, got query.ListUjianFilter) {
				t.Helper()
				assert.Equal(t, "uas", got.Search)
				assert.Equal(t, 50, got.Limit)
				assert.Equal(t, 0, got.Offset)
				assert.Equal(t, query.BERLANGSUNG, got.KategoriUjian)
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := ujian_service.NewGetujianService(tc.repo)
			got, err := svc.GetAllUjianService(ctx, tc.filter)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantCalled, tc.repo.getAllCalled)
			assert.Equal(t, tc.wantItems, got)
			if tc.wantCalled && tc.assertFilter != nil {
				tc.assertFilter(t, tc.repo.gotFilter)
			}
		})
	}
}

func TestGetUjianService_GetUjianByID_BranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	expected := buildValidListUjian()

	tests := []struct {
		name       string
		idUjian    ujian.ID
		repo       *fakeGetUjianRepo
		wantErr    error
		wantCalled bool
		wantItem   ujian.ListUjian
	}{
		{
			name:       "branch 5 -> id ujian tidak valid",
			idUjian:    0,
			repo:       &fakeGetUjianRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name:       "branch 6 -> repo get ujian by id gagal",
			idUjian:    7,
			repo:       &fakeGetUjianRepo{getByIDErr: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
		},
		{
			name:       "branch 7 -> berhasil get ujian by id",
			idUjian:    7,
			repo:       &fakeGetUjianRepo{getByIDRet: expected},
			wantCalled: true,
			wantItem:   expected,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := ujian_service.NewGetujianService(tc.repo)
			got, err := svc.GetUjianByIdService(ctx, tc.idUjian)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.wantCalled, tc.repo.getByIDCalled)
			if tc.wantCalled {
				assert.Equal(t, tc.idUjian, tc.repo.gotID)
			}
			assert.Equal(t, tc.wantItem, got)
		})
	}
}

func TestGetUjianService_FilterSanitizer(t *testing.T) {
	t.Parallel()

	repo := &fakeGetUjianRepo{getAllRet: []ujian.ListUjian{buildValidListUjian()}}
	svc := ujian_service.NewGetujianService(repo)
	date := " 2026-06-10 "
	year := " 2026 "

	items, err := svc.GetAllUjianService(context.Background(), query.ListUjianFilter{
		Search:        "  mapel  ",
		Limit:         0,
		Offset:        -1,
		TanggalUjian:  &date,
		Tahun:         &year,
		KategoriUjian: query.SELESAI,
	})

	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "mapel", repo.gotFilter.Search)
	assert.Equal(t, 20, repo.gotFilter.Limit)
	assert.Equal(t, 0, repo.gotFilter.Offset)
	require.NotNil(t, repo.gotFilter.TanggalUjian)
	assert.Equal(t, "2026-06-10", *repo.gotFilter.TanggalUjian)
	require.NotNil(t, repo.gotFilter.Tahun)
	assert.Equal(t, "2026", *repo.gotFilter.Tahun)
}
