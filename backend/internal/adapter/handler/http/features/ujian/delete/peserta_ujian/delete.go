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

type DeletePesertaUjianHandler struct {
	svc *ujian_service.DeletePesertaUjianService
}

func NewDeletePesertaUjianHandler(svc *ujian_service.DeletePesertaUjianService) *DeletePesertaUjianHandler {
	return &DeletePesertaUjianHandler{svc: svc}
}

func (h *DeletePesertaUjianHandler) DeletePesertaUjian(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodDelete {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	rawIDPesertaUjian := strings.TrimSpace(ps.ByName("idPesertaUjian"))
	idPesertaUjian, err := strconv.Atoi(rawIDPesertaUjian)
	if err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id peserta ujian")
		return
	}

	data := DeletePesertaUjianRequest{IDPesertaUjian: idPesertaUjian}
	if err := ValidateInputDeletePesertaUjianRequest(data); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id peserta ujian")
		return
	}

	if err := h.svc.DeletePesertaUjianService(r.Context(), ujian.ID(data.IDPesertaUjian)); err != nil {
		logger.Error(r.Context(), "failed deleting peserta ujian", "layer", "adapter.http.handler", "op", "ujian.delete_peserta", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		case errors.Is(err, coreerror.ErrDeleteRestricted):
			httpResponse.WriteErr(w, http.StatusBadRequest, "DELETE_RESTRICTED", "delete restricted : constraint violation")
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id peserta ujian")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed delete peserta ujian")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "success")
}
