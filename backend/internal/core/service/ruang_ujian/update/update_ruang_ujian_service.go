package ruangujian_service

import (
	"context"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ruangujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ruang_ujian"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type UpdateRuangUjianRepo struct {
	ruangRepo ruangujian_repo.RuangUjianRepo
}

func NewUpdateRuangUjianService(ruangRepo ruangujian_repo.RuangUjianRepo) *UpdateRuangUjianRepo {
	return &UpdateRuangUjianRepo{ruangRepo: ruangRepo}
}

func(r *UpdateRuangUjianRepo)UpdateRuangUjian(ctx context.Context, idRuangan int, ruangUjian updatepatch.UpdateRuangUjianPatch) error{
	logger := corelog.FromContext(ctx)

	

	if idRuangan <= 0 {
		logger.Error(ctx,"failed update ruang ujian","layer","core.service","op","ruangujian.update","err",coreerror.ErrMissingId)
		return coreerror.ErrMissingId
	}

	if ruangUjian.KodeRuang != nil {
		s := strings.TrimSpace(*ruangUjian.KodeRuang)
		if s == ""  {
			logger.Error(ctx, "failed updating ruangujian", "layer", "core.service", "op", "ruangujian.update_ruang_ujian.KodeRuang", "err", coreerror.ErrMissingField)
		} 

		s = strings.ToUpper(s)
		ruangUjian.KodeRuang = &s
	}

	if ruangUjian.NamaRuang != nil {
		s := strings.TrimSpace(*ruangUjian.NamaRuang)
		if s == ""  {
			logger.Error(ctx, "failed updating ruangujian", "layer", "core.service", "op", "ruangujian.update_ruang_ujian.NamaRuangan", "err", coreerror.ErrMissingField)
		} 
		ruangUjian.NamaRuang = &s
	}

	err := r.ruangRepo.UpdateRuangUjian(ctx,idRuangan,ruangUjian)
	if err != nil {
		logger.Error(ctx,"failed update ruang ujian","layer","core.service","op","ruangujian.update","err",err)
		return err
	}
	return nil
}