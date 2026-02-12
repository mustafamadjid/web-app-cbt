package kelas_service

import (
	"context"
	"errors"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	kelas_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/kelas"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type UpdateKelasService struct {
	updateRepo kelas_repo.KelasRepository
}

func NewUpdateKelasService(kelasRepo kelas_repo.KelasRepository) *UpdateKelasService {
	return &UpdateKelasService{updateRepo: kelasRepo}
}

func (s *UpdateKelasService) UpdateNamaKelas(ctx context.Context, idNamaKelas int, dataUpdate updatepatch.NamaKelasPatch) error {
	logger := corelog.FromContext(ctx)

	if idNamaKelas <= 0 {
		return coreerror.ErrMissingId
	}

	if dataUpdate.NamaKelas != nil {
		s := strings.TrimSpace(*dataUpdate.NamaKelas)

		if s == "" {
			logger.Error(ctx, "failed updating nama kelas", "layer", "core.service", "op", "kelas.update_nama_kelas.NamaKelas", "err", coreerror.ErrMissingField)
			return coreerror.ErrMissingField
		}
		dataUpdate.NamaKelas = &s
	}

	if dataUpdate.IdTingkatKelas != nil {
		if *dataUpdate.IdTingkatKelas == 0 {
			logger.Error(ctx, "failed updating nama kelas", "layer", "core.service", "op", "kelas.update_nama_kelas.IdTingkatKelas", "err", coreerror.ErrMissingField)
			return coreerror.ErrMissingField
		}
	}

	if err := s.updateRepo.UpdateNamaKelas(ctx, idNamaKelas, dataUpdate); err != nil {
		logger.Error(ctx, "failed updating nama kelas", "layer", "core.service", "op", "kelas.update_nama_kela.UpdateNamaKelas", "err", err)

		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			return coreerror.ErrNotFound
		default:
			return err
		}
	}

	return nil
}
