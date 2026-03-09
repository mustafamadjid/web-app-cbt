package ujian_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

type DeleteUjianService struct {
	ujianRepo UjianRepository
}

func NewDeleteUjianService(ujianRepo UjianRepository) *DeleteUjianService {
	return &DeleteUjianService{
		ujianRepo: ujianRepo,
	}
}

func (r *DeleteUjianService) DeleteUjianService(ctx context.Context, id ujian.ID) error {
	logger := corelog.FromContext(ctx)

	if id <= 0 {
		logger.Error(ctx, "failed deleting ujian",
			"layer", "core.service",
			"op", "ujian.delete_ujian.DeleteUjianService",
			"err", coreerror.ErrMissingId,
		)
		return coreerror.ErrMissingId
	}

	if err := r.ujianRepo.DeleteUjian(ctx, id); err != nil {
		logger.Error(ctx, "failed deleting ujian",
			"layer", "core.service",
			"op", "ujian.delete_ujian.DeleteUjianService",
			"err", err,
		)
		return err
	}

	return nil
}
