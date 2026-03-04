package matapelajaran_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"strings"
)

func sanitizeKodeMapelPatch(mapel *updatepatch.UpdateMapelPatch) error {
	if mapel.KodeMapel == nil {
		return nil
	}
	kodeMapel := strings.TrimSpace(*mapel.KodeMapel)
	if kodeMapel == "" {
		return coreerror.ErrMissingField
	}
	kodeMapel = strings.ToUpper(kodeMapel)
	mapel.KodeMapel = &kodeMapel
	return nil
}
func sanitizeNamaMapelPatch(mapel *updatepatch.UpdateMapelPatch) error {
	if mapel.NamaMapel == nil {
		return nil
	}
	namaMapel := strings.TrimSpace(*mapel.NamaMapel)
	if namaMapel == "" {
		return coreerror.ErrMissingField
	}
	mapel.NamaMapel = &namaMapel
	return nil
}
func sanitizeDeskripsiMapelPatch(mapel *updatepatch.UpdateMapelPatch) error {
	if mapel.Deskripsi == nil {
		return nil
	}
	deskripsi := strings.TrimSpace(*mapel.Deskripsi)
	if deskripsi == "" {
		return coreerror.ErrMissingField
	}
	mapel.Deskripsi = &deskripsi
	return nil
}
