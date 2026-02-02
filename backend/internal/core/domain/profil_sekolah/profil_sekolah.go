package profil_sekolah

import "time"

type IDProfil int

type ProfilSekolah struct {
	IDProfil      IDProfil
	EmailSekolah  string
	NoTelpSekolah string
	KepalaSekolah string
	WakaSekolah   string
	NamaSekolah   string
	AlamatSekolah string
	LogoSekolah   *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
