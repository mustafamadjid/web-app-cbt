package sesi_service

import (
	"context"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"

	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	sesi_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/sesi"
)

type CreateSesiService struct {
	sesiRepo sesi_repo.SesiRepository
}

func NewCreateSesiService(sesiRepo sesi_repo.SesiRepository) *CreateSesiService {
	return &CreateSesiService{sesiRepo: sesiRepo}
}

func(r *CreateSesiService)CreateSesiService(ctx context.Context,sesi sesi.Sesi) error {
	logger := corelog.FromContext(ctx)

	sesi.NamaSesi = strings.TrimSpace(sesi.NamaSesi)
	sesi.KodeSesi = strings.TrimSpace(sesi.KodeSesi)
	sesi.KodeSesi = strings.ToUpper(sesi.KodeSesi)

	if len(sesi.NamaSesi) == 0 || sesi.NamaSesi == "" {
		logger.Error(ctx,"failed create sesi","layer","core.service","op","sesi.create","err",coreerror.ErrMissingField)
		return coreerror.ErrMissingField
	}

	if len(sesi.KodeSesi) == 0 || sesi.KodeSesi == "" {
		logger.Error(ctx,"failed create sesi","layer","core.service","op","sesi.create","err",coreerror.ErrMissingField)
		return coreerror.ErrMissingField
	}

	exist, err := r.sesiRepo.ExistByKodeSesi(ctx,sesi.KodeSesi)
	if err != nil {
		logger.Error(ctx,"failed check exist sesi","layer","core.service","op","sesi.create","err",err)
		return err
	}

	if exist {
		logger.Error(ctx,"failed create sesi","layer","core.service","op","sesi.create","err",err)
		return coreerror.ErrSesiUjianExist
	}

	if err := r.sesiRepo.CreateSesi(ctx,sesi); err != nil {
		logger.Error(ctx,"failed create sesi","layer","core.service","op","sesi.create","err",err)
		return err
	}
	return nil
}