package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/cookie"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

	"github.com/mustafamadjid/web-app-cbt/internal/core/service/auth_service"
)

type AuthHandler struct {
	svc 	*auth_service.AuthService
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


func NewAuthHandler(svc *auth_service.AuthService, cookies cookie.CookieConfig, accessTTL time.Duration, refreshTTL time.Duration) *AuthHandler {
	return &AuthHandler{
		svc: svc,
		cookies: cookies,
		accessTTL: accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (h *AuthHandler) Login(write http.ResponseWriter, req *http.Request) {
	var reqBody LoginRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest,"BAD_REQUEST","Bad request")
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

