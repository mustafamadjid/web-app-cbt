package httpx

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

type LoginResponse struct {
	IdPengguna user.ID `json:"id_pengguna"`
	Username   string `json:"username"`
}

type AuthMeResponse struct {
	IdPengguna user.ID 	`json:"id_pengguna"`
	Username   string   `json:"username"`
	Role       user.Role `json:"role"`
}