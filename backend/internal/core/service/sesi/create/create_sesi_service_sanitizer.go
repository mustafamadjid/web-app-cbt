package sesi_service

import (
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	"strings"
)

func sanitizeCreateSesi(data sesi.Sesi) sesi.Sesi {
	data.NamaSesi = strings.TrimSpace(data.NamaSesi)
	data.KodeSesi = strings.TrimSpace(data.KodeSesi)
	data.KodeSesi = strings.ToUpper(data.KodeSesi)
	return data
}
