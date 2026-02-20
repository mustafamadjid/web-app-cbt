package updatepatch

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"

type PengumumanUpdatePatch struct {
	IdPengguna               *pengumuman.ID
	JudulPengumuman          *string
	IsiPengumuman            *string
	TanggalRilisPengumuman   *string
	TanggalSelesaiPengumuman *string
	DokumenPengumuman        *string
}