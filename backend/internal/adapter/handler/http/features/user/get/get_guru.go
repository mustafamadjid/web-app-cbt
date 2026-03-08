package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/get"
)

type GetGuruHandler struct {
	svc *user_service.GetGuruService
}

func NewGetGuruHandler(svc *user_service.GetGuruService) *GetGuruHandler {
	return &GetGuruHandler{svc: svc}
}

func (h *GetGuruHandler) ListGuru(write http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(req.Context())
	if req.Method != http.MethodGet {
		httpResponse.WriteErr(write, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	filters, err := parseListGuruFilters(req)
	if err != nil {
		logger.Info(req.Context(), "invalid guru filters", "layer", "adapter.http.handler", "op", "user.list_guru", "err", err)
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	items, err := h.svc.ListGuru(req.Context(), filters)
	if err != nil {
		logger.Error(req.Context(), "failed listing guru", "layer", "adapter.http.handler", "op", "user.list_guru", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		default:
			httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get list guru")
		}
		return
	}

	httpResponse.WriteOK(write, http.StatusOK, toGuruResponses(items), "Success")
}

func (h *GetGuruHandler) GetGuruByID(write http.ResponseWriter, req *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(req.Context())
	if req.Method != http.MethodGet {
		httpResponse.WriteErr(write, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	rawID := strings.TrimSpace(params.ByName("id"))
	if rawID == "" {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : id is required")
		return
	}

	parsedID, err := strconv.Atoi(rawID)
	if err != nil || parsedID <= 0 {
		logger.Info(req.Context(), "invalid guru id", "layer", "adapter.http.handler", "op", "user.get_guru", "err", err)
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : invalid id")
		return
	}

	result, err := h.svc.FindProfilGuruByID(req.Context(), user.ID(parsedID))
	if err != nil {
		logger.Error(req.Context(), "failed getting guru", "layer", "adapter.http.handler", "op", "user.get_guru", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(write, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get guru")
		}
		return
	}

	httpResponse.WriteOK(write, http.StatusOK, toGuruResponseFromData(result), "Success")
}
