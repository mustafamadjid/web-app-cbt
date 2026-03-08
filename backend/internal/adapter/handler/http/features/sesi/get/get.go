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
	sesi_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/get"
)

type GetSesiHandler struct {
	svc *sesi_service.GetSesiService
}

func NewGetSesiHandler(svc *sesi_service.GetSesiService) *GetSesiHandler {
	return &GetSesiHandler{svc: svc}
}

func (h *GetSesiHandler) ListSesi(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	filters, err := parseListSesiFilters(r)
	if err != nil {
		logger.Info(r.Context(), "invalid sesi filters", "layer", "adapter.http.handler", "op", "sesi.list", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	items, err := h.svc.GetSesiService(r.Context(), filters)
	if err != nil {
		logger.Error(r.Context(), "failed listing sesi", "layer", "adapter.http.handler", "op", "sesi.list", "err", err)
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get sesi")
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toListSesiResponse(items), "Success")
}

func (h *GetSesiHandler) GetSesiByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	rawID := strings.TrimSpace(params.ByName("idSesi"))
	idSesi, err := strconv.Atoi(rawID)
	if err != nil || idSesi <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id sesi")
		return
	}

	item, err := h.svc.GetSesiByIdService(r.Context(), idSesi)
	if err != nil {
		logger.Error(r.Context(), "failed get sesi by id", "layer", "adapter.http.handler", "op", "sesi.get_by_id", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get sesi")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toSesiResponse(item), "Success")
}

func (h *GetSesiHandler) GetSesiByKode(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	kodeSesi := strings.TrimSpace(params.ByName("kodeSesi"))
	if err := validator.ValidateInputSafe(kodeSesi, "kode_sesi"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	item, err := h.svc.GetSesiByKodeService(r.Context(), kodeSesi)
	if err != nil {
		logger.Error(r.Context(), "failed get sesi by kode", "layer", "adapter.http.handler", "op", "sesi.get_by_kode", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingField):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: kode sesi is required")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get sesi")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toSesiResponse(item), "Success")
}
