package httpx

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminRevokeRequest struct {
	SessionId string `json:"session_id"`
}