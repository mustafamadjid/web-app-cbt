package httpx

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	listessayungraded_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/grading_ujian/list/list_ujian_essay_ungraded"
)

type ListUjianEssayUngradedHandler struct {
	svc *listessayungraded_service.ListUjianEssayUngradedService
}

func NewListUjianEssayUngradedHandler(svc *listessayungraded_service.ListUjianEssayUngradedService) *ListUjianEssayUngradedHandler {
	return &ListUjianEssayUngradedHandler{svc: svc}
}

func (h *ListUjianEssayUngradedHandler) ListUjianEssayUngraded(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	req, err := parseListUjianEssayUngradedRequest(r)
	if err != nil {
		logger.Info(r.Context(), "invalid essay ungraded ujian filters", "layer", "adapter.http.handler", "op", "ujian.grading.list_essay_ungraded", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		return
	}

	items, err := h.svc.ListUjianEssayUngraded(r.Context(), toListUjianEssayUngradedFilter(req))
	if err != nil {
		logger.Error(r.Context(), "failed listing essay ungraded ujian", "layer", "adapter.http.handler", "op", "ujian.grading.list_essay_ungraded", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get list ujian essay ungraded")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toListUjianEssayUngradedResponses(items), "Success")
}

func toListUjianEssayUngradedFilter(req ListUjianEssayUngradedRequest) query.ListUjianEssayUngradedFilter {
	return query.ListUjianEssayUngradedFilter{
		Search:         req.Search,
		Limit:          req.Limit,
		Offset:         req.Offset,
		TanggalUjian:   req.TanggalUjian,
		Tahun:          req.Tahun,
		Bulan:          req.Bulan,
		TingkatKelasID: req.TingkatKelasID,
		NamaKelasID:    req.NamaKelasID,
		MapelID:        req.MapelID,
		SesiID:         req.SesiID,
	}
}
