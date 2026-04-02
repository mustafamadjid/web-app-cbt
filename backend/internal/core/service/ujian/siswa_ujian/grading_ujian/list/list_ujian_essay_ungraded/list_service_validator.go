package gradingujian_service

import (
	"errors"
	"strconv"
	"time"

	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

var (
	errInvalidTanggalUjian = errors.New("invalid tanggal ujian")
	errInvalidTahun        = errors.New("invalid tahun")
	errInvalidBulan        = errors.New("invalid bulan")
	errInvalidTingkatKelas = errors.New("invalid tingkat kelas")
	errInvalidNamaKelas    = errors.New("invalid nama kelas")
	errInvalidIDMapel      = errors.New("invalid id mapel")
	errInvalidSesi         = errors.New("invalid sesi")
)

func validateListUjianEssayUngradedFilter(filter query.ListUjianEssayUngradedFilter) error {
	if filter.TanggalUjian != nil {
		if *filter.TanggalUjian == "" {
			return errInvalidTanggalUjian
		}
		if _, err := time.Parse("2006-01-02", *filter.TanggalUjian); err != nil {
			return errInvalidTanggalUjian
		}
	}

	if filter.Tahun != nil {
		if *filter.Tahun == "" {
			return errInvalidTahun
		}
		if len(*filter.Tahun) != 4 {
			return errInvalidTahun
		}

		tahunInt, err := strconv.Atoi(*filter.Tahun)
		if err != nil || tahunInt <= 0 {
			return errInvalidTahun
		}
	}

	if filter.Bulan != nil {
		if *filter.Bulan == "" {
			return errInvalidBulan
		}

		bulanInt, err := strconv.Atoi(*filter.Bulan)
		if err != nil || bulanInt < 1 || bulanInt > 12 {
			return errInvalidBulan
		}
	}

	if filter.TingkatKelasID != nil && *filter.TingkatKelasID <= 0 {
		return errInvalidTingkatKelas
	}
	if filter.NamaKelasID != nil && *filter.NamaKelasID <= 0 {
		return errInvalidNamaKelas
	}
	if filter.MapelID != nil && *filter.MapelID <= 0 {
		return errInvalidIDMapel
	}
	if filter.SesiID != nil && *filter.SesiID <= 0 {
		return errInvalidSesi
	}

	return nil
}
