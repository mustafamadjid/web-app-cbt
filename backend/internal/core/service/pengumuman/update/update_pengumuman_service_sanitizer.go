package pengumuman_service

import (
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"strings"
)

func sanitizeUpdatePengumumanPatch(payload updatepatch.PengumumanUpdatePatch) updatepatch.PengumumanUpdatePatch {
	if payload.JudulPengumuman != nil {
		judulPengumuman := strings.TrimSpace(*payload.JudulPengumuman)
		payload.JudulPengumuman = &judulPengumuman
	}

	if payload.IsiPengumuman != nil {
		isiPengumuman := strings.TrimSpace(*payload.IsiPengumuman)
		payload.IsiPengumuman = &isiPengumuman
	}

	if payload.TanggalRilisPengumuman != nil {
		tanggalRilisPengumuman := strings.TrimSpace(*payload.TanggalRilisPengumuman)
		payload.TanggalRilisPengumuman = &tanggalRilisPengumuman
	}

	if payload.TanggalSelesaiPengumuman != nil {
		tanggalSelesaiPengumuman := strings.TrimSpace(*payload.TanggalSelesaiPengumuman)
		payload.TanggalSelesaiPengumuman = &tanggalSelesaiPengumuman
	}

	if payload.DokumenPengumuman != nil {
		dokumenPengumuman := strings.TrimSpace(*payload.DokumenPengumuman)
		payload.DokumenPengumuman = &dokumenPengumuman
	}

	return payload
}
