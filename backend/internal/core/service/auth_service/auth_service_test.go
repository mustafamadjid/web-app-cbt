package auth_service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
	// "github.com/mustafamadjid/web-app-cbt/internal/core/service/auth_service"
	"github.com/stretchr/testify/assert"
)

type fakeAuthUserRepo struct {
	byUsername      map[string]user.Pengguna
	findUsernameErr error
}

type fakeUserRepo struct {
	byID        map[user.ID]user.Pengguna
	findUserErr error
}

type fakeHasher struct {
	ok bool
}

type fakeSessionRepo struct {
	sessionNextID    string
	store            map[string]session.Session
	createSessionErr error
	revokeSessionErr error
	hasActiveErr     error
}

type fakeAccessToken struct {
	GenerateAccessTokenErr error
}

type fakeRefreshToken struct {
	GenerateRefreshTokenErr error
	VerifyRefreshTokenErr   error
}

func (fakeRepo *fakeAuthUserRepo) FindByUsername(ctx context.Context, username string) (user.Pengguna, error) {
	if fakeRepo.findUsernameErr != nil {
		return user.Pengguna{}, fakeRepo.findUsernameErr
	}
	u, ok := fakeRepo.byUsername[username]
	if !ok {
		return user.Pengguna{}, coreerror.ErrNotFound
	}
	return u, nil
}

func (fakeRepo *fakeUserRepo) FindUserByID(ctx context.Context, id user.ID) (user.Pengguna, error) {
	if fakeRepo.findUserErr != nil {
		return user.Pengguna{}, fakeRepo.findUserErr
	}
	if fakeRepo.byID == nil {
		return user.Pengguna{}, coreerror.ErrNotFound
	}
	u, ok := fakeRepo.byID[id]
	if !ok {
		return user.Pengguna{}, coreerror.ErrNotFound
	}
	return u, nil
}

func (fakeRepo *fakeUserRepo) UserExistByUsername(ctx context.Context, username string) (bool, error) {
	return false, nil
}

func (fakeRepo *fakeUserRepo) CreateUser(ctx context.Context, pengguna user.Pengguna) (user.ID, error) {
	return 0, nil
}

func (fakeRepo *fakeUserRepo) UpdateUser(ctx context.Context, idPengguna user.ID, pengguna outuser.UpdatePenggunaPatch) error {
	return nil
}

func (fakeRepo *fakeUserRepo) DeleteUser(ctx context.Context, id user.ID) error {
	return nil
}

