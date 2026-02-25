package pengumuman_service

import (
	"context"
	"errors"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	delete_file_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/delete_file_system"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	pengumuman_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/pengumuman"
)

type DeletePengumumanService struct {
	repo pengumuman_repo.PengumumanRepo
	deleteFile delete_file_repo.DeleteFileRepo
}

func NewDeletePengumumanService(repo pengumuman_repo.PengumumanRepo, deleteFile delete_file_repo.DeleteFileRepo) *DeletePengumumanService {
	return &DeletePengumumanService{
		repo: repo,
		deleteFile: deleteFile,
	}
}

func(r *DeletePengumumanService)DeletePengumumanService(ctx context.Context,idPengumuman pengumuman.ID)error{
	logger := corelog.FromContext(ctx)

	if idPengumuman <= 0 {
		logger.Error(ctx, "failed deleting pengumuman", "layer", "core.service", "op", "pengumuman.delete_pengumuman.DeletePengumumanService", "err", coreerror.ErrMissingId)
		return coreerror.ErrMissingId
	}

	pengumumanData,err := r.repo.GetPengumumanById(ctx,idPengumuman)
	if err != nil {
		logger.Error(ctx, "failed deleting pengumuman", "layer", "core.service", "op", "pengumuman.delete_pengumuman.DeletePengumumanService", "err", err)
		return err
	}
	if err := r.deleteFile.DeleteFile(ctx,pengumumanData.DokumenPengumuman); err != nil {
		logger.Error(ctx, "failed deleting pengumuman", "layer", "core.service", "op", "pengumuman.delete_file_foto", "pengumuman_id", idPengumuman, "err", err)
		return err
	}

	if err := r.repo.DeletePengumuman(ctx,idPengumuman); err != nil {
		logger.Error(ctx, "failed deleting pengumuman", "layer", "core.service", "op", "pengumuman.delete_pengumuman.DeletePengumumanService", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrDeleteRestricted):
			return coreerror.ErrDeleteRestricted
		default:
			return err
		}
	}
	return nil
}