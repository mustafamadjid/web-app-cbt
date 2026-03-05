package httpx

import (
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/create"
)

type CreatePesertaUjianHandler struct {
	svc *ujian_service.CreatePesertaUjianService
}

func NewCreatePesertaUjianHandler(svc *ujian_service.CreatePesertaUjianService) *CreatePesertaUjianHandler {
	return &CreatePesertaUjianHandler{svc: svc}
}

func(s *CreatePesertaUjianHandler)CreatePesertaUjian(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodPost {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	if err := httphelper.JsonHeaderBodyValidator(w,r); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	var dataRequest CreatePesertaUjianRequest
	if err := httphelper.JSONDecoder(r, &dataRequest); err != nil {
		switch {
		case errors.Is(err, coreerror.ErrInvalidRequestBody):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		case strings.Contains(err.Error(), "parsing time") || strings.Contains(err.Error(), "cannot parse"):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid time format")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	if err:= ValidateInputIdRequestPesertUjian(dataRequest); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	data := ujian.PesertaUjian{
		IdJadwalUjian: ujian.ID(dataRequest.IdJadwalUjian),
		IdSiswa: ujian.ID(dataRequest.IdSiswa),
		WaktuMulai: dataRequest.WaktuMulai,
	}

	_,err := s.svc.CreatePesertaUjianService(r.Context(), data)
	if err != nil {
		logger.Error(r.Context(), "failed create peserta ujian", "layer", "core.service", "op", "ujian.create_peserta", "err", err)
		switch {
		case errors.Is(err,coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: missing id")
			return
		case errors.Is(err,coreerror.ErrMissingField):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: missing field")
			return
		case errors.Is(err,coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid input")
			return
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
			return
		}
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "Success",)
}