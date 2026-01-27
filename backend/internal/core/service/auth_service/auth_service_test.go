package auth_service_test

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/login"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/auth_service"
	"github.com/stretchr/testify/assert"
)

type fakeUserRepo struct {
	byUsername      map[string]user.Pengguna
	findUsernameErr error
}

type fakeHasher struct {
	ok bool
}

type fakeSessionRepo struct {
	sessionNextID    string
	store            map[string]session.Session
	createSessionErr error
	revokeSessionErr error
}

type fakeAccessToken struct {
	GenerateAccessTokenErr error
}

type fakeRefreshToken struct {
	GenerateRefreshTokenErr error
	VerifyRefreshTokenErr   error
}

func (fakeRepo *fakeUserRepo) FindByUsername(ctx context.Context, username string) (user.Pengguna, error) {
	if fakeRepo.findUsernameErr != nil {
		return user.Pengguna{}, fakeRepo.findUsernameErr
	}
	u, ok := fakeRepo.byUsername[username]
	if !ok {
		return user.Pengguna{}, coreerror.ErrNotFound
	}

	return u, nil
}

func (fakeHash *fakeHasher) ComparePaswordAndHashed(hash string, plain string) bool {
	return fakeHash.ok
}

func (fakeHash *fakeHasher) GenerateHash(plain string) (string, error) {
	return "hashed password", nil
}

func (fakeSession *fakeSessionRepo) CreateSession(ctx context.Context, userID user.ID, expiresAt time.Time) (string, error) {
	if fakeSession.createSessionErr != nil {
		return "", fakeSession.createSessionErr
	}
	if fakeSession.store == nil {
		fakeSession.store = map[string]session.Session{}
	}

	if fakeSession.sessionNextID == "" {
		fakeSession.sessionNextID = "session_user_1"
	}

	id := fakeSession.sessionNextID
	fakeSession.store[id] = session.Session{
		SessionID: id,
		UserID:    userID,
		Revoked:   false,
		ExpiresAt: expiresAt,
	}
	return id, nil
}

func (fakeSession *fakeSessionRepo) GetSession(ctx context.Context, sessionID string) (session.Session, error) {
	ss, ok := fakeSession.store[sessionID]
	if !ok {
		return session.Session{}, coreerror.ErrNotFound
	}
	return ss, nil
}

func (fakeSession *fakeSessionRepo) RevokeSession(ctx context.Context, sessionID string) error {
	if fakeSession.revokeSessionErr != nil {
		return fakeSession.revokeSessionErr
	}
	ss, ok := fakeSession.store[sessionID]
	if !ok {
		return coreerror.ErrNotFound
	}
	ss.Revoked = true
	fakeSession.store[sessionID] = ss
	return nil
}

func (fakeSession *fakeSessionRepo) RevokeSessionAllbyUser(ctx context.Context, userID user.ID, now time.Time) error {
	for ssid, sess := range fakeSession.store {
		if sess.UserID == userID && !sess.Revoked && now.Before(sess.ExpiresAt) {
			sess.Revoked = true
			fakeSession.store[ssid] = sess
		}
	}
	return nil
}

func (f *fakeAccessToken) GenerateAccessToken(userID user.ID, role user.Role, username string, tokenDuration time.Duration) (string, error) {
	if f.GenerateAccessTokenErr != nil {
		return "", f.GenerateAccessTokenErr
	}

	return fmt.Sprintf("ACCESS TOKEN :|%v|%v|%s", userID, role, username), nil
}

func (f *fakeAccessToken) VerifyAccessToken(token string, now time.Time) (userID user.ID, role user.Role, username string, err error) {
	const prefix = "ACCESS TOKEN :|"
	if !strings.HasPrefix(token, prefix) {
		return 0, "", "", coreerror.ErrInvalidToken
	}

	parts := strings.Split(token[len(prefix):], "|")
	if len(parts) != 3 {
		return 0, "", "", coreerror.ErrInvalidToken
	}

	idValue, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", "", coreerror.ErrInvalidToken
	}

	return user.ID(idValue), user.Role(parts[1]), parts[2], nil
}

