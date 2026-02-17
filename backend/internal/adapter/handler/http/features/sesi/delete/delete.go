package httpx

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	sesi_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/delete"
)

type DeleteSesiHandler struct {
	svc *sesi_service.DeleteSesiService
}

func NewDeleteSesiHandler(svc *sesi_service.DeleteSesiService) *DeleteSesiHandler {
	return &DeleteSesiHandler{svc: svc}
}

func (h *DeleteSesiHandler) DeleteSesi(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodDelete {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idSesi, err := strconv.Atoi(ps.ByName("idSesi"))
	if err != nil || idSesi <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id sesi")
		return
	}

	if err := h.svc.DeleteSesiService(r.Context(), idSesi); err != nil {
		logger.Error(r.Context(), "failed deleting sesi", "layer", "adapter.http.handler", "op", "sesi.delete", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		case errors.Is(err, coreerror.ErrDeleteRestricted):
			httpResponse.WriteErr(w, http.StatusBadRequest, "DELETE_RESTRICTED", "delete restricted : constraint violation")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed delete sesi")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "success")
}
