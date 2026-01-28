package user_service

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

type CreateGuruRes struct {
	IdPengguna   user.ID
	IdProfilGuru user.ID
}

type CreateSiswaRes struct {
	IdPengguna    user.ID
	IdProfilSiswa user.ID
}
