package kelas_service

import (
	"context"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	kelas_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/kelas"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
)

type GetKelasService struct {
	kelasSvc kelas_repo.KelasRepository
}

func NewGetKelasService(kelasSvc kelas_repo.KelasRepository) *GetKelasService {
	return &GetKelasService{
		kelasSvc: kelasSvc,
	}
}

func (s *GetKelasService) GetFullKelas(ctx context.Context, filter query.ListKelasFilter) ([]kelas.FullKelasData, error) {
	logger := corelog.FromContext(ctx)

	filter = sanitizeListKelasFilter(filter)

	items, err := s.kelasSvc.GetKelas(ctx, filter)
	if err != nil {
		logger.Error(ctx, "failed get kelas", "layer", "core.service", "op", "kelas.get", "err", err)
		return nil, err
	}
	return items, nil
}

func (s *GetKelasService) GetKelasById(ctx context.Context, idTingkatKelas int, idNamaKelas int) (kelas.KelasData, error) {
	return s.kelasSvc.GetKelasById(ctx, idTingkatKelas, idNamaKelas)
}
