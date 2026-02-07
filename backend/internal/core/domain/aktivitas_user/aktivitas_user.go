package aktivitas_user

import (
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type AktivitasID string
type Action string

const (
	LOGIN  Action = "LOGIN"
	LOGOUT Action = "LOGOUT"
	CREATE Action = "CREATE"
	UPDATE Action = "UPDATE"
	DELETE Action = "DELETE"
)

type AktivitasUser struct {
	IdAktivitas AktivitasID
	IdPengguna  user.ID
	Username    string
	Role        user.Role
	Action      Action
	Description string
	IpAddress   string
	CreatedAt   time.Time
}

func (action Action) ValidAction() bool {
	switch action {
	case LOGIN, LOGOUT, CREATE, UPDATE, DELETE:
		return true
	default:
		return false
	}
}

func ValidIpAddress(ip string) bool {
	if ip == "" {
		return false
	}

	if len(ip) > 15 {
		return false
	}

	return true
}
