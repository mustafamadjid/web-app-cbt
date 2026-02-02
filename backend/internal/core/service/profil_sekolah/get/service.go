package profil_sekolah_service

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/profil_sekolah"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/profil_sekolah"
)

type GetProfilSekolahService struct {
	repo out.ProfilSekolahRepository
}

func NewGetProfilSekolahService(repo out.ProfilSekolahRepository) *GetProfilSekolahService {
	return &GetProfilSekolahService{repo: repo}
}

func (s *GetProfilSekolahService) GetProfilSekolah(ctx context.Context) (profil_sekolah.ProfilSekolah, error) {
	return s.repo.GetProfilSekolah(ctx)
}
