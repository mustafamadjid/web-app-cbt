package ruangujian_service

import (
	"context"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ruangujian_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/ruang_ujian"
)

type CreateRuangUjianService struct {
	ruangRepo ruangujian_repo.RuangUjianRepo
}

func NewRuangUjianService(ruangRepo ruangujian_repo.RuangUjianRepo) *CreateRuangUjianService {
	return &CreateRuangUjianService{
		ruangRepo: ruangRepo,
	}
}

func(r *CreateRuangUjianService)CreateRuangUjianService(ctx context.Context, ruangUjian ruangujian.RuangUjian)error{
	logger := corelog.FromContext(ctx)

	ruangUjian.KodeRuang = strings.TrimSpace(ruangUjian.KodeRuang)
	ruangUjian.KodeRuang = strings.ToUpper(ruangUjian.KodeRuang)

	ruangUjian.NamaRuangan = strings.TrimSpace(ruangUjian.NamaRuangan)

	exist, err := r.ruangRepo.ExistByKodeRuang(ctx,ruangUjian.KodeRuang)
	if err != nil {
		logger.Error(ctx, "failed creating ruang ujian", "layer", "core.service", "op", "ruangujian.create_ruangujian.existByKodeRuang", "err", err)
		return err
	}

	if exist {
		return coreerror.ErrKodeRuangUjianExist
	}

	if err := r.ruangRepo.CreateRuangUjian(ctx,ruangUjian); err != nil {
		logger.Error(ctx, "failed creating ruang ujian", "layer", "core.service", "op", "ruangujian.create_ruangujian", "err", err)
		return err
	}

	return nil
}
