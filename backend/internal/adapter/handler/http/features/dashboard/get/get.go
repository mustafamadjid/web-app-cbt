package httpx

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	dashboard_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/dashboard"
)

type GetDashboardStatistikHandler struct {
	svc *dashboard_service.DashboardService
}

func NewGetDashboardStatistikHandler(svc *dashboard_service.DashboardService) *GetDashboardStatistikHandler {
	return &GetDashboardStatistikHandler{svc: svc}
}

func (h *GetDashboardStatistikHandler) GetDashboardStatistik(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	statistik, err := h.svc.GetDashboardStatistik(r.Context())
	if err != nil {
		logger.Error(r.Context(), "failed get dashboard statistik", "layer", "adapter.http.handler", "op", "dashboard.statistik.get", "err", err)
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get dashboard statistik")
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toGetDashboardStatistikResponse(statistik), "Success")
}
