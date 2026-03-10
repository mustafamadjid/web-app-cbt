package httpx

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	sesi_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/update"
)

type UpdateSesiHandler struct {
	svc *sesi_service.UpdateSesiService
}

func NewUpdateSesiHandler(svc *sesi_service.UpdateSesiService) *UpdateSesiHandler {
	return &UpdateSesiHandler{svc: svc}
}

func (h *UpdateSesiHandler) UpdateSesi(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodPatch {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idSesi, err := strconv.Atoi(ps.ByName("idSesi"))
	if err != nil || idSesi <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id sesi")
		return
	}

	if err := httphelper.JsonHeaderBodyValidator(w, r); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	var dataRequest UpdateSesiRequest
	if err := httphelper.JSONDecoder(r, &dataRequest); err != nil {
		switch {
		case errors.Is(err, coreerror.ErrInvalidRequestBody):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	if dataRequest.KodeSesi != nil {
		value, err := validator.ValidateRequiredPrintableText(*dataRequest.KodeSesi, "kode_sesi")
		if err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
			return
		}
		*dataRequest.KodeSesi = value
	}

	if dataRequest.NamaSesi != nil {
		value, err := validator.ValidateRequiredPrintableText(*dataRequest.NamaSesi, "nama_sesi")
		if err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
			return
		}
		*dataRequest.NamaSesi = value
	}

	patch := updatepatch.UpdateSesiPatch{
		KodeSesi: dataRequest.KodeSesi,
		NamaSesi: dataRequest.NamaSesi,
	}

	if err := h.svc.UpdateSesiService(r.Context(), idSesi, patch); err != nil {
		logger.Error(r.Context(), "failed updating sesi", "layer", "adapter.http.handler", "op", "sesi.update", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		case errors.Is(err, coreerror.ErrMissingId), errors.Is(err, coreerror.ErrMissingField), errors.Is(err, coreerror.ErrNoFieldToUpdate):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid field")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed update sesi")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "success")
}
