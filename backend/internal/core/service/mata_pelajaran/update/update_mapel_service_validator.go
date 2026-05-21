package matapelajaran_service

import (
	"context"
	"errors"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

func (r *UpdateMapelRepo) validateUpdateMapelPatch(ctx context.Context, idMapel int, mapel updatepatch.UpdateMapelPatch) error {
	logger := corelog.FromContext(ctx)

	if err := validateUpdateMapelID(idMapel); err != nil {
		return err
	}

	if err := validateMapelPatch(mapel); err != nil {
		return err
	}

	if err := validateMapelIdKelasPatch(mapel); err != nil {
		logger.Error(ctx, "failed updating mapel", "layer", "core.service", "op", "matapelajaran.update_mapel.IdKelas", "err", err)
		return err
	}

	if err := validateMapelRequiredFieldsPatch(mapel); err != nil {
		return err
	}

	if err := r.validateKodeMapelUniqueness(ctx, mapel); err != nil {
		if errors.Is(err, coreerror.ErrKodeMapelExist) {
			return err
		}
		logger.Error(ctx, "failed updating mapel", "layer", "core.service", "op", "matapelajaran.update_mapel.existKodeMapel", "err", err)
		return err
	}

	return nil
}

func validateUpdateMapelID(idMapel int) error {
	if idMapel <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}

func validateMapelPatch(mapel updatepatch.UpdateMapelPatch) error {
	if mapel.IdKelas == nil && mapel.KodeMapel == nil && mapel.NamaMapel == nil && mapel.Deskripsi == nil {
		return coreerror.ErrNoFieldToUpdate
	}
	return nil
}

func validateMapelIdKelasPatch(mapel updatepatch.UpdateMapelPatch) error {
	if mapel.IdKelas != nil && *mapel.IdKelas <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}

func validateMapelRequiredFieldsPatch(mapel updatepatch.UpdateMapelPatch) error {
	if isEmptyStringPatch(mapel.KodeMapel) || isEmptyStringPatch(mapel.NamaMapel) || isEmptyStringPatch(mapel.Deskripsi) {
		return coreerror.ErrMissingField
	}
	return nil
}

func isEmptyStringPatch(value *string) bool {
	return value != nil && *value == ""
}

func (r *UpdateMapelRepo) validateKodeMapelUniqueness(ctx context.Context, mapel updatepatch.UpdateMapelPatch) error {
	if mapel.KodeMapel == nil {
		return nil
	}
	exist, err := r.mapelRepo.ExistKodeMapel(ctx, *mapel.KodeMapel)
	if err != nil {
		return err
	}
	if exist {
		return coreerror.ErrKodeMapelExist
	}
	return nil
}
