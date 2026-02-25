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
	getByIDFn     func(context.Context, pengumuman.ID) (pengumuman.Pengumuman, error)
	getByIDCalled bool
	gotGetByID    pengumuman.ID
	deleteFn      func(context.Context, pengumuman.ID) error
	deleteCalled  bool
	gotID         pengumuman.ID
}

type fakeDeleteFileRepo struct {
	deleteFn     func(context.Context, string) error
	deleteCalled bool
	gotPath      string
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

func (f *fakeDeletePengumumanRepo) GetPengumumanById(ctx context.Context, id pengumuman.ID) (pengumuman.Pengumuman, error) {
	f.getByIDCalled = true
	f.gotGetByID = id
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return pengumuman.Pengumuman{DokumenPengumuman: "dokumen-lama.pdf"}, nil
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

func (f *fakeDeleteFileRepo) DeleteFile(ctx context.Context, filePath string) error {
	f.deleteCalled = true
	f.gotPath = filePath
	if f.deleteFn != nil {
		return f.deleteFn(ctx, filePath)
	}
	return nil
}

func TestDeletePengumumanService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	deleteFileErr := errors.New("delete file error")

	tests := []struct {
		name                 string
		idPengumuman         pengumuman.ID
		repo                 *fakeDeletePengumumanRepo
		deleteFile           *fakeDeleteFileRepo
		expectErr            error
		wantGetByIDCalled    bool
		wantDeleteFileCalled bool
		wantDeleteRepoCalled bool
	}{
		{
			name:                 "branch 1 -> id pengumuman tidak valid",
			idPengumuman:         0,
			repo:                 &fakeDeletePengumumanRepo{},
			deleteFile:           &fakeDeleteFileRepo{},
			expectErr:            coreerror.ErrMissingId,
			wantGetByIDCalled:    false,
			wantDeleteFileCalled: false,
			wantDeleteRepoCalled: false,
		},
		{
			name:         "branch 2 -> get by id gagal",
			idPengumuman: 10,
			repo: &fakeDeletePengumumanRepo{
				getByIDFn: func(context.Context, pengumuman.ID) (pengumuman.Pengumuman, error) {
					return pengumuman.Pengumuman{}, repoErr
				},
			},
			deleteFile:           &fakeDeleteFileRepo{},
			expectErr:            repoErr,
			wantGetByIDCalled:    true,
			wantDeleteFileCalled: false,
			wantDeleteRepoCalled: false,
		},
		{
			name:         "branch 3 -> delete file gagal",
			idPengumuman: 10,
			repo:         &fakeDeletePengumumanRepo{},
			deleteFile: &fakeDeleteFileRepo{
				deleteFn: func(context.Context, string) error { return deleteFileErr },
			},
			expectErr:            deleteFileErr,
			wantGetByIDCalled:    true,
			wantDeleteFileCalled: true,
			wantDeleteRepoCalled: false,
		},
		{
			name:         "branch 4 -> delete restricted",
			idPengumuman: 10,
			repo: &fakeDeletePengumumanRepo{
				deleteFn: func(context.Context, pengumuman.ID) error { return coreerror.ErrDeleteRestricted },
			},
			deleteFile:           &fakeDeleteFileRepo{},
			expectErr:            coreerror.ErrDeleteRestricted,
			wantGetByIDCalled:    true,
			wantDeleteFileCalled: true,
			wantDeleteRepoCalled: true,
		},
		{
			name:         "branch 5 -> repo error lain",
			idPengumuman: 10,
			repo: &fakeDeletePengumumanRepo{
				deleteFn: func(context.Context, pengumuman.ID) error { return repoErr },
			},
			deleteFile:           &fakeDeleteFileRepo{},
			expectErr:            repoErr,
			wantGetByIDCalled:    true,
			wantDeleteFileCalled: true,
			wantDeleteRepoCalled: true,
		},
		{
			name:                 "happy path -> delete pengumuman berhasil",
			idPengumuman:         11,
			repo:                 &fakeDeletePengumumanRepo{},
			deleteFile:           &fakeDeleteFileRepo{},
			expectErr:            nil,
			wantGetByIDCalled:    true,
			wantDeleteFileCalled: true,
			wantDeleteRepoCalled: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := pengumuman_service.NewDeletePengumumanService(tt.repo, tt.deleteFile)
			err := svc.DeletePengumumanService(ctx, tt.idPengumuman)

			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantGetByIDCalled, tt.repo.getByIDCalled)
			if tt.wantGetByIDCalled {
				assert.Equal(t, tt.idPengumuman, tt.repo.gotGetByID)
			}

			assert.Equal(t, tt.wantDeleteFileCalled, tt.deleteFile.deleteCalled)
			if tt.wantDeleteFileCalled {
				assert.Equal(t, "dokumen-lama.pdf", tt.deleteFile.gotPath)
			}

			assert.Equal(t, tt.wantDeleteRepoCalled, tt.repo.deleteCalled)
			if tt.wantDeleteRepoCalled {
				assert.Equal(t, tt.idPengumuman, tt.repo.gotID)
			}
		})
	}
}
