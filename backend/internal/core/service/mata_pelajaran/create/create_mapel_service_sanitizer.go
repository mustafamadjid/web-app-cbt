package matapelajaran_service

import (
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	"strings"
)

func sanitizeMapel(mapel matapelajaran.MataPelajaran) matapelajaran.MataPelajaran {
	mapel.KodeMapel = strings.TrimSpace(mapel.KodeMapel)
	mapel.KodeMapel = strings.ToUpper(mapel.KodeMapel)
	mapel.NamaMapel = strings.TrimSpace(mapel.NamaMapel)
	mapel.Deskripsi = strings.TrimSpace(mapel.Deskripsi)
	return mapel
}
