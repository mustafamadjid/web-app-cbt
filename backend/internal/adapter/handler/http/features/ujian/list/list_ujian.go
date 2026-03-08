package httpx

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/get"
)

type ListUjianHandler struct {
	svc *ujian_service.GetUjianService
}

func NewListUjianHandler(svc *ujian_service.GetUjianService) *ListUjianHandler {
	return &ListUjianHandler{svc: svc}
}

func (h *ListUjianHandler) ListUjian(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	req, err := parseListUjianRequest(r)
	if err != nil {
		logger.Info(r.Context(), "invalid ujian filters", "layer", "adapter.http.handler", "op", "ujian.list", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		return
	}

	filter := query.ListUjianFilter{
		Search:         req.Search,
		Limit:          req.Limit,
		Offset:         req.Offset,
		TanggalUjian:   req.Tanggal,
		Tahun:          req.Tahun,
		TingkatKelasID: req.TingkatKelasID,
		TingkatKelas:   req.TingkatKelas,
		RuangUjian:     req.RuangUjianID,
	}

	items, err := h.svc.GetAllUjianService(r.Context(), filter)
	if err != nil {
		logger.Error(r.Context(), "failed listing ujian", "layer", "adapter.http.handler", "op", "ujian.list", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get jadwal ujian")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toListUjianResponses(items), "Success")
}
