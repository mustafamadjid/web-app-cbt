package httpx

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

type LoginResponse struct {
	IdPengguna user.ID
	Username   string
}