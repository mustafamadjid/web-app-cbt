package pengumuman_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	pengumuman_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/update"
)

type fakeUpdatePengumumanRepo struct {
	getByIDFn     func(context.Context, pengumuman.ID) (pengumuman.Pengumuman, error)
	getByIDCalled bool
	gotGetByID    pengumuman.ID
	updateFn      func(context.Context, pengumuman.ID, updatepatch.PengumumanUpdatePatch) error
	updateCalled  bool
	gotID         pengumuman.ID
	gotPatch      updatepatch.PengumumanUpdatePatch
}

type fakeDeleteFileRepo struct {
	deleteFn     func(context.Context, string) error
	deleteCalled bool
	gotPath      string
}

func (f *fakeUpdatePengumumanRepo) GetPengumumanActive(context.Context) ([]pengumuman.Pengumuman, error) {
	return nil, nil
}

func (f *fakeUpdatePengumumanRepo) GetPengumumanNonActive(context.Context) ([]pengumuman.Pengumuman, error) {
	return nil, nil
}

func (f *fakeUpdatePengumumanRepo) GetPengumumanIncoming(context.Context) ([]pengumuman.Pengumuman, error) {
	return nil, nil
}

func (f *fakeUpdatePengumumanRepo) GetPengumumanById(ctx context.Context, id pengumuman.ID) (pengumuman.Pengumuman, error) {
	f.getByIDCalled = true
	f.gotGetByID = id
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return pengumuman.Pengumuman{DokumenPengumuman: "dokumen-lama.pdf"}, nil
}

func (f *fakeUpdatePengumumanRepo) CreatePengumuman(context.Context, pengumuman.Pengumuman) error {
	return nil
}

func (f *fakeUpdatePengumumanRepo) UpdatePengumuman(ctx context.Context, idPengumuman pengumuman.ID, patch updatepatch.PengumumanUpdatePatch) error {
	f.updateCalled = true
	f.gotID = idPengumuman
	f.gotPatch = patch
	if f.updateFn != nil {
		return f.updateFn(ctx, idPengumuman, patch)
	}
	return nil
}

