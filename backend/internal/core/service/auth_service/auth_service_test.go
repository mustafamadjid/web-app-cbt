package auth_service_test

import (
	"context"
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
)

type fakeUserRepo struct {
	byUsername map[string] user.Pengguna
	findUsernameErr error
}

type fakeHasher struct {
	ok bool
}

type fakeSessionRepo struct {
	sessionNextID string
	store map[string] session.Session
	createSessionErr error
}

type fakeAccessToken struct {
	GenerateAccessTokenErr error
}

type fakeRefreshToken struct {
	GenerateRefreshTokenErr error
}

func (fakeRepo *fakeUserRepo) FindByUsername(ctx context.Context, username string) (user.Pengguna, error) {
	u,ok := fakeRepo.byUsername[username]
	if !ok {
		return user.Pengguna{}, coreerror.ErrNotFound
	}

	return  u,nil
}

func (fakeHash *fakeHasher) ComparePaswordAndHashed (hash string,plain string) bool {
	return fakeHash.ok
}

func (fakeHash *fakeHasher) GenerateHash(plain string) (string,error) {
	return "hashed password",nil
}

func (fakeSession *fakeSessionRepo) CreateSession(ctx context.Context, userID user.ID, expiresAt time.Time) (string,error){
	if fakeSession.store == nil {
		fakeSession.store = map[string] session.Session{}
	}

	if fakeSession.sessionNextID == "" {
		fakeSession.sessionNextID = "session_user_1";
	}

	id := fakeSession.sessionNextID
	fakeSession.store[id] = session.Session{
		SessionID: id,
		UserID: userID,
		Revoked: false,
		ExpiresAt: expiresAt,
	}
	return id,nil
}

func (fakeSession *fakeSessionRepo) GetSession(ctx context.Context, sessionID string) (session.Session, error){
	ss,ok := fakeSession.store[sessionID]
	if !ok {
		return session.Session{}, coreerror.ErrNotFound
	}
	return ss,nil
}

func (fakeSession *fakeSessionRepo) RevokeSession(ctx context.Context, sessionID string) error {
	ss,ok := fakeSession.store[sessionID]
	if !ok {
		return coreerror.ErrNotFound
	}
	ss.Revoked = true
	fakeSession.store[sessionID] = ss
	return nil
}

func (fakeSession *fakeSessionRepo) RevokeSessionAllbyUser(ctx context.Context, userID user.ID,now time.Time) error {
	for ssid,sess := range fakeSession.store {
		if sess.UserID == userID && !sess.Revoked && now.Before(sess.ExpiresAt) {
			sess.Revoked = true
			fakeSession.store[ssid] = sess
		}
	}
	return  nil
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

	idValue,err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", "", coreerror.ErrInvalidToken
	}

	return user.ID(idValue), user.Role(parts[1]), parts[2], nil
}


func (fakeRefreshToken *fakeRefreshToken) GenerateRefreshToken(sessionID string, tokenDuration time.Duration) (string,error){
	if fakeRefreshToken.GenerateRefreshTokenErr != nil {
		return "",fakeRefreshToken.GenerateRefreshTokenErr
	}
	return "REFRESH TOKEN : " + sessionID ,nil
}

func (fakeRefreshToken *fakeRefreshToken) VerifyRefreshToken(token string, now time.Time) (sessionID string,err error){
	const prefix = "REFRESH TOKEN : "
	if len(token) <= len(prefix) || token[:len(prefix)] != prefix {
		return "", coreerror.ErrInvalidToken
	}
	return token[len(prefix):], nil
}

func TestLogin_Success_ReturnToken(t *testing.T) {
	users := &fakeUserRepo{
		byUsername: map[string]user.Pengguna {
			"myadmin": {ID: 1, Username: "myadmin", PasswordHashed: "X", Role: "ADMIN", StatusAkun: user.AKTIF, NamaLengkap: "My Admin", Email: "admin@example.com", JenisKelamin: "LAKI_LAKI", NoHp: "081234567891"},
		},
	}

	hasher := &fakeHasher{ok: true}
	sessions := &fakeSessionRepo{sessionNextID: "sess_1"}
	accessTokens := &fakeAccessToken{}
	refreshTokens := &fakeRefreshToken{}

	uc := auth_service.NewAuthService(users,hasher,sessions,accessTokens,refreshTokens)

	res, err := uc.Login(context.Background(), login.LoginCmd{
		Username: "myadmin", Password: "password",
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if res.AccessToken != "ACCESS TOKEN :|1|ADMIN|myadmin" {
		t.Fatalf("expected ACCESS TOKEN, got %s", res.AccessToken)
	}
	if res.RefreshToken != "REFRESH TOKEN : sess_1" {
		t.Fatalf("expected refresh for sess_1, got %s", res.RefreshToken)
	}
}