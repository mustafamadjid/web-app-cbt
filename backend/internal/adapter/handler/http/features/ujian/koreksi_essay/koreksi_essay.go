package httpx

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	essaygrading_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/grading/essay_grading"
)

type KoreksiEssayHandler struct {
	svc *essaygrading_service.EssayGradingUjianService
}

func NewKoreksiEssayHandler(svc *essaygrading_service.EssayGradingUjianService) *KoreksiEssayHandler {
	return &KoreksiEssayHandler{svc: svc}
}

func (h *KoreksiEssayHandler) KoreksiEssay(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodPatch {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	actor, ok := middleware.ActorFromContext(r.Context())
	if !ok {
		logger.Error(r.Context(), "missing actor in context", "layer", "adapter.http.handler", "op", "ujian.koreksi_essay", "err", "actor_not_found")
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	if err := httphelper.JsonHeaderBodyValidator(w, r); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: content type must be application/json")
		return
	}

	var req KoreksiEssayRequest
	if err := httphelper.JSONDecoder(r, &req); err != nil {
		switch {
		case errors.Is(err, coreerror.ErrInvalidRequestBody):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}
	req, err := sanitizeAndValidateKoreksiEssayRequest(req)
	if err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid essay grading payload")
		return
	}

	if err := h.svc.EssayGrading(r.Context(), toKoreksiEssayPayload(req), ujian.ID(actor.IdPengguna)); err != nil {
		logger.Error(r.Context(), "failed koreksi essay ujian", "layer", "adapter.http.handler", "op", "ujian.koreksi_essay", "graded_by", actor.IdPengguna, "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId),
			errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid essay grading payload")
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
