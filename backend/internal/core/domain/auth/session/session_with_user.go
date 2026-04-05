package session

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

type SessionWithUser struct {
	Session  Session
	Pengguna user.Pengguna
}
