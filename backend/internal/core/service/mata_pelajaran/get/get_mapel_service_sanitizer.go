package matapelajaran_service

import (
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
	"strings"
)

func sanitizeAndValidateListMapelFilter(filter query.ListMapelFilter) (query.ListMapelFilter, error) {
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
	if filter.NamaMapel != nil {
		if *filter.NamaMapel == "" {
			return filter, errInvalidNamaMapelFilter
		}
		namaMapel := strings.TrimSpace(*filter.NamaMapel)
		filter.NamaMapel = &namaMapel
	}
	if filter.TingkatKelas != nil && *filter.TingkatKelas <= 0 {
		return filter, errInvalidTingkatKelas
	}
	return filter, nil
}
