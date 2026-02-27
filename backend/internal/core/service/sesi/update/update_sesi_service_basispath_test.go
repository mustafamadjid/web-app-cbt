package sesi_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	sesi_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/update"
	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/update/fake_test"
	"github.com/stretchr/testify/assert"
)

func ptrString(s string) *string { return &s }

func TestUpdateSesiService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("update sesi error")

	tests := []struct {
		name       string
		idSesi     int
		patch      updatepatch.UpdateSesiPatch
		repo       *fake_test.FakeSesiRepo
		wantErr    error
		wantUpdate bool
	}{
		{
			name:       "Path 1 -> idSesi == 0",
			idSesi:     0,
			patch:      updatepatch.UpdateSesiPatch{},
			repo:       &fake_test.FakeSesiRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantUpdate: false,
		},
		{
			name:       "Path 2 -> NamaSesi not nil tapi kosong",
			idSesi:     1,
			patch:      updatepatch.UpdateSesiPatch{NamaSesi: ptrString("  ")},
			repo:       &fake_test.FakeSesiRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantUpdate: false,
		},
		{
			name:       "Path 3 -> KodeSesi not nil tapi kosong",
			idSesi:     1,
			patch:      updatepatch.UpdateSesiPatch{KodeSesi: ptrString("  ")},
			repo:       &fake_test.FakeSesiRepo{},
			wantErr:    coreerror.ErrMissingField,
			wantUpdate: false,
		},
		{
			name:       "Path 4 -> UpdateSesi error",
			idSesi:     1,
			patch:      updatepatch.UpdateSesiPatch{NamaSesi: ptrString("Sesi 1")},
			repo:       &fake_test.FakeSesiRepo{UpdateSesiErr: repoErr},
			wantErr:    repoErr,
			wantUpdate: true,
		},
		{
			name:       "Path 5 -> happy path berhasil update",
			idSesi:     1,
			patch:      updatepatch.UpdateSesiPatch{NamaSesi: ptrString("  Sesi 1  "), KodeSesi: ptrString("  sesi01  ")},
			repo:       &fake_test.FakeSesiRepo{},
			wantErr:    nil,
			wantUpdate: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := sesi_service.NewUpdateSesiService(tt.repo)
			err := svc.UpdateSesiService(ctx, tt.idSesi, tt.patch)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantUpdate, tt.repo.UpdateSesiCalled)
		})
	}
}
