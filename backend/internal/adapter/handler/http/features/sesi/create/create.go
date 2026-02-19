package httpx

import (
	"errors"
	"net/http"
	"strings"

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

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

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

	if strings.TrimSpace(dataRequest.NamaSesi) == "" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: nama sesi is required")
		return
	}
	if err := validator.ValidateInputSafe(dataRequest.NamaSesi, "nama_sesi"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	if strings.TrimSpace(dataRequest.KodeSesi) == "" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: kode sesi is required")
		return
	}
	if err := validator.ValidateInputSafe(dataRequest.KodeSesi, "kode_sesi"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

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
