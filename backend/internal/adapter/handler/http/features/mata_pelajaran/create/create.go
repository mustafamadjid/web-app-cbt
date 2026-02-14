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

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	var dataRequest CreateMapelRequest
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

	if dataRequest.IdKelas <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: id kelas is required")
		return
	}

	if strings.TrimSpace(dataRequest.KodeMapel) == "" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: kode mapel is required")
		return
	}
	if err := validator.ValidateInputSafe(dataRequest.KodeMapel, "kode_mapel"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	if strings.TrimSpace(dataRequest.NamaMapel) == "" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: nama mapel is required")
		return
	}
	if err := validator.ValidateInputSafe(dataRequest.NamaMapel, "nama_mapel"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	if strings.TrimSpace(dataRequest.Deskripsi) == "" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: deskripsi is required")
		return
	}
	if err := validator.ValidateInputSafe(dataRequest.Deskripsi, "deskripsi"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

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

	httpResponse.WriteOKNoData(w,http.StatusOK,"success")
}
