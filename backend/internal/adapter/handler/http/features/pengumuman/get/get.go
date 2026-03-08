package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	pengumuman_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/get"
)

type GetPengumumanHandler struct {
	svc *pengumuman_service.GetPengumumanService
}

func NewGetPengumumanHandler(svc *pengumuman_service.GetPengumumanService) *GetPengumumanHandler {
	return &GetPengumumanHandler{svc: svc}
}

func (h *GetPengumumanHandler) GetPengumumanActive(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	items, err := h.svc.GetPengumumanActiveService(r.Context())
	if err != nil {
		logger.Error(r.Context(), "failed getting active pengumuman", "layer", "adapter.http.handler", "op", "pengumuman.get_active", "err", err)
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get pengumuman")
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toPengumumanResponses(items), "Success")
}

func (h *GetPengumumanHandler) GetPengumumanNonActive(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	items, err := h.svc.GetPengumumanNonActiveService(r.Context())
	if err != nil {
		logger.Error(r.Context(), "failed getting non active pengumuman", "layer", "adapter.http.handler", "op", "pengumuman.get_non_active", "err", err)
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get pengumuman")
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toPengumumanResponses(items), "Success")
}

func (h *GetPengumumanHandler) GetPengumumanIncoming(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	items, err := h.svc.GetPengumumanIncomingService(r.Context())
	if err != nil {
		logger.Error(r.Context(), "failed getting incoming pengumuman", "layer", "adapter.http.handler", "op", "pengumuman.get_incoming", "err", err)
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get pengumuman")
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toPengumumanResponses(items), "Success")
}

func (h *GetPengumumanHandler) GetPengumumanByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())

	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	idPengumuman, err := strconv.Atoi(strings.TrimSpace(ps.ByName("idPengumuman")))
	if err != nil || idPengumuman <= 0 {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id pengumuman")
		return
	}

	item, err := h.svc.GetPengumumanByIdService(r.Context(), pengumuman.ID(idPengumuman))
	if err != nil {
		logger.Error(r.Context(), "failed getting pengumuman by id", "layer", "adapter.http.handler", "op", "pengumuman.get_by_id", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: missing id")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get pengumuman")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toPengumumanResponse(item), "Success")
}
