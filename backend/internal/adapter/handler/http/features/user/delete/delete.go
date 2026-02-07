package httpx

import (
	"encoding/json"
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

func (h *DeleteHandler) DeleteUsers(write http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(req.Context())
	if req.Method != http.MethodDelete {
		httpResponse.WriteErr(write, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	var reqBody DeleteUsersRequest
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&reqBody); err != nil {
		logger.Error(req.Context(), "failed decoding delete users request", "op", "user.delete_many", "err", err)
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request")
		return
	}

	if len(reqBody.Ids) == 0 {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : ids is required")
		return
	}

	actor, ok := middleware.ActorFromContext(req.Context())
	if !ok {
		logger.Error(req.Context(), "missing actor in context", "op", "user.delete_many", "err", "actor_not_found")
		httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	ids := make([]user.ID, 0, len(reqBody.Ids))
	seen := make(map[user.ID]struct{}, len(reqBody.Ids))
	for _, rawID := range reqBody.Ids {
		if rawID <= 0 {
			httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : invalid id")
			return
		}
		id := user.ID(rawID)
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : ids is required")
		return
	}

	affected, err := h.svc.DeleteMany(req.Context(), ids)
	if err != nil {
		logger.Error(req.Context(), "failed deleting users", "op", "user.delete_many", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(write, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed delete users")
		}
		return
	}

	if h.aktivitasUser != nil {
		aktivitasCmd := aktivitas_user_service.AktivitasUserCmd{
			IdPengguna:  actor.IdPengguna,
			Action:      aktivitas_user.DELETE,
			Description: "Menghapus beberapa pengguna",
			IpAddress:   httphelper.GetClientIP(req),
		}
		if err := h.aktivitasUser.CreateAktivitasUserService(req.Context(), aktivitasCmd); err != nil {
			logger.Error(req.Context(), "failed creating aktivitas user", "op", "user.delete_many.activity", "err", err)
		}
	}

	httpResponse.WriteOK(write, http.StatusOK, map[string]int64{"deleted": affected}, "Success")
}
