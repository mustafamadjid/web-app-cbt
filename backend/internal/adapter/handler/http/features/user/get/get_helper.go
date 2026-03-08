package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
)

func parseListGuruFilters(req *http.Request) (query.ListGuruFilter, error) {
	values := req.URL.Query()
	filters := query.ListGuruFilter{}

	filters.Search = strings.TrimSpace(values.Get("q"))
	if filters.Search == "" {
		filters.Search = strings.TrimSpace(values.Get("search"))
	}
	if err := validator.ValidateInputSafe(filters.Search, "search"); err != nil {
		return query.ListGuruFilter{}, err
	}

	if statusRaw := strings.TrimSpace(values.Get("status")); statusRaw != "" {
		if err := validator.ValidateInputSafe(statusRaw, "status"); err != nil {
			return query.ListGuruFilter{}, err
		}
		status := user.StatusAkun(strings.ToUpper(statusRaw))
		filters.Status = &status
	}

	if limitRaw := strings.TrimSpace(values.Get("limit")); limitRaw != "" {
		limit, err := strconv.Atoi(limitRaw)
		if err != nil {
			return query.ListGuruFilter{}, errors.New("limit must be a number")
		}
		filters.Limit = limit
	}

	if offsetRaw := strings.TrimSpace(values.Get("offset")); offsetRaw != "" {
		offset, err := strconv.Atoi(offsetRaw)
		if err != nil {
			return query.ListGuruFilter{}, errors.New("offset must be a number")
		}
		filters.Offset = offset
	}

	filters.SortBy = strings.TrimSpace(values.Get("sort_by"))
	if err := validator.ValidateInputSafe(filters.SortBy, "sort_by"); err != nil {
		return query.ListGuruFilter{}, err
	}

	if sortDescRaw := strings.TrimSpace(values.Get("sort_desc")); sortDescRaw != "" {
		sortDesc, err := strconv.ParseBool(sortDescRaw)
		if err != nil {
			return query.ListGuruFilter{}, errors.New("sort_desc must be a boolean")
		}
		filters.SortDesc = sortDesc
	}

	if bidangRaw := strings.TrimSpace(values.Get("bidang")); bidangRaw != "" {
		if err := validator.ValidateInputSafe(bidangRaw, "bidang"); err != nil {
			return query.ListGuruFilter{}, err
		}
		filters.Bidang = &bidangRaw
	} else if bidangStudiRaw := strings.TrimSpace(values.Get("bidang_studi")); bidangStudiRaw != "" {
		if err := validator.ValidateInputSafe(bidangStudiRaw, "bidang_studi"); err != nil {
			return query.ListGuruFilter{}, err
		}
		filters.Bidang = &bidangStudiRaw
	}

	return filters, nil
}

func parseListSiswaFilters(req *http.Request) (query.ListSiswaFilter, error) {
	values := req.URL.Query()
	filters := query.ListSiswaFilter{}

	filters.Search = strings.TrimSpace(values.Get("q"))
	if filters.Search == "" {
		filters.Search = strings.TrimSpace(values.Get("search"))
	}
	if err := validator.ValidateInputSafe(filters.Search, "search"); err != nil {
		return query.ListSiswaFilter{}, err
	}

	if statusRaw := strings.TrimSpace(values.Get("status")); statusRaw != "" {
		if err := validator.ValidateInputSafe(statusRaw, "status"); err != nil {
			return query.ListSiswaFilter{}, err
		}
		status := user.StatusAkun(strings.ToUpper(statusRaw))
		filters.Status = &status
	}

	if limitRaw := strings.TrimSpace(values.Get("limit")); limitRaw != "" {
		limit, err := strconv.Atoi(limitRaw)
		if err != nil {
			return query.ListSiswaFilter{}, errors.New("limit must be a number")
		}
		filters.Limit = limit
	}

	if offsetRaw := strings.TrimSpace(values.Get("offset")); offsetRaw != "" {
		offset, err := strconv.Atoi(offsetRaw)
		if err != nil {
			return query.ListSiswaFilter{}, errors.New("offset must be a number")
		}
		filters.Offset = offset
	}

	filters.SortBy = strings.TrimSpace(values.Get("sort_by"))
	if err := validator.ValidateInputSafe(filters.SortBy, "sort_by"); err != nil {
		return query.ListSiswaFilter{}, err
	}

	if sortDescRaw := strings.TrimSpace(values.Get("sort_desc")); sortDescRaw != "" {
		sortDesc, err := strconv.ParseBool(sortDescRaw)
		if err != nil {
			return query.ListSiswaFilter{}, errors.New("sort_desc must be a boolean")
		}
		filters.SortDesc = sortDesc
	}

	if angkatanRaw := strings.TrimSpace(values.Get("angkatan")); angkatanRaw != "" {
		angkatan, err := strconv.Atoi(angkatanRaw)
		if err != nil {
			return query.ListSiswaFilter{}, errors.New("angkatan must be a number")
		}
		filters.Angkatan = &angkatan
	}

	if tingkatKelasRaw := strings.TrimSpace(values.Get("tingkat_kelas")); tingkatKelasRaw != "" {
		tingkatKelas, err := strconv.Atoi(tingkatKelasRaw)
		if err != nil {
			return query.ListSiswaFilter{}, errors.New("tingkat_kelas must be a number")
		}
		filters.TingkatKelas = &tingkatKelas
	}

	if jenisKelaminRaw := strings.TrimSpace(values.Get("jenis_kelamin")); jenisKelaminRaw != "" {
		jenisKelamin, err := strconv.Atoi(jenisKelaminRaw)
		if err != nil {
			return query.ListSiswaFilter{}, errors.New("jenis_kelamin must be a number")
		}
		filters.JenisKelamin = &jenisKelamin
	}

	if idTingkatKelasRaw := strings.TrimSpace(values.Get("id_tingkat_kelas")); idTingkatKelasRaw != "" {
		idTingkatKelas, err := strconv.Atoi(idTingkatKelasRaw)
		if err != nil {
			return query.ListSiswaFilter{}, errors.New("id_tingkat_kelas must be a number")
		}
		filters.IdTingkatKelas = &idTingkatKelas
	}

	if idNamaKelasRaw := strings.TrimSpace(values.Get("id_nama_kelas")); idNamaKelasRaw != "" {
		idNamaKelas, err := strconv.Atoi(idNamaKelasRaw)
		if err != nil {
			return query.ListSiswaFilter{}, errors.New("id_nama_kelas must be a number")
		}
		filters.IdNamaKelas = &idNamaKelas
	}

	return filters, nil
}

