package ruangujian_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ruang_ujian"
	ruangujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/delete"
	"github.com/stretchr/testify/assert"
)

type fakeDeleteRepo struct {
	deleteFn func(context.Context, int) error
}

func (f *fakeDeleteRepo) GetRuangUjian(context.Context, query.ListRuangUjianFilter) ([]ruangujian.RuangUjian, error) {
	return nil, nil
}
func (f *fakeDeleteRepo) GetRuangUjianById(context.Context, int) (ruangujian.RuangUjian, error) {
	return ruangujian.RuangUjian{}, nil
}
func (f *fakeDeleteRepo) GetRuangUjianByKode(context.Context, string) (ruangujian.RuangUjian, error) {
	return ruangujian.RuangUjian{}, nil
}
func (f *fakeDeleteRepo) ExistByKodeRuang(context.Context, string) (bool, error)        { return false, nil }
func (f *fakeDeleteRepo) CreateRuangUjian(context.Context, ruangujian.RuangUjian) error { return nil }
func (f *fakeDeleteRepo) UpdateRuangUjian(context.Context, int, updatepatch.UpdateRuangUjianPatch) error {
	return nil
}
func (f *fakeDeleteRepo) DeleteRuangUjian(ctx context.Context, id int) error {
	if f.deleteFn != nil {
		return f.deleteFn(ctx, id)
	}
	return nil
}

func TestDeleteRuangUjianService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")

	tests := []struct {
		name      string
		idRuangan int
		repo      *fakeDeleteRepo
		expectErr error
	}{
		{name: "branch 1 -> id tidak valid", idRuangan: 0, repo: &fakeDeleteRepo{}, expectErr: coreerror.ErrMissingId},
		{
			name:      "branch 2 -> delete restricted",
			idRuangan: 10,
			repo:      &fakeDeleteRepo{deleteFn: func(context.Context, int) error { return coreerror.ErrDeleteRestricted }},
			expectErr: coreerror.ErrDeleteRestricted,
		},
		{
			name:      "branch 3 -> repo error lain",
			idRuangan: 10,
			repo:      &fakeDeleteRepo{deleteFn: func(context.Context, int) error { return repoErr }},
			expectErr: repoErr,
		},
		{
			name:      "loop coverage -> delete berhasil",
			idRuangan: 10,
			repo:      &fakeDeleteRepo{deleteFn: func(context.Context, int) error { return nil }},
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			svc := ruangujian_service.NewDeleteRuangUjianService(tc.repo)
			err := svc.DeleteRuangUjian(ctx, tc.idRuangan)
			assert.ErrorIs(t, err, tc.expectErr)
		})
	}
}
