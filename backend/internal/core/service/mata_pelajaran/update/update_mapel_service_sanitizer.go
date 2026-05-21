package matapelajaran_service

import (
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	"strings"
)

func sanitizeUpdateMapelPatch(mapel updatepatch.UpdateMapelPatch) updatepatch.UpdateMapelPatch {
	if mapel.KodeMapel != nil {
		kodeMapel := strings.TrimSpace(*mapel.KodeMapel)
		kodeMapel = strings.ToUpper(kodeMapel)
		mapel.KodeMapel = &kodeMapel
	}

	if mapel.NamaMapel != nil {
		namaMapel := strings.TrimSpace(*mapel.NamaMapel)
		mapel.NamaMapel = &namaMapel
	}

	if mapel.Deskripsi != nil {
		deskripsi := strings.TrimSpace(*mapel.Deskripsi)
		mapel.Deskripsi = &deskripsi
	}

	return mapel
}
