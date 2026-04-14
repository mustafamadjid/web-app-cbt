package tests

import (
	"context"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	logport "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	ratelimiterport "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/rate_limiter"
)

type fakeAccessTokenService struct {
	token    string
	userID   user.ID
	role     user.Role
	username string
	err      error
}

func (f *fakeAccessTokenService) GenerateAccessToken(user.ID, user.Role, string, time.Duration) (string, error) {
	return "", nil
}

func (f *fakeAccessTokenService) VerifyAccessToken(token string, now time.Time) (user.ID, user.Role, string, error) {
	f.token = token
	return f.userID, f.role, f.username, f.err
}

type fakeRefreshTokenService struct{}

func (f *fakeRefreshTokenService) GenerateRefreshToken(string, time.Duration) (string, error) {
	return "", nil
}

func (f *fakeRefreshTokenService) VerifyRefreshToken(string, time.Time) (string, error) {
	return "", nil
}

var _ out.AccessTokenService = (*fakeAccessTokenService)(nil)
var _ out.RefreshTokenService = (*fakeRefreshTokenService)(nil)

type fakeSessionRepository struct {
	hasActive bool
	err       error
	userID    user.ID
	calls     int
}

func (f *fakeSessionRepository) GetSession(context.Context, string) (session.Session, error) {
	return session.Session{}, nil
}

func (f *fakeSessionRepository) GetSessionByUserId(context.Context, user.ID) (session.Session, error) {
	return session.Session{}, nil
}

func (f *fakeSessionRepository) GetAllActiveSession(context.Context) ([]session.SessionWithUser, error) {
	return nil, nil
}

func (f *fakeSessionRepository) CreateSession(context.Context, user.ID, user.Role, time.Time) (string, error) {
	return "", nil
}

func (f *fakeSessionRepository) RevokeSession(context.Context, string) error {
	return nil
}

func (f *fakeSessionRepository) RevokeSessionAllbyUser(context.Context, user.ID) error {
	return nil
}

func (f *fakeSessionRepository) RevokeExpiredSessions(context.Context, user.ID) (bool, error) {
	return false, nil
}

func (f *fakeSessionRepository) HasActiveSession(_ context.Context, userID user.ID) (bool, error) {
	f.calls++
	f.userID = userID
	return f.hasActive, f.err
}

var _ out.SessionRepository = (*fakeSessionRepository)(nil)

type fakeRateLimiter struct {
	key        string
	allowed    bool
	retryAfter time.Duration
	err        error
	calls      int
}

func (f *fakeRateLimiter) Allow(_ context.Context, key string) (bool, time.Duration, error) {
	f.key = key
	f.calls++
	return f.allowed, f.retryAfter, f.err
}

var _ ratelimiterport.RateLimiter = (*fakeRateLimiter)(nil)

type capturedLogCall struct {
	msg   string
	attrs []any
}

type capturedLogger struct {
	withAttrs  []any
	infoCalls  []capturedLogCall
	errorCalls []capturedLogCall
}

func (l *capturedLogger) With(attrs ...any) logport.Logger {
	l.withAttrs = append([]any(nil), attrs...)
	return l
}

func (l *capturedLogger) Info(_ context.Context, msg string, attrs ...any) {
	l.infoCalls = append(l.infoCalls, capturedLogCall{
		msg:   msg,
		attrs: append([]any(nil), attrs...),
	})
}

func (l *capturedLogger) Error(_ context.Context, msg string, attrs ...any) {
	l.errorCalls = append(l.errorCalls, capturedLogCall{
		msg:   msg,
		attrs: append([]any(nil), attrs...),
	})
}

var _ logport.Logger = (*capturedLogger)(nil)
