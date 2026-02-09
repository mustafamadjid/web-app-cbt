package auth_service

import (
	"context"
	"errors"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	fake_test "github.com/mustafamadjid/web-app-cbt/internal/core/service/auth_service/fake_test"
	// "github.com/mustafamadjid/web-app-cbt/internal/core/service/auth_service"
	"github.com/stretchr/testify/assert"
)

var (
	ErrSessionDown     = errors.New("session db down")
	ErrSignAccessFail  = errors.New("sign access fail")
	ErrSignRefreshFail = errors.New("sign refresh fail")
	ErrRevokeFailed    = errors.New("session revoke failed")
)

func TestLogin_Success_ReturnToken(t *testing.T) {
	authUsers := &fake_test.FakeAuthUserRepo{
		ByUsername: map[string]user.Pengguna{
			"myadmin": {ID: 1, Username: "myadmin", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF, NamaLengkap: "My Admin", Email: "admin@example.com", JenisKelamin: "LAKI_LAKI", NoHp: "081234567891"},
		},
	}
	users := &fake_test.FakeUserRepo{}
	hasher := &fake_test.FakeHasher{Ok: true}
	sessions := &fake_test.FakeSessionRepo{SessionNextID: "sess_1"}
	accessTokens := &fake_test.FakeAccessToken{}
	refreshTokens := &fake_test.FakeRefreshToken{}
	uc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)

	res, err := uc.Login(context.Background(), LoginCmd{Username: "myadmin", Password: "password"})

	assert.NoError(t, err)
	assert.Equal(t, "ACCESS TOKEN :|1|ADMIN|myadmin", res.AccessToken)
	assert.Equal(t, "REFRESH TOKEN : sess_1", res.RefreshToken)
}

