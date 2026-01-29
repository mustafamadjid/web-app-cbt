package httpx

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

type CreateGuruResponse struct {
	IdPengguna   user.ID `json:"id_pengguna"`
	IdProfilGuru user.ID `json:"id_profil_guru"`
}

type CreateSiswaResponse struct {
	IdPengguna    user.ID `json:"id_pengguna"`
	IdProfilSiswa user.ID `json:"id_profil_siswa"`
}
