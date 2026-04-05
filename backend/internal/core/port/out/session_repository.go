package out

import (
	"context"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type SessionRepository interface {
	GetSession(ctx context.Context, sessionID string) (session.Session, error)
	GetSessionByUserId(ctx context.Context, userId user.ID) (session.Session, error)
	GetAllActiveSession(ctx context.Context) ([]session.SessionWithUser, error)

	CreateSession(ctx context.Context, userID user.ID, role user.Role, expiresAt time.Time) (string, error)
	RevokeSession(ctx context.Context, sessionID string) error
	RevokeSessionAllbyUser(ctx context.Context, userID user.ID) error
	RevokeExpiredSessions(ctx context.Context, userID user.ID) (bool, error)

	HasActiveSession(ctx context.Context, userID user.ID) (bool, error)
}