func toGuruResponses(items []query.GuruListItem) []GuruResponseItem {
	response := make([]GuruResponseItem, 0, len(items))
	for _, item := range items {
		response = append(response, toGuruResponseFromList(item))
	}

	return response
}

func toGuruResponseFromList(item query.GuruListItem) GuruResponseItem {
	return GuruResponseItem{
		IdPengguna:   item.IdPengguna,
		Username:     item.Username,
		NoHp:         item.NoHp,
		Email:        item.Email,
		NamaLengkap:  item.NamaLengkap,
		StatusAkun:   item.StatusAkun,
		Nip:          item.Nip,
		Jabatan:      item.Jabatan,
		BidangStudi:  item.BidangStudi,
		Foto:         item.Foto,
		Role:         item.Role,
		JenisKelamin: item.JenisKelamin,
	}
}

func toGuruResponseFromData(item user.DataGuru) GuruResponseItem {
	return GuruResponseItem{
		IdPengguna:   item.IdPengguna,
		Username:     item.Username,
		NoHp:         item.NoHp,
		Email:        user.Email(item.Email),
		NamaLengkap:  item.NamaLengkap,
		Nip:          user.NIP(item.Nip),
		Jabatan:      item.Jabatan,
		BidangStudi:  item.BidangStudi,
		Foto:         item.Foto,
		JenisKelamin: item.JenisKelamin,
		Role:         item.Role,
		StatusAkun:   item.StatusAkun,
	}
}

func toSiswaResponses(items []query.SiswaListItem) []SiswaResponseItem {
	response := make([]SiswaResponseItem, 0, len(items))
	for _, item := range items {
		response = append(response, toSiswaResponseFromList(item))
	}

	return response
}

func toSiswaResponseFromList(item query.SiswaListItem) SiswaResponseItem {
	return SiswaResponseItem{
		IdPengguna:   item.IdPengguna,
		Role:         user.SISWA,
		Username:     item.Username,
		NoHp:         item.NoHp,
		Email:        item.Email,
		NamaLengkap:  item.NamaLengkap,
		StatusAkun:   item.StatusAkun,
		NoAbsen:      item.NoAbsen,
		Angkatan:     item.Angkatan,
		TempatLahir:  item.TempatLahir,
		TanggalLahir: httphelper.FormatDateOnly(item.TanggalLahir),
		NamaKelas:    item.NamaKelas,
		TingkatKelas: item.TingkatKelas,
		Foto:         item.Foto,
		JenisKelamin: item.JenisKelamin,
		Nisn:         item.Nisn,
	}
}

func toSiswaResponseFromData(item user.DataSiswa) SiswaResponseItem {
	return SiswaResponseItem{
		IdPengguna:   item.IdPengguna,
		Role:         item.Role,
		Username:     item.Username,
		NoHp:         item.NoHp,
		Email:        user.Email(item.Email),
		NamaLengkap:  item.NamaLengkap,
		StatusAkun:   item.StatusAkun,
		Nisn:         item.Nisn,
		NoAbsen:      item.NoAbsen,
		Angkatan:     item.Angkatan,
		TempatLahir:  item.TempatLahir,
		TanggalLahir: httphelper.FormatDateOnly(item.TanggalLahir),
		NamaKelas:    item.NamaKelas,
		TingkatKelas: item.TingkatKelas,
		Foto:         item.Foto,
		JenisKelamin: item.JenisKelamin,
	}
}
