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
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	sesi_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/update"
)

type UpdateSesiHandler struct {
	svc *sesi_service.UpdateSesiService
}

func NewUpdateSesiHandler(svc *sesi_service.UpdateSesiService) *UpdateSesiHandler {
	return &UpdateSesiHandler{svc: svc}
}

func (h *UpdateSesiHandler) UpdateSesi(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodPatch {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idSesi, err := strconv.Atoi(ps.ByName("idSesi"))
	if err != nil || idSesi <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id sesi")
		return
	}

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	var dataRequest UpdateSesiRequest
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

	if dataRequest.KodeSesi == nil && dataRequest.NamaSesi == nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: no fields to update")
		return
	}

	patch := updatepatch.UpdateSesiPatch{}
	if dataRequest.KodeSesi != nil {
		kodeSesi := strings.TrimSpace(*dataRequest.KodeSesi)
		if kodeSesi == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: kode sesi is required")
			return
		}
		if err := validator.ValidateInputSafe(kodeSesi, "kode_sesi"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
			return
		}
		patch.KodeSesi = &kodeSesi
	}

	if dataRequest.NamaSesi != nil {
		namaSesi := strings.TrimSpace(*dataRequest.NamaSesi)
		if namaSesi == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: nama sesi is required")
			return
		}
		if err := validator.ValidateInputSafe(namaSesi, "nama_sesi"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
			return
		}
		patch.NamaSesi = &namaSesi
	}

	if err := h.svc.UpdateSesiService(r.Context(), idSesi, patch); err != nil {
		logger.Error(r.Context(), "failed updating sesi", "layer", "adapter.http.handler", "op", "sesi.update", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		case errors.Is(err, coreerror.ErrMissingId), errors.Is(err, coreerror.ErrMissingField):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid field")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed update sesi")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "success")
}
