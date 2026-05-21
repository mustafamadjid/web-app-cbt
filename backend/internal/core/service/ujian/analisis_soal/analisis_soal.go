package analisissoal_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	analisissoal_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian/analisis_soal"
)

type AnalisisSoalService struct {
	repo analisissoal_repo.AnalisisSoalInterface
}

type AnalisisSoalrepo = AnalisisSoalService

func NewAnalisisSoalService(repo analisissoal_repo.AnalisisSoalInterface) *AnalisisSoalService {
	return &AnalisisSoalService{repo: repo}
}

func (r *AnalisisSoalService) GetListAnalisisSoal(ctx context.Context, idJadwalUjian int) ([]ujian.AnalisisSoal, error) {
	if idJadwalUjian == 0 {
		return nil, coreerror.ErrMissingId
	}

	data, err := r.repo.GetListAnalisisSoal(ctx, ujian.ID(idJadwalUjian))
	if err != nil {
		return nil, err
	}
	return data, nil
}
