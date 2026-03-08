package httpx

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	aktivitas_user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/aktivitas_user"
)

type AktivitasUserHandler struct {
	svc *aktivitas_user_service.AktivitasUserService
}

func NewAktivitasUserHandler(svc *aktivitas_user_service.AktivitasUserService) *AktivitasUserHandler {
	return &AktivitasUserHandler{svc: svc}
}

func (h *AktivitasUserHandler) GetAktivitasUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	data, err := h.svc.GetAktivitasUserService(r.Context())
	if err != nil {
		logger.Error(r.Context(), "failed getting aktivitas user", "layer", "adapter.http.handler", "op", "aktivitas_user.get", "err", err)
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get aktivitas user")
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toAktivitasUserResponses(data), "success")
}
