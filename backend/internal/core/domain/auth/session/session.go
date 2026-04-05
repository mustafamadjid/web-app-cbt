package session

import (
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type Session struct {
	SessionID string
	UserID    user.ID
	Role      user.Role
	Revoked   bool
	ExpiresAt time.Time
}
