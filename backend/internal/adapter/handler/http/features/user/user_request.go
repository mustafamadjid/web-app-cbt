package httpx

import (
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type CreateGuruRequest struct {
	Username     string
	Email        string
	Password     string
	NamaLengkap  string
	JenisKelamin string
	NoHp         string
	Foto         string
	Nip          string
	Jabatan      string
	BidangStudi  string
}

type CreateSiswaRequest struct {
	Username       string
	Email          string
	Password       string
	NamaLengkap    string
	JenisKelamin   string
	NoHp           string
	Foto           string
	IdTingkatKelas user.ID
	IdNamaKelas    user.ID
	Nisn           string
	NoAbsen        int
	Angkatan       int
	TempatLahir    string
	TanggalLahir   time.Time
}

type UpdateGuruRequest struct {
	IdPengguna   user.ID
	Username     *string
	Email        *string
	NamaLengkap  *string
	JenisKelamin *string
	NoHp         *string
	Foto         *string
	StatusAkun   *string
	Role         *string
	Nip          *string
	Jabatan      *string
	BidangStudi  *string
}

type UpdateSiswaRequest struct {
	IdPengguna     user.ID
	Username       *string
	Email          *string
	NamaLengkap    *string
	JenisKelamin   *string
	NoHp           *string
	Foto           *string
	StatusAkun     *string
	Role           *string
	IdTingkatKelas *user.ID
	IdNamaKelas    *user.ID
	Nisn           *string
	NoAbsen        *int
	Angkatan       *int
	TempatLahir    *string
	TanggalLahir   *time.Time
}

type DeleteUserRequest struct {
	IdPengguna user.ID `json:"id_pengguna"`
}

type DeleteUsersRequest struct {
	Ids []int `json:"ids"`
}
