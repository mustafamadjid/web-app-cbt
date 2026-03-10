package httpx

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/sesi"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	sesi_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/create"
)

type CreateSesiHandler struct {
	svc *sesi_service.CreateSesiService
}

func NewCreateSesiHandler(svc *sesi_service.CreateSesiService) *CreateSesiHandler {
	return &CreateSesiHandler{svc: svc}
}

func (h *CreateSesiHandler) CreateSesiHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodPost {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	if err := httphelper.JsonHeaderBodyValidator(w, r); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	var dataRequest CreateSesiRequest
	if err := httphelper.JSONDecoder(r, &dataRequest); err != nil {
		switch {
		case errors.Is(err, coreerror.ErrInvalidRequestBody):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	if dataRequest.NamaSesi == "" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: nama sesi is required")
		return
	}
	namaSesi, err := validator.ValidateRequiredPrintableText(dataRequest.NamaSesi, "nama_sesi")
	if err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	dataRequest.NamaSesi = namaSesi

	if dataRequest.KodeSesi == "" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: kode sesi is required")
		return
	}
	kodeSesi, err := validator.ValidateRequiredPrintableText(dataRequest.KodeSesi, "kode_sesi")
	if err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	dataRequest.KodeSesi = kodeSesi

	if err := h.svc.CreateSesiService(r.Context(), sesi.Sesi{
		NamaSesi: dataRequest.NamaSesi,
		KodeSesi: dataRequest.KodeSesi,
	}); err != nil {
		logger.Error(r.Context(), "failed creating sesi", "layer", "adapter.http.handler", "op", "sesi.create", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrSesiUjianExist):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: kode sesi already exist")
		case errors.Is(err, coreerror.ErrInvalidInput), errors.Is(err, coreerror.ErrInvalidInputSafe), errors.Is(err, coreerror.ErrMissingField):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed create sesi")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "success")
}
