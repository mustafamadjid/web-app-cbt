package ujian_service

import (
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	"strconv"
	"strings"
	"time"
)

func sanitizeAndValidateListUjianFilter(filter query.ListUjianFilter) (query.ListUjianFilter, error) {
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
		if tanggalUjian == "" {
			return filter, errInvalidTanggalUjian
		}
		if _, err := time.Parse("2006-01-02", tanggalUjian); err != nil {
			return filter, errInvalidTanggalUjian
		}
		filter.TanggalUjian = &tanggalUjian
	}
	if filter.Tahun != nil {
		tahun := strings.TrimSpace(*filter.Tahun)
		if tahun == "" {
			return filter, errInvalidTahun
		}
		if len(tahun) != 4 {
			return filter, errInvalidTahun
		}
		tahunInt, err := strconv.Atoi(tahun)
		if err != nil || tahunInt <= 0 {
			return filter, errInvalidTahun
		}
		filter.Tahun = &tahun
	}
	if filter.TingkatKelasID != nil && *filter.TingkatKelasID <= 0 {
		return filter, errInvalidTingkatKelas
	}
	if filter.TingkatKelas != nil && *filter.TingkatKelas <= 0 {
		return filter, errInvalidTingkatKelas
	}
	if filter.RuangUjian != nil && *filter.RuangUjian <= 0 {
		return filter, errInvalidRuangUjian
	}
	return filter, nil
}

func sanitizeAndValidateListUjianItems(items []ujian.ListUjian) ([]ujian.ListUjian, error) {
	sanitized := make([]ujian.ListUjian, 0, len(items))

	for _, item := range items {
		item = sanitizeListUjianItem(item)
		if err := validateListUjian(item); err != nil {
			return nil, err
		}
		sanitized = append(sanitized, item)
	}

	return sanitized, nil
}

func sanitizeJadwalUjian(item ujian.JadwalUjian) ujian.JadwalUjian {
	item.Token = strings.ToUpper(strings.TrimSpace(item.Token))
	item.StatusUjian = ujian.StatusUjian(strings.ToUpper(strings.TrimSpace(string(item.StatusUjian))))
	return item
}

func sanitizeListUjianItem(item ujian.ListUjian) ujian.ListUjian {
	item.NamaUjian = strings.TrimSpace(item.NamaUjian)
	item.PembuatUsername = strings.TrimSpace(item.PembuatUsername)
	item.NamaPengawas = strings.TrimSpace(item.NamaPengawas)
	item.PengawasUsername = strings.TrimSpace(item.PengawasUsername)
	item.NamaSesi = strings.TrimSpace(item.NamaSesi)
	item.NamaRuangan = strings.TrimSpace(item.NamaRuangan)

	if item.NamaKelas != nil {
		namaKelas := strings.TrimSpace(*item.NamaKelas)
		if namaKelas == "" {
			item.NamaKelas = nil
		} else {
			item.NamaKelas = &namaKelas
		}
	}

	return item
}
