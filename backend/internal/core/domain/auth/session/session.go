package session

import (
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type Session struct {
	SessionID string
	UserID    user.ID
	Revoked   bool
	ExpiresAt time.Time
}