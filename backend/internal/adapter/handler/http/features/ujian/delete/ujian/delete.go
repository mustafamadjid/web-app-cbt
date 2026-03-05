package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/delete"
)

type DeleteUjianHandler struct {
	svc *ujian_service.DeleteUjianService
}

func NewDeleteUjianHandler(svc *ujian_service.DeleteUjianService) *DeleteUjianHandler {
	return &DeleteUjianHandler{svc: svc}
}

func (h *DeleteUjianHandler) DeleteUjian(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodDelete {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	rawIDUjian := strings.TrimSpace(ps.ByName("idUjian"))
	idUjian, err := strconv.Atoi(rawIDUjian)
	if err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id ujian")
		return
	}

	data := DeleteUjianRequest{IDUjian: idUjian}
	if err := ValidateInputDeleteUjianRequest(data); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id ujian")
		return
	}

	if err := h.svc.DeleteUjianService(r.Context(), ujian.ID(data.IDUjian)); err != nil {
		logger.Error(r.Context(), "failed deleting ujian", "layer", "adapter.http.handler", "op", "ujian.delete", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		case errors.Is(err, coreerror.ErrDeleteRestricted):
			httpResponse.WriteErr(w, http.StatusBadRequest, "DELETE_RESTRICTED", "delete restricted : constraint violation")
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id ujian")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed delete ujian")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "success")
}
