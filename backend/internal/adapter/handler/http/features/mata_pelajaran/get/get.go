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
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/mata_pelajaran"
	mapel_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/get"
)

type GetMapelHandler struct {
	svc *mapel_service.GetMapelRepo
}

func NewGetMapelHandler(svc *mapel_service.GetMapelRepo) *GetMapelHandler {
	return &GetMapelHandler{svc: svc}
}

func (h *GetMapelHandler) ListMapel(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	filters, err := parseListMapelFilters(r)
	if err != nil {
		logger.Info(r.Context(), "invalid mapel filters", "layer", "adapter.http.handler", "op", "mata_pelajaran.list", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	items, err := h.svc.GetMapelService(r.Context(), filters)
	if err != nil {
		logger.Error(r.Context(), "failed listing mapel", "layer", "adapter.http.handler", "op", "mata_pelajaran.list", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get mata pelajaran")
		}
		return
	}

	response := ListMapelResponse{Items: make([]MapelResponse, 0, len(items))}
	for _, item := range items {
		response.Items = append(response.Items, MapelResponse{
			IDMapel:   int(item.IdMapel),
			IDKelas:   int(item.IdKelas),
			KodeMapel: item.KodeMapel,
			NamaMapel: item.NamaMapel,
			Deskripsi: item.Deskripsi,
		})
	}

	httpResponse.WriteOK(w, http.StatusOK, response, "Success")
}

func (h *GetMapelHandler) GetMapelByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	rawID := strings.TrimSpace(params.ByName("idMapel"))
	idMapel, err := strconv.Atoi(rawID)
	if err != nil || idMapel <= 0 {
		logger.Info(r.Context(), "invalid id mapel", "layer", "adapter.http.handler", "op", "mata_pelajaran.get_by_id", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id mapel")
		return
	}

	item, err := h.svc.GetMapelById(r.Context(), idMapel)
	if err != nil {
		logger.Error(r.Context(), "failed get mapel by id", "layer", "adapter.http.handler", "op", "mata_pelajaran.get_by_id", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get mata pelajaran")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, MapelResponse{
		IDMapel:   int(item.IdMapel),
		IDKelas:   int(item.IdKelas),
		KodeMapel: item.KodeMapel,
		NamaMapel: item.NamaMapel,
		Deskripsi: item.Deskripsi,
	}, "Success")
}

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
