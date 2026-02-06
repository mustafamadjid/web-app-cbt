package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/cookie"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/service/auth_service"

	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	authin "github.com/mustafamadjid/web-app-cbt/internal/core/port/in/auth_port_in"
	out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	aktivitas_user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/aktivitas_user"
)

type AuthHandler struct {
	svc           authin.AuthUsecase
	cookies       cookie.CookieConfig
	accessTTL     time.Duration
	refreshTTL    time.Duration
	aktivitasUser *aktivitas_user_service.AktivitasUserService
	accessTokens  out.AccessTokenService
}

func NewAuthHandler(svc authin.AuthUsecase, cookies cookie.CookieConfig, accessTTL time.Duration, refreshTTL time.Duration, aktivitasUser *aktivitas_user_service.AktivitasUserService, accessTokens out.AccessTokenService) *AuthHandler {
	return &AuthHandler{
		svc:           svc,
		cookies:       cookies,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
		aktivitasUser: aktivitasUser,
		accessTokens:  accessTokens,
	}
}

func (h *AuthHandler) Login(write http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(req.Context())
	var reqBody LoginRequest

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&reqBody); err != nil {
		logger.Error(req.Context(), "failed decoding login request", "op", "auth.login", "err", err)
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request")
		return
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9._]{3,13}$`)
	if !re.MatchString(reqBody.Username) {
		logger.Info(req.Context(), "invalid login input", "op", "auth.login", "err", "invalid_username")
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : invalid character. It must be 3-16 characters long and contain only letters, numbers, and underscores.")
		return
	}

	reqCmd := auth_service.LoginCmd{
		Username: reqBody.Username,
		Password: reqBody.Password,
	}

	now := time.Now()

	res, err := h.svc.Login(req.Context(), reqCmd)
	if err != nil {
		logger.Error(req.Context(), "login failed", "op", "auth.login", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrInvalidCreds):
			httpResponse.WriteErr(write, http.StatusUnauthorized, "INVALID_CREDENTIALS", "unauthorized : invalid credentials")
		case errors.Is(err, coreerror.ErrHasSession):
			httpResponse.WriteErr(write, http.StatusUnauthorized, "HAS_SESSION", "unauthorized : already has session")
		default:
			httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}

	responseData := LoginResponse{
		IdPengguna: res.IdPengguna,
		Username:   res.Username,
	}

	cookie.SetAccessCookie(write, h.cookies, res.AccessToken, now.Add(h.accessTTL))
	cookie.SetRefreshCookie(write, h.cookies, res.RefreshToken, now.Add(h.refreshTTL))

	if h.aktivitasUser != nil {
		aktivitasCmd := aktivitas_user_service.AktivitasUserCmd{
			IdPengguna:  res.IdPengguna,
			Action:      aktivitas_user.LOGIN,
			Description: "Login pengguna",
			IpAddress:   httphelper.GetClientIP(req),
		}
		if err := h.aktivitasUser.CreateAktivitasUserService(req.Context(), aktivitasCmd); err != nil {
			logger.Error(req.Context(), "failed creating aktivitas user", "op", "auth.login.activity", "err", err)
		}
	}

	httpResponse.WriteOK(write, http.StatusOK, responseData, "success")
}

func (h *AuthHandler) Logout(write http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(req.Context())
	c, err := req.Cookie(h.cookies.RefreshName)
	if err != nil || c.Value == "" {
		logger.Info(req.Context(), "missing refresh token", "op", "auth.logout", "err", "missing_refresh_token")
		httpResponse.WriteErr(write, http.StatusUnauthorized, "INVALID_TOKEN", "unauthorized : invalid token")
		return
	}

	var userID user.ID
	if h.accessTokens != nil {
		if accessCookie, err := req.Cookie(h.cookies.AccessName); err == nil && accessCookie.Value != "" {
			uid, _, _, err := h.accessTokens.VerifyAccessToken(accessCookie.Value, time.Now())
			if err == nil {
				userID = uid
			}
		}
	}

	if err := h.svc.Logout(req.Context(), c.Value, time.Now()); err != nil {
		logger.Error(req.Context(), "logout failed", "op", "auth.logout", "err", err)
	}

	cookie.ClearAuthCookies(write, h.cookies)

	if h.aktivitasUser != nil && userID > 0 {
		aktivitasCmd := aktivitas_user_service.AktivitasUserCmd{
			IdPengguna:  userID,
			Action:      aktivitas_user.LOGOUT,
			Description: "Logout pengguna",
			IpAddress:   httphelper.GetClientIP(req),
		}
		if err := h.aktivitasUser.CreateAktivitasUserService(req.Context(), aktivitasCmd); err != nil {
			logger.Error(req.Context(), "failed creating aktivitas user", "op", "auth.logout.activity", "err", err)
		}
	}

	httpResponse.WriteOKNoData(write, http.StatusOK, "success")
}

func (h *AuthHandler) Refresh(write http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(req.Context())
	c, err := req.Cookie(h.cookies.RefreshName)
	if err != nil || c.Value == "" {
		logger.Info(req.Context(), "missing refresh token", "op", "auth.refresh", "err", "missing_refresh_token")
		httpResponse.WriteErr(write, http.StatusUnauthorized, "INVALID_TOKEN", "unauthorized : invalid token refresh token empty")
		return
	}

	newAccessToken, err := h.svc.RefreshAccessToken(req.Context(), c.Value, h.accessTTL)
	if err != nil {
		logger.Error(req.Context(), "refresh token failed", "op", "auth.refresh", "err", err)
		httpResponse.WriteErr(write, http.StatusUnauthorized, "INVALID_TOKEN", "unauthorized : invalid token refresh token failed to generate")
		return
	}

	now := time.Now()
	cookie.SetAccessCookie(write, h.cookies, newAccessToken, now.Add(h.accessTTL))
	httpResponse.WriteOKNoData(write, http.StatusOK, "success")

}

func (h *AuthHandler) AdminRevokeUser(write http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(req.Context())
	if req.Method != http.MethodPut {
		httpResponse.WriteErr(write, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	var reqBody AdminRevokeRequest

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&reqBody); err != nil {
		logger.Error(req.Context(), "failed decoding revoke request", "op", "auth.admin_revoke", "err", err)
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request")
		return
	}

	if err := h.svc.AdminRevokingSession(req.Context(), reqBody.SessionId); err != nil {
		logger.Error(req.Context(), "failed revoking session", "op", "auth.admin_revoke", "session_id", reqBody.SessionId, "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(write, http.StatusNotFound, "NOT_FOUND", "Session not found")
		default:
			httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		}
		return
	}
	httpResponse.WriteOKNoData(write, http.StatusOK, "success")
}


func (h *AuthHandler) AuthMe(write http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(req.Context())
	if req.Method != http.MethodGet {
		httpResponse.WriteErr(write, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	actor,ok := middleware.ActorFromContext(req.Context())
	if !ok {
		logger.Error(req.Context(), "failed getting actor", "op", "auth.auth_me", "err", "actor_not_found")
		httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
		return
	}
	responseData := AuthMeResponse {
		IdPengguna: actor.IdPengguna,
		Username: actor.Username,
		Role: actor.Role,
	}

	httpResponse.WriteOK(write, http.StatusOK, responseData, "success")
}