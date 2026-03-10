package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
)

func parseListKelasFilters(r *http.Request) (query.ListKelasFilter, error) {
	values := r.URL.Query()
	filters := query.ListKelasFilter{}

	filters.Search = strings.TrimSpace(values.Get("q"))
	if filters.Search == "" {
		filters.Search = strings.TrimSpace(values.Get("search"))
	}
	search, err := validator.ValidateOptionalPrintableText(filters.Search, "search")
	if err != nil {
		return query.ListKelasFilter{}, err
	}
	filters.Search = search

	if tingkatKelasRaw := strings.TrimSpace(values.Get("tingkat_kelas")); tingkatKelasRaw != "" {
		tingkatKelas, err := strconv.Atoi(tingkatKelasRaw)
		if err != nil {
			return query.ListKelasFilter{}, errors.New("tingkat_kelas must be a number")
		}
		filters.TingkatKelas = &tingkatKelas
	}

	if limitRaw := strings.TrimSpace(values.Get("limit")); limitRaw != "" {
		limit, err := strconv.Atoi(limitRaw)
		if err != nil {
			return query.ListKelasFilter{}, errors.New("limit must be a number")
		}
		filters.Limit = limit
	}

	if offsetRaw := strings.TrimSpace(values.Get("offset")); offsetRaw != "" {
		offset, err := strconv.Atoi(offsetRaw)
		if err != nil {
			return query.ListKelasFilter{}, errors.New("offset must be a number")
		}
		filters.Offset = offset
	}

	return filters, nil
}

func toFullKelasResponse(items []kelas.FullKelasData) FullKelasResponse {
	response := FullKelasResponse{}
	for _, item := range items {
		for _, tingkat := range item.ItemsTingkatKelas {
			response.ItemsTingkatKelas = append(response.ItemsTingkatKelas, TingkatKelasResponse{
				IDTingkatKelas: int(tingkat.IdTingkatKelas),
				TingkatKelas:   tingkat.TingkatKelas,
			})
		}
		for _, nama := range item.ItemsNamaKelas {
			response.ItemsNamaKelas = append(response.ItemsNamaKelas, NamaKelasResponse{
				IDNamaKelas:    int(nama.IdNamaKelas),
				IDTingkatKelas: int(nama.IdTingkatKelas),
				NamaKelas:      nama.NamaKelas,
			})
		}
	}

	return response
}

func toDataKelasResponse(item kelas.KelasData) DataKelasResponse {
	return DataKelasResponse{
		ItemsTingkatKelas: TingkatKelasResponse{
			IDTingkatKelas: int(item.ItemsTingkatKelas.IdTingkatKelas),
			TingkatKelas:   item.ItemsTingkatKelas.TingkatKelas,
		},
		ItemsNamaKelas: NamaKelasResponse{
			IDNamaKelas:    int(item.ItemsNamaKelas.IdNamaKelas),
			IDTingkatKelas: int(item.ItemsNamaKelas.IdTingkatKelas),
			NamaKelas:      item.ItemsNamaKelas.NamaKelas,
		},
	}
}
