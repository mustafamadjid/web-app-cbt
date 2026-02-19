package user_service

import (
	"context"
	"errors"
	"testing"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
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
}

func (f *fakePasswordHasher) GenerateHash(plain string) (string, error) {
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
	tests := []struct {
		name          string
		hasherErr     error
		repoErr       error
		expectErr     bool
		expectRepoHit bool
	}{
		{
			name:          "Branch 1 -> hashing failed",
			hasherErr:     errors.New("hash failed"),
			expectErr:     true,
			expectRepoHit: false,
		},
		{
			name:          "Branch 2 -> username taken",
			repoErr:       errors.New("repo failed"),
			expectErr:     true,
			expectRepoHit: true,
		},
		{
			name:          "Branch 3 -> success",
			expectErr:     false,
			expectRepoHit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeResetPasswordRepo{resetErr: tt.repoErr}
			hasher := &fakePasswordHasher{hash: "hashed-pass", hashErr: tt.hasherErr}
			svc := NewResetPasswordService(repo, hasher)

			err := svc.ResetPasswordService(context.Background(), user.ID(10), "new-pass")
			if tt.expectErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if repo.called != tt.expectRepoHit {
				t.Fatalf("unexpected repo call state, got %v want %v", repo.called, tt.expectRepoHit)
			}
			if tt.expectRepoHit {
				if repo.receivedID != user.ID(10) {
					t.Fatalf("unexpected user id, got %v", repo.receivedID)
				}
				if repo.receivedPlain != "hashed-pass" {
					t.Fatalf("unexpected hashed password, got %q", repo.receivedPlain)
				}
			}
		})
	}
}
