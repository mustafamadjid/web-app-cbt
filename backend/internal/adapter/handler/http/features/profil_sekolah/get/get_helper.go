package httpx

import (
	profilsekolahresponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/profil_sekolah"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/profil_sekolah"
)

func toProfilSekolahResponse(profil profil_sekolah.ProfilSekolah) profilsekolahresponse.ProfilSekolahResponse {
	return profilsekolahresponse.ProfilSekolahResponse{
		IDProfil:      int(profil.IDProfil),
		EmailSekolah:  profil.EmailSekolah,
		NoTelpSekolah: profil.NoTelpSekolah,
		KepalaSekolah: profil.KepalaSekolah,
		WakaSekolah:   profil.WakaSekolah,
		NamaSekolah:   profil.NamaSekolah,
		AlamatSekolah: profil.AlamatSekolah,
		LogoSekolah:   profil.LogoSekolah,
		CreatedAt:     httphelper.FormatRFC3339(profil.CreatedAt),
		UpdatedAt:     httphelper.FormatRFC3339(profil.UpdatedAt),
	}
}
