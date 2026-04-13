package auth_service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	auth_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/auth_service"
	"github.com/stretchr/testify/assert"
)

func TestAuthServiceLoginBasisPath(t *testing.T) {
	t.Parallel()

	const (
		username = "admin"
		password = "rahasia"
	)

	cmd := auth_service.LoginCmd{
		Username: username,
		Password: password,
	}

	baseUser := user.Pengguna{
		ID:             1,
		Username:       username,
		PasswordHashed: "hashed-password",
		Role:           user.ADMIN,
		StatusAkun:     user.AKTIF,
	}
	nonAdminUser := baseUser
	nonAdminUser.Role = user.GURU

	testErr := errors.New("test error")

	newDeps := func() (
		*FakeAuthUserRepo,
		*FakeUserRepo,
		*FakeHasher,
		*FakeSessionRepo,
		*FakeAccessToken,
		*FakeRefreshToken,
	) {
		return &FakeAuthUserRepo{
				ByUsername: map[string]user.Pengguna{},
			},
			&FakeUserRepo{},
			&FakeHasher{Ok: true},
			&FakeSessionRepo{Store: map[string]session.Session{}, SessionNextID: "session_user_1"},
			&FakeAccessToken{},
			&FakeRefreshToken{}
	}

	testCases := []struct {
		name      string
		cmd       auth_service.LoginCmd
		setup     func(*FakeAuthUserRepo, *FakeHasher, *FakeSessionRepo, *FakeAccessToken, *FakeRefreshToken)
		wantErr   error
		assertion func(*testing.T, auth_service.LoginRes, *FakeAuthUserRepo, *FakeHasher, *FakeSessionRepo, *FakeAccessToken, *FakeRefreshToken)
	}{
		{
			name:    "Path 1 -> Panjang username tidak valid",
			cmd:     auth_service.LoginCmd{Username: "adm", Password: password},
			wantErr: coreerror.ErrUsernameLengthInvalid,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.False(t, authRepo.FindCalled)
				assert.False(t, hasher.CompareCalled)
				assert.False(t, sessionRepo.RevokeExpiredCalled)
			},
		},
		{
			name: "Path 2 -> Username gagal ditemukan",
			setup: func(authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				authRepo.ByUsername = map[string]user.Pengguna{}
			},
			wantErr: coreerror.ErrInvalidCreds,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.True(t, authRepo.FindCalled)
				assert.False(t, hasher.CompareCalled)
				assert.False(t, sessionRepo.RevokeExpiredCalled)
			},
		},
		{
			name: "Path 3 -> Status akun tidak aktif",
			setup: func(authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				inactiveUser := baseUser
				inactiveUser.StatusAkun = user.NONAKTIF
				authRepo.ByUsername[username] = inactiveUser
			},
			wantErr: coreerror.ErrInvalidCreds,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.True(t, authRepo.FindCalled)
				assert.False(t, hasher.CompareCalled)
				assert.False(t, sessionRepo.RevokeExpiredCalled)
			},
		},
		{
			name: "Path 4 -> Password tidak cocok",
			setup: func(authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				authRepo.ByUsername[username] = baseUser
				hasher.Ok = false
			},
			wantErr: coreerror.ErrInvalidCreds,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.True(t, hasher.CompareCalled)
				assert.False(t, sessionRepo.RevokeExpiredCalled)
			},
		},
		{
			name: "Path 5 -> Revoke expired session gagal",
			setup: func(authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				authRepo.ByUsername[username] = baseUser
				sessionRepo.RevokeExpiredErr = testErr
			},
			wantErr: testErr,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.True(t, sessionRepo.RevokeExpiredCalled)
				assert.False(t, sessionRepo.HasActiveCalled)
			},
		},
		{
			name: "Path 6 -> Cek active session gagal",
			setup: func(authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				authRepo.ByUsername[username] = nonAdminUser
				sessionRepo.HasActiveErr = testErr
			},
			wantErr: testErr,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.True(t, sessionRepo.RevokeExpiredCalled)
				assert.True(t, sessionRepo.HasActiveCalled)
				assert.False(t, sessionRepo.CreateSessionCalled)
			},
		},
		{
			name: "Path 7 -> User masih memiliki session aktif",
			setup: func(authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				authRepo.ByUsername[username] = nonAdminUser
				sessionRepo.Store["existing_session"] = session.Session{
					SessionID: "existing_session",
					UserID:    nonAdminUser.ID,
					Role:      nonAdminUser.Role,
					Revoked:   false,
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}
			},
			wantErr: coreerror.ErrHasSession,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.True(t, sessionRepo.HasActiveCalled)
				assert.False(t, sessionRepo.CreateSessionCalled)
			},
		},
		{
			name: "Path 8 -> Admin boleh login walau masih memiliki session aktif",
			setup: func(authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				authRepo.ByUsername[username] = baseUser
				sessionRepo.Store["existing_session"] = session.Session{
					SessionID: "existing_session",
					UserID:    baseUser.ID,
					Role:      baseUser.Role,
					Revoked:   false,
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}
			},
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				assert.Equal(t, baseUser.ID, got.IdPengguna)
				assert.Equal(t, baseUser.Username, got.Username)
				assert.Equal(t, "ACCESS TOKEN :|1|ADMIN|admin", got.AccessToken)
				assert.Equal(t, "REFRESH TOKEN : session_user_1", got.RefreshToken)
				assert.False(t, sessionRepo.HasActiveCalled)
				assert.True(t, sessionRepo.CreateSessionCalled)
				assert.True(t, accessToken.GenerateCalled)
				assert.True(t, refreshToken.GenerateCalled)
				assert.Len(t, sessionRepo.Store, 2)
			},
		},
		{
			name: "Path 9 -> Create session gagal",
			setup: func(authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				authRepo.ByUsername[username] = baseUser
				sessionRepo.CreateSessionErr = testErr
			},
			wantErr: testErr,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.True(t, sessionRepo.CreateSessionCalled)
				assert.False(t, accessToken.GenerateCalled)
			},
		},
		{
			name: "Path 10 -> Generate access token gagal",
			setup: func(authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				authRepo.ByUsername[username] = baseUser
				accessToken.GenerateAccessTokenErr = testErr
			},
			wantErr: testErr,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.True(t, sessionRepo.CreateSessionCalled)
				assert.True(t, accessToken.GenerateCalled)
				assert.False(t, refreshToken.GenerateCalled)
				_, exists := sessionRepo.Store["session_user_1"]
				assert.True(t, exists)
			},
		},
		{
			name: "Path 11 -> Generate refresh token gagal",
			setup: func(authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				authRepo.ByUsername[username] = baseUser
				refreshToken.GenerateRefreshTokenErr = testErr
			},
			wantErr: testErr,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.True(t, sessionRepo.CreateSessionCalled)
				assert.True(t, accessToken.GenerateCalled)
				assert.True(t, refreshToken.GenerateCalled)
				_, exists := sessionRepo.Store["session_user_1"]
				assert.True(t, exists)
			},
		},
		{
			name: "Path 12 -> Login berhasil",
			setup: func(authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				authRepo.ByUsername[username] = baseUser
			},
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *FakeAuthUserRepo, hasher *FakeHasher, sessionRepo *FakeSessionRepo, accessToken *FakeAccessToken, refreshToken *FakeRefreshToken) {
				assert.Equal(t, baseUser.ID, got.IdPengguna)
				assert.Equal(t, baseUser.Username, got.Username)
				assert.Equal(t, "ACCESS TOKEN :|1|ADMIN|admin", got.AccessToken)
				assert.Equal(t, "REFRESH TOKEN : session_user_1", got.RefreshToken)
				assert.True(t, sessionRepo.CreateSessionCalled)
				assert.True(t, accessToken.GenerateCalled)
				assert.True(t, refreshToken.GenerateCalled)

				sess, exists := sessionRepo.Store["session_user_1"]
				assert.True(t, exists)
				assert.Equal(t, baseUser.ID, sess.UserID)
				assert.Equal(t, baseUser.Role, sess.Role)
				assert.False(t, sess.Revoked)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			authRepo, userRepo, hasher, sessionRepo, accessToken, refreshToken := newDeps()
			if tc.setup != nil {
				tc.setup(authRepo, hasher, sessionRepo, accessToken, refreshToken)
			}

			service := auth_service.NewAuthService(
				authRepo,
				userRepo,
				hasher,
				sessionRepo,
				accessToken,
				refreshToken,
				14*24*time.Hour,
			)

			inputCmd := cmd
			if tc.cmd != (auth_service.LoginCmd{}) {
				inputCmd = tc.cmd
			}

			got, err := service.Login(context.Background(), inputCmd)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}

			tc.assertion(t, got, authRepo, hasher, sessionRepo, accessToken, refreshToken)
		})
	}
}
