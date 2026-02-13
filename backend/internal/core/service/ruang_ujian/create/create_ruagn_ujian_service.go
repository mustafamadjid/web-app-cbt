package ruangujian_service

import (
	"context"

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
}
