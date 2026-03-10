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
	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/auth_service/fake_test"
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

	testErr := errors.New("test error")

	newDeps := func() (
		*fake_test.FakeAuthUserRepo,
		*fake_test.FakeUserRepo,
		*fake_test.FakeHasher,
		*fake_test.FakeSessionRepo,
		*fake_test.FakeAccessToken,
		*fake_test.FakeRefreshToken,
	) {
		return &fake_test.FakeAuthUserRepo{
				ByUsername: map[string]user.Pengguna{},
			},
			&fake_test.FakeUserRepo{},
			&fake_test.FakeHasher{Ok: true},
			&fake_test.FakeSessionRepo{Store: map[string]session.Session{}, SessionNextID: "session_user_1"},
			&fake_test.FakeAccessToken{},
			&fake_test.FakeRefreshToken{}
	}

	testCases := []struct {
		name      string
		cmd       auth_service.LoginCmd
		setup     func(*fake_test.FakeAuthUserRepo, *fake_test.FakeHasher, *fake_test.FakeSessionRepo, *fake_test.FakeAccessToken, *fake_test.FakeRefreshToken)
		wantErr   error
		assertion func(*testing.T, auth_service.LoginRes, *fake_test.FakeAuthUserRepo, *fake_test.FakeHasher, *fake_test.FakeSessionRepo, *fake_test.FakeAccessToken, *fake_test.FakeRefreshToken)
	}{
		{
			name:    "Path 1 -> Panjang username tidak valid",
			cmd:     auth_service.LoginCmd{Username: "adm", Password: password},
			wantErr: coreerror.ErrUsernameLengthInvalid,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.False(t, authRepo.FindCalled)
				assert.False(t, hasher.CompareCalled)
				assert.False(t, sessionRepo.RevokeExpiredCalled)
			},
		},
		{
			name: "Path 2 -> Username gagal ditemukan",
			setup: func(authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				authRepo.ByUsername = map[string]user.Pengguna{}
			},
			wantErr: coreerror.ErrInvalidCreds,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.True(t, authRepo.FindCalled)
				assert.False(t, hasher.CompareCalled)
				assert.False(t, sessionRepo.RevokeExpiredCalled)
			},
		},
		{
			name: "Path 3 -> Status akun tidak aktif",
			setup: func(authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				inactiveUser := baseUser
				inactiveUser.StatusAkun = user.NONAKTIF
				authRepo.ByUsername[username] = inactiveUser
			},
			wantErr: coreerror.ErrInvalidCreds,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.True(t, authRepo.FindCalled)
				assert.False(t, hasher.CompareCalled)
				assert.False(t, sessionRepo.RevokeExpiredCalled)
			},
		},
		{
			name: "Path 4 -> Password tidak cocok",
			setup: func(authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				authRepo.ByUsername[username] = baseUser
				hasher.Ok = false
			},
			wantErr: coreerror.ErrInvalidCreds,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.True(t, hasher.CompareCalled)
				assert.False(t, sessionRepo.RevokeExpiredCalled)
			},
		},
		{
			name: "Path 5 -> Revoke expired session gagal",
			setup: func(authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				authRepo.ByUsername[username] = baseUser
				sessionRepo.RevokeExpiredErr = testErr
			},
			wantErr: testErr,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.True(t, sessionRepo.RevokeExpiredCalled)
				assert.False(t, sessionRepo.HasActiveCalled)
			},
		},
		{
			name: "Path 6 -> Cek active session gagal",
			setup: func(authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				authRepo.ByUsername[username] = baseUser
				sessionRepo.HasActiveErr = testErr
			},
			wantErr: testErr,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.True(t, sessionRepo.RevokeExpiredCalled)
				assert.True(t, sessionRepo.HasActiveCalled)
				assert.False(t, sessionRepo.CreateSessionCalled)
			},
		},
		{
			name: "Path 7 -> User masih memiliki session aktif",
			setup: func(authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				authRepo.ByUsername[username] = baseUser
				sessionRepo.Store["existing_session"] = session.Session{
					SessionID: "existing_session",
					UserID:    baseUser.ID,
					Revoked:   false,
					ExpiresAt: time.Now().Add(1 * time.Hour),
				}
			},
			wantErr: coreerror.ErrHasSession,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.True(t, sessionRepo.HasActiveCalled)
				assert.False(t, sessionRepo.CreateSessionCalled)
			},
		},
		{
			name: "Path 8 -> Create session gagal",
			setup: func(authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				authRepo.ByUsername[username] = baseUser
				sessionRepo.CreateSessionErr = testErr
			},
			wantErr: testErr,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.True(t, sessionRepo.CreateSessionCalled)
				assert.False(t, accessToken.GenerateCalled)
			},
		},
		{
			name: "Path 9 -> Generate access token gagal",
			setup: func(authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				authRepo.ByUsername[username] = baseUser
				accessToken.GenerateAccessTokenErr = testErr
			},
			wantErr: testErr,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.True(t, sessionRepo.CreateSessionCalled)
				assert.True(t, accessToken.GenerateCalled)
				assert.False(t, refreshToken.GenerateCalled)
				_, exists := sessionRepo.Store["session_user_1"]
				assert.True(t, exists)
			},
		},
		{
			name: "Path 10 -> Generate refresh token gagal",
			setup: func(authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				authRepo.ByUsername[username] = baseUser
				refreshToken.GenerateRefreshTokenErr = testErr
			},
			wantErr: testErr,
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				assert.Equal(t, auth_service.LoginRes{}, got)
				assert.True(t, sessionRepo.CreateSessionCalled)
				assert.True(t, accessToken.GenerateCalled)
				assert.True(t, refreshToken.GenerateCalled)
				_, exists := sessionRepo.Store["session_user_1"]
				assert.True(t, exists)
			},
		},
		{
			name: "Path 11 -> Login berhasil",
			setup: func(authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
				authRepo.ByUsername[username] = baseUser
			},
			assertion: func(t *testing.T, got auth_service.LoginRes, authRepo *fake_test.FakeAuthUserRepo, hasher *fake_test.FakeHasher, sessionRepo *fake_test.FakeSessionRepo, accessToken *fake_test.FakeAccessToken, refreshToken *fake_test.FakeRefreshToken) {
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
