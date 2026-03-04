package sesi_service

import (
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
	"strings"
)

func sanitizeListSesiFilter(filter query.ListSesiFilter) query.ListSesiFilter {
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Limit > 50 {
		filter.Limit = 50
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	return filter
}
func sanitizeKodeSesi(kodeSesi string) string {
	kodeSesi = strings.TrimSpace(kodeSesi)
	kodeSesi = strings.ToUpper(kodeSesi)
	return kodeSesi
}
