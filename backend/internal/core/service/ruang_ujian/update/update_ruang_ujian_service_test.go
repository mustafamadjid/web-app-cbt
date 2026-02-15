package ruangujian_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ruang_ujian"
	ruangujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/update"
	"github.com/stretchr/testify/assert"
)

type fakeUpdateRepo struct {
	existFn  func(context.Context, string) (bool, error)
	updateFn func(context.Context, int, updatepatch.UpdateRuangUjianPatch) error
}

func (f *fakeUpdateRepo) GetRuangUjian(context.Context, query.ListRuangUjianFilter) ([]ruangujian.RuangUjian, error) {
	return nil, nil
}
func (f *fakeUpdateRepo) GetRuangUjianById(context.Context, int) (ruangujian.RuangUjian, error) {
	return ruangujian.RuangUjian{}, nil
}
func (f *fakeUpdateRepo) GetRuangUjianByKode(context.Context, string) (ruangujian.RuangUjian, error) {
	return ruangujian.RuangUjian{}, nil
}
func (f *fakeUpdateRepo) ExistByKodeRuang(ctx context.Context, kode string) (bool, error) {
	if f.existFn != nil {
		return f.existFn(ctx, kode)
	}
	return false, nil
}
func (f *fakeUpdateRepo) CreateRuangUjian(context.Context, ruangujian.RuangUjian) error { return nil }
func (f *fakeUpdateRepo) UpdateRuangUjian(ctx context.Context, id int, patch updatepatch.UpdateRuangUjianPatch) error {
	if f.updateFn != nil {
		return f.updateFn(ctx, id, patch)
	}
	return nil
}
func (f *fakeUpdateRepo) DeleteRuangUjian(context.Context, int) error { return nil }

func TestUpdateRuangUjianService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")

	empty := "   "
	kode := " a-01 "
	nama := " Ruang Utama "

	tests := []struct {
		name      string
		idRuangan int
		patch     updatepatch.UpdateRuangUjianPatch
		repo      *fakeUpdateRepo
		expectErr error
	}{
		{name: "branch 1 -> id tidak valid", idRuangan: 0, patch: updatepatch.UpdateRuangUjianPatch{}, repo: &fakeUpdateRepo{}, expectErr: coreerror.ErrMissingId},
		{
			name:      "branch 2 -> gagal cek kode ruang",
			idRuangan: 1,
			patch:     updatepatch.UpdateRuangUjianPatch{KodeRuang: &empty},
			repo: &fakeUpdateRepo{existFn: func(_ context.Context, got string) (bool, error) {
				assert.Equal(t, "", got)
				return false, repoErr
			}},
			expectErr: repoErr,
		},
		{
			name:      "branch 3 -> kode sudah ada",
			idRuangan: 1,
			patch:     updatepatch.UpdateRuangUjianPatch{KodeRuang: &kode},
			repo: &fakeUpdateRepo{existFn: func(_ context.Context, got string) (bool, error) {
				assert.Equal(t, "A-01", got)
				return true, nil
			}},
			expectErr: coreerror.ErrKodeRuangUjianExist,
		},
		{
			name:      "branch 4 -> gagal update",
			idRuangan: 2,
			patch:     updatepatch.UpdateRuangUjianPatch{KodeRuang: &kode, NamaRuang: &empty},
			repo: &fakeUpdateRepo{
				existFn: func(context.Context, string) (bool, error) { return false, nil },
				updateFn: func(_ context.Context, id int, patch updatepatch.UpdateRuangUjianPatch) error {
					assert.Equal(t, 2, id)
					assert.Equal(t, "A-01", *patch.KodeRuang)
					assert.Equal(t, "", *patch.NamaRuang)
					return repoErr
				},
			},
			expectErr: repoErr,
		},
		{
			name:      "loop coverage -> normalisasi patch dan berhasil",
			idRuangan: 3,
			patch:     updatepatch.UpdateRuangUjianPatch{KodeRuang: &kode, NamaRuang: &nama},
			repo: &fakeUpdateRepo{
				existFn: func(context.Context, string) (bool, error) { return false, nil },
				updateFn: func(_ context.Context, id int, patch updatepatch.UpdateRuangUjianPatch) error {
					assert.Equal(t, 3, id)
					assert.Equal(t, "A-01", *patch.KodeRuang)
					assert.Equal(t, "Ruang Utama", *patch.NamaRuang)
					return nil
				},
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			svc := ruangujian_service.NewUpdateRuangUjianService(tc.repo)
			err := svc.UpdateRuangUjian(ctx, tc.idRuangan, tc.patch)
			assert.ErrorIs(t, err, tc.expectErr)
		})
	}
}
