package matapelajaran_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	mapel_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/update"
	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/update/fake_test"
	"github.com/stretchr/testify/assert"
)

func ptrString(s string) *string                { return &s }
func ptrMapelID(id matapelajaran.ID) *matapelajaran.ID { return &id }

func TestUpdateMapelService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	existErr := errors.New("exist kode error")
	genericErr := errors.New("generic update error")

	tests := []struct {
		name       string
		idMapel    int
		patch      updatepatch.UpdateMapelPatch
		repo       *fake_test.FakeMapelRepo
		wantErr    error
		wantUpdate bool
	}{
		{
			name:       "Path 1 -> idMapel <= 0",
			idMapel:    0,
			patch:      updatepatch.UpdateMapelPatch{},
			repo:       &fake_test.FakeMapelRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantUpdate: false,
		},
		{
			name:       "Path 2 -> IdKelas not nil tapi <= 0",
			idMapel:    1,
			patch:      updatepatch.UpdateMapelPatch{IdKelas: ptrMapelID(0)},
			repo:       &fake_test.FakeMapelRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantUpdate: false,
		},
		{
			name:       "Path 3 -> KodeMapel not nil tapi kosong",
			idMapel:    1,
			patch:      updatepatch.UpdateMapelPatch{KodeMapel: ptrString("  ")},
			repo:       &fake_test.FakeMapelRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantUpdate: false,
		},
		{
			name:       "Path 4 -> ExistKodeMapel error",
			idMapel:    1,
			patch:      updatepatch.UpdateMapelPatch{KodeMapel: ptrString("MTK01")},
			repo:       &fake_test.FakeMapelRepo{ExistKodeMapelErr: existErr},
			wantErr:    existErr,
			wantUpdate: false,
		},
		{
			name:       "Path 5 -> KodeMapel sudah exist",
			idMapel:    1,
			patch:      updatepatch.UpdateMapelPatch{KodeMapel: ptrString("MTK01")},
			repo:       &fake_test.FakeMapelRepo{ExistKodeMapelRet: true},
			wantErr:    coreerror.ErrKodeMapelExist,
			wantUpdate: false,
		},
		{
			name:       "Path 6 -> NamaMapel not nil tapi kosong",
			idMapel:    1,
			patch:      updatepatch.UpdateMapelPatch{NamaMapel: ptrString("  ")},
			repo:       &fake_test.FakeMapelRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantUpdate: false,
		},
		{
			name:       "Path 7 -> Deskripsi not nil tapi kosong",
			idMapel:    1,
			patch:      updatepatch.UpdateMapelPatch{Deskripsi: ptrString("  ")},
			repo:       &fake_test.FakeMapelRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantUpdate: false,
		},
		{
			name:       "Path 8 -> UpdateMapel err ErrNotFound",
			idMapel:    1,
			patch:      updatepatch.UpdateMapelPatch{NamaMapel: ptrString("Matematika")},
			repo:       &fake_test.FakeMapelRepo{UpdateMapelErr: coreerror.ErrNotFound},
			wantErr:    coreerror.ErrNotFound,
			wantUpdate: true,
		},
		{
			name:       "Path 9 -> UpdateMapel generic error",
			idMapel:    1,
			patch:      updatepatch.UpdateMapelPatch{NamaMapel: ptrString("Matematika")},
			repo:       &fake_test.FakeMapelRepo{UpdateMapelErr: genericErr},
			wantErr:    genericErr,
			wantUpdate: true,
		},
		{
			name:       "Path 10 -> happy path berhasil update",
			idMapel:    1,
			patch:      updatepatch.UpdateMapelPatch{NamaMapel: ptrString("  Matematika  "), Deskripsi: ptrString("  Desc  ")},
			repo:       &fake_test.FakeMapelRepo{},
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
