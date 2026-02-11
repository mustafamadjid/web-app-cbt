package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	kelas_domain "github.com/mustafamadjid/web-app-cbt/internal/core/domain/kelas"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/kelas"
	kelas_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/get"
)

type GetKelasHandler struct {
	svc *kelas_service.GetKelasService
}

func NewGetKelasHandler(svc *kelas_service.GetKelasService) *GetKelasHandler {
	return &GetKelasHandler{svc: svc}
}

func (h *GetKelasHandler) ListKelas(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	filters, err := parseListKelasFilters(r)
	if err != nil {
		logger.Info(r.Context(), "invalid kelas filters", "layer", "adapter.http.handler", "op", "kelas.list", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	items, err := h.svc.GetFullKelas(r.Context(), filters)
	if err != nil {
		logger.Error(r.Context(), "failed listing kelas", "layer", "adapter.http.handler", "op", "kelas.list", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get kelas")
		}
		return
	}

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

	httpResponse.WriteOK(w, http.StatusOK, response, "Success")
}

func parseListKelasFilters(r *http.Request) (query.ListKelasFilter, error) {
	values := r.URL.Query()
	filters := query.ListKelasFilter{}

	filters.Search = strings.TrimSpace(values.Get("q"))
	if filters.Search == "" {
		filters.Search = strings.TrimSpace(values.Get("search"))
	}
	if err := validator.ValidateInputSafe(filters.Search, "search"); err != nil {
		return query.ListKelasFilter{}, err
	}

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

func (h *GetKelasHandler) GetKelasByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	rawID := strings.TrimSpace(params.ByName("id"))
	if rawID == "" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "Bad request : id is required")
		return
	}

	parsedID, err := strconv.Atoi(rawID)
	if err != nil || parsedID <= 0 {
		logger.Info(r.Context(), "invalid kelas id", "layer", "adapter.http.handler", "op", "kelas.get_by_id", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "Bad request : invalid id")
		return
	}

	item, err := h.svc.GetKelasByID(r.Context(), kelas_domain.ID(parsedID))
	if err != nil {
		logger.Error(r.Context(), "failed getting kelas by id", "layer", "adapter.http.handler", "op", "kelas.get_by_id", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get kelas")
		}
		return
	}

	response := DataKelasResponse{
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

	httpResponse.WriteOK(w, http.StatusOK, response, "Success")
}
