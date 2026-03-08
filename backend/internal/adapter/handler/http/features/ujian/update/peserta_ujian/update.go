package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/update"
)

type UpdatePesertaUjianHandler struct {
	svc *ujian_service.UpdatePesertaUjianService
}

func NewUpdatePesertaUjianHandler(svc *ujian_service.UpdatePesertaUjianService) *UpdatePesertaUjianHandler {
	return &UpdatePesertaUjianHandler{svc: svc}
}

func (h *UpdatePesertaUjianHandler) UpdatePesertaUjian(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodPatch {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	rawIDPesertaUjian := strings.TrimSpace(ps.ByName("idPesertaUjian"))
	idPesertaUjian, err := strconv.Atoi(rawIDPesertaUjian)
	if err != nil || idPesertaUjian <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id peserta ujian")
		return
	}

	if err := httphelper.JsonHeaderBodyValidator(w, r); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}

	var dataRequest UpdatePesertaUjianRequest
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

	if err := ValidateInputIDRequestUpdatePesertaUjian(dataRequest); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id")
		return
	}

	patch := updatepatch.UpdatePesertaUjianPatch{
		IdJadwalUjian: toIDPesertaUjianPointer(dataRequest.IdJadwalUjian),
		IdSiswa:       toIDPesertaUjianPointer(dataRequest.IdSiswa),
		WaktuMulai:    dataRequest.WaktuMulai,
		WaktuSubmit:   dataRequest.WaktuSubmit,
		NilaiUjian:    dataRequest.NilaiUjian,
	}

	if err := h.svc.UpdatePesertaUjianService(r.Context(), ujian.ID(idPesertaUjian), patch); err != nil {
		logger.Error(r.Context(), "failed updating peserta ujian", "layer", "adapter.http.handler", "op", "ujian.update_peserta", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		case errors.Is(err, coreerror.ErrMissingId),
			errors.Is(err, coreerror.ErrMissingField),
			errors.Is(err, coreerror.ErrNoFieldToUpdate),
			errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid field")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed update peserta ujian")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "success")
}
