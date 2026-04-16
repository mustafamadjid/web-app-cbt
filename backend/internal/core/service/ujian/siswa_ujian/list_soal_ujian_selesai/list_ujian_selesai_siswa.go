package siswaujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
)

type ListUjianSelesaiSiswaService struct {
	repo ujian_repo.ListUjianRepository
}

func NewListUjianSelesaiSiswaService(repo ujian_repo.ListUjianRepository) *ListUjianSelesaiSiswaService {
	return &ListUjianSelesaiSiswaService{repo: repo}
}

func (s *ListUjianSelesaiSiswaService) ListUjianSelesaiSiswa(ctx context.Context, idSiswa int) ([]ujian.ListUjian, error) {
	logger := corelog.FromContext(ctx)

	if idSiswa <= 0 {
		logger.Error(ctx, "invalid id siswa", "layer", "core.service", "op", "siswa_ujian.list_ujian_selesai", "err", coreerror.ErrMissingId)
		return nil, coreerror.ErrMissingId
	}

	items, err := s.repo.GetAllUjianSubmittedByIdSiswa(ctx, idSiswa)
	if err != nil {
		logger.Error(ctx, "failed list ujian selesai siswa", "layer", "core.service", "op", "siswa_ujian.list_ujian_selesai", "err", err)
		return nil, err
	}

	return items, nil
}
