package siswaujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ujian"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

type ListUjianSiswaService struct {
	repo ujian_repo.UjianSiswaRepository
}

func NewListUjianSiswaService(repo ujian_repo.UjianSiswaRepository) *ListUjianSiswaService {
	return &ListUjianSiswaService{
		repo: repo,
	}
}

func (s *ListUjianSiswaService) ListUjianSiswa(ctx context.Context, idSiswa int, filter query.ListUjianFilter) ([]ujian.ListUjian, error) {
	logger := corelog.FromContext(ctx)

	if err := validateListUjianSiswaID(idSiswa); err != nil {
		logger.Error(ctx, "failed list ujian siswa", "layer", "core.service", "op", "siswa_ujian.list", "err", err)
		return nil, err
	}

	var err error
	filter, err = sanitizeAndValidateListUjianSiswaFilter(filter)
	if err != nil {
		logger.Error(ctx, "failed list ujian siswa", "layer", "core.service", "op", "siswa_ujian.list.filter", "err", err)
		return nil, coreerror.ErrInvalidInput
	}

	items, err := s.repo.ListUjianSiswa(ctx, idSiswa, filter)
	if err != nil {
		logger.Error(ctx, "failed list ujian siswa", "layer", "core.service", "op", "siswa_ujian.list", "err", err)
		return nil, err
	}

	return items, nil
}
