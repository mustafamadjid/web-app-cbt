package user_service_test

import (
	"context"
	"errors"
	"testing"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/delete"
	"github.com/stretchr/testify/assert"
)

type fakeDeleteFileRepo struct {
	DeleteErr error

	DeleteCalled bool
	LastPath     string
}

func (f *fakeDeleteFileRepo) DeleteFile(ctx context.Context, filePath string) error {
	f.DeleteCalled = true
	f.LastPath = filePath
	if f.DeleteErr != nil {
		return f.DeleteErr
	}
	return nil
}

func TestDeleteUserService_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := user.ID(123)

	tests := []struct {
		name            string
		id              user.ID
		repo            *FakeDeleteUserRepo
		deleteFile      *fakeDeleteFileRepo
		wantErr         error
		wantFindCalled  bool
		wantDeleteFile  bool
		wantDeleteRepo  bool
		wantDeletePath  string
		wantFindID      user.ID
		wantDeleteID    user.ID
	}{
		{
			name:           "Path 1 -> id tidak valid mengembalikan ErrMissingId",
			id:             0,
			repo:           &FakeDeleteUserRepo{},
			deleteFile:     &fakeDeleteFileRepo{},
			wantErr:        coreerror.ErrMissingId,
			wantFindCalled: false,
			wantDeleteFile: false,
			wantDeleteRepo: false,
		},
		{
			name:           "Path 2 -> find user gagal",
			id:             userID,
			repo:           &FakeDeleteUserRepo{FindErr: errors.New("find failed")},
			deleteFile:     &fakeDeleteFileRepo{},
			wantErr:        errors.New("find failed"),
			wantFindCalled: true,
			wantDeleteFile: false,
			wantDeleteRepo: false,
			wantFindID:     userID,
		},
		{
			name: "Path 3 -> hapus file foto gagal tapi delete user tetap berjalan",
			id:   userID,
			repo: &FakeDeleteUserRepo{
				FindResult: user.Pengguna{ID: userID, Foto: "foto-lama.png"},
			},
			deleteFile:     &fakeDeleteFileRepo{DeleteErr: errors.New("delete file failed")},
			wantFindCalled: true,
			wantDeleteFile: true,
			wantDeleteRepo: true,
			wantDeletePath: "foto-lama.png",
			wantFindID:     userID,
			wantDeleteID:   userID,
		},
		{
			name: "Path 4 -> repository delete gagal",
			id:   userID,
			repo: &FakeDeleteUserRepo{
				DeleteErr:  errors.New("delete failed"),
				FindResult: user.Pengguna{ID: userID, Foto: "foto-lama.png"},
			},
			deleteFile:     &fakeDeleteFileRepo{},
			wantErr:        errors.New("delete failed"),
			wantFindCalled: true,
			wantDeleteFile: true,
			wantDeleteRepo: true,
			wantDeletePath: "foto-lama.png",
			wantFindID:     userID,
			wantDeleteID:   userID,
		},
		{
			name: "Path 5 -> delete user berhasil",
			id:   userID,
			repo: &FakeDeleteUserRepo{
				FindResult: user.Pengguna{ID: userID, Foto: "foto-lama.png"},
			},
			deleteFile:     &fakeDeleteFileRepo{},
			wantFindCalled: true,
			wantDeleteFile: true,
			wantDeleteRepo: true,
			wantDeletePath: "foto-lama.png",
			wantFindID:     userID,
			wantDeleteID:   userID,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			service := user_service.NewDeleteUserService(tc.repo, tc.deleteFile)

			err := service.Delete(ctx, tc.id)

			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantFindCalled, tc.repo.FindCalled)
			assert.Equal(t, tc.wantDeleteFile, tc.deleteFile.DeleteCalled)
			assert.Equal(t, tc.wantDeleteRepo, tc.repo.DeleteCalled)
			if tc.wantFindCalled {
				assert.Equal(t, tc.wantFindID, tc.repo.LastFindID)
			}
			if tc.wantDeleteFile {
				assert.Equal(t, tc.wantDeletePath, tc.deleteFile.LastPath)
			}
			if tc.wantDeleteRepo {
				assert.Equal(t, tc.wantDeleteID, tc.repo.LastID)
			}
		})
	}
}

func TestDeleteUserService_DeleteMany(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userIDs := []user.ID{101, 202, 303}

	tests := []struct {
		name            string
		repo            *FakeDeleteUserRepo
		wantAffected    int64
		wantErr         error
		wantDeleteUsers bool
	}{
		{
			name:            "Path 1 -> repository delete many gagal",
			repo:            &FakeDeleteUserRepo{DeleteErr: errors.New("delete many failed")},
			wantAffected:    0,
			wantErr:         errors.New("delete many failed"),
			wantDeleteUsers: true,
		},
		{
			name:            "Path 2 -> delete many berhasil",
			repo:            &FakeDeleteUserRepo{},
			wantAffected:    int64(len(userIDs)),
			wantDeleteUsers: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deleteFile := &fakeDeleteFileRepo{}
			service := user_service.NewDeleteUserService(tc.repo, deleteFile)

			affected, err := service.DeleteMany(ctx, userIDs)

			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.EqualValues(t, tc.wantAffected, affected)
			assert.Equal(t, tc.wantDeleteUsers, tc.repo.DeleteUsersCalled)
			assert.Equal(t, userIDs, tc.repo.LastIDs)
			assert.False(t, deleteFile.DeleteCalled)
		})
	}
}
