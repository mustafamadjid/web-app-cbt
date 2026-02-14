package httpx

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ruangujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/delete"
)

type DeleteRuangUjianHandler struct {
	svc *ruangujian_service.DeleteRuangUjianService
}

func NewDeleteRuangUjianHandler(svc *ruangujian_service.DeleteRuangUjianService) *DeleteRuangUjianHandler {
	return &DeleteRuangUjianHandler{svc: svc}
}

func(h *DeleteRuangUjianHandler)DeleteRuangUjian(w http.ResponseWriter,r *http.Request, ps httprouter.Params)  {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodDelete {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idRuangan, err := strconv.Atoi(ps.ByName("idRuangan"))
	if err != nil || idRuangan <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id ruangan")
		return
	}

	if err := h.svc.DeleteRuangUjian(r.Context(), idRuangan); err != nil {
		logger.Error(r.Context(), "failed to delete ruang ujian", "layer", "adapter.http.handler", "op", "ruang_ujian.delete", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "not found")
			return
		case errors.Is(err, coreerror.ErrDeleteRestricted):
			httpResponse.WriteErr(w, http.StatusBadRequest, "DELETE_RESTRICTED", "delete restricted : constraint violation")
			return
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed delete ruang ujian")
			return
		}
	}

	httpResponse.WriteOKNoData(w,http.StatusOK,"success")

}