package pengumuman_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	pengumuman_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/date_validation"
	"strings"
)

func sanitizeJudulPengumumanPatch(payload *updatepatch.PengumumanUpdatePatch) error {
	if payload.JudulPengumuman == nil {
		return nil
	}
	judulPengumuman := strings.TrimSpace(*payload.JudulPengumuman)
	if judulPengumuman == "" {
		return coreerror.ErrMissingField
	}
	payload.JudulPengumuman = &judulPengumuman
	return nil
}
func sanitizeIsiPengumumanPatch(payload *updatepatch.PengumumanUpdatePatch) error {
	if payload.IsiPengumuman == nil {
		return nil
	}
	isiPengumuman := strings.TrimSpace(*payload.IsiPengumuman)
	if isiPengumuman == "" {
		return coreerror.ErrMissingField
	}
	payload.IsiPengumuman = &isiPengumuman
	return nil
}
func sanitizeTanggalRilisPengumumanPatch(payload *updatepatch.PengumumanUpdatePatch) error {
	if payload.TanggalRilisPengumuman == nil {
		return nil
	}
	tanggalRilisPengumuman := strings.TrimSpace(*payload.TanggalRilisPengumuman)
	if tanggalRilisPengumuman == "" {
		return coreerror.ErrMissingField
	}
	if err := pengumuman_service.ValidateDate(tanggalRilisPengumuman); err != nil {
		return err
	}
	payload.TanggalRilisPengumuman = &tanggalRilisPengumuman
	return nil
}
func sanitizeTanggalSelesaiPengumumanPatch(payload *updatepatch.PengumumanUpdatePatch) error {
	if payload.TanggalSelesaiPengumuman == nil {
		return nil
	}
	tanggalSelesaiPengumuman := strings.TrimSpace(*payload.TanggalSelesaiPengumuman)
	if tanggalSelesaiPengumuman == "" {
		return coreerror.ErrMissingField
	}
	if err := pengumuman_service.ValidateDate(tanggalSelesaiPengumuman); err != nil {
		return err
	}
	payload.TanggalSelesaiPengumuman = &tanggalSelesaiPengumuman
	return nil
}
func sanitizeDokumenPengumumanPatch(payload *updatepatch.PengumumanUpdatePatch) error {
	if payload.DokumenPengumuman == nil {
		return nil
	}
	dokumenPengumuman := strings.TrimSpace(*payload.DokumenPengumuman)
	if dokumenPengumuman == "" {
		return coreerror.ErrMissingField
	}
	payload.DokumenPengumuman = &dokumenPengumuman
	return nil
}
