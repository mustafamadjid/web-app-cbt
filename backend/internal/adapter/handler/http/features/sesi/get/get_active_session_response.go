package httpx

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

type ActiveSessionResponse struct {
	SessionID  string    `json:"session_id"`
	IdPengguna user.ID   `json:"id_pengguna"`
	Role       user.Role `json:"role"`
	Revoked    bool      `json:"revoked"`
	ExpiresAt  string    `json:"expires_at"`
}

type ActiveSessionUserResponse struct {
	IdPengguna   user.ID         `json:"id_pengguna"`
	Username     string          `json:"username"`
	Email        *user.Email     `json:"email"`
	NamaLengkap  string          `json:"nama_lengkap"`
	JenisKelamin string          `json:"jenis_kelamin"`
	NoHp         *string         `json:"no_hp"`
	Role         user.Role       `json:"role"`
	StatusAkun   user.StatusAkun `json:"status_akun"`
	Foto         string          `json:"foto_profil"`
}

type ActiveLoginSessionResponse struct {
	Session  ActiveSessionResponse     `json:"session"`
	Pengguna ActiveSessionUserResponse `json:"pengguna"`
}

type ListActiveLoginSessionResponse struct {
	Items []ActiveLoginSessionResponse `json:"items"`
}
