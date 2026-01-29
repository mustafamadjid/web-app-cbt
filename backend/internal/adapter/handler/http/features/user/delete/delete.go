package httpx

import (
	"errors"
	"net/http"
	"strconv"

	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/delete"
)

type DeleteHandler struct {
	svc *user_service.DeleteUserService
}

func NewDeleteUserHandler(svc *user_service.DeleteUserService) *DeleteHandler {
	return &DeleteHandler{svc: svc}
}

func (h *DeleteHandler) DeleteGuru(write http.ResponseWriter, req *http.Request) {
	h.deleteUser(write, req)
}

func (h *DeleteHandler) DeleteSiswa(write http.ResponseWriter, req *http.Request) {
	h.deleteUser(write, req)
}

func (h *DeleteHandler) deleteUser(write http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		httpResponse.WriteErr(write, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	id := req.URL.Query().Get("id")
	if id == "" {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : id is required")
		return
	}

	parsedId, err := strconv.Atoi(id)
	if err != nil || parsedId <= 0 {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : invalid id")
		return
	}

	uid := user.ID(parsedId)
	

	if err := h.svc.Delete(req.Context(), uid); err != nil {
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(write, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed delete user")
		}
		return
	}

	httpResponse.WriteOKNoData(write, http.StatusOK, "Success")
}
