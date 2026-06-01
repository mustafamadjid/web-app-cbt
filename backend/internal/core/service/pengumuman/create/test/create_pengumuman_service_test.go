package pengumuman_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	pengumuman_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/create"
)

type fakeCreatePengumumanRepo struct {
	createFn     func(context.Context, pengumuman.Pengumuman) error
	createCalled bool
	gotCreate    pengumuman.Pengumuman
}

func (f *fakeCreatePengumumanRepo) GetPengumumanActive(context.Context) ([]pengumuman.Pengumuman, error) {
	return nil, nil
}

func (f *fakeCreatePengumumanRepo) GetPengumumanNonActive(context.Context) ([]pengumuman.Pengumuman, error) {
	return nil, nil
}

func (f *fakeCreatePengumumanRepo) GetPengumumanIncoming(context.Context) ([]pengumuman.Pengumuman, error) {
	return nil, nil
}

func (f *fakeCreatePengumumanRepo) GetPengumumanById(context.Context, pengumuman.ID) (pengumuman.Pengumuman, error) {
	return pengumuman.Pengumuman{}, nil
}

func (f *fakeCreatePengumumanRepo) CreatePengumuman(ctx context.Context, payload pengumuman.Pengumuman) error {
	f.createCalled = true
	f.gotCreate = payload
	if f.createFn != nil {
		return f.createFn(ctx, payload)
	}
	return nil
}

func (f *fakeCreatePengumumanRepo) UpdatePengumuman(context.Context, pengumuman.ID, updatepatch.PengumumanUpdatePatch) error {
	return nil
}

func (f *fakeCreatePengumumanRepo) DeletePengumuman(context.Context, pengumuman.ID) error {
	return nil
}

func TestCreatePengumumanService(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")
	trimmedPayload := pengumuman.Pengumuman{
		IdPengguna:               1,
		JudulPengumuman:          "Pengumuman Penting",
		IsiPengumuman:            "Isi pengumuman",
		TanggalRilisPengumuman:   "2026-01-01",
		TanggalSelesaiPengumuman: "2026-01-31",
		DokumenPengumuman:        "dokumen.pdf",
	}

	tests := []struct {
		name           string
		payload        pengumuman.Pengumuman
		repo           *fakeCreatePengumumanRepo
		wantErr        error
		wantRepoCalled bool
		wantPayload    *pengumuman.Pengumuman
	}{
		{
			name: "branch 1 -> id pengguna tidak valid",
			payload: pengumuman.Pengumuman{
				IdPengguna:               0,
				JudulPengumuman:          "judul",
				IsiPengumuman:            "isi",
				TanggalRilisPengumuman:   "2026-01-01",
				TanggalSelesaiPengumuman: "2026-01-31",
			},
			repo:           &fakeCreatePengumumanRepo{},
			wantErr:        coreerror.ErrMissingId,
			wantRepoCalled: false,
		},
		{
			name: "branch 2 -> format tanggal rilis tidak valid",
			payload: pengumuman.Pengumuman{
				IdPengguna:               1,
				JudulPengumuman:          "judul",
				IsiPengumuman:            "isi",
				TanggalRilisPengumuman:   "2026-15-01",
				TanggalSelesaiPengumuman: "2026-01-31",
			},
			repo:           &fakeCreatePengumumanRepo{},
			wantErr:        coreerror.ErrInvalidDateFormat,
			wantRepoCalled: false,
		},
		{
			name: "branch 3 -> format tanggal selesai tidak valid",
			payload: pengumuman.Pengumuman{
				IdPengguna:               1,
				JudulPengumuman:          "judul",
				IsiPengumuman:            "isi",
				TanggalRilisPengumuman:   "2026-01-01",
				TanggalSelesaiPengumuman: "2026-99-31",
			},
			repo:           &fakeCreatePengumumanRepo{},
			wantErr:        coreerror.ErrInvalidDateFormat,
			wantRepoCalled: false,
		},
		{
			name: "branch 4 -> repo gagal create",
			payload: pengumuman.Pengumuman{
				IdPengguna:               1,
				JudulPengumuman:          "  Pengumuman Penting  ",
				IsiPengumuman:            "  Isi pengumuman  ",
				TanggalRilisPengumuman:   " 2026-01-01 ",
				TanggalSelesaiPengumuman: " 2026-01-31 ",
				DokumenPengumuman:        " dokumen.pdf ",
			},
			repo:           &fakeCreatePengumumanRepo{createFn: func(context.Context, pengumuman.Pengumuman) error { return repoErr }},
			wantErr:        repoErr,
			wantRepoCalled: true,
			wantPayload:    &trimmedPayload,
		},
		{
			name: "branch 5 -> normalisasi payload dan berhasil",
			payload: pengumuman.Pengumuman{
				IdPengguna:               1,
				JudulPengumuman:          "  Pengumuman Penting  ",
				IsiPengumuman:            "  Isi pengumuman  ",
				TanggalRilisPengumuman:   " 2026-01-01 ",
				TanggalSelesaiPengumuman: " 2026-01-31 ",
				DokumenPengumuman:        " dokumen.pdf ",
			},
			repo:           &fakeCreatePengumumanRepo{},
			wantErr:        nil,
			wantRepoCalled: true,
			wantPayload:    &trimmedPayload,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := pengumuman_service.NewCreatePengumumanRepo(tt.repo)
			err := svc.CreatePengumuman(ctx, tt.payload)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantRepoCalled, tt.repo.createCalled)
			if tt.wantRepoCalled && tt.wantPayload != nil {
				assert.Equal(t, *tt.wantPayload, tt.repo.gotCreate)
			}
		})
	}
}
