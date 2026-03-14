package siswaujian_service

import (
	"context"
	"math/rand/v2"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
)

type ListSoalUjianSiswaService struct {
	repo ujian_repo.SoalUjianRepository
}

func NewListSoalUjianSiswaService(repo ujian_repo.SoalUjianRepository) *ListSoalUjianSiswaService {
	return &ListSoalUjianSiswaService{
		repo: repo,
	}
}

func (r *ListSoalUjianSiswaService) ListSoalUjianSiswa(ctx context.Context, idJadwalUjian ujian.ID) ([]ujian.SoalUjianSiswa, error) {
	logger := corelog.FromContext(ctx)

	if idJadwalUjian <= 0 {
		logger.Error(ctx, "failed list soal ujian", "layer", "core.service", "op", "ujian.list_soal", "err", coreerror.ErrMissingId)
		return nil, coreerror.ErrMissingId
	}

	soal, acakSoal, err := r.repo.GetSoalUjianByBankSoalForSiswa(ctx, idJadwalUjian)
	if err != nil {
		logger.Error(ctx, "failed list soal ujian", "layer", "core.service", "op", "ujian.list_soal", "err", err)
		return nil, err
	}

	lastElement := len(soal) - 1

	if acakSoal && len(soal) > 1 {
		for i := lastElement; i > 0; i-- {
			j := rand.IntN(i + 1)
			soal[i], soal[j] = soal[j], soal[i]
		}

		return soal, nil
	}

	return soal, nil
}
