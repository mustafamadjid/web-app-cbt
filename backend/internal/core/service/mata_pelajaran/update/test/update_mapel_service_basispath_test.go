package matapelajaran_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	mapel_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/update"
	"github.com/stretchr/testify/assert"
)

func ptrString(s string) *string                       { return &s }
func ptrMapelID(id matapelajaran.ID) *matapelajaran.ID { return &id }

func TestUpdateMapelService_BranchCoverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	existErr := errors.New("exist kode error")
	genericErr := errors.New("generic update error")

	tests := []struct {
		name       string
		idMapel    int
		patch      updatepatch.UpdateMapelPatch
		repo       *FakeMapelRepo
		wantErr    error
		wantUpdate bool
	}{
		{
			name:       "Branch 1 -> id mapel tidak valid",
			idMapel:    0,
			patch:      updatepatch.UpdateMapelPatch{},
			repo:       &FakeMapelRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantUpdate: false,
		},
		{
			name:       "Branch 2 -> id kelas tidak valid",
			idMapel:    1,
			patch:      updatepatch.UpdateMapelPatch{IdKelas: ptrMapelID(0)},
			repo:       &FakeMapelRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantUpdate: false,
		},
		{
			name:       "Branch 3 -> kode mapel kosong",
			idMapel:    1,
			patch:      updatepatch.UpdateMapelPatch{KodeMapel: ptrString("  ")},
			repo:       &FakeMapelRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantUpdate: false,
		},
		{
			name:       "Branch 4 -> cek kode mapel gagal",
			idMapel:    1,
			patch:      updatepatch.UpdateMapelPatch{KodeMapel: ptrString("MTK01")},
			repo:       &FakeMapelRepo{ExistKodeMapelErr: existErr},
			wantErr:    existErr,
			wantUpdate: false,
		},
		{
			name:       "Branch 5 -> kode mapel sudah ada",
			idMapel:    1,
			patch:      updatepatch.UpdateMapelPatch{KodeMapel: ptrString("MTK01")},
			repo:       &FakeMapelRepo{ExistKodeMapelRet: true},
			wantErr:    coreerror.ErrKodeMapelExist,
			wantUpdate: false,
		},
		{
			name:       "Branch 6 -> nama mapel kosong",
			idMapel:    1,
			patch:      updatepatch.UpdateMapelPatch{NamaMapel: ptrString("  ")},
			repo:       &FakeMapelRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantUpdate: false,
		},
		{
			name:       "Branch 7 -> deskripsi kosong",
			idMapel:    1,
			patch:      updatepatch.UpdateMapelPatch{Deskripsi: ptrString("  ")},
			repo:       &FakeMapelRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantUpdate: false,
		},
		{
			name:       "Branch 8 -> mapel yang diupdate tidak ditemukan",
			idMapel:    1,
			patch:      updatepatch.UpdateMapelPatch{NamaMapel: ptrString("Matematika")},
			repo:       &FakeMapelRepo{UpdateMapelErr: coreerror.ErrNotFound},
			wantErr:    coreerror.ErrNotFound,
			wantUpdate: true,
		},
		{
			name:       "Branch 9 -> repo update gagal",
			idMapel:    1,
			patch:      updatepatch.UpdateMapelPatch{NamaMapel: ptrString("Matematika")},
			repo:       &FakeMapelRepo{UpdateMapelErr: genericErr},
			wantErr:    genericErr,
			wantUpdate: true,
		},
		{
			name:       "Branch 10 -> update mapel berhasil",
			idMapel:    1,
			patch:      updatepatch.UpdateMapelPatch{NamaMapel: ptrString("  Matematika  "), Deskripsi: ptrString("  Desc  ")},
			repo:       &FakeMapelRepo{},
			wantErr:    nil,
			wantUpdate: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mapel_service.NewUpdateMapelService(tt.repo)
			err := svc.UpdateMapelService(ctx, tt.idMapel, tt.patch)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantUpdate, tt.repo.UpdateMapelCalled)
		})
	}
}
