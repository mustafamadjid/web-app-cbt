package ruangujian_service

import (
	"context"
	"errors"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ruangujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ruang_ujian"
)

type DeleteRuangUjianService struct {
	ruangRepo ruangujian_repo.RuangUjianRepo
}

func NewDeleteRuangUjianService(ruangRepo ruangujian_repo.RuangUjianRepo) *DeleteRuangUjianService {
	return &DeleteRuangUjianService{ruangRepo: ruangRepo}
}

func(r *DeleteRuangUjianService)DeleteRuangUjian(ctx context.Context, idRuangan int) error {
	logger := corelog.FromContext(ctx)

	if idRuangan <= 0 {
		logger.Error(ctx,"failed delete ruang ujian","layer","core.service","op","ruangujian.delete","err",coreerror.ErrMissingId)
		return coreerror.ErrMissingId
	}

	if err := r.ruangRepo.DeleteRuangUjian(ctx,idRuangan); err != nil {
		logger.Error(ctx,"failed delete ruang ujian","layer","core.service","op","ruangujian.delete","err",err)
		switch {
		case errors.Is(err, coreerror.ErrDeleteRestricted):
			return coreerror.ErrDeleteRestricted
		default:
			return err
		}
	}
	return nil
}