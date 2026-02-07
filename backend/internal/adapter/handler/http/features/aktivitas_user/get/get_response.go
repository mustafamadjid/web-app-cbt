package httpx

import (
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type AktivitasUserResponse struct {
	IdAktivitas aktivitas_user.AktivitasID `json:"id_aktivitas"`
	IdPengguna  user.ID                    `json:"id_pengguna"`
	Username    string                     `json:"username"`
	Role        user.Role                  `json:"role"`
	Action      aktivitas_user.Action      `json:"action"`
	Description string                     `json:"description"`
	IpAddress   string                     `json:"ip_address"`
	CreatedAt   string                     `json:"created_at"`
}
