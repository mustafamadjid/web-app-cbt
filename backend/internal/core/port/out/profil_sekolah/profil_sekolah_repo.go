package profil_sekolah

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/profil_sekolah"
)

type ProfilSekolahRepository interface {
	UpdateProfilSekolah(ctx context.Context, idProfil profil_sekolah.IDProfil, profil profil_sekolah.ProfilSekolah) error
	GetProfilSekolah(ctx context.Context) (profil_sekolah.ProfilSekolah, error)
}
