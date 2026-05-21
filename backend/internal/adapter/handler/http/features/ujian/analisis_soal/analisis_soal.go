package httpx

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	analisissoal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/analisis_soal"
)

type AnalisisSoalHandler struct {
	svc *analisissoal_service.AnalisisSoalService
}

func NewAnalisisSoalHandler(svc *analisissoal_service.AnalisisSoalService) *AnalisisSoalHandler {
	return &AnalisisSoalHandler{svc: svc}
}

func (h *AnalisisSoalHandler) GetListAnalisisSoal(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	req, err := parseGetListAnalisisSoalRequest(params)
	if err != nil {
		logger.Info(r.Context(), "invalid get list analisis soal request", "layer", "adapter.http.handler", "op", "ujian.analisis_soal.list", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	items, err := h.svc.GetListAnalisisSoal(r.Context(), req.IDJadwalUjian)
	if err != nil {
		logger.Error(r.Context(), "failed get list analisis soal", "layer", "adapter.http.handler", "op", "ujian.analisis_soal.list", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "id jadwal ujian tidak valid")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "analisis soal not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get list analisis soal")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toGetListAnalisisSoalResponse(req.IDJadwalUjian, items), "Success")
}
