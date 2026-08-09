package ruangujian_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ruang_ujian"
	ruangujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/create"
	"github.com/stretchr/testify/assert"
)

type fakeCreateRepo struct {
	existFn  func(ctx context.Context, kodeRuang string) (bool, error)
	createFn func(ctx context.Context, ruangUjian ruangujian.RuangUjian) error
}

func (f *fakeCreateRepo) GetRuangUjian(context.Context, query.ListRuangUjianFilter) ([]ruangujian.RuangUjian, error) {
	return nil, nil
}
func (f *fakeCreateRepo) GetRuangUjianById(context.Context, int) (ruangujian.RuangUjian, error) {
	return ruangujian.RuangUjian{}, nil
}
func (f *fakeCreateRepo) GetRuangUjianByKode(context.Context, string) (ruangujian.RuangUjian, error) {
	return ruangujian.RuangUjian{}, nil
}
func (f *fakeCreateRepo) ExistByKodeRuang(ctx context.Context, kodeRuang string) (bool, error) {
	if f.existFn != nil {
		return f.existFn(ctx, kodeRuang)
	}
	return false, nil
}
func (f *fakeCreateRepo) CreateRuangUjian(ctx context.Context, ruangUjian ruangujian.RuangUjian) error {
	if f.createFn != nil {
		return f.createFn(ctx, ruangUjian)
	}
	return nil
}
func (f *fakeCreateRepo) UpdateRuangUjian(context.Context, int, updatepatch.UpdateRuangUjianPatch) error {
	return nil
}
func (f *fakeCreateRepo) DeleteRuangUjian(context.Context, int) error { return nil }

func TestCreateRuangUjianService(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("repo error")
	ctx := context.Background()

	tests := []struct {
		name      string
		repo      *fakeCreateRepo
		payload   ruangujian.RuangUjian
		expectErr error
	}{
		{
			name:      "branch 1 -> gagal cek kode ruang",
			repo:      &fakeCreateRepo{existFn: func(context.Context, string) (bool, error) { return false, repoErr }},
			payload:   ruangujian.RuangUjian{KodeRuang: " a1 ", NamaRuangan: " Ruang A "},
			expectErr: repoErr,
		},
		{
			name:      "branch 2 -> kode sudah ada",
			repo:      &fakeCreateRepo{existFn: func(context.Context, string) (bool, error) { return true, nil }},
			payload:   ruangujian.RuangUjian{KodeRuang: "a1", NamaRuangan: "Ruang A"},
			expectErr: coreerror.ErrKodeRuangUjianExist,
		},
		{
			name: "branch 3 -> gagal create",
			repo: &fakeCreateRepo{
				existFn:  func(context.Context, string) (bool, error) { return false, nil },
				createFn: func(context.Context, ruangujian.RuangUjian) error { return repoErr },
			},
			payload:   ruangujian.RuangUjian{KodeRuang: "a1", NamaRuangan: "Ruang A"},
			expectErr: repoErr,
		},
		{
			name: "happy path -> normalisasi input dan berhasil",
			repo: &fakeCreateRepo{
				existFn: func(_ context.Context, kode string) (bool, error) {
					assert.Equal(t, "A-01", kode)
					return false, nil
				},
				createFn: func(_ context.Context, ruang ruangujian.RuangUjian) error {
					assert.Equal(t, "A-01", ruang.KodeRuang)
					assert.Equal(t, "Lab Utama", ruang.NamaRuangan)
					return nil
				},
			},
			payload: ruangujian.RuangUjian{KodeRuang: " a-01 ", NamaRuangan: " Lab Utama "},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			svc := ruangujian_service.NewRuangUjianService(tc.repo)
			err := svc.CreateRuangUjianService(ctx, tc.payload)
			assert.ErrorIs(t, err, tc.expectErr)
		})
	}
}

func TestCreateRuangUjianService_KodeSudahAda(t *testing.T) {
	repo := &fakeCreateRepo{
		existFn: func(context.Context, string) (bool, error) {
			return true, nil
		},
	}
	svc := ruangujian_service.NewRuangUjianService(repo)

	err := svc.CreateRuangUjianService(context.Background(),
		ruangujian.RuangUjian{KodeRuang: "A-01"})

	assert.ErrorIs(t, err, coreerror.ErrKodeRuangUjianExist)
}
// func TestCreateRuangUjianService_RepoGagal(t *testing.T) {
// 	repo := &fakeCreateRepo{
		
// 	}
// }
