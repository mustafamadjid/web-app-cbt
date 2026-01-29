package httpx

import (
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

type CreateGuruRequest struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	NamaLengkap  string `json:"nama_lengkap"`
	JenisKelamin string `json:"jenis_kelamin"`
	NoHp         string `json:"no_hp"`
	Foto         string `json:"foto"`
	Nip          string `json:"nip"`
	Jabatan      string `json:"jabatan"`
	BidangStudi  string `json:"bidang_studi"`
}

type CreateSiswaRequest struct {
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	NamaLengkap    string    `json:"nama_lengkap"`
	JenisKelamin   string    `json:"jenis_kelamin"`
	NoHp           string    `json:"no_hp"`
	Foto           string    `json:"foto"`
	IdTingkatKelas user.ID   `json:"id_tingkat_kelas"`
	IdNamaKelas    user.ID   `json:"id_nama_kelas"`
	Nisn           string    `json:"nisn"`
	NoAbsen        int       `json:"no_absen"`
	Angkatan       int       `json:"angkatan"`
	TempatLahir    string    `json:"tempat_lahir"`
	TanggalLahir   time.Time `json:"tanggal_lahir"`
}

type UpdateGuruRequest struct {
	IdPengguna   user.ID `json:"id_pengguna"`
	Username     *string `json:"username,omitempty"`
	Email        *string `json:"email,omitempty"`
	NamaLengkap  *string `json:"nama_lengkap,omitempty"`
	JenisKelamin *string `json:"jenis_kelamin,omitempty"`
	NoHp         *string `json:"no_hp,omitempty"`
	Foto         *string `json:"foto,omitempty"`
	StatusAkun   *string `json:"status_akun,omitempty"`
	Role         *string `json:"role,omitempty"`
	Nip          *string `json:"nip,omitempty"`
	Jabatan      *string `json:"jabatan,omitempty"`
	BidangStudi  *string `json:"bidang_studi,omitempty"`
}

type UpdateSiswaRequest struct {
	IdPengguna     user.ID    `json:"id_pengguna"`
	Username       *string    `json:"username,omitempty"`
	Email          *string    `json:"email,omitempty"`
	NamaLengkap    *string    `json:"nama_lengkap,omitempty"`
	JenisKelamin   *string    `json:"jenis_kelamin,omitempty"`
	NoHp           *string    `json:"no_hp,omitempty"`
	Foto           *string    `json:"foto,omitempty"`
	StatusAkun     *string    `json:"status_akun,omitempty"`
	Role           *string    `json:"role,omitempty"`
	IdTingkatKelas *user.ID   `json:"id_tingkat_kelas,omitempty"`
	IdNamaKelas    *user.ID   `json:"id_nama_kelas,omitempty"`
	Nisn           *string    `json:"nisn,omitempty"`
	NoAbsen        *int       `json:"no_absen,omitempty"`
	Angkatan       *int       `json:"angkatan,omitempty"`
	TempatLahir    *string    `json:"tempat_lahir,omitempty"`
	TanggalLahir   *time.Time `json:"tanggal_lahir,omitempty"`
}

type DeleteUserRequest struct {
	IdPengguna user.ID `json:"id_pengguna"`
}
