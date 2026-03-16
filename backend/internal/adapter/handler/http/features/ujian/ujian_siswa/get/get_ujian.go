package httpx

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/ujian_penjadwalan/get"
)

type GetUjianHandler struct {
	svc *ujian_service.GetUjianService
}

func NewGetUjianHandler(svc *ujian_service.GetUjianService) *GetUjianHandler {
	return &GetUjianHandler{svc: svc}
}

func (h *GetUjianHandler) GetUjianById(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idUjianStr := params.ByName("idUjian")
	idUjian, err := strconv.ParseInt(idUjianStr, 10, 64)
	if err != nil || idUjian <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "id ujian tidak valid")
		return
	}

	item, err := h.svc.GetUjianByIdService(r.Context(), ujian.ID(idUjian))
	if err != nil {
		logger.Error(r.Context(), "failed get ujian by id", "layer", "adapter.http.handler", "op", "ujian.get_by_id", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "id ujian tidak valid")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "ujian not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get ujian by id")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, ToUjianResponse(item), "Success")
}
