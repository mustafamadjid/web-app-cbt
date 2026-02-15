package ruangujian_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ruang_ujian"
	ruangujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/get"
	"github.com/stretchr/testify/assert"
)

type fakeGetRepo struct {
	getListFn   func(context.Context, query.ListRuangUjianFilter) ([]ruangujian.RuangUjian, error)
	getByIDFn   func(context.Context, int) (ruangujian.RuangUjian, error)
	getByKodeFn func(context.Context, string) (ruangujian.RuangUjian, error)
}

func (f *fakeGetRepo) GetRuangUjian(ctx context.Context, filter query.ListRuangUjianFilter) ([]ruangujian.RuangUjian, error) {
	if f.getListFn != nil {
		return f.getListFn(ctx, filter)
	}
	return nil, nil
}
func (f *fakeGetRepo) GetRuangUjianById(ctx context.Context, id int) (ruangujian.RuangUjian, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return ruangujian.RuangUjian{}, nil
}
func (f *fakeGetRepo) GetRuangUjianByKode(ctx context.Context, kode string) (ruangujian.RuangUjian, error) {
	if f.getByKodeFn != nil {
		return f.getByKodeFn(ctx, kode)
	}
	return ruangujian.RuangUjian{}, nil
}
func (f *fakeGetRepo) ExistByKodeRuang(context.Context, string) (bool, error)        { return false, nil }
func (f *fakeGetRepo) CreateRuangUjian(context.Context, ruangujian.RuangUjian) error { return nil }
func (f *fakeGetRepo) UpdateRuangUjian(context.Context, int, updatepatch.UpdateRuangUjianPatch) error {
	return nil
}
func (f *fakeGetRepo) DeleteRuangUjian(context.Context, int) error { return nil }

func TestGetRuangUjian(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repoErr := errors.New("repo error")
	expected := []ruangujian.RuangUjian{{IdRuangan: 1, KodeRuang: "A1", NamaRuangan: "Lab"}}

	tests := []struct {
		name      string
		filter    query.ListRuangUjianFilter
		repo      *fakeGetRepo
		want      []ruangujian.RuangUjian
		expectErr error
	}{
		{
			name:   "branch 1 -> gagal get list",
			filter: query.ListRuangUjianFilter{Limit: 0, Offset: -1, Search: "  abc  "},
			repo: &fakeGetRepo{getListFn: func(_ context.Context, got query.ListRuangUjianFilter) ([]ruangujian.RuangUjian, error) {
				assert.Equal(t, 20, got.Limit)
				assert.Equal(t, 0, got.Offset)
				assert.Equal(t, "abc", got.Search)
				return nil, repoErr
			}},
			expectErr: repoErr,
		},
		{
			name:   "loop coverage -> limit lebih dari 50 dan berhasil",
			filter: query.ListRuangUjianFilter{Limit: 99, Offset: 3, Search: " test "},
			repo: &fakeGetRepo{getListFn: func(_ context.Context, got query.ListRuangUjianFilter) ([]ruangujian.RuangUjian, error) {
				assert.Equal(t, 50, got.Limit)
				assert.Equal(t, 3, got.Offset)
				assert.Equal(t, "test", got.Search)
				return expected, nil
			}},
			want: expected,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			svc := ruangujian_service.NewGetRuangUjianService(tc.repo)
			got, err := svc.GetRuangUjian(ctx, tc.filter)
			assert.ErrorIs(t, err, tc.expectErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGetRuangUjianById(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repoErr := errors.New("repo error")
	expected := ruangujian.RuangUjian{IdRuangan: 1, KodeRuang: "A1", NamaRuangan: "Lab"}

	tests := []struct {
		name      string
		id        int
		repo      *fakeGetRepo
		want      ruangujian.RuangUjian
		expectErr error
	}{
		{name: "branch 1 -> id tidak valid", id: 0, repo: &fakeGetRepo{}, expectErr: coreerror.ErrMissingId},
		{name: "branch 2 -> repo error", id: 1, repo: &fakeGetRepo{getByIDFn: func(context.Context, int) (ruangujian.RuangUjian, error) { return ruangujian.RuangUjian{}, repoErr }}, expectErr: repoErr},
		{name: "loop coverage -> berhasil", id: 1, repo: &fakeGetRepo{getByIDFn: func(context.Context, int) (ruangujian.RuangUjian, error) { return expected, nil }}, want: expected},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			svc := ruangujian_service.NewGetRuangUjianService(tc.repo)
			got, err := svc.GetRuangUjianById(ctx, tc.id)
			assert.ErrorIs(t, err, tc.expectErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGetRuangUjianByKode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repoErr := errors.New("repo error")
	expected := ruangujian.RuangUjian{IdRuangan: 2, KodeRuang: "B1", NamaRuangan: "Ruang B"}

	tests := []struct {
		name      string
		kode      string
		repo      *fakeGetRepo
		want      ruangujian.RuangUjian
		expectErr error
	}{
		{name: "branch 1 -> kode tidak ditemukan", kode: "   ", repo: &fakeGetRepo{}, expectErr: coreerror.ErrMissingField},
		{name: "branch 2 -> repo error", kode: " b1 ", repo: &fakeGetRepo{getByKodeFn: func(_ context.Context, got string) (ruangujian.RuangUjian, error) {
			assert.Equal(t, "B1", got)
			return ruangujian.RuangUjian{}, repoErr
		}}, expectErr: repoErr},
		{name: "loop coverage -> normalisasi kode dan berhasil", kode: " b1 ", repo: &fakeGetRepo{getByKodeFn: func(_ context.Context, got string) (ruangujian.RuangUjian, error) {
			assert.Equal(t, "B1", got)
			return expected, nil
		}}, want: expected},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			svc := ruangujian_service.NewGetRuangUjianService(tc.repo)
			got, err := svc.GetRuangUjianByKode(ctx, tc.kode)
			assert.ErrorIs(t, err, tc.expectErr)
			assert.Equal(t, tc.want, got)
		})
	}
}
