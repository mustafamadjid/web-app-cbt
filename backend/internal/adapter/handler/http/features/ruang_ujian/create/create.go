package httpx

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ruangujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ruang_ujian"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ruangujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/create"
)

type CreateRuangUjianHandler struct {
	svc *ruangujian_service.CreateRuangUjianService
}

func NewCreateRuangUjianHandler(svc *ruangujian_service.CreateRuangUjianService) *CreateRuangUjianHandler {
	return &CreateRuangUjianHandler{svc: svc}
}

func (h *CreateRuangUjianHandler) CreateRuangUjian(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodPost {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	if err := httphelper.JsonHeaderBodyValidator(w, r); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	var dataRequest CreateRuangUjianrequest
	if err := httphelper.JSONDecoder(r, &dataRequest); err != nil {
		switch {
		case errors.Is(err, coreerror.ErrInvalidRequestBody):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	if dataRequest.NamaRuangan == "" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		return
	}
	if dataRequest.KodeRuang == "" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		return
	}

	if err := validator.ValidateInputSafe(dataRequest.NamaRuangan, "nama_ruangan"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	if err := validator.ValidateInputSafe(dataRequest.KodeRuang, "kode_ruang"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	data := ruangujian.RuangUjian{
		NamaRuangan: dataRequest.NamaRuangan,
		KodeRuang:   dataRequest.KodeRuang,
	}

	if err := h.svc.CreateRuangUjianService(r.Context(), data); err != nil {
		logger.Error(r.Context(), "failed creating ruang ujian", "layer", "adapter.http.handler", "op", "ruang_ujian.create", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrKodeRuangUjianExist):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: kode ruang ujian already exist")
			return
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid input")
			return
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
			return
		}
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "success")
}
