package httpx

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	statistikujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/statistik_ujian"
)

type GetStatistikUjianHandler struct {
	svc *statistikujian_service.StatistikUjianService
}

func NewGetStatistikUjianHandler(svc *statistikujian_service.StatistikUjianService) *GetStatistikUjianHandler {
	return &GetStatistikUjianHandler{svc: svc}
}

func (h *GetStatistikUjianHandler) GetStatistikUjian(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	req, err := parseGetStatistikUjianRequest(params)
	if err != nil {
		logger.Info(r.Context(), "invalid get statistik ujian request", "layer", "adapter.http.handler", "op", "ujian.statistik.get", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	item, err := h.svc.GetStatistikUjian(r.Context(), req.IDJadwalUjian)
	if err != nil {
		logger.Error(r.Context(), "failed get statistik ujian", "layer", "adapter.http.handler", "op", "ujian.statistik.get", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "id jadwal ujian tidak valid")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "statistik ujian not found")
		case errors.Is(err, coreerror.ErrConflict):
			httpResponse.WriteErr(w, http.StatusConflict, "CONFLICT", "statistik ujian conflict")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get statistik ujian")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toGetStatistikUjianResponse(item), "Success")
}
