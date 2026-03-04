package kelas_service

import (
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
	"strings"
)

func sanitizeListKelasFilter(filter query.ListKelasFilter) query.ListKelasFilter {
	filter.Search = strings.TrimSpace(filter.Search)
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 50 {
		filter.Limit = 50
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	return filter
}
