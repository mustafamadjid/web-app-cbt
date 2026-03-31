package user_service_test

import (
	"context"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type FakeSessionRepo struct {
	RevokeAllErr error

	RevokeAllCalled bool
	LastUserID      user.ID
}

func (r *FakeSessionRepo) GetSession(ctx context.Context, sessionID string) (session.Session, error) {
	return session.Session{}, nil
}

func (r *FakeSessionRepo) GetSessionByUserId(ctx context.Context, userID user.ID) (session.Session, error) {
	return session.Session{}, nil
}

func (r *FakeSessionRepo) GetAllActiveSession(ctx context.Context) ([]session.Session, error) {
	return nil, nil
}

func (r *FakeSessionRepo) CreateSession(ctx context.Context, userID user.ID, expiresAt time.Time) (string, error) {
	return "", nil
}

func (r *FakeSessionRepo) RevokeSession(ctx context.Context, sessionID string) error {
	return nil
}

func (r *FakeSessionRepo) RevokeSessionAllbyUser(ctx context.Context, userID user.ID) error {
	r.RevokeAllCalled = true
	r.LastUserID = userID
	if r.RevokeAllErr != nil {
		return r.RevokeAllErr
	}
	return nil
}

func (r *FakeSessionRepo) RevokeExpiredSessions(ctx context.Context, userID user.ID) (bool, error) {
	return false, nil
}

func (r *FakeSessionRepo) HasActiveSession(ctx context.Context, userID user.ID) (bool, error) {
	return false, nil
}
