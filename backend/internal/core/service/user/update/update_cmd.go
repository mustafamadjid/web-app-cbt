package user_service

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

type UpdateGuruCmd struct {
	IdPengguna user.ID

	Username     *string
	Email        *string
	NamaLengkap  *string
	JenisKelamin *string
	NoHp         *string
	Foto         *string
	StatusAkun   *user.StatusAkun
	Role         *user.Role

	Nip         *string
	Jabatan     *string
	BidangStudi *string
}
