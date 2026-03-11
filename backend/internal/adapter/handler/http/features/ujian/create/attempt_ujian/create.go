package httpx

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	siswaujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian"
)

type AttemptUjianHandler struct {
	svc *siswaujian_service.AttemptUjianService
}

func NewAttemptUjianHandler(svc *siswaujian_service.AttemptUjianService) *AttemptUjianHandler {
	return &AttemptUjianHandler{svc: svc}
}

func (h *AttemptUjianHandler) AttemptUjian(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodPost {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	if err := httphelper.JsonHeaderBodyValidator(w, r); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	var req AttemptUjianRequest
	if err := httphelper.JSONDecoder(r, &req); err != nil {
		switch {
		case errors.Is(err, coreerror.ErrInvalidRequestBody):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	err := h.svc.AttemptUjian(
		r.Context(),
		req.IdSiswa,
		req.IdJadwalUjian,
		req.TokenUjian,
		req.WaktuMulai,
	)
	if err != nil {
		logger.Error(r.Context(), "failed create attempt ujian", "layer", "adapter.http.handler", "op", "ujian.attempt.create", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: missing id")
		case errors.Is(err, coreerror.ErrMissingTokenUjian):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: missing token ujian")
		case errors.Is(err, coreerror.ErrTimeEmpty):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: waktu_mulai is required")
		case errors.Is(err, coreerror.ErrPesertaNotAllowedToAttemptJadwal):
			httpResponse.WriteErr(w, http.StatusForbidden, "SISWA_NOT_ALLOWED", "peserta ujian is not allowed to attempt this ujian")
		case errors.Is(err, coreerror.ErrTokenUjianInvalid):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_TOKEN_UJIAN", "bad request: invalid token ujian")
		case errors.Is(err, coreerror.ErrWaktuAttemptPesertaInvalid):
			httpResponse.WriteErr(w, http.StatusBadRequest, "UJIAN_ATTEMPT_TIME_EXPIRED", "bad request: expired attempt time")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		case errors.Is(err, coreerror.ErrConflict):
			httpResponse.WriteErr(w, http.StatusConflict, "CONFLICT", "conflict")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "Success")
}
