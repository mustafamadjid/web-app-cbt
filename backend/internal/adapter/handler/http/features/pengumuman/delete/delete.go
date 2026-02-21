package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	pengumuman_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/delete"
)

type DeletePengumumanHandler struct {
	svc *pengumuman_service.DeletePengumumanService
}

func NewDeletePengumumanHandler(svc *pengumuman_service.DeletePengumumanService) *DeletePengumumanHandler {
	return &DeletePengumumanHandler{svc: svc}
}

func (h *DeletePengumumanHandler) DeletePengumuman(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodDelete {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idPengumuman, err := strconv.Atoi(strings.TrimSpace(ps.ByName("idPengumuman")))
	if err != nil || idPengumuman <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id pengumuman")
		return
	}

	if err := h.svc.DeletePengumumanService(r.Context(), pengumuman.ID(idPengumuman)); err != nil {
		logger.Error(r.Context(), "failed deleting pengumuman", "layer", "adapter.http.handler", "op", "pengumuman.delete", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: missing id")
		case errors.Is(err, coreerror.ErrDeleteRestricted):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "delete restricted : constraint violation")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed delete pengumuman")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "success")
}
