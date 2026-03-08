package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/sesi"
)

func parseListSesiFilters(r *http.Request) (query.ListSesiFilter, error) {
	values := r.URL.Query()
	filters := query.ListSesiFilter{}

	filters.Search = strings.TrimSpace(values.Get("q"))
	if filters.Search == "" {
		filters.Search = strings.TrimSpace(values.Get("search"))
	}
	if err := validator.ValidateInputSafe(filters.Search, "search"); err != nil {
		return query.ListSesiFilter{}, err
	}

	if limitRaw := strings.TrimSpace(values.Get("limit")); limitRaw != "" {
		limit, err := strconv.Atoi(limitRaw)
		if err != nil {
			return query.ListSesiFilter{}, errors.New("limit must be a number")
		}
		filters.Limit = limit
	}

	if offsetRaw := strings.TrimSpace(values.Get("offset")); offsetRaw != "" {
		offset, err := strconv.Atoi(offsetRaw)
		if err != nil {
			return query.ListSesiFilter{}, errors.New("offset must be a number")
		}
		filters.Offset = offset
	}

	return filters, nil
}

func toListSesiResponse(items []sesi.Sesi) ListSesiResponse {
	response := ListSesiResponse{Items: make([]SesiResponse, 0, len(items))}
	for _, item := range items {
		response.Items = append(response.Items, toSesiResponse(item))
	}

	return response
}

func toSesiResponse(item sesi.Sesi) SesiResponse {
	return SesiResponse{
		IdSesi:   int(item.IdSesi),
		KodeSesi: item.KodeSesi,
		NamaSesi: item.NamaSesi,
	}
}
