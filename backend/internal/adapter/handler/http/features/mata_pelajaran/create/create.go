package httpx

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	mapel_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/create"
)

type CreateMapelHandler struct {
	svc *mapel_service.CreateMapelRepo
}

func NewCreateMapelHandler(svc *mapel_service.CreateMapelRepo) *CreateMapelHandler {
	return &CreateMapelHandler{svc: svc}
}

func (h *CreateMapelHandler) CreateMapel(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodPost {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	if err := httphelper.JsonHeaderBodyValidator(w, r); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	var dataRequest CreateMapelRequest
	if err := httphelper.JSONDecoder(r, &dataRequest); err != nil {
		switch {
		case errors.Is(err, coreerror.ErrInvalidRequestBody):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	if dataRequest.IdKelas <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: id kelas is required")
		return
	}

	if dataRequest.KodeMapel == "" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: kode mapel is required")
		return
	}
	kodeMapel, err := validator.ValidateRequiredPrintableText(dataRequest.KodeMapel, "kode_mapel")
	if err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	dataRequest.KodeMapel = kodeMapel

	if dataRequest.NamaMapel == "" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: nama mapel is required")
		return
	}
	namaMapel, err := validator.ValidateRequiredPrintableText(dataRequest.NamaMapel, "nama_mapel")
	if err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	dataRequest.NamaMapel = namaMapel

	if dataRequest.Deskripsi == "" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: deskripsi is required")
		return
	}
	deskripsi, err := validator.ValidateRequiredPrintableText(dataRequest.Deskripsi, "deskripsi")
	if err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}
	dataRequest.Deskripsi = deskripsi

	cmdData := matapelajaran.MataPelajaran{
		IdKelas:   matapelajaran.ID(dataRequest.IdKelas),
		KodeMapel: dataRequest.KodeMapel,
		NamaMapel: dataRequest.NamaMapel,
		Deskripsi: dataRequest.Deskripsi,
	}

	if err := h.svc.CreateMapelService(r.Context(), cmdData); err != nil {
		logger.Error(r.Context(), "failed creating mapel", "layer", "adapter.http.handler", "op", "mata_pelajaran.create", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrKodeMapelExist):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: kode mapel already exist")
			return
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid input")
			return
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed create mapel")
			return
		}
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "success")
}
