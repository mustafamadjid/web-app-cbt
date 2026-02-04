package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	aktivitas_user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/aktivitas_user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/delete"
)

type DeleteHandler struct {
	svc           *user_service.DeleteUserService
	aktivitasUser *aktivitas_user_service.AktivitasUserService
}

func NewDeleteUserHandler(svc *user_service.DeleteUserService, aktivitasUser *aktivitas_user_service.AktivitasUserService) *DeleteHandler {
	return &DeleteHandler{svc: svc, aktivitasUser: aktivitasUser}
}

func (h *DeleteHandler) DeleteUser(write http.ResponseWriter, req *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(req.Context())
	if req.Method != http.MethodDelete {
		httpResponse.WriteErr(write, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	rawId := strings.TrimSpace(params.ByName("id"))
	if rawId == "" {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : id is required")
		return
	}

	parsedId, err := strconv.Atoi(rawId)
	if err != nil || parsedId <= 0 {
		logger.Info(req.Context(), "invalid user id", "op", "user.delete", "err", err)
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : invalid id")
		return
	}

	uid := user.ID(parsedId)

	actor, ok := middleware.ActorFromContext(req.Context())
	if !ok {
		logger.Error(req.Context(), "missing actor in context", "op", "user.delete", "err", "actor_not_found")
		httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	if err := h.svc.Delete(req.Context(), uid); err != nil {
		logger.Error(req.Context(), "failed deleting user", "op", "user.delete", "user_id", uid, "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(write, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed delete user")
		}
		return
	}

	if h.aktivitasUser != nil {
		aktivitasCmd := aktivitas_user_service.AktivitasUserCmd{
			IdPengguna:  actor.IdPengguna,
			Action:      aktivitas_user.DELETE,
			Description: "Menghapus pengguna",
			IpAddress:   httphelper.GetClientIP(req),
		}
		if err := h.aktivitasUser.CreateAktivitasUserService(req.Context(), aktivitasCmd); err != nil {
			logger.Error(req.Context(), "failed creating aktivitas user", "op", "user.delete.activity", "err", err)
		}
	}

	httpResponse.WriteOKNoData(write, http.StatusOK, "Success")
}
