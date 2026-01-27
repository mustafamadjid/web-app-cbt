package user

type ID int
type Role string
type StatusAkun string

const (
	ADMIN Role = "ADMIN"
	GURU Role = "GURU"
	SISWA Role = "SISWA"
)

const (
	AKTIF StatusAkun = "AKTIF"
	NONAKTIF StatusAkun = "NONAKTIF"
)

type Pengguna struct {
	ID ID
	Username string
	Email string
	PasswordHashed string
	NamaLengkap string
	JenisKelamin string
	NoHp string
	Role Role
	StatusAkun StatusAkun
	Foto string
}

type ProfilSiswa struct{
	ID ID
	IdPengguna ID
	IdTingkatKelas ID
	IdNamaKelas ID
	nisn string
	NoAbsen int
	Angkatan int
	TempatLahir string
	TanggalLahir string
}
type ProfilGuru struct{
	ID ID
	IdPengguna ID
	Nip string
	Jabatan string
	BidangStudi string
}


func (role Role) ValidRole() bool {
	switch role {
	case ADMIN, GURU, SISWA:
		return true
	default:
		return false
	}
}