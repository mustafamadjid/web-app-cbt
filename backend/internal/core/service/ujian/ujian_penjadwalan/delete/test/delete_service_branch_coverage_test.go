package ujian_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/delete"
	"github.com/stretchr/testify/assert"
)

type fakeDeleteUjianRepo struct {
	deleteErr    error
	deleteCalled bool
	gotDeleteID  ujian.ID
}

func (*fakeDeleteUjianRepo) GetIdUjianByAttempt(context.Context, ujian.ID) (ujian.ID, error) {
	return 0, nil
}

func (*fakeDeleteUjianRepo) CreateUjian(context.Context, ujian.PenjadwalanUjian) error {
	return nil
}

func (*fakeDeleteUjianRepo) UpdateUjian(context.Context, ujian.ID, updatepatch.UpdatePenjadwalanUjian) error {
	return nil
}

func (f *fakeDeleteUjianRepo) DeleteUjian(_ context.Context, id ujian.ID) error {
	f.deleteCalled = true
	f.gotDeleteID = id
	return f.deleteErr
}

func TestDeleteUjianService_BasisPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repoErr := errors.New("repo error")

	tests := []struct {
		name       string
		idUjian    ujian.ID
		repo       *fakeDeleteUjianRepo
		wantErr    error
		wantCalled bool
	}{
		{
			name:       "Path 1 -> id ujian tidak valid",
			idUjian:    0,
			repo:       &fakeDeleteUjianRepo{},
			wantErr:    coreerror.ErrMissingId,
			wantCalled: false,
		},
		{
			name:       "Path 2 -> repo delete ujian gagal",
			idUjian:    15,
			repo:       &fakeDeleteUjianRepo{deleteErr: repoErr},
			wantErr:    repoErr,
			wantCalled: true,
		},
		{
			name:       "Path 3 -> berhasil delete ujian",
			idUjian:    15,
			repo:       &fakeDeleteUjianRepo{},
			wantCalled: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := ujian_service.NewDeleteUjianService(tc.repo)
			err := svc.DeleteUjianService(ctx, tc.idUjian)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantCalled, tc.repo.deleteCalled)
			if tc.wantCalled {
				assert.Equal(t, tc.idUjian, tc.repo.gotDeleteID)
			}
		})
	}
}
