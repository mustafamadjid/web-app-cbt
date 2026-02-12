package auth_service_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type FakeAuthUserRepo struct {
	ByUsername      map[string]user.Pengguna
	FindUsernameErr error
}

type FakeUserRepo struct {
	ByID        map[user.ID]user.Pengguna
	FindUserErr error
}

type FakeHasher struct {
	Ok bool
}

type FakeSessionRepo struct {
	SessionNextID    string
	Store            map[string]session.Session
	CreateSessionErr error
	RevokeSessionErr error
	HasActiveErr     error

	// NEW: for RevokeExpiredSessions
	RevokeExpiredErr error
}

type FakeAccessToken struct {
	GenerateAccessTokenErr error
}

type FakeRefreshToken struct {
	GenerateRefreshTokenErr error
	VerifyRefreshTokenErr   error
}

type FakeRefreshTTl struct {
	TTL time.Duration
}

func (fakeRepo *FakeAuthUserRepo) FindByUsername(ctx context.Context, username string) (user.Pengguna, error) {
	if fakeRepo.FindUsernameErr != nil {
		return user.Pengguna{}, fakeRepo.FindUsernameErr
	}
	u, ok := fakeRepo.ByUsername[username]
	if !ok {
		return user.Pengguna{}, coreerror.ErrNotFound
	}
	return u, nil
}

func (fakeRepo *FakeUserRepo) FindUserByID(ctx context.Context, id user.ID) (user.Pengguna, error) {
	if fakeRepo.FindUserErr != nil {
		return user.Pengguna{}, fakeRepo.FindUserErr
	}
	if fakeRepo.ByID == nil {
		return user.Pengguna{}, coreerror.ErrNotFound
	}
	u, ok := fakeRepo.ByID[id]
	if !ok {
		return user.Pengguna{}, coreerror.ErrNotFound
	}
	return u, nil
}

func (fakeRepo *FakeUserRepo) UserExistByUsername(ctx context.Context, username string) (bool, error) {
	return false, nil
}

func (fakeRepo *FakeUserRepo) CreateUser(ctx context.Context, pengguna user.Pengguna) (user.ID, error) {
	return 0, nil
}

func (fakeRepo *FakeUserRepo) UpdateUser(ctx context.Context, idPengguna user.ID, pengguna updatepatch.Pengguna) error {
	return nil
}

func (fakeRepo *FakeUserRepo) DeleteUser(ctx context.Context, id user.ID) error {
	return nil
}

func (fakeRepo *FakeUserRepo) DeleteUsers(ctx context.Context, ids []user.ID) (int64, error) {
	return int64(len(ids)), nil
}

func (fakeRepo *FakeUserRepo) ListUser(ctx context.Context) ([]user.Pengguna, error) {
	return []user.Pengguna{}, nil
}

func (fakeHash *FakeHasher) ComparePaswordAndHashed(hash string, plain string) bool {
	return fakeHash.Ok
}

func (fakeHash *FakeHasher) GenerateHash(plain string) (string, error) {
	return "hashed password", nil
}

func (fakeSession *FakeSessionRepo) CreateSession(ctx context.Context, userID user.ID, expiresAt time.Time) (string, error) {
	if fakeSession.CreateSessionErr != nil {
		return "", fakeSession.CreateSessionErr
	}
	if fakeSession.Store == nil {
		fakeSession.Store = map[string]session.Session{}
	}
	if fakeSession.SessionNextID == "" {
		fakeSession.SessionNextID = "session_user_1"
	}
	id := fakeSession.SessionNextID
	fakeSession.Store[id] = session.Session{
		SessionID: id,
		UserID:    userID,
		Revoked:   false,
		ExpiresAt: expiresAt,
	}
	return id, nil
}

func (fakeSession *FakeSessionRepo) GetSession(ctx context.Context, sessionID string) (session.Session, error) {
	ss, ok := fakeSession.Store[sessionID]
	if !ok {
		return session.Session{}, coreerror.ErrNotFound
	}
	return ss, nil
}

func (fakeSession *FakeSessionRepo) RevokeSession(ctx context.Context, sessionID string) error {
	if fakeSession.RevokeSessionErr != nil {
		return fakeSession.RevokeSessionErr
	}
	ss, ok := fakeSession.Store[sessionID]
	if !ok {
		return coreerror.ErrNotFound
	}
	ss.Revoked = true
	fakeSession.Store[sessionID] = ss
	return nil
}

