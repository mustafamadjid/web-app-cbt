package httpx

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	siswaujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/waktu_selesai"
)

type GetWaktuSelesaiUjianHandler struct {
	svc *siswaujian_service.GetWaktuSelesaiService
}

func NewGetWaktuSelesaiUjianHandler(svc *siswaujian_service.GetWaktuSelesaiService) *GetWaktuSelesaiUjianHandler {
	return &GetWaktuSelesaiUjianHandler{svc: svc}
}

func (h *GetWaktuSelesaiUjianHandler) GetWaktuSelesaiUjian(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	req, err := parseGetWaktuSelesaiUjianRequest(params)
	if err != nil {
		logger.Info(r.Context(), "invalid get waktu selesai ujian request", "layer", "adapter.http.handler", "op", "siswa_ujian.get_waktu_selesai", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", err.Error())
		return
	}

	waktuSelesai, err := h.svc.GetWaktuSelesai(r.Context(), req.IDJadwalUjian)
	if err != nil {
		logger.Error(r.Context(), "failed get waktu selesai ujian", "layer", "adapter.http.handler", "op", "siswa_ujian.get_waktu_selesai", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "id jadwal ujian tidak valid")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "jadwal ujian not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get waktu selesai ujian")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toGetWaktuSelesaiUjianResponse(req.IDJadwalUjian, waktuSelesai), "Success")
}
