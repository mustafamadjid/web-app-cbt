package ruangujian_service

import (
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"strings"
)

func sanitizeKodeRuangPatch(ruangUjian *updatepatch.UpdateRuangUjianPatch) {
	if ruangUjian.KodeRuang == nil {
		return
	}
	kodeRuang := strings.TrimSpace(*ruangUjian.KodeRuang)
	kodeRuang = strings.ToUpper(kodeRuang)
	ruangUjian.KodeRuang = &kodeRuang
}
func sanitizeNamaRuangPatch(ruangUjian *updatepatch.UpdateRuangUjianPatch) {
	if ruangUjian.NamaRuang == nil {
		return
	}
	namaRuang := strings.TrimSpace(*ruangUjian.NamaRuang)
	ruangUjian.NamaRuang = &namaRuang
}
