package httpx

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	activeattempt_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/active_attempt"
)

type GetActiveAttemptUjianHandler struct {
	svc *activeattempt_service.GetActiveAttemptUjianService
}

func NewGetActiveAttemptUjianHandler(svc *activeattempt_service.GetActiveAttemptUjianService) *GetActiveAttemptUjianHandler {
	return &GetActiveAttemptUjianHandler{svc: svc}
}

func (h *GetActiveAttemptUjianHandler) GetActiveAttemptUjian(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	actor, ok := middleware.ActorFromContext(r.Context())
	if !ok {
		logger.Error(r.Context(), "missing actor in context", "layer", "adapter.http.handler", "op", "siswa_ujian.get_active_attempt", "err", "actor_not_found")
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	req, err := parseGetActiveAttemptUjianRequest(r)
	if err != nil {
		logger.Info(r.Context(), "invalid get active attempt ujian request", "layer", "adapter.http.handler", "op", "siswa_ujian.get_active_attempt", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	item, err := h.svc.GetActiveAttemptUjian(r.Context(), int(actor.IdPengguna), req.IDJadwalUjian)
	if err != nil {
		logger.Error(r.Context(), "failed get active attempt ujian", "layer", "adapter.http.handler", "op", "siswa_ujian.get_active_attempt", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "id jadwal ujian tidak valid")
		case errors.Is(err, sql.ErrNoRows):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "active attempt not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get active attempt ujian")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toGetActiveAttemptUjianResponse(item), "Success")
}
