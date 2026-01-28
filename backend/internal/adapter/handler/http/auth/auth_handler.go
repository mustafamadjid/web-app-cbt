package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/cookie"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/response_envelope"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
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

	now := time.Now()
	
}

