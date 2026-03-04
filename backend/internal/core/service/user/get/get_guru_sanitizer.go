package user_service

import (
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
	"strings"
)

func sanitizeListGuruFilter(filter query.ListGuruFilter) query.ListGuruFilter {
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
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}
	if filter.Bidang != nil {
		trimmedBidang := strings.TrimSpace(*filter.Bidang)
		filter.Bidang = &trimmedBidang
	}
	return filter
}
