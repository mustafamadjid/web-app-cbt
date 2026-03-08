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

	if err := httphelper.JsonHeaderBodyValidator(w, r); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	var dataRequest UpdateMapelRequest
	if err := httphelper.JSONDecoder(r, &dataRequest); err != nil {
		switch {
		case errors.Is(err, coreerror.ErrInvalidRequestBody):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	if dataRequest.KodeMapel != nil {
		if err := validator.ValidateInputSafe(*dataRequest.KodeMapel, "kode_mapel"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "invalid input")
			return
		}
	}

	if dataRequest.NamaMapel != nil {
		if err := validator.ValidateInputSafe(*dataRequest.NamaMapel, "nama_mapel"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "invalid input")
			return
		}
	}

	if dataRequest.Deskripsi != nil {
		if err := validator.ValidateInputSafe(*dataRequest.Deskripsi, "deskripsi"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "invalid input")
			return
		}
	}

	patch := updatepatch.UpdateMapelPatch{
		IdKelas:   toMapelIDPointer(dataRequest.IdKelas),
		KodeMapel: dataRequest.KodeMapel,
		NamaMapel: dataRequest.NamaMapel,
		Deskripsi: dataRequest.Deskripsi,
	}

	if err := h.svc.UpdateMapelService(r.Context(), idMapel, patch); err != nil {
		logger.Error(r.Context(), "failed updating mapel", "layer", "adapter.http.handler", "op", "mata_pelajaran.update", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data mata pelajaran tidak ditemukan")
			return
		case errors.Is(err, coreerror.ErrMissingId), errors.Is(err, coreerror.ErrMissingField), errors.Is(err, coreerror.ErrKodeMapelExist), errors.Is(err, coreerror.ErrNoFieldToUpdate):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid field")
			return
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed update mata pelajaran")
			return
		}
	}

	httpResponse.WriteOK(w, http.StatusOK, UpdateMapelResponse{Success: true}, "success update mata pelajaran")
}
