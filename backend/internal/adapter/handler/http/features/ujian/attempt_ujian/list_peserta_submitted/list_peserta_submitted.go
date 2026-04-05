package httpx

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	attemptujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/list_peserta_submitted"
)

type ListPesertaUjianSubmittedHandler struct {
	svc *attemptujian_service.PesertaUjianSubmittedService
}

func NewListPesertaUjianSubmittedHandler(svc *attemptujian_service.PesertaUjianSubmittedService) *ListPesertaUjianSubmittedHandler {
	return &ListPesertaUjianSubmittedHandler{svc: svc}
}

func (h *ListPesertaUjianSubmittedHandler) ListPesertaUjianSubmitted(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	req, err := parseListPesertaUjianSubmittedRequest(params)
	if err != nil {
		logger.Info(r.Context(), "invalid list peserta ujian submitted request", "layer", "adapter.http.handler", "op", "ujian.attempt.list_peserta_submitted", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	items, err := h.svc.ListPesertaUjianSubmitted(r.Context(), req.IDJadwalUjian)
	if err != nil {
		logger.Error(r.Context(), "failed list peserta ujian submitted", "layer", "adapter.http.handler", "op", "ujian.attempt.list_peserta_submitted", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "id jadwal ujian tidak valid")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "jadwal ujian not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get list peserta ujian submitted")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toListPesertaUjianSubmittedResponses(items), "Success")
}
