package user_service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/reset_password"
	"github.com/stretchr/testify/assert"
)

type fakeResetPasswordRepo struct {
	resetErr      error
	called        bool
	receivedID    user.ID
	receivedPlain string
}

func (f *fakeResetPasswordRepo) ResetPassword(ctx context.Context, idPengguna user.ID, password string) error {
	f.called = true
	f.receivedID = idPengguna
	f.receivedPlain = password
	if f.resetErr != nil {
		return f.resetErr
	}
	return nil
}

type fakePasswordHasher struct {
	hash    string
	hashErr error
	called  bool
}

func (f *fakePasswordHasher) GenerateHash(plain string) (string, error) {
	f.called = true
	if f.hashErr != nil {
		return "", f.hashErr
	}
	if f.hash != "" {
		return f.hash, nil
	}
	return "hashed-" + plain, nil
}

func (f *fakePasswordHasher) ComparePaswordAndHashed(hash string, plain string) bool {
	return true
}

func TestResetPasswordService_ResetPasswordService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		hasherErr        error
		repoErr          error
		wantErr          error
		wantRepoHit      bool
		wantHasherCalled bool
	}{
		{
			name:             "Path 1 -> hashing failed",
			hasherErr:        errors.New("hash failed"),
			wantErr:          errors.New("hash failed"),
			wantRepoHit:      false,
			wantHasherCalled: true,
		},
		{
			name:             "Path 2 -> repo reset password gagal",
			repoErr:          errors.New("repo failed"),
			wantErr:          errors.New("repo failed"),
			wantRepoHit:      true,
			wantHasherCalled: true,
		},
		{
			name:             "Path 3 -> reset password berhasil",
			wantRepoHit:      true,
			wantHasherCalled: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &fakeResetPasswordRepo{resetErr: tt.repoErr}
			hasher := &fakePasswordHasher{hash: "hashed-pass", hashErr: tt.hasherErr}
			svc := user_service.NewResetPasswordService(repo, hasher)

			err := svc.ResetPasswordService(context.Background(), user.ID(10), "new-pass")

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantRepoHit, repo.called)
			assert.Equal(t, tt.wantHasherCalled, hasher.called)
			if tt.wantRepoHit {
				assert.Equal(t, user.ID(10), repo.receivedID)
				assert.Equal(t, "hashed-pass", repo.receivedPlain)
			}
		})
	}
}
