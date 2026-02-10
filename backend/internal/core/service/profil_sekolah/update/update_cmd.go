package profil_sekolah_service

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/profil_sekolah"

type UpdateProfilSekolahCmd struct {
	IDProfil      profil_sekolah.IDProfil
	EmailSekolah  *string
	NoTelpSekolah *string
	KepalaSekolah *string
	WakaSekolah   *string
	NamaSekolah   *string
	AlamatSekolah *string
	LogoSekolah   *string
}