package sesi_service

import (
	"context"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	sesi_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/sesi"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

type UpdateSesiService struct {
	sesiRepo sesi_repo.SesiRepository
}

func NewUpdateSesiService(sesiRepo sesi_repo.SesiRepository) *UpdateSesiService {
	return &UpdateSesiService{sesiRepo: sesiRepo}
}

func(r *UpdateSesiService)UpdateSesiService(ctx context.Context,idSesi int,sesi updatepatch.UpdateSesiPatch) error{
	logger := corelog.FromContext(ctx)

	if idSesi == 0 {
		return coreerror.ErrMissingId
	}

	if sesi.NamaSesi != nil {
		s := strings.TrimSpace(*sesi.NamaSesi)
		if s == "" {
			return coreerror.ErrMissingField
		}

		sesi.NamaSesi = &s
	}

	if sesi.KodeSesi != nil {
		s := strings.TrimSpace(*sesi.KodeSesi)
		if s == "" {
			return coreerror.ErrMissingField
		}

		s = strings.ToUpper(s)
		sesi.KodeSesi = &s
	}

	if err := r.sesiRepo.UpdateSesi(ctx,idSesi,sesi); err != nil {
		logger.Error(ctx,"failed update sesi","layer","core.service","op","sesi.update","err",err)
		return err
	}
	return nil
}