func (fakeRepo *fakeUserRepo) ListUser(ctx context.Context) ([]user.Pengguna, error) {
	return []user.Pengguna{}, nil
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

func (fakeSession *fakeSessionRepo) GetSessionByUserId(ctx context.Context, userId user.ID) (session.Session, error) {
	for _, sess := range fakeSession.store {
		if sess.UserID == userId {
			return sess, nil
		}
	}
	return session.Session{}, coreerror.ErrNotFound
}

func (fakeSession *fakeSessionRepo)GetAllActiveSession(ctx context.Context) ([]session.Session, error) {
	var sessions []session.Session
	for _, sess := range fakeSession.store {
		if !sess.Revoked {
			sessions = append(sessions, sess)
		}
	}
	return sessions, nil
}

func (fakeSession *fakeSessionRepo)HasActiveSession(ctx context.Context, userID user.ID) (bool, error) {
	if fakeSession.hasActiveErr != nil {
		return false, fakeSession.hasActiveErr
	}
	for _, sess := range fakeSession.store {
		if sess.UserID == userID && !sess.Revoked {
			return true, nil
		}
	}
	return false, nil
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
	ErrSessionDown     = errors.New("session db down")
	ErrSignAccessFail  = errors.New("sign access fail")
	ErrSignRefreshFail = errors.New("sign refresh fail")
	ErrRevokeFailed    = errors.New("session revoke failed")
)

func TestLogin_Success_ReturnToken(t *testing.T) {
	authUsers := &fakeAuthUserRepo{
		byUsername: map[string]user.Pengguna{
			"myadmin": {ID: 1, Username: "myadmin", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF, NamaLengkap: "My Admin", Email: "admin@example.com", JenisKelamin: "LAKI_LAKI", NoHp: "081234567891"},
		},
	}
	users := &fakeUserRepo{}
	hasher := &fakeHasher{ok: true}
	sessions := &fakeSessionRepo{sessionNextID: "sess_1"}
	accessTokens := &fakeAccessToken{}
	refreshTokens := &fakeRefreshToken{}
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
				authUsers := &fakeAuthUserRepo{byUsername: map[string]user.Pengguna{}}
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, LoginCmd{Username: "myadmin", Password: "password"}
			},
			wantErr: coreerror.ErrInvalidCreds,
		},
		{
			name: "P2 status tidak aktif -> invalid creds",
			setup: func() (*AuthService, LoginCmd) {
				authUsers := &fakeAuthUserRepo{
					byUsername: map[string]user.Pengguna{
						"myadmin": {ID: 1, Username: "myadmin", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.NONAKTIF},
					},
				}
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, LoginCmd{Username: "myadmin", Password: "password"}
			},
			wantErr: coreerror.ErrInvalidCreds,
		},
		{
			name: "P3 wrong password -> invalid creds",
			setup: func() (*AuthService, LoginCmd) {
				authUsers := &fakeAuthUserRepo{
					byUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: false}
				sessions := &fakeSessionRepo{}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, LoginCmd{Username: "admin@example.com", Password: "wrong"}
			},
			wantErr: coreerror.ErrInvalidCreds,
		},
		{
			name: "P4 create session error -> propagated",
			setup: func() (*AuthService, LoginCmd) {
				authUsers := &fakeAuthUserRepo{
					byUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{createSessionErr: ErrSessionDown}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, LoginCmd{Username: "admin@example.com", Password: "password"}
			},
			wantErr: ErrSessionDown,
		},
		{
			name: "P5 has active session error -> propagated",
			setup: func() (*AuthService, LoginCmd) {
				authUsers := &fakeAuthUserRepo{
					byUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{hasActiveErr: ErrSessionDown}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, LoginCmd{Username: "admin@example.com", Password: "password"}
			},
			wantErr: ErrSessionDown,
		},
		{
			name: "P6 has active session -> has session",
			setup: func() (*AuthService, LoginCmd) {
				authUsers := &fakeAuthUserRepo{
					byUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{
					store: map[string]session.Session{
						"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)},
					},
				}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, LoginCmd{Username: "admin@example.com", Password: "password"}
			},
			wantErr: coreerror.ErrHasSession,
		},
		{
			name: "P7 access token error -> propagated",
			setup: func() (*AuthService, LoginCmd) {
				authUsers := &fakeAuthUserRepo{
					byUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{sessionNextID: "sess_1"}
				accessTokens := &fakeAccessToken{GenerateAccessTokenErr: ErrSignAccessFail}
				refreshTokens := &fakeRefreshToken{}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, LoginCmd{Username: "admin@example.com", Password: "password"}
			},
			wantErr: ErrSignAccessFail,
		},
		{
			name: "P8 refresh token error -> propagated",
			setup: func() (*AuthService, LoginCmd) {
				authUsers := &fakeAuthUserRepo{
					byUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{sessionNextID: "sess_1"}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{GenerateRefreshTokenErr: ErrSignRefreshFail}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, LoginCmd{Username: "admin@example.com", Password: "password"}
			},
			wantErr: ErrSignRefreshFail,
		},
		{
			name: "P9 success",
			setup: func() (*AuthService, LoginCmd) {
				authUsers := &fakeAuthUserRepo{
					byUsername: map[string]user.Pengguna{
						"admin@example.com": {ID: 1, Username: "admin@example.com", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF},
					},
				}
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{sessionNextID: "sess_99"}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
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
				authUsers := &fakeAuthUserRepo{}
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{VerifyRefreshTokenErr: coreerror.ErrInvalidToken}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, "bad-token"
			},
			wantErr: coreerror.ErrInvalidToken,
		},
		{
			name: "P2 revoke session error -> propagated",
			setup: func() (*AuthService, string) {
				authUsers := &fakeAuthUserRepo{}
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{
					store: map[string]session.Session{
						"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)},
					},
					revokeSessionErr: ErrRevokeFailed,
				}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				svc := NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
				return svc, "REFRESH TOKEN : sess_1"
			},
			wantErr: ErrRevokeFailed,
		},
		{
			name: "P3 success",
			setup: func() (*AuthService, string) {
				authUsers := &fakeAuthUserRepo{}
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{
					store: map[string]session.Session{
						"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)},
					},
				}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
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
				authUsers := &fakeAuthUserRepo{}
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				return NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
			},
			token:   "",
			wantErr: coreerror.ErrNoTokenProvided,
		},
		{
			name: "P2 refresh token invalid -> propagated",
			setup: func() *AuthService {
				authUsers := &fakeAuthUserRepo{}
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{VerifyRefreshTokenErr: coreerror.ErrInvalidToken}
				return NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
			},
			token:   "bad-token",
			wantErr: coreerror.ErrInvalidToken,
		},
		{
			name: "P3 session not found -> propagated",
			setup: func() *AuthService {
				authUsers := &fakeAuthUserRepo{}
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{store: map[string]session.Session{}}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				return NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
			},
			token:   "REFRESH TOKEN : sess_1",
			wantErr: coreerror.ErrNotFound,
		},
		{
			name: "P4 session expired -> session expired",
			setup: func() *AuthService {
				authUsers := &fakeAuthUserRepo{}
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{
					store: map[string]session.Session{
						"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(-time.Hour)},
					},
				}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				return NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
			},
			token:   "REFRESH TOKEN : sess_1",
			wantErr: coreerror.ErrSessionExpired,
		},
		{
			name: "P5 session revoked -> session expired",
			setup: func() *AuthService {
				authUsers := &fakeAuthUserRepo{}
				users := &fakeUserRepo{}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{
					store: map[string]session.Session{
						"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour), Revoked: true},
					},
				}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				return NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
			},
			token:   "REFRESH TOKEN : sess_1",
			wantErr: coreerror.ErrSessionExpired,
		},
		{
			name: "P6 user not found -> propagated",
			setup: func() *AuthService {
				authUsers := &fakeAuthUserRepo{}
				users := &fakeUserRepo{byID: map[user.ID]user.Pengguna{}}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{
					store: map[string]session.Session{
						"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)},
					},
				}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
				return NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
			},
			token:   "REFRESH TOKEN : sess_1",
			wantErr: coreerror.ErrNotFound,
		},
		{
			name: "P7 access token error -> propagated",
			setup: func() *AuthService {
				authUsers := &fakeAuthUserRepo{}
				users := &fakeUserRepo{
					byID: map[user.ID]user.Pengguna{
						1: {ID: 1, Username: "admin@example.com", Role: "ADMIN"},
					},
				}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{
					store: map[string]session.Session{
						"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)},
					},
				}
				accessTokens := &fakeAccessToken{GenerateAccessTokenErr: ErrSignAccessFail}
				refreshTokens := &fakeRefreshToken{}
				return NewAuthService(authUsers, users, hasher, sessions, accessTokens, refreshTokens)
			},
			token:   "REFRESH TOKEN : sess_1",
			wantErr: ErrSignAccessFail,
		},
		{
			name: "P8 success",
			setup: func() *AuthService {
				authUsers := &fakeAuthUserRepo{}
				users := &fakeUserRepo{
					byID: map[user.ID]user.Pengguna{
						1: {ID: 1, Username: "admin@example.com", Role: "ADMIN"},
					},
				}
				hasher := &fakeHasher{ok: true}
				sessions := &fakeSessionRepo{
					store: map[string]session.Session{
						"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)},
					},
				}
				accessTokens := &fakeAccessToken{}
				refreshTokens := &fakeRefreshToken{}
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
	sessions := &fakeSessionRepo{
		store: map[string]session.Session{
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
	sessions := &fakeSessionRepo{
		store: map[string]session.Session{
			"sess_1": {SessionID: "sess_1", UserID: 1, ExpiresAt: time.Now().Add(time.Hour)},
		},
	}

	_, err := sessions.GetSessionByUserId(context.Background(), 3)

	assert.Error(t, err)
	assert.ErrorIs(t, err, coreerror.ErrNotFound)
}
