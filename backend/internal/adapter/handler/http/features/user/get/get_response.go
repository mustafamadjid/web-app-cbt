package httpx

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

type GuruResponseItem struct {
	IdPengguna   user.ID         `json:"id_pengguna"`
	Role         user.Role       `json:"role"`
	Username     string          `json:"username"`
	NoHp         string          `json:"no_hp"`
	Email        user.Email      `json:"email"`
	NamaLengkap  string          `json:"nama_lengkap"`
	StatusAkun   user.StatusAkun `json:"status_akun"`
	Nip          user.NIP        `json:"nip"`
	Jabatan      string          `json:"jabatan"`
	BidangStudi  string          `json:"bidang_studi"`
	Foto         string          `json:"foto_profil"`
	JenisKelamin string          `json:"jenis_kelamin"`
}

type SiswaResponseItem struct {
	IdPengguna   user.ID         `json:"id_pengguna"`
	Role         user.Role       `json:"role"`
	Username     string          `json:"username"`
	NoHp         string          `json:"no_hp"`
	Email        user.Email      `json:"email"`
	NamaLengkap  string          `json:"nama_lengkap"`
	StatusAkun   user.StatusAkun `json:"status_akun"`
	Nisn         string          `json:"nisn"`
	NoAbsen      int             `json:"no_absen"`
	Angkatan     int             `json:"angkatan"`
	TempatLahir  string          `json:"tempat_lahir"`
	TanggalLahir string          `json:"tanggal_lahir"`
	NamaKelas    string          `json:"nama_kelas"`
	TingkatKelas int             `json:"tingkat_kelas"`
	Foto         string          `json:"foto_profil"`
	JenisKelamin string          `json:"jenis_kelamin"`
}
