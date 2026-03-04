package ruangujian_service

import (
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ruang_ujian"
	"strings"
)

func sanitizeListRuangUjianFilter(filter query.ListRuangUjianFilter) query.ListRuangUjianFilter {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 50 {
		filter.Limit = 50
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	filter.Search = strings.TrimSpace(filter.Search)
	return filter
}
func sanitizeKodeRuang(kodeRuang string) string {
	kodeRuang = strings.TrimSpace(kodeRuang)
	kodeRuang = strings.ToUpper(kodeRuang)
	return kodeRuang
}
