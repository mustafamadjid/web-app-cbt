package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	ruangujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ruang_ujian/update"
)

type UpdateRuangUjianHandler struct {
	svc *ruangujian_service.UpdateRuangUjianService
}

func NewUpdateRuangUjianHandler(svc *ruangujian_service.UpdateRuangUjianService) *UpdateRuangUjianHandler {
	return &UpdateRuangUjianHandler{svc: svc}
}

func (h *UpdateRuangUjianHandler) UpdateRuangUjian(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodPatch {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idRuangan, err := strconv.Atoi(ps.ByName("idRuangan"))
	if err != nil || idRuangan <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id ruangan")
		return
	}

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	var dataRequest UpdateRuangUjianRequest

	if err := httphelper.JSONDecoder(r, &dataRequest); err != nil {
		switch {
		case errors.Is(err, coreerror.ErrInvalidRequestBody):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	if dataRequest.KodeRuang == nil && dataRequest.NamaRuangan == nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		return
	}

	patch := updatepatch.UpdateRuangUjianPatch{}

	if dataRequest.NamaRuangan != nil {
		namaRuangan := strings.TrimSpace(*dataRequest.NamaRuangan)
		if namaRuangan == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: missing fields")
			return
		}

		if err := validator.ValidateInputSafe(namaRuangan, "nama_ruang"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "invalid input")
			return
		}

		patch.NamaRuang = &namaRuangan
	}

	if dataRequest.KodeRuang != nil {
		kodeRuang := strings.TrimSpace(*dataRequest.KodeRuang)
		if kodeRuang == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: missing fields")
			return
		}

		if err := validator.ValidateInputSafe(kodeRuang, "kode_ruang"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "invalid input")
			return
		}

		patch.KodeRuang = &kodeRuang
	}

	if err := h.svc.UpdateRuangUjian(r.Context(), idRuangan, patch); err != nil {
		logger.Error(r.Context(), "failed update ruang ujian", "layer", "core.service", "op", "ruangujian.update", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrKodeRuangUjianExist):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: kode ruang ujian already exist")
			return
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: missing id")
			return
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
			return
		}
	}
	httpResponse.WriteOKNoData(w, http.StatusOK, "success")
}
