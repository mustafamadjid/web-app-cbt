package httpx

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	mapel_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/delete"
)

type DeleteMapelHandler struct {
	svc *mapel_service.DeleteMapelRepo
}

func NewDeleteMapelHandler(svc *mapel_service.DeleteMapelRepo) *DeleteMapelHandler {
	return &DeleteMapelHandler{svc: svc}
}

func (h *DeleteMapelHandler) DeleteMapel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodDelete {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idMapel, err := strconv.Atoi(ps.ByName("idMapel"))
	if err != nil || idMapel <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id mapel")
		return
	}

	if err := h.svc.DeleteMapelService(r.Context(), idMapel); err != nil {
		logger.Error(r.Context(), "failed to delete mapel", "layer", "adapter.http.handler", "op", "mata_pelajaran.delete", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "not found")
			return
		case errors.Is(err, coreerror.ErrDeleteRestricted):
			httpResponse.WriteErr(w, http.StatusBadRequest, "DELETE_RESTRICTED", "delete restricted : constraint violation")
			return
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed delete mapel")
			return
		}
	}

	httpResponse.WriteOKNoData(w,http.StatusOK,"success")
}
