package auth_service

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

type LoginCmd struct {
	Username string
	Password string
}

type LoginRes struct {
	IdPengguna   user.ID
	Username	 string
	AccessToken  string
	RefreshToken string
}