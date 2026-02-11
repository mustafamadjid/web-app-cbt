package updatepatch

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

type Pengguna struct {
	NamaLengkap  *string
	Email        *user.Email
	NoHp         *string
	Foto         *string
	StatusAkun   *user.StatusAkun
	Role         *user.Role
	JenisKelamin *string
}
