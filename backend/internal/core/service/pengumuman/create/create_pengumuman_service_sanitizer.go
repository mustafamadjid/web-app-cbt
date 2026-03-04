package pengumuman_service

import (
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	"strings"
)

func sanitizeCreatePengumuman(data pengumuman.Pengumuman) pengumuman.Pengumuman {
	data.IsiPengumuman = strings.TrimSpace(data.IsiPengumuman)
	data.JudulPengumuman = strings.TrimSpace(data.JudulPengumuman)
	data.TanggalRilisPengumuman = strings.TrimSpace(data.TanggalRilisPengumuman)
	data.TanggalSelesaiPengumuman = strings.TrimSpace(data.TanggalSelesaiPengumuman)
	data.DokumenPengumuman = strings.TrimSpace(data.DokumenPengumuman)
	return data
}
