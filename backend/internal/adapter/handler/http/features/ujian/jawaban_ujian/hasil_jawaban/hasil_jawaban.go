package httpx

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	hasiljawaban_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/jawaban_ujian/hasil_jawaban"
)

type HasilJawabanUjianHandler struct {
	svc *hasiljawaban_service.HasilJawabanUjianService
}

func NewHasilJawabanUjianHandler(svc *hasiljawaban_service.HasilJawabanUjianService) *HasilJawabanUjianHandler {
	return &HasilJawabanUjianHandler{svc: svc}
}

func (h *HasilJawabanUjianHandler) ListHasilJawabanUjianByAttempt(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idAttempt, err := parseHasilJawabanUjianIDAttempt(params)
	if err != nil {
		logger.Info(r.Context(), "invalid list hasil jawaban ujian request", "layer", "adapter.http.handler", "op", "ujian.jawaban.list_hasil_by_attempt_id", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id attempt")
		return
	}

	items, err := h.svc.ListHasilJawabanUjianByAttempt(r.Context(), idAttempt)
	if err != nil {
		logger.Error(r.Context(), "failed list hasil jawaban ujian by attempt id", "layer", "adapter.http.handler", "op", "ujian.jawaban.list_hasil_by_attempt_id", "attempt_id", idAttempt, "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId),
			errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id attempt")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toHasilJawabanUjianResponse(idAttempt, items), "Success")
}
