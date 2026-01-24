package auth_test

import (
	"context"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
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

type fakeToken struct {
	GenerateAccessTokenErr error
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

func (fakeToken *fakeToken) GenerateAccessToken(userID int, role string, tokenDuration time.Duration) (string,error){
	if fakeToken.GenerateAccessTokenErr != nil {
		return "",fakeToken.GenerateAccessTokenErr
	}
	return "ACCESS TOKEN",nil
}

func (faketoken *fakeToken) GenerateRefreshToken(sessionID string, tokenDuration time.Duration) (string,error){
	if faketoken.GenerateRefreshTokenErr != nil {
		return "",faketoken.GenerateRefreshTokenErr
	}
	return "REFRESH TOKEN : " + sessionID ,nil
}

func (faketoken *fakeToken) VerifyRefreshToken(token string) (sessionID string,err error){
	const prefix = "BEARER "
	if len(token) <= len(prefix) || token[:len(prefix)] != prefix {
		return "", coreerror.ErrInvalidToken
	}
	return token[len(prefix):], nil
}
