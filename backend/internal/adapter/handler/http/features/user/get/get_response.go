package httpx

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

type GuruResponseItem struct {
	IdPengguna   user.ID          `json:"id_pengguna"`
	Role         user.Role        `json:"role"`
	Username     string           `json:"username"`
	NoHp         string           `json:"no_hp"`
	Email        user.Email       `json:"email"`
	NamaLengkap  string           `json:"nama_lengkap"`
	StatusAkun   user.StatusAkun  `json:"status_akun"`
	Nip          user.NIP         `json:"nip"`
	Jabatan      string           `json:"jabatan"`
	BidangStudi  string           `json:"bidang_studi"`
	Foto         string           `json:"foto_profil"`
	JenisKelamin string           `json:"jenis_kelamin"`
}