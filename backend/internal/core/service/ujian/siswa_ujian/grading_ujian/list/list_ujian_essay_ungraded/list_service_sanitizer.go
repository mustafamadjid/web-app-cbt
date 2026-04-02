package gradingujian_service

import (
	"strings"

	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

func sanitizeAndValidateListUjianEssayUngradedFilter(filter query.ListUjianEssayUngradedFilter) (query.ListUjianEssayUngradedFilter, error) {
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

	if filter.TanggalUjian != nil {
		tanggalUjian := strings.TrimSpace(*filter.TanggalUjian)
		filter.TanggalUjian = &tanggalUjian
	}
	if filter.Tahun != nil {
		tahun := strings.TrimSpace(*filter.Tahun)
		filter.Tahun = &tahun
	}
	if filter.Bulan != nil {
		bulan := strings.TrimSpace(*filter.Bulan)
		filter.Bulan = &bulan
	}

	if err := validateListUjianEssayUngradedFilter(filter); err != nil {
		return filter, err
	}

	return filter, nil
}
