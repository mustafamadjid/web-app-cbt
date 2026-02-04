package aktivitas_user_service

import (

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type AktivitasUserCmd struct {
	IdPengguna  user.ID
	Action      aktivitas_user.Action
	Description string
	IpAddress   string
}