func (fakeSession *FakeSessionRepo) RevokeSessionAllbyUser(ctx context.Context, userID user.ID, now time.Time) error {
	for ssid, sess := range fakeSession.Store {
		if sess.UserID == userID && !sess.Revoked && now.Before(sess.ExpiresAt) {
			sess.Revoked = true
			fakeSession.Store[ssid] = sess
		}
	}
	return nil
}

func (fakeSession *FakeSessionRepo) GetSessionByUserId(ctx context.Context, userId user.ID) (session.Session, error) {
	for _, sess := range fakeSession.Store {
		if sess.UserID == userId {
			return sess, nil
		}
	}
	return session.Session{}, coreerror.ErrNotFound
}

func (fakeSession *FakeSessionRepo) GetAllActiveSession(ctx context.Context) ([]session.Session, error) {
	var sessions []session.Session
	for _, sess := range fakeSession.Store {
		if !sess.Revoked {
			sessions = append(sessions, sess)
		}
	}
	return sessions, nil
}

func (fakeSession *FakeSessionRepo) HasActiveSession(ctx context.Context, userID user.ID) (bool, error) {
	if fakeSession.HasActiveErr != nil {
		return false, fakeSession.HasActiveErr
	}
	for _, sess := range fakeSession.Store {
		if sess.UserID == userID && !sess.Revoked {
			return true, nil
		}
	}
	return false, nil
}

// NEW: RevokeExpiredSessions
// Meniru query DB:
//   UPDATE sessions SET revoked_at=NOW()
//   WHERE id_pengguna=$1 AND revoked_at IS NULL AND expires_at <= NOW()
//
// Di fake kita pakai:
// - revoked_at diganti dengan Revoked bool
// - NOW() diganti time.Now() (karena signature tidak menerima now)
func (fakeSession *FakeSessionRepo) RevokeExpiredSessions(ctx context.Context, userID user.ID) (bool, error) {
	if fakeSession.RevokeExpiredErr != nil {
		return false, fakeSession.RevokeExpiredErr
	}
	if fakeSession.Store == nil {
		return false, nil
	}

	now := time.Now()
	changed := false

	for id, sess := range fakeSession.Store {
		if sess.UserID != userID {
			continue
		}
		if sess.Revoked {
			continue
		}
		if sess.ExpiresAt.After(now) {
			continue // masih aktif
		}

		sess.Revoked = true
		fakeSession.Store[id] = sess
		changed = true
	}

	return changed, nil
}

func (f *FakeAccessToken) GenerateAccessToken(userID user.ID, role user.Role, username string, tokenDuration time.Duration) (string, error) {
	if f.GenerateAccessTokenErr != nil {
		return "", f.GenerateAccessTokenErr
	}
	return fmt.Sprintf("ACCESS TOKEN :|%v|%v|%s", userID, role, username), nil
}

func (f *FakeAccessToken) VerifyAccessToken(token string, now time.Time) (userID user.ID, role user.Role, username string, err error) {
	const prefix = "ACCESS TOKEN :|"
	if !strings.HasPrefix(token, prefix) {
		return 0, "", "", coreerror.ErrInvalidToken
	}

	parts := strings.Split(strings.TrimPrefix(token, prefix), "|")
	if len(parts) != 3 {
		return 0, "", "", coreerror.ErrInvalidToken
	}

	id, convErr := strconv.Atoi(parts[0])
	if convErr != nil {
		return 0, "", "", coreerror.ErrInvalidToken
	}
	return user.ID(id), user.Role(parts[1]), parts[2], nil
}

func (f *FakeRefreshToken) GenerateRefreshToken(sessionID string, tokenDuration time.Duration) (string, error) {
	if f.GenerateRefreshTokenErr != nil {
		return "", f.GenerateRefreshTokenErr
	}
	return fmt.Sprintf("REFRESH TOKEN : %s", sessionID), nil
}

func (f *FakeRefreshToken) VerifyRefreshToken(token string, now time.Time) (sessionID string, err error) {
	if f.VerifyRefreshTokenErr != nil {
		return "", f.VerifyRefreshTokenErr
	}
	const prefix = "REFRESH TOKEN : "
	if !strings.HasPrefix(token, prefix) {
		return "", coreerror.ErrInvalidToken
	}

	sessionID = strings.TrimPrefix(token, prefix)
	if sessionID == "" {
		return "", coreerror.ErrInvalidToken
	}
	return sessionID, nil
}