func (fakeRefreshToken *fakeRefreshToken) GenerateRefreshToken(sessionID string, tokenDuration time.Duration) (string, error) {
	if fakeRefreshToken.GenerateRefreshTokenErr != nil {
		return "", fakeRefreshToken.GenerateRefreshTokenErr
	}
	return "REFRESH TOKEN : " + sessionID, nil
}

func (fakeRefreshToken *fakeRefreshToken) VerifyRefreshToken(token string, now time.Time) (sessionID string, err error) {
	if fakeRefreshToken.VerifyRefreshTokenErr != nil {
		return "", fakeRefreshToken.VerifyRefreshTokenErr
	}
	const prefix = "REFRESH TOKEN : "
	if len(token) <= len(prefix) || token[:len(prefix)] != prefix {
		return "", coreerror.ErrInvalidToken
	}
	return token[len(prefix):], nil
}

var (
	ErrSessionDown    = errors.New("session db down")
	ErrSignAccessFail = errors.New("sign access fail")
	ErrSignRefreshFail = errors.New("sign refresh fail")
)

func TestLogin_Success_ReturnToken(t *testing.T) {
	users := &fakeUserRepo{
		byUsername: map[string]user.Pengguna{
			"myadmin": {ID: 1, Username: "myadmin", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF, NamaLengkap: "My Admin", Email: "admin@example.com", JenisKelamin: "LAKI_LAKI", NoHp: "081234567891"},
		},
	}

	hasher := &fakeHasher{ok: true}
	sessions := &fakeSessionRepo{sessionNextID: "sess_1"}
	accessTokens := &fakeAccessToken{}
	refreshTokens := &fakeRefreshToken{}

	uc := auth_service.NewAuthService(users, hasher, sessions, accessTokens, refreshTokens)

	res, err := uc.Login(context.Background(), login.LoginCmd{
		Username: "myadmin", Password: "password",
	})

	assert.NoError(t, err)
	assert.Equal(t, "ACCESS TOKEN :|1|ADMIN|myadmin", res.AccessToken)
	assert.Equal(t, "REFRESH TOKEN : sess_1", res.RefreshToken)
}

