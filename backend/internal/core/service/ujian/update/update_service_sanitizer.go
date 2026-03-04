package ujian_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"strings"
)

func sanitizeNamaUjianPatch(payload *updatepatch.UpdateUjianPatch) error {
	if payload.NamaUjian == nil {
		return nil
	}
	namaUjian := strings.TrimSpace(*payload.NamaUjian)
	if namaUjian == "" {
		return coreerror.ErrMissingField
	}
	if len(namaUjian) > 100 {
		return coreerror.ErrInvalidInput
	}
	payload.NamaUjian = &namaUjian
	return nil
}
func sanitizeDeskripsiUjianPatch(payload *updatepatch.UpdateUjianPatch) error {
	if payload.DeskripsiUjian == nil {
		return nil
	}
	deskripsiUjian := strings.TrimSpace(*payload.DeskripsiUjian)
	if deskripsiUjian == "" {
		return coreerror.ErrMissingField
	}
	payload.DeskripsiUjian = &deskripsiUjian
	return nil
}
func sanitizeStatusUjianPatch(payload *updatepatch.UpdateJadwalUjianPatch) error {
	if payload.StatusUjian == nil {
		return nil
	}
	statusUjian := strings.TrimSpace(string(*payload.StatusUjian))
	if statusUjian == "" {
		return coreerror.ErrMissingField
	}
	statusUjian = strings.ToUpper(statusUjian)
	status := ujian.StatusUjian(statusUjian)
	payload.StatusUjian = &status
	return nil
}
func sanitizeTokenJadwalUjianPatch(payload *updatepatch.UpdateJadwalUjianPatch) error {
	if payload.Token == nil {
		return nil
	}
	token := strings.ToUpper(strings.TrimSpace(*payload.Token))
	if token == "" {
		return coreerror.ErrMissingField
	}
	payload.Token = &token
	return nil
}
func sanitizeJawabanEssayPatch(payload *updatepatch.UpdateJawabanUjianSiswaPatch) error {
	if payload.JawabanEssay == nil {
		return nil
	}
	jawabanEssay := strings.TrimSpace(*payload.JawabanEssay)
	if jawabanEssay == "" {
		return coreerror.ErrMissingField
	}
	payload.JawabanEssay = &jawabanEssay
	return nil
}
