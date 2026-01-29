package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/cookie"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	authin "github.com/mustafamadjid/web-app-cbt/internal/core/port/in/auth_port_in"

	"github.com/mustafamadjid/web-app-cbt/internal/core/service/auth_service"
)

type AuthHandler struct {
	svc 	authin.AuthUsecase
	cookies	cookie.CookieConfig
	accessTTL time.Duration
	refreshTTL time.Duration
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	IdPengguna user.ID
	Username string
}


func NewAuthHandler(svc authin.AuthUsecase, cookies cookie.CookieConfig, accessTTL time.Duration, refreshTTL time.Duration) *AuthHandler {
	return &AuthHandler{
		svc: svc,
		cookies: cookies,
		accessTTL: accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (h *AuthHandler) Login(write http.ResponseWriter, req *http.Request) {
	var reqBody LoginRequest

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&reqBody); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest,"BAD_REQUEST","Bad request")
		return
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	if !re.MatchString(reqBody.Username) {
		httpResponse.WriteErr(write, http.StatusBadRequest,"BAD_REQUEST","Bad request : invalid character")
		return
	}

	reqCmd := auth_service.LoginCmd{
		Username: reqBody.Username,
		Password: reqBody.Password,
	}

	now := time.Now()

	res, err := h.svc.Login(req.Context(), reqCmd)
	if err != nil {
		switch{
		case errors.Is(err,coreerror.ErrInvalidCreds):
			httpResponse.WriteErr( write, http.StatusUnauthorized,"INVALID_CREDENTIALS","unauthorized : invalid credentials")
		default:
			httpResponse.WriteErr( write, http.StatusInternalServerError,"INTERNAL_SERVER_ERROR","internal server error")
		}
		return
	}

	responseData := LoginResponse{
		IdPengguna: res.IdPengguna,
		Username: res.Username,
	}

	cookie.SetAccessCookie(write, h.cookies, res.AccessToken, now.Add(h.accessTTL))
	cookie.SetRefreshCookie(write, h.cookies, res.RefreshToken, now.Add(h.refreshTTL))

	httpResponse.WriteOK(write, http.StatusOK, responseData, "success")
}

func (h *AuthHandler) Logout(write http.ResponseWriter, req *http.Request) {
	c,err := req.Cookie(h.cookies.RefreshName)
	if err != nil || c.Value == "" {
		httpResponse.WriteErr( write, http.StatusUnauthorized,"INVALID_TOKEN","unauthorized : invalid token")
		return
	}

	_= h.svc.Logout(req.Context(),c.Value,time.Now())

	cookie.ClearAuthCookies(write,h.cookies)

	httpResponse.WriteOKNoData(write, http.StatusOK, "success")
}

func (h *AuthHandler)Refresh(write http.ResponseWriter, req *http.Request) {
	c,err := req.Cookie(h.cookies.RefreshName)
	if err != nil || c.Value == "" {
		httpResponse.WriteErr( write, http.StatusUnauthorized,"INVALID_TOKEN","unauthorized : invalid token")
		return
	}

	newAccessToken,err := h.svc.RefreshAccessToken(req.Context(),c.Value,h.accessTTL)
	if err != nil  {
		httpResponse.WriteErr( write, http.StatusUnauthorized,"INVALID_TOKEN","unauthorized : invalid token")
		return
	}

	now := time.Now()
	cookie.SetAccessCookie(write, h.cookies,newAccessToken,now)
	httpResponse.WriteOKNoData(write, http.StatusOK,"success")

}




















// func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
// 	c, err := r.Cookie(h.cookies.RefreshName)
// 	if err != nil || c.Value == "" {
// 		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
// 		return
// 	}

// 	now := time.Now()
// 	sid, err := h.refresh.VerifyRefreshToken(c.Value, now)
// 	if err != nil {
// 		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
// 		return
// 	}

// 	sess, err := h.sessions.GetSession(r.Context(), sid)
// 	if err != nil {
// 		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
// 		return
// 	}

// 	if sess.Revoked || now.After(sess.ExpiresAt) {
// 		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
// 		return
// 	}

// 	at, err := h.access.GenerateAccessToken(sess.UserID, sess.Role, sess.Username, h.accessTTL)
// 	if err != nil {
// 		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "internal error"})
// 		return
// 	}

// 	SetAccessCookie(w, h.cookies, at, now.Add(h.accessTTL))
// 	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
// }