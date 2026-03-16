package httpx

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	savejawaban_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/jawaban_ujian/save_jawaban"
)

type SaveJawabanUjianHandler struct {
	svc *savejawaban_service.JawabanUjianService
}

func NewSaveJawabanUjianHandler(svc *savejawaban_service.JawabanUjianService) *SaveJawabanUjianHandler {
	return &SaveJawabanUjianHandler{svc: svc}
}

func (h *SaveJawabanUjianHandler) SaveJawabanUjian(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodPost {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	if err := httphelper.JsonHeaderBodyValidator(w, r); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: content type must be application/json")
		return
	}

	var req SaveJawabanRequest
	if err := httphelper.JSONDecoder(r, &req); err != nil {
		switch {
		case errors.Is(err, coreerror.ErrInvalidRequestBody):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	idAttempt, jawaban := toSaveJawabanPayload(req)
	if err := h.svc.SaveJawabanUjian(r.Context(), idAttempt, jawaban); err != nil {
		logger.Error(r.Context(), "failed save jawaban ujian", "layer", "adapter.http.handler", "op", "ujian.jawaban.save", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId),
			errors.Is(err, coreerror.ErrMissingJawabanEssayAndPilgan),
			errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid jawaban ujian payload")
		case errors.Is(err, coreerror.ErrConflict):
			httpResponse.WriteErr(w, http.StatusConflict, "CONFLICT", "conflict")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "Success")
}
