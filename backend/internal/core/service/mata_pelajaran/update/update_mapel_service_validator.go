package matapelajaran_service

import (
	"context"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

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
