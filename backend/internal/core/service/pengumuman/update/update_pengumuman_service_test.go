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
	updateFn     func(context.Context, pengumuman.ID, updatepatch.PengumumanUpdatePatch) error
	updateCalled bool
	gotID        pengumuman.ID
	gotPatch     updatepatch.PengumumanUpdatePatch
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

func (f *fakeUpdatePengumumanRepo) GetPengumumanById(context.Context, pengumuman.ID) (pengumuman.Pengumuman, error) {
	return pengumuman.Pengumuman{}, nil
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

func TestUpdatePengumumanService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	idPenggunaValid := pengumuman.ID(9)

	tests := []struct {
		name           string
		idPengumuman   pengumuman.ID
		patch          updatepatch.PengumumanUpdatePatch
		repo           *fakeUpdatePengumumanRepo
		wantErr        error
		wantRepoCalled bool
	}{
		{
			name:           "branch 1 -> id pengumuman tidak valid",
			idPengumuman:   0,
			patch:          updatepatch.PengumumanUpdatePatch{},
			repo:           &fakeUpdatePengumumanRepo{},
			wantErr:        coreerror.ErrMissingId,
			wantRepoCalled: false,
		},
		{
			name:         "branch 2 -> id pengguna tidak valid",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna: ptrPengumumanID(0),
			},
			repo:           &fakeUpdatePengumumanRepo{},
			wantErr:        coreerror.ErrMissingId,
			wantRepoCalled: false,
		},
		{
			name:         "branch 3 -> judul pengumuman kosong setelah trim",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:      &idPenggunaValid,
				JudulPengumuman: ptrString("   "),
			},
			repo:           &fakeUpdatePengumumanRepo{},
			wantErr:        coreerror.ErrMissingField,
			wantRepoCalled: false,
		},
		{
			name:         "branch 4 -> isi pengumuman kosong setelah trim",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:    &idPenggunaValid,
				IsiPengumuman: ptrString("   "),
			},
			repo:           &fakeUpdatePengumumanRepo{},
			wantErr:        coreerror.ErrMissingField,
			wantRepoCalled: false,
		},
		{
			name:         "branch 5 -> tanggal rilis kosong setelah trim",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:             &idPenggunaValid,
				TanggalRilisPengumuman: ptrString("   "),
			},
			repo:           &fakeUpdatePengumumanRepo{},
			wantErr:        coreerror.ErrMissingField,
			wantRepoCalled: false,
		},
		{
			name:         "branch 6 -> tanggal rilis tidak valid",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:             &idPenggunaValid,
				TanggalRilisPengumuman: ptrString("2026-88-01"),
			},
			repo:           &fakeUpdatePengumumanRepo{},
			wantErr:        coreerror.ErrInvalidDateFormat,
			wantRepoCalled: false,
		},
		{
			name:         "branch 7 -> tanggal selesai kosong setelah trim",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:               &idPenggunaValid,
				TanggalSelesaiPengumuman: ptrString("   "),
			},
			repo:           &fakeUpdatePengumumanRepo{},
			wantErr:        coreerror.ErrMissingField,
			wantRepoCalled: false,
		},
		{
			name:         "branch 8 -> tanggal selesai tidak valid",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:               &idPenggunaValid,
				TanggalSelesaiPengumuman: ptrString("2026-99-31"),
			},
			repo:           &fakeUpdatePengumumanRepo{},
			wantErr:        coreerror.ErrInvalidDateFormat,
			wantRepoCalled: false,
		},
		{
			name:         "branch 9 -> dokumen pengumuman kosong setelah trim",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:        &idPenggunaValid,
				DokumenPengumuman: ptrString("   "),
			},
			repo:           &fakeUpdatePengumumanRepo{},
			wantErr:        coreerror.ErrMissingField,
			wantRepoCalled: false,
		},
		{
			name:         "branch 10 -> repo gagal update",
			idPengumuman: 10,
			patch: updatepatch.PengumumanUpdatePatch{
				IdPengguna:               &idPenggunaValid,
				JudulPengumuman:          ptrString("Judul"),
				IsiPengumuman:            ptrString("Isi"),
				TanggalRilisPengumuman:   ptrString("2026-01-01"),
				TanggalSelesaiPengumuman: ptrString("2026-01-31"),
				DokumenPengumuman:        ptrString("dokumen.pdf"),
			},
			repo:           &fakeUpdatePengumumanRepo{updateFn: func(context.Context, pengumuman.ID, updatepatch.PengumumanUpdatePatch) error { return repoErr }},
			wantErr:        repoErr,
			wantRepoCalled: true,
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
			repo:           &fakeUpdatePengumumanRepo{},
			wantErr:        nil,
			wantRepoCalled: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := pengumuman_service.NewUpdatePengumumanService(tt.repo)
			err := svc.UpdatePengumumanService(ctx, tt.idPengumuman, tt.patch)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantRepoCalled, tt.repo.updateCalled)
			if tt.wantRepoCalled {
				assert.Equal(t, tt.idPengumuman, tt.repo.gotID)
				assert.Equal(t, tt.patch, tt.repo.gotPatch)
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
