package httpx

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	KodeUnik string `json:"kode_unik"`
}

type AdminRevokeRequest struct {
	SessionId string `json:"session_id"`
}