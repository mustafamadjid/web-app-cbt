package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ruang_ujian"
)

func parseListRuangUjianFilters(r *http.Request) (query.ListRuangUjianFilter, error) {
	values := r.URL.Query()
	filters := query.ListRuangUjianFilter{}

	filters.Search = strings.TrimSpace(values.Get("q"))
	if filters.Search == "" {
		filters.Search = strings.TrimSpace(values.Get("search"))
	}
	search, err := validator.ValidateOptionalPrintableText(filters.Search, "search")
	if err != nil {
		return query.ListRuangUjianFilter{}, err
	}
	filters.Search = search

	if limitRaw := strings.TrimSpace(values.Get("limit")); limitRaw != "" {
		limit, err := strconv.Atoi(limitRaw)
		if err != nil {
			return query.ListRuangUjianFilter{}, errors.New("limit must be a number")
		}
		filters.Limit = limit
	}

	if offsetRaw := strings.TrimSpace(values.Get("offset")); offsetRaw != "" {
		offset, err := strconv.Atoi(offsetRaw)
		if err != nil {
			return query.ListRuangUjianFilter{}, errors.New("offset must be a number")
		}
		filters.Offset = offset
	}

	return filters, nil
}

func toRuangUjianResponses(items []ruangujian.RuangUjian) []RuangUjianResponse {
	response := make([]RuangUjianResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toRuangUjianResponse(item))
	}

	return response
}

func toRuangUjianResponse(item ruangujian.RuangUjian) RuangUjianResponse {
	return RuangUjianResponse{
		IdRuangan:   int(item.IdRuangan),
		NamaRuangan: item.NamaRuangan,
		KodeRuang:   item.KodeRuang,
	}
}