func (f *fakeUpdatePengumumanRepo) DeletePengumuman(context.Context, pengumuman.ID) error {
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

func TestUpdatePengumumanService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	deleteFileErr := errors.New("delete file error")
	idPenggunaValid := pengumuman.ID(9)

	tests := []struct {
		name                 string
		idPengumuman         pengumuman.ID
		patch                updatepatch.PengumumanUpdatePatch
		repo                 *fakeUpdatePengumumanRepo
		deleteFile           *fakeDeleteFileRepo
		wantErr              error
		wantGetByIDCalled    bool
		wantDeleteFileCalled bool
		wantRepoCalled       bool
		wantPatch            updatepatch.PengumumanUpdatePatch
	}{
		{
			name:                 "branch 1 -> id pengumuman tidak valid",
			idPengumuman:         0,
			patch:                updatepatch.PengumumanUpdatePatch{},
			repo:                 &fakeUpdatePengumumanRepo{},
			deleteFile:           &fakeDeleteFileRepo{},
			wantErr:              coreerror.ErrMissingId,
			wantGetByIDCalled:    false,
			wantDeleteFileCalled: false,
			wantRepoCalled:       false,
		},
		{
			name:         "branch 2 -> id pengguna tidak valid",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna: ptrPengumumanID(0),
			},
			repo:                 &fakeUpdatePengumumanRepo{},
			deleteFile:           &fakeDeleteFileRepo{},
			wantErr:              coreerror.ErrMissingId,
			wantGetByIDCalled:    false,
			wantDeleteFileCalled: false,
			wantRepoCalled:       false,
		},
		{
			name:         "branch 3 -> judul pengumuman kosong setelah trim",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:      &idPenggunaValid,
				JudulPengumuman: ptrString("   "),
			},
			repo:                 &fakeUpdatePengumumanRepo{},
			deleteFile:           &fakeDeleteFileRepo{},
			wantErr:              coreerror.ErrMissingField,
			wantGetByIDCalled:    false,
			wantDeleteFileCalled: false,
			wantRepoCalled:       false,
		},
		{
			name:         "branch 4 -> isi pengumuman kosong setelah trim",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:    &idPenggunaValid,
				IsiPengumuman: ptrString("   "),
			},
			repo:                 &fakeUpdatePengumumanRepo{},
			deleteFile:           &fakeDeleteFileRepo{},
			wantErr:              coreerror.ErrMissingField,
			wantGetByIDCalled:    false,
			wantDeleteFileCalled: false,
			wantRepoCalled:       false,
		},
		{
			name:         "branch 5 -> tanggal rilis kosong setelah trim",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:             &idPenggunaValid,
				TanggalRilisPengumuman: ptrString("   "),
			},
			repo:                 &fakeUpdatePengumumanRepo{},
			deleteFile:           &fakeDeleteFileRepo{},
			wantErr:              coreerror.ErrMissingField,
			wantGetByIDCalled:    false,
			wantDeleteFileCalled: false,
			wantRepoCalled:       false,
		},
		{
			name:         "branch 6 -> tanggal rilis tidak valid",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:             &idPenggunaValid,
				TanggalRilisPengumuman: ptrString("2026-88-01"),
			},
			repo:                 &fakeUpdatePengumumanRepo{},
			deleteFile:           &fakeDeleteFileRepo{},
			wantErr:              coreerror.ErrInvalidDateFormat,
			wantGetByIDCalled:    false,
			wantDeleteFileCalled: false,
			wantRepoCalled:       false,
		},
		{
			name:         "branch 7 -> tanggal selesai kosong setelah trim",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:               &idPenggunaValid,
				TanggalSelesaiPengumuman: ptrString("   "),
			},
			repo:                 &fakeUpdatePengumumanRepo{},
			deleteFile:           &fakeDeleteFileRepo{},
			wantErr:              coreerror.ErrMissingField,
			wantGetByIDCalled:    false,
			wantDeleteFileCalled: false,
			wantRepoCalled:       false,
		},
		{
			name:         "branch 8 -> tanggal selesai tidak valid",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:               &idPenggunaValid,
				TanggalSelesaiPengumuman: ptrString("2026-99-31"),
			},
			repo:                 &fakeUpdatePengumumanRepo{},
			deleteFile:           &fakeDeleteFileRepo{},
			wantErr:              coreerror.ErrInvalidDateFormat,
			wantGetByIDCalled:    false,
			wantDeleteFileCalled: false,
			wantRepoCalled:       false,
		},
		{
			name:         "branch 9 -> dokumen pengumuman kosong setelah trim",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:        &idPenggunaValid,
				DokumenPengumuman: ptrString("   "),
			},
			repo:                 &fakeUpdatePengumumanRepo{},
			deleteFile:           &fakeDeleteFileRepo{},
			wantErr:              coreerror.ErrMissingField,
			wantGetByIDCalled:    false,
			wantDeleteFileCalled: false,
			wantRepoCalled:       false,
		},
		{
			name:         "branch 10 -> get by id gagal saat update dokumen",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:               &idPenggunaValid,
				JudulPengumuman:          ptrString("Judul"),
				IsiPengumuman:            ptrString("Isi"),
				TanggalRilisPengumuman:   ptrString("2026-01-01"),
				TanggalSelesaiPengumuman: ptrString("2026-01-31"),
				DokumenPengumuman:        ptrString("dokumen.pdf"),
			},
			repo: &fakeUpdatePengumumanRepo{
				getByIDFn: func(context.Context, pengumuman.ID) (pengumuman.Pengumuman, error) {
					return pengumuman.Pengumuman{}, repoErr
				},
			},
			deleteFile:           &fakeDeleteFileRepo{},
			wantErr:              repoErr,
			wantGetByIDCalled:    true,
			wantDeleteFileCalled: false,
			wantRepoCalled:       false,
		},
		{
			name:         "branch 11 -> delete file gagal saat update dokumen",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:               &idPenggunaValid,
				JudulPengumuman:          ptrString("Judul"),
				IsiPengumuman:            ptrString("Isi"),
				TanggalRilisPengumuman:   ptrString("2026-01-01"),
				TanggalSelesaiPengumuman: ptrString("2026-01-31"),
				DokumenPengumuman:        ptrString("dokumen.pdf"),
			},
			repo: &fakeUpdatePengumumanRepo{},
			deleteFile: &fakeDeleteFileRepo{
				deleteFn: func(context.Context, string) error { return deleteFileErr },
			},
			wantErr:              deleteFileErr,
			wantGetByIDCalled:    true,
			wantDeleteFileCalled: true,
			wantRepoCalled:       false,
		},
		{
			name:         "branch 12 -> repo gagal update",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:               &idPenggunaValid,
				JudulPengumuman:          ptrString("Judul"),
				IsiPengumuman:            ptrString("Isi"),
				TanggalRilisPengumuman:   ptrString("2026-01-01"),
				TanggalSelesaiPengumuman: ptrString("2026-01-31"),
				DokumenPengumuman:        ptrString("dokumen.pdf"),
			},
			repo: &fakeUpdatePengumumanRepo{
				updateFn: func(context.Context, pengumuman.ID, updatepatch.PengumumanUpdatePatch) error { return repoErr },
			},
			deleteFile:           &fakeDeleteFileRepo{},
			wantErr:              repoErr,
			wantGetByIDCalled:    true,
			wantDeleteFileCalled: true,
			wantRepoCalled:       true,
			wantPatch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:               &idPenggunaValid,
				JudulPengumuman:          ptrString("Judul"),
				IsiPengumuman:            ptrString("Isi"),
				TanggalRilisPengumuman:   ptrString("2026-01-01"),
				TanggalSelesaiPengumuman: ptrString("2026-01-31"),
				DokumenPengumuman:        ptrString("dokumen.pdf"),
			},
		},
		{
			name:         "happy path -> validasi seluruh field dan update berhasil",
			idPengumuman: 11,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:               &idPenggunaValid,
				JudulPengumuman:          ptrString("  Judul Baru  "),
				IsiPengumuman:            ptrString("  Isi Baru  "),
				TanggalRilisPengumuman:   ptrString(" 2026-02-01 "),
				TanggalSelesaiPengumuman: ptrString(" 2026-02-28 "),
				DokumenPengumuman:        ptrString("  dok-baru.pdf  "),
			},
			repo:                 &fakeUpdatePengumumanRepo{},
			deleteFile:           &fakeDeleteFileRepo{},
			wantErr:              nil,
			wantGetByIDCalled:    true,
			wantDeleteFileCalled: true,
			wantRepoCalled:       true,
			wantPatch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:               &idPenggunaValid,
				JudulPengumuman:          ptrString("Judul Baru"),
				IsiPengumuman:            ptrString("Isi Baru"),
				TanggalRilisPengumuman:   ptrString("2026-02-01"),
				TanggalSelesaiPengumuman: ptrString("2026-02-28"),
				DokumenPengumuman:        ptrString("dok-baru.pdf"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := pengumuman_service.NewUpdatePengumumanService(tt.repo, tt.deleteFile)
			err := svc.UpdatePengumumanService(ctx, tt.idPengumuman, tt.patch)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
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

			assert.Equal(t, tt.wantRepoCalled, tt.repo.updateCalled)
			if tt.wantRepoCalled {
				assert.Equal(t, tt.idPengumuman, tt.repo.gotID)
				assert.Equal(t, tt.wantPatch, tt.repo.gotPatch)
			}
		})
	}
}

func ptrString(s string) *string {
	return &s
}

func ptrPengumumanID(id pengumuman.ID) *pengumuman.ID {
	return &id
}
