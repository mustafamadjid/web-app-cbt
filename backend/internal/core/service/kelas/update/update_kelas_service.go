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

	if err := validateUpdateNamaKelasID(idNamaKelas); err != nil {
		return err
	}

	if err := sanitizeNamaKelasPatch(&dataUpdate); err != nil {
		logger.Error(ctx, "failed updating nama kelas", "layer", "core.service", "op", "kelas.update_nama_kelas.NamaKelas", "err", err)
		return err
	}

	if err := validateIdTingkatKelasPatch(dataUpdate); err != nil {
		logger.Error(ctx, "failed updating nama kelas", "layer", "core.service", "op", "kelas.update_nama_kelas.IdTingkatKelas", "err", err)
		return err
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

// -----------------------
// Sanitizer and validator
// -----------------------

func validateUpdateNamaKelasID(idNamaKelas int) error {
	if idNamaKelas <= 0 {
		return coreerror.ErrMissingId
	}

	return nil
}

func sanitizeNamaKelasPatch(dataUpdate *updatepatch.NamaKelasPatch) error {
	if dataUpdate.NamaKelas != nil {
		trimmedNamaKelas := strings.TrimSpace(*dataUpdate.NamaKelas)
		if trimmedNamaKelas == "" {
			return coreerror.ErrMissingField
		}
		dataUpdate.NamaKelas = &trimmedNamaKelas
	}

	return nil
}

func validateIdTingkatKelasPatch(dataUpdate updatepatch.NamaKelasPatch) error {
	if dataUpdate.IdTingkatKelas != nil && *dataUpdate.IdTingkatKelas == 0 {
		return coreerror.ErrMissingField
	}

	return nil
}
