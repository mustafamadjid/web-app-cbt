package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	matapelajaran "github.com/mustafamadjid/web-app-cbt/internal/core/domain/mata_pelajaran"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	mapel_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/mata_pelajaran/update"
)

type UpdateMapelHandler struct {
	svc *mapel_service.UpdateMapelRepo
}

func NewUpdateMapelHandler(svc *mapel_service.UpdateMapelRepo) *UpdateMapelHandler {
	return &UpdateMapelHandler{svc: svc}
}

func (h *UpdateMapelHandler) UpdateMapel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodPatch {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idMapel, err := strconv.Atoi(ps.ByName("idMapel"))
	if err != nil || idMapel <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id mapel")
		return
	}

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	var dataRequest UpdateMapelRequest
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

	if dataRequest.IdKelas == nil && dataRequest.KodeMapel == nil && dataRequest.NamaMapel == nil && dataRequest.Deskripsi == nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: no fields to update")
		return
	}

	patch := updatepatch.UpdateMapelPatch{}
	if dataRequest.IdKelas != nil {
		if *dataRequest.IdKelas <= 0 {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: id kelas is required")
			return
		}
		idKelas := matapelajaran.ID(*dataRequest.IdKelas)
		patch.IdKelas = &idKelas
	}

	if dataRequest.KodeMapel != nil {
		kodeMapel := strings.TrimSpace(*dataRequest.KodeMapel)
		if kodeMapel == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: kode mapel is required")
			return
		}
		if err := validator.ValidateInputSafe(kodeMapel, "kode_mapel"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "invalid input")
			return
		}
		patch.KodeMapel = &kodeMapel
	}

	if dataRequest.NamaMapel != nil {
		namaMapel := strings.TrimSpace(*dataRequest.NamaMapel)
		if namaMapel == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: nama mapel is required")
			return
		}
		if err := validator.ValidateInputSafe(namaMapel, "nama_mapel"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "invalid input")
			return
		}
		patch.NamaMapel = &namaMapel
	}

	if dataRequest.Deskripsi != nil {
		deskripsi := strings.TrimSpace(*dataRequest.Deskripsi)
		if deskripsi == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: deskripsi is required")
			return
		}
		if err := validator.ValidateInputSafe(deskripsi, "deskripsi"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "invalid input")
			return
		}
		patch.Deskripsi = &deskripsi
	}

	if err := h.svc.UpdateMapelService(r.Context(), idMapel, patch); err != nil {
		logger.Error(r.Context(), "failed updating mapel", "layer", "adapter.http.handler", "op", "mata_pelajaran.update", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data mata pelajaran tidak ditemukan")
			return
		case errors.Is(err, coreerror.ErrMissingId), errors.Is(err, coreerror.ErrMissingField), errors.Is(err, coreerror.ErrKodeMapelExist):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid field")
			return
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed update mata pelajaran")
			return
		}
	}

	httpResponse.WriteOK(w, http.StatusOK, UpdateMapelResponse{Success: true}, "success update mata pelajaran")
}
