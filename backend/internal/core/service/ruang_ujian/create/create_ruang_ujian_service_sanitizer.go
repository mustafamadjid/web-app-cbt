package ruangujian_service

import (
	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	"strings"
)

func sanitizeCreateRuangUjian(data ruangujian.RuangUjian) ruangujian.RuangUjian {
	data.KodeRuang = strings.TrimSpace(data.KodeRuang)
	data.KodeRuang = strings.ToUpper(data.KodeRuang)
	data.NamaRuangan = strings.TrimSpace(data.NamaRuangan)
	return data
}
