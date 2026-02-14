package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
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

func (h *CreateRuangUjianHandler)CreateRuangUian(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodPost {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct,"application/json"){
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	var dataRequest CreateRuangUjianrequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&dataRequest); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		return
	}

	if dec.Decode(&struct{}{}) != io.EOF {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		return
	}
	
	if strings.TrimSpace(dataRequest.NamaRuangan) == ""{
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		return
	}
	if strings.TrimSpace(dataRequest.KodeRuang) == ""{
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		return
	}

	if err := validator.ValidateInputSafe(dataRequest.NamaRuangan,"nama_ruangan"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	if err := validator.ValidateInputSafe(dataRequest.KodeRuang,"kode_ruang"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	data := ruangujian.RuangUjian {
		NamaRuangan: dataRequest.NamaRuangan,
		KodeRuang: dataRequest.KodeRuang,
	}

	if err := h.svc.CreateRuangUjianService(r.Context(), data); err != nil {
		logger.Error(r.Context(), "failed creating ruang ujian", "layer", "adapter.http.handler", "op", "ruang_ujian.create", "err", err)
		switch {
		case errors.Is(err,coreerror.ErrKodeRuangUjianExist):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: kode ruang ujian already exist")
			return
		case errors.Is(err,coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
			return
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err.Error())
			return
		}
	}

	httpResponse.WriteOKNoData(w,http.StatusOK,"success")
}