package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
)

func parseListMapelFilters(r *http.Request) (query.ListMapelFilter, error) {
	values := r.URL.Query()
	filters := query.ListMapelFilter{}

	filters.Search = strings.TrimSpace(values.Get("q"))
	if filters.Search == "" {
		filters.Search = strings.TrimSpace(values.Get("search"))
	}
	if err := validator.ValidateInputSafe(filters.Search, "search"); err != nil {
		return query.ListMapelFilter{}, err
	}

	if namaMapel := strings.TrimSpace(values.Get("nama_mapel")); namaMapel != "" {
		filters.NamaMapel = &namaMapel
	}

	if tingkatKelasRaw := strings.TrimSpace(values.Get("tingkat_kelas")); tingkatKelasRaw != "" {
		tingkatKelas, err := strconv.Atoi(tingkatKelasRaw)
		if err != nil {
			return query.ListMapelFilter{}, errors.New("tingkat_kelas must be a number")
		}
		filters.TingkatKelas = &tingkatKelas
	}

	if limitRaw := strings.TrimSpace(values.Get("limit")); limitRaw != "" {
		limit, err := strconv.Atoi(limitRaw)
		if err != nil {
			return query.ListMapelFilter{}, errors.New("limit must be a number")
		}
		filters.Limit = limit
	}

	if offsetRaw := strings.TrimSpace(values.Get("offset")); offsetRaw != "" {
		offset, err := strconv.Atoi(offsetRaw)
		if err != nil {
			return query.ListMapelFilter{}, errors.New("offset must be a number")
		}
		filters.Offset = offset
	}

	return filters, nil
}

func toListMapelResponse(items []matapelajaran.MataPelajaran) ListMapelResponse {
	response := ListMapelResponse{Items: make([]MapelResponse, 0, len(items))}
	for _, item := range items {
		response.Items = append(response.Items, toMapelResponse(item))
	}

	return response
}

func toMapelResponse(item matapelajaran.MataPelajaran) MapelResponse {
	return MapelResponse{
		IDMapel:   int(item.IdMapel),
		IDKelas:   int(item.IdKelas),
		KodeMapel: item.KodeMapel,
		NamaMapel: item.NamaMapel,
		Deskripsi: item.Deskripsi,
	}
}
