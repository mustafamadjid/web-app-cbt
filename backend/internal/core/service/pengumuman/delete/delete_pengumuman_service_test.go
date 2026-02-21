package pengumuman_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	pengumuman_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/delete"
)

type fakeDeletePengumumanRepo struct {
	deleteFn     func(context.Context, pengumuman.ID) error
	deleteCalled bool
	gotID        pengumuman.ID
}

func (f *fakeDeletePengumumanRepo) GetPengumumanActive(context.Context) ([]pengumuman.Pengumuman, error) {
	return nil, nil
}

func (f *fakeDeletePengumumanRepo) GetPengumumanNonActive(context.Context) ([]pengumuman.Pengumuman, error) {
	return nil, nil
}

func (f *fakeDeletePengumumanRepo) GetPengumumanIncoming(context.Context) ([]pengumuman.Pengumuman, error) {
	return nil, nil
}

func (f *fakeDeletePengumumanRepo) GetPengumumanById(context.Context, pengumuman.ID) (pengumuman.Pengumuman, error) {
	return pengumuman.Pengumuman{}, nil
}

func (f *fakeDeletePengumumanRepo) CreatePengumuman(context.Context, pengumuman.Pengumuman) error {
	return nil
}

func (f *fakeDeletePengumumanRepo) UpdatePengumuman(context.Context, pengumuman.ID, updatepatch.PengumumanUpdatePatch) error {
	return nil
}

func (f *fakeDeletePengumumanRepo) DeletePengumuman(ctx context.Context, idPengumuman pengumuman.ID) error {
	f.deleteCalled = true
	f.gotID = idPengumuman
	if f.deleteFn != nil {
		return f.deleteFn(ctx, idPengumuman)
	}
	return nil
}

func TestDeletePengumumanService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")

	tests := []struct {
		name           string
		idPengumuman   pengumuman.ID
		repo           *fakeDeletePengumumanRepo
		expectErr      error
		wantRepoCalled bool
	}{
		{
			name:           "branch 1 -> id pengumuman tidak valid",
			idPengumuman:   0,
			repo:           &fakeDeletePengumumanRepo{},
			expectErr:      coreerror.ErrMissingId,
			wantRepoCalled: false,
		},
		{
			name:           "branch 2 -> delete restricted",
			idPengumuman:   10,
			repo:           &fakeDeletePengumumanRepo{deleteFn: func(context.Context, pengumuman.ID) error { return coreerror.ErrDeleteRestricted }},
			expectErr:      coreerror.ErrDeleteRestricted,
			wantRepoCalled: true,
		},
		{
			name:           "branch 3 -> repo error lain",
			idPengumuman:   10,
			repo:           &fakeDeletePengumumanRepo{deleteFn: func(context.Context, pengumuman.ID) error { return repoErr }},
			expectErr:      repoErr,
			wantRepoCalled: true,
		},
		{
			name:           "happy path -> delete pengumuman berhasil",
			idPengumuman:   11,
			repo:           &fakeDeletePengumumanRepo{},
			expectErr:      nil,
			wantRepoCalled: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := pengumuman_service.NewDeletePengumumanService(tt.repo)
			err := svc.DeletePengumumanService(ctx, tt.idPengumuman)

			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantRepoCalled, tt.repo.deleteCalled)
			if tt.wantRepoCalled {
				assert.Equal(t, tt.idPengumuman, tt.repo.gotID)
			}
		})
	}
}
