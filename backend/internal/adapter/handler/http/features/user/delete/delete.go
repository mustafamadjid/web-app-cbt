package httpx

import (
	"encoding/json"
	"errors"
	"net/http"

	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
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

	var payload DeleteUserRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : invalid request body")
		return
	}

	if payload.IdPengguna == 0 {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : id_pengguna is required")
		return
	}

	if err := h.svc.Delete(req.Context(), payload.IdPengguna); err != nil {
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
