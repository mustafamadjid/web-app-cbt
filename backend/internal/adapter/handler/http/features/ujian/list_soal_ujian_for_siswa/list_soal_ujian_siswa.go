package httpx

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/soal_ujian/siswa"
)

type ListSoalUjianSiswaHandler struct {
	svc *ujian_service.ListSoalUjianSiswaService
}

func NewListSoalUjianSiswaHandler(svc *ujian_service.ListSoalUjianSiswaService) *ListSoalUjianSiswaHandler {
	return &ListSoalUjianSiswaHandler{svc: svc}
}

func (h *ListSoalUjianSiswaHandler) ListSoalUjianSiswa(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	req, err := parseListSoalUjianSiswaRequest(params)
	if err != nil {
		logger.Info(r.Context(), "invalid list soal ujian siswa request", "layer", "adapter.http.handler", "op", "siswa_soal_ujian.list", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	items, err := h.svc.ListSoalUjianSiswa(r.Context(), ujian.ID(req.IDJadwalUjian))
	if err != nil {
		logger.Error(r.Context(), "failed list soal ujian siswa", "layer", "adapter.http.handler", "op", "siswa_soal_ujian.list", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "MISSING_ID", "missing id")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get soal ujian siswa")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toListSoalUjianSiswaResponses(items), "Success")
}