func TestLogin_BasisPaths(t *testing.T) {
	type tc struct {
		name        string
		setup       func() (*auth_service.AuthService, login.LoginCmd)
		wantErr     error
		wantAccess  string
		wantRefresh string
	}

	tests := []tc{
		{
			name: "P1 user not found -> invalid creds",
			setup: func() (*auth_service.AuthService, login.LoginCmd) {
				users := &fakeUserRepo{byUsername: map[string]user.Pengguna{}}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				svc := auth_service.NewAuthService(users, hasher, sessions, accessTokens, refreshTokens)
				return svc, login.LoginCmd{Username: "myadmin", Password: "password"}
			},
			wantErr: coreerror.ErrInvalidCreds,
		},
		{
			name: "P2 status tidak aktif -> invalid creds",
			setup: func() (*auth_service.AuthService, login.LoginCmd) {
				users := &fakeUserRepo{
					byUsername: map[string]user.Pengguna{
						"myadmin": {ID: 1, Username: "myadmin", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.NONAKTIF},
					},
				}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				svc := auth_service.NewAuthService(users, hasher, sessions, accessTokens, refreshTokens)
				return svc, login.LoginCmd{Username: "myadmin", Password: "password"}
			},
			wantErr: coreerror.ErrInvalidCreds,
		},
		{
			name: "P3 wrong password -> invalid creds",
			setup: func() (*auth_service.AuthService, login.LoginCmd) {
				users := &fakeUserRepo{
					byUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				hasher := &fakeHasher{ok: false}
				sessions := &fakeSessionRepo{}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				svc := auth_service.NewAuthService(users, hasher, sessions, accessTokens, refreshTokens)
				return svc, login.LoginCmd{Username: "admin@example.com", Password: "wrong"}
			},
			wantErr: coreerror.ErrInvalidCreds,
		},
		{
			name: "P4 create session error -> propagated",
			setup: func() (*auth_service.AuthService, login.LoginCmd) {
				users := &fakeUserRepo{
					byUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{createSessionErr: fmt.Errorf("session db down")}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				svc := auth_service.NewAuthService(users, hasher, sessions, accessTokens, refreshTokens)
				return svc, login.LoginCmd{Username: "admin@example.com", Password: "password"}
			},
			wantErr: fmt.Errorf("session db down"),
		},
		{
			name: "P5 access token error -> propagated",
			setup: func() (*auth_service.AuthService, login.LoginCmd) {
				users := &fakeUserRepo{
					byUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{sessionNextID: "sess_1"}
				accessTokens := &fakeAccessToken{GenerateAccessTokenErr: fmt.Errorf("sign access fail")}
				refreshTokens := &fakeRefreshToken{}
				svc := auth_service.NewAuthService(users, hasher, sessions, accessTokens, refreshTokens)
				return svc, login.LoginCmd{Username: "admin@example.com", Password: "password"}
			},
			wantErr: fmt.Errorf("sign access fail"),
		},
		{
			name: "P6 refresh token error -> propagated",
			setup: func() (*auth_service.AuthService, login.LoginCmd) {
				users := &fakeUserRepo{
					byUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{sessionNextID: "sess_1"}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{GenerateRefreshTokenErr: fmt.Errorf("sign refresh fail")}
				svc := auth_service.NewAuthService(users, hasher, sessions, accessTokens, refreshTokens)
				return svc, login.LoginCmd{Username: "admin@example.com", Password: "password"}
			},
			wantErr: fmt.Errorf("sign refresh fail"),
		},
		{
			name: "P7 success",
			setup: func() (*auth_service.AuthService, login.LoginCmd) {
				users := &fakeUserRepo{
					byUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{sessionNextID: "sess_99"}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				svc := auth_service.NewAuthService(users, hasher, sessions, accessTokens, refreshTokens)
				return svc, login.LoginCmd{Username: "admin@example.com", Password: "password"}
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
			assert.EqualError(t, err, tt.wantErr.Error())
		})
	}
}

func TestLogout_BasisPaths(t *testing.T) {
	type tc struct {
		name    string
		setup   func() (*auth_service.AuthService, string)
		wantErr error
	}

	tests := []tc{
		{
			name: "P1 invalid refresh token -> invalid token",
			setup: func() (*auth_service.AuthService, string) {
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{VerifyRefreshTokenErr: coreerror.ErrInvalidToken}
				svc := auth_service.NewAuthService(users, hasher, sessions, accessTokens, refreshTokens)
				return svc, "bad-token"
			},
			wantErr: coreerror.ErrInvalidToken,
		},
		{
			name: "P2 revoke session error -> propagated",
			setup: func() (*auth_service.AuthService, string) {
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{
					store: map[string]session.Session{
						"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)},
					},
					revokeSessionErr: fmt.Errorf("session revoke failed"),
				}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				svc := auth_service.NewAuthService(users, hasher, sessions, accessTokens, refreshTokens)
				return svc, "REFRESH TOKEN : sess_1"
			},
			wantErr: fmt.Errorf("session revoke failed"),
		},
		{
			name: "P3 success",
			setup: func() (*auth_service.AuthService, string) {
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{
					store: map[string]session.Session{
						"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)},
					},
				}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				svc := auth_service.NewAuthService(users, hasher, sessions, accessTokens, refreshTokens)
				return svc, "REFRESH TOKEN : sess_1"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, token := tt.setup()
			err := svc.Logout(context.Background(), token, time.Now())
			if tt.wantErr == nil {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
			assert.EqualError(t, err, tt.wantErr.Error())
		})
	}
}
