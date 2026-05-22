package pengumuman_service

import (
	"context"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	pengumuman_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/date_validation"
)

func validateUpdatePengumumanPatch(ctx context.Context, idPengumuman pengumuman.ID, payload updatepatch.PengumumanUpdatePatch) error {
	logger := corelog.FromContext(ctx)

	if err := validateUpdatePengumumanID(idPengumuman); err != nil {
		logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", coreerror.ErrMissingId)
		return err
	}

	if err := validatePengumumanUpdateUserID(payload); err != nil {
		logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", coreerror.ErrMissingId)
		return err
	}

	if err := validatePengumumanUpdatePatch(payload); err != nil {
		logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", err)
		return err
	}

	if err := validatePengumumanRequiredFieldsPatch(payload); err != nil {
		logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", err)
		return err
	}

	if err := validatePengumumanUpdateDatePatch(payload); err != nil {
		logger.Error(ctx, "failed updating pengumuman", "layer", "core.service", "op", "pengumuman.update_pengumuman.UpdatePengumumanService", "err", err)
		return err
	}

	return nil
}

func validateUpdatePengumumanID(idPengumuman pengumuman.ID) error {
	if idPengumuman <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}

func validatePengumumanUpdateUserID(payload updatepatch.PengumumanUpdatePatch) error {
	if payload.IdPengguna != nil && *payload.IdPengguna <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}

func validatePengumumanUpdatePatch(payload updatepatch.PengumumanUpdatePatch) error {
	if payload.JudulPengumuman == nil && payload.IsiPengumuman == nil && payload.TanggalRilisPengumuman == nil && payload.TanggalSelesaiPengumuman == nil && payload.DokumenPengumuman == nil {
		return coreerror.ErrNoFieldToUpdate
	}
	return nil
}

func validatePengumumanRequiredFieldsPatch(payload updatepatch.PengumumanUpdatePatch) error {
	if isEmptyStringPatch(payload.JudulPengumuman) ||
		isEmptyStringPatch(payload.IsiPengumuman) ||
		isEmptyStringPatch(payload.TanggalRilisPengumuman) ||
		isEmptyStringPatch(payload.TanggalSelesaiPengumuman) ||
		isEmptyStringPatch(payload.DokumenPengumuman) {
		return coreerror.ErrMissingField
	}
	return nil
}

func validatePengumumanUpdateDatePatch(payload updatepatch.PengumumanUpdatePatch) error {
	if payload.TanggalRilisPengumuman != nil {
		if err := pengumuman_service.ValidateDate(*payload.TanggalRilisPengumuman); err != nil {
			return err
		}
	}

	if payload.TanggalSelesaiPengumuman != nil {
		if err := pengumuman_service.ValidateDate(*payload.TanggalSelesaiPengumuman); err != nil {
			return err
		}
	}

	return nil
}

func isEmptyStringPatch(value *string) bool {
	return value != nil && *value == ""
}
