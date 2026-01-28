package user_service

import (
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

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

type UpdateSiswaCmd struct {
	IdPengguna user.ID

	Username     *string
	Email        *string
	NamaLengkap  *string
	JenisKelamin *string
	NoHp         *string
	Foto         *string
	StatusAkun   *user.StatusAkun
	Role         *user.Role

	IdTingkatKelas *user.ID
	IdNamaKelas    *user.ID
	Nisn           *string
	NoAbsen        *int
	Angkatan       *int
	TempatLahir    *string
	TanggalLahir   *time.Time
}
