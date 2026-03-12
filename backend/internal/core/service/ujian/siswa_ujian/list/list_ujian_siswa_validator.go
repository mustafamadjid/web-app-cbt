package siswaujian_service

import (
	"errors"
	"strconv"
	"time"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
)

func validateListUjianSiswaID(idSiswa int) error {
	if idSiswa <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}

var (
	errInvalidTanggalUjian  = errors.New("invalid tanggal ujian")
	errInvalidTahun         = errors.New("invalid tahun")
	errInvalidBulan         = errors.New("invalid bulan")
	errInvalidTingkatKelas  = errors.New("invalid tingkat kelas")
	errInvalidRuangUjian    = errors.New("invalid ruang ujian")
	errInvalidIDMapel       = errors.New("invalid id mapel")
	errInvalidKategoriUjian = errors.New("invalid kategori ujian")
)

func validateListUjianSiswaFilter(filter query.ListUjianFilter) error {
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
	if filter.TingkatKelas != nil && *filter.TingkatKelas <= 0 {
		return errInvalidTingkatKelas
	}
	if filter.RuangUjian != nil && *filter.RuangUjian <= 0 {
		return errInvalidRuangUjian
	}
	if filter.IDMapel != nil && *filter.IDMapel <= 0 {
		return errInvalidIDMapel
	}

	switch filter.KategoriUjian {
	case "", query.MENDATANG, query.BERLANGSUNG, query.SELESAI, query.DIBATALKAN:
		return nil
	default:
		return errInvalidKategoriUjian
	}
}
