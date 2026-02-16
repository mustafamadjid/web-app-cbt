package sesi_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	sesi_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/sesi"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"

)

type DeleteSesiService struct {
	sesiRepo sesi_repo.SesiRepository
}

func NewDeleteSesiService(sesiRepo sesi_repo.SesiRepository) *DeleteSesiService {
	return &DeleteSesiService{sesiRepo: sesiRepo}
}

func(r *DeleteSesiService)DeleteSesiService(ctx context.Context,idSesi int) error {
	logger := corelog.FromContext(ctx)
	
	if idSesi == 0 {
		return coreerror.ErrMissingId
	}

	if err := r.sesiRepo.DeleteSesi(ctx,idSesi); err != nil {
		logger.Error(ctx,"failed delete sesi","layer","core.service","op","sesi.delete","err",err)
		return err
	}
	return nil
}