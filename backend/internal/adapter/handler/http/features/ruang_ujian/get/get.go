package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ruangujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/get"
)

type GetRuangUjianHandler struct {
	svc *ruangujian_service.GetRuangUjianRepo
}

func NewGetRuangUjianHandler(svc *ruangujian_service.GetRuangUjianRepo) *GetRuangUjianHandler {
	return &GetRuangUjianHandler{svc: svc}
}

func (h *GetRuangUjianHandler) GetRuangUjian(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	filters, err := parseListRuangUjianFilters(r)
	if err != nil {
		logger.Info(r.Context(), "invalid ruang ujian filters", "layer", "adapter.http.handler", "op", "ruang_ujian.list", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	items, err := h.svc.GetRuangUjian(r.Context(), filters)
	if err != nil {
		logger.Error(r.Context(), "failed listing ruang ujian", "layer", "adapter.http.handler", "op", "ruang_ujian.get_ruang_ujian", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get ruang ujian")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toRuangUjianResponses(items), "Success")
}

func (h *GetRuangUjianHandler) GetRuangUjianByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	rawId := strings.TrimSpace(params.ByName("IdRuangan"))
	IdRuangan, err := strconv.Atoi(rawId)
	if err != nil || IdRuangan <= 0 {
		logger.Info(r.Context(), "invalid id ruangan", "layer", "adapter.http.handler", "op", "ruang_ujian.get_by_id", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id ruangan")
		return
	}

	item, err := h.svc.GetRuangUjianById(r.Context(), IdRuangan)
	if err != nil {
		logger.Error(r.Context(), "failed get ruang ujian by id", "layer", "adapter.http.handler", "op", "ruang_ujian.get_by_id", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "MISSING_ID", "missing id")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get ruang ujian")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toRuangUjianResponse(item), "Success")
}

func (h *GetRuangUjianHandler) GetRuangUjianByKode(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	rawKode := strings.TrimSpace(params.ByName("KodeRuang"))

	item, err := h.svc.GetRuangUjianByKode(r.Context(), rawKode)
	if err != nil {
		logger.Error(r.Context(), "failed get ruang ujian by id", "layer", "adapter.http.handler", "op", "ruang_ujian.get_by_id", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingField):
			httpResponse.WriteErr(w, http.StatusBadRequest, "MISSING_FIELD", "missing field")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get ruang ujian")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toRuangUjianResponse(item), "Success")
}
