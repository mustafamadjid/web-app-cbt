package siswaujian_service

import (
	"strings"

	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

func sanitizeAndValidateListUjianSiswaFilter(filter query.ListUjianFilter) (query.ListUjianFilter, error) {
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

	switch strings.ToLower(strings.TrimSpace(string(filter.KategoriUjian))) {
	case "":
		filter.KategoriUjian = ""
	case string(query.MENDATANG):
		filter.KategoriUjian = query.MENDATANG
	case string(query.BERLANGSUNG):
		filter.KategoriUjian = query.BERLANGSUNG
	case string(query.SELESAI):
		filter.KategoriUjian = query.SELESAI
	default:
		return filter, errInvalidKategoriUjian
	}

	if err := validateListUjianSiswaFilter(filter); err != nil {
		return filter, err
	}

	return filter, nil
}