func TestLogin_BasisPaths(t *testing.T) {
	type tc struct {
		name        string
		setup       func() (*AuthService, LoginCmd)
		wantErr     error
		wantAccess  string
		wantRefresh string
	}

	tests := []tc{
		{
			name: "P1 user not found -> invalid creds",
			setup: func() (*AuthService, LoginCmd) {
				authUsers := &fake_test.FakeAuthUserRepo{ByUsername: map[string]user.Pengguna{}}
				users := &fake_test.FakeUserRepo{}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{}
				accessTokens := &fake_test.FakeAccessToken{}
				refreshTokens := &fake_test.FakeRefreshToken{}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, LoginCmd{Username: "myadmin", Password: "password"}
			},
			wantErr: coreerror.ErrInvalidCreds,
		},
		{
			name: "P2 status tidak aktif -> invalid creds",
			setup: func() (*AuthService, LoginCmd) {
				authUsers := &fake_test.FakeAuthUserRepo{
					ByUsername: map[string]user.Pengguna{
						"myadmin": {ID: 1, Username: "myadmin", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.NONAKTIF},
					},
				}
				users := &fake_test.FakeUserRepo{}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{}
				accessTokens := &fake_test.FakeAccessToken{}
				refreshTokens := &fake_test.FakeRefreshToken{}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, LoginCmd{Username: "myadmin", Password: "password"}
			},
			wantErr: coreerror.ErrInvalidCreds,
		},
		{
			name: "P3 wrong password -> invalid creds",
			setup: func() (*AuthService, LoginCmd) {
				authUsers := &fake_test.FakeAuthUserRepo{
					ByUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				users := &fake_test.FakeUserRepo{}
				hasher := &fake_test.FakeHasher{Ok: false}
				sessions := &fake_test.FakeSessionRepo{}
				accessTokens := &fake_test.FakeAccessToken{}
				refreshTokens := &fake_test.FakeRefreshToken{}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, LoginCmd{Username: "admin@example.com", Password: "wrong"}
			},
			wantErr: coreerror.ErrInvalidCreds,
		},
		{
			name: "P4 create session error -> propagated",
			setup: func() (*AuthService, LoginCmd) {
				authUsers := &fake_test.FakeAuthUserRepo{
					ByUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				users := &fake_test.FakeUserRepo{}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{CreateSessionErr: ErrSessionDown}
				accessTokens := &fake_test.FakeAccessToken{}
				refreshTokens := &fake_test.FakeRefreshToken{}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, LoginCmd{Username: "admin@example.com", Password: "password"}
			},
			wantErr: ErrSessionDown,
		},
		{
			name: "P5 has active session error -> propagated",
			setup: func() (*AuthService, LoginCmd) {
				authUsers := &fake_test.FakeAuthUserRepo{
					ByUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				users := &fake_test.FakeUserRepo{}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{HasActiveErr: ErrSessionDown}
				accessTokens := &fake_test.FakeAccessToken{}
				refreshTokens := &fake_test.FakeRefreshToken{}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, LoginCmd{Username: "admin@example.com", Password: "password"}
			},
			wantErr: ErrSessionDown,
		},
		{
			name: "P6 has active session -> has session",
			setup: func() (*AuthService, LoginCmd) {
				authUsers := &fake_test.FakeAuthUserRepo{
					ByUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				users := &fake_test.FakeUserRepo{}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{
					Store: map[string]session.Session{
						"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)},
					},
				}
				accessTokens := &fake_test.FakeAccessToken{}
				refreshTokens := &fake_test.FakeRefreshToken{}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, LoginCmd{Username: "admin@example.com", Password: "password"}
			},
			wantErr: coreerror.ErrHasSession,
		},
		{
			name: "P7 access token error -> propagated",
			setup: func() (*AuthService, LoginCmd) {
				authUsers := &fake_test.FakeAuthUserRepo{
					ByUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				users := &fake_test.FakeUserRepo{}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{SessionNextID: "sess_1"}
				accessTokens := &fake_test.FakeAccessToken{GenerateAccessTokenErr: ErrSignAccessFail}
				refreshTokens := &fake_test.FakeRefreshToken{}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, LoginCmd{Username: "admin@example.com", Password: "password"}
			},
			wantErr: ErrSignAccessFail,
		},
		{
			name: "P8 refresh token error -> propagated",
			setup: func() (*AuthService, LoginCmd) {
				authUsers := &fake_test.FakeAuthUserRepo{
					ByUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				users := &fake_test.FakeUserRepo{}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{SessionNextID: "sess_1"}
				accessTokens := &fake_test.FakeAccessToken{}
				refreshTokens := &fake_test.FakeRefreshToken{GenerateRefreshTokenErr: ErrSignRefreshFail}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, LoginCmd{Username: "admin@example.com", Password: "password"}
			},
			wantErr: ErrSignRefreshFail,
		},
		{
			name: "P9 success",
			setup: func() (*AuthService, LoginCmd) {
				authUsers := &fake_test.FakeAuthUserRepo{
					ByUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				users := &fake_test.FakeUserRepo{}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{SessionNextID: "sess_99"}
				accessTokens := &fake_test.FakeAccessToken{}
				refreshTokens := &fake_test.FakeRefreshToken{}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, LoginCmd{Username: "admin@example.com", Password: "password"}
			},
			wantAccess:  "ACCESS TOKEN :|1|ADMIN|admin@example.com",
			wantRefresh: "REFRESH TOKEN : sess_99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, cmd := tt.setup()
			res, err := svc.Login(context.Background(), cmd)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantAccess, res.AccessToken)
				assert.Equal(t, tt.wantRefresh, res.RefreshToken)
				return
			}

			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestLogout_BasisPaths(t *testing.T) {
	type tc struct {
		name    string
		setup   func() (*AuthService, string)
		wantErr error
	}

	tests := []tc{
		{
			name: "P1 invalid refresh token -> invalid token",
			setup: func() (*AuthService, string) {
				authUsers := &fake_test.FakeAuthUserRepo{}
				users := &fake_test.FakeUserRepo{}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{}
				accessTokens := &fake_test.FakeAccessToken{}
				refreshTokens := &fake_test.FakeRefreshToken{VerifyRefreshTokenErr: coreerror.ErrInvalidToken}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, "bad-token"
			},
			wantErr: coreerror.ErrInvalidToken,
		},
		{
			name: "P2 revoke session error -> propagated",
			setup: func() (*AuthService, string) {
				authUsers := &fake_test.FakeAuthUserRepo{}
				users := &fake_test.FakeUserRepo{}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{
					Store: map[string]session.Session{
						"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)},
					},
					RevokeSessionErr: ErrRevokeFailed,
				}
				accessTokens := &fake_test.FakeAccessToken{}
				refreshTokens := &fake_test.FakeRefreshToken{}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, "REFRESH TOKEN : sess_1"
			},
			wantErr: ErrRevokeFailed,
		},
		{
			name: "P3 success",
			setup: func() (*AuthService, string) {
				authUsers := &fake_test.FakeAuthUserRepo{}
				users := &fake_test.FakeUserRepo{}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{
					Store: map[string]session.Session{
						"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)},
					},
				}
				accessTokens := &fake_test.FakeAccessToken{}
				refreshTokens := &fake_test.FakeRefreshToken{}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, "REFRESH TOKEN : sess_1"
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc, token := tt.setup()
			err := svc.Logout(context.Background(), token, time.Now())

			if tt.wantErr == nil {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestRefreshAccessToken_BasisPaths(t *testing.T) {
	type tc struct {
		name    string
		setup   func() *AuthService
		token   string
		wantErr error
		wantTok string
	}

	tests := []tc{
		{
			name: "P1 empty token -> no token provided",
			setup: func() *AuthService {
				authUsers := &fake_test.FakeAuthUserRepo{}
				users := &fake_test.FakeUserRepo{}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{}
				accessTokens := &fake_test.FakeAccessToken{}
				refreshTokens := &fake_test.FakeRefreshToken{}
				return NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
			},
			token:   "",
			wantErr: coreerror.ErrNoTokenProvided,
		},
		{
			name: "P2 refresh token invalid -> propagated",
			setup: func() *AuthService {
				authUsers := &fake_test.FakeAuthUserRepo{}
				users := &fake_test.FakeUserRepo{}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{}
				accessTokens := &fake_test.FakeAccessToken{}
				refreshTokens := &fake_test.FakeRefreshToken{VerifyRefreshTokenErr: coreerror.ErrInvalidToken}
				return NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
			},
			token:   "bad-token",
			wantErr: coreerror.ErrInvalidToken,
		},
		{
			name: "P3 session not found -> propagated",
			setup: func() *AuthService {
				authUsers := &fake_test.FakeAuthUserRepo{}
				users := &fake_test.FakeUserRepo{}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{Store: map[string]session.Session{}}
				accessTokens := &fake_test.FakeAccessToken{}
				refreshTokens := &fake_test.FakeRefreshToken{}
				return NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
			},
			token:   "REFRESH TOKEN : sess_1",
			wantErr: coreerror.ErrNotFound,
		},
		{
			name: "P4 session expired -> session expired",
			setup: func() *AuthService {
				authUsers := &fake_test.FakeAuthUserRepo{}
				users := &fake_test.FakeUserRepo{}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{
					Store: map[string]session.Session{
						"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(-time.Hour)},
					},
				}
				accessTokens := &fake_test.FakeAccessToken{}
				refreshTokens := &fake_test.FakeRefreshToken{}
				return NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
			},
			token:   "REFRESH TOKEN : sess_1",
			wantErr: coreerror.ErrSessionExpired,
		},
		{
			name: "P5 session revoked -> session expired",
			setup: func() *AuthService {
				authUsers := &fake_test.FakeAuthUserRepo{}
				users := &fake_test.FakeUserRepo{}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{
					Store: map[string]session.Session{
						"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour), Revoked: true},
					},
				}
				accessTokens := &fake_test.FakeAccessToken{}
				refreshTokens := &fake_test.FakeRefreshToken{}
				return NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
			},
			token:   "REFRESH TOKEN : sess_1",
			wantErr: coreerror.ErrSessionExpired,
		},
		{
			name: "P6 user not found -> propagated",
			setup: func() *AuthService {
				authUsers := &fake_test.FakeAuthUserRepo{}
				users := &fake_test.FakeUserRepo{ByID: map[user.ID]user.Pengguna{}}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{
					Store: map[string]session.Session{
						"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)},
					},
				}
				accessTokens := &fake_test.FakeAccessToken{}
				refreshTokens := &fake_test.FakeRefreshToken{}
				return NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
			},
			token:   "REFRESH TOKEN : sess_1",
			wantErr: coreerror.ErrNotFound,
		},
		{
			name: "P7 access token error -> propagated",
			setup: func() *AuthService {
				authUsers := &fake_test.FakeAuthUserRepo{}
				users := &fake_test.FakeUserRepo{
					ByID: map[user.ID]user.Pengguna{
						1: {ID: 1, Username: "admin@example.com", Role: "ADMIN"},
					},
				}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{
					Store: map[string]session.Session{
						"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)},
					},
				}
				accessTokens := &fake_test.FakeAccessToken{GenerateAccessTokenErr: ErrSignAccessFail}
				refreshTokens := &fake_test.FakeRefreshToken{}
				return NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
			},
			token:   "REFRESH TOKEN : sess_1",
			wantErr: ErrSignAccessFail,
		},
		{
			name: "P8 success",
			setup: func() *AuthService {
				authUsers := &fake_test.FakeAuthUserRepo{}
				users := &fake_test.FakeUserRepo{
					ByID: map[user.ID]user.Pengguna{
						1: {ID: 1, Username: "admin@example.com", Role: "ADMIN"},
					},
				}
				hasher := &fake_test.FakeHasher{Ok: true}
				sessions := &fake_test.FakeSessionRepo{
					Store: map[string]session.Session{
						"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)},
					},
				}
				accessTokens := &fake_test.FakeAccessToken{}
				refreshTokens := &fake_test.FakeRefreshToken{}
				return NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
			},
			token:   "REFRESH TOKEN : sess_1",
			wantTok: "ACCESS TOKEN :|1|ADMIN|admin@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.setup()
			token, err := svc.RefreshAccessToken(context.Background(), tt.token, time.Minute*15)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantTok, token)
				return
			}

			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestSessionRepository_GetSessionByUserId(t *testing.T) {
	sessions := &fake_test.FakeSessionRepo{
		Store: map[string]session.Session{
			"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)},
			"sess_2": {SessionID: "sess_2", UserID: 2, ExpiresAt: time.Now().Add(time.Hour)},
		},
	}

	found, err := sessions.GetSessionByUserId(context.Background(), 2)

	assert.NoError(t, err)
	assert.Equal(t, "sess_2", found.SessionID)
	assert.Equal(t, user.ID(2), found.UserID)
}

func TestSessionRepository_GetSessionByUserId_NotFound(t *testing.T) {
	sessions := &fake_test.FakeSessionRepo{
		Store: map[string]session.Session{
			"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)},
		},
	}

	_, err := sessions.GetSessionByUserId(context.Background(), 3)

	assert.Error(t, err)
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}
