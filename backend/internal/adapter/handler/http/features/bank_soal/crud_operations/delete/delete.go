package httpx

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/bank_soal"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	bank_soal_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/bank_soal/delete"
)

type DeleteBankSoalHandler struct {
	svc *bank_soal_service.DeleteBankSoalService
}

func NewDeleteBankSoalHandler(svc *bank_soal_service.DeleteBankSoalService) *DeleteBankSoalHandler {
	return &DeleteBankSoalHandler{svc: svc}
}

func (h *DeleteBankSoalHandler) DeleteBankSoal(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodDelete {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idBankSoal, err := strconv.Atoi(ps.ByName("idBankSoal"))
	if err != nil || idBankSoal <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id bank soal")
		return
	}

	if err := h.svc.DeleteBankSoalService(r.Context(), bank_soal.ID(idBankSoal)); err != nil {
		logger.Error(r.Context(), "failed deleting bank soal", "layer", "adapter.http.handler", "op", "bank_soal.delete", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id bank soal")
		case errors.Is(err, coreerror.ErrDeleteRestricted):
			httpResponse.WriteErr(w, http.StatusBadRequest, "DELETE_RESTRICTED", "delete restricted : constraint violation")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "Success")
}
