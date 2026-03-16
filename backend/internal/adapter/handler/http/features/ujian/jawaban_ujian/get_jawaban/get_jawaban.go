package httpx

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	getjawaban_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/jawaban_ujian/get_jawaban"
)

type GetJawabanUjianHandler struct {
	svc *getjawaban_service.SiswaGetJawabanUjianService
}

func NewGetJawabanUjianHandler(svc *getjawaban_service.SiswaGetJawabanUjianService) *GetJawabanUjianHandler {
	return &GetJawabanUjianHandler{svc: svc}
}

func (h *GetJawabanUjianHandler) GetJawabanUjian(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	actor, ok := middleware.ActorFromContext(r.Context())
	if !ok {
		logger.Error(r.Context(), "missing actor in context", "layer", "adapter.http.handler", "op", "ujian.jawaban.get_by_attempt_id_siswa", "err", "actor_not_found")
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	req, err := parseGetJawabanUjianRequest(params)
	if err != nil {
		logger.Info(r.Context(), "invalid get jawaban ujian request", "layer", "adapter.http.handler", "op", "ujian.jawaban.get_by_attempt_id_siswa", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id attempt")
		return
	}

	items, err := h.svc.GetJawabanUjianByAttemptId(r.Context(), int(actor.IdPengguna), ujian.ID(req.IDAttempt))
	if err != nil {
		logger.Error(r.Context(), "failed get jawaban ujian by attempt id", "layer", "adapter.http.handler", "op", "ujian.jawaban.get_by_attempt_id_siswa", "attempt_id", req.IDAttempt, "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id attempt")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toGetJawabanUjianResponse(req.IDAttempt, items), "Success")
}
