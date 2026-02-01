package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/cookie"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	authin "github.com/mustafamadjid/web-app-cbt/internal/core/port/in/auth_port_in"

	"github.com/mustafamadjid/web-app-cbt/internal/core/service/auth_service"
)

type AuthHandler struct {
	svc        authin.AuthUsecase
	cookies    cookie.CookieConfig
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewAuthHandler(svc authin.AuthUsecase, cookies cookie.CookieConfig, accessTTL time.Duration, refreshTTL time.Duration) *AuthHandler {
	return &AuthHandler{
		svc:        svc,
		cookies:    cookies,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (h *AuthHandler) Login(write http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var reqBody LoginRequest

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&reqBody); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request")
		return
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9._]{3,13}$`)
	if !re.MatchString(reqBody.Username) {
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

	httpResponse.WriteOK(write, http.StatusOK, responseData, "success")
}

func (h *AuthHandler) Logout(write http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	c, err := req.Cookie(h.cookies.RefreshName)
	if err != nil || c.Value == "" {
		httpResponse.WriteErr(write, http.StatusUnauthorized, "INVALID_TOKEN", "unauthorized : invalid token")
		return
	}

	_ = h.svc.Logout(req.Context(), c.Value, time.Now())

	cookie.ClearAuthCookies(write, h.cookies)

	httpResponse.WriteOKNoData(write, http.StatusOK, "success")
}

func (h *AuthHandler) Refresh(write http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	c, err := req.Cookie(h.cookies.RefreshName)
	if err != nil || c.Value == "" {
		httpResponse.WriteErr(write, http.StatusUnauthorized, "INVALID_TOKEN", "unauthorized : invalid token")
		return
	}

	newAccessToken, err := h.svc.RefreshAccessToken(req.Context(), c.Value, h.accessTTL)
	if err != nil {
		httpResponse.WriteErr(write, http.StatusUnauthorized, "INVALID_TOKEN", "unauthorized : invalid token")
		return
	}

	now := time.Now()
	cookie.SetAccessCookie(write, h.cookies, newAccessToken, now)
	httpResponse.WriteOKNoData(write, http.StatusOK, "success")

}


func (h *AuthHandler)AdminRevokeUser(write http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Method != http.MethodPut {
		httpResponse.WriteErr(write, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	var reqBody AdminRevokeRequest

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&reqBody); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request")
		return
	}

	if err := h.svc.AdminRevokingSession(req.Context(),reqBody.SessionId); err != nil {
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