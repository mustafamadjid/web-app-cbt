package httpx

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	attemptujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/attempt_ujian/update"
)

type ExpireAttemptUjianHandler struct {
	svc *attemptujian_service.ExpireAttemptUjianService
}

func NewExpireAttemptUjianHandler(svc *attemptujian_service.ExpireAttemptUjianService) *ExpireAttemptUjianHandler {
	return &ExpireAttemptUjianHandler{svc: svc}
}

func (h *ExpireAttemptUjianHandler) ExpireAttemptUjian(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodPatch {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	if _, ok := middleware.ActorFromContext(r.Context()); !ok {
		logger.Error(r.Context(), "missing actor in context", "layer", "adapter.http.handler", "op", "ujian.attempt.expire_admin", "err", "actor_not_found")
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	idAttempt, err := parseExpireAttemptUjianID(ps)
	if err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id attempt")
		return
	}

	if err := h.svc.ExpireAttemptUjian(r.Context(), idAttempt); err != nil {
		logger.Error(r.Context(), "failed expire attempt ujian", "layer", "adapter.http.handler", "op", "ujian.attempt.expire_admin", "attempt_id", idAttempt, "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId),
			errors.Is(err, coreerror.ErrMissingField),
			errors.Is(err, coreerror.ErrNoFieldToUpdate),
			errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid attempt update payload")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		case errors.Is(err, coreerror.ErrSiswaHasActiveAttempt):
			httpResponse.WriteErr(w, http.StatusConflict, "DOUBLE_ATTEMPT_NOT_ALLOWED", "double attempt not allowed")
		case errors.Is(err, coreerror.ErrConflict):
			httpResponse.WriteErr(w, http.StatusConflict, "CONFLICT", "conflict")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "Success")
}
