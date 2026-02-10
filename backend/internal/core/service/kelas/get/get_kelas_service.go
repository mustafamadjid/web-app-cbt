package kelas_service

import (
	"context"
	"strings"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"

	kelas_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/kelas"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type GetKelasService struct {
	kelasSvc kelas_repo.KelasRepository
}

func NewGetKelasService(kelasSvc kelas_repo.KelasRepository) *GetKelasService {
	return &GetKelasService{
		kelasSvc: kelasSvc,
	}
}

func (s *GetKelasService)GetFullKelas(ctx context.Context, filter query.ListKelasFilter) ([]kelas.FullKelasData, error) {
	logger := corelog.FromContext(ctx)

	filter.Search = strings.TrimSpace(filter.Search)

	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	if filter.Limit > 50 {
		filter.Limit = 50
	}

	if filter.Offset < 0 {
		filter.Offset = 0
	}

	items, err := s.kelasSvc.GetKelas(ctx,filter)
	if err != nil {
		logger.Error(ctx, "failed get kelas", "layer", "core.service", "op", "kelas.get", "err", err)
		return nil, err
	}
	return items, nil
}
