package cookie


import (
	"net/http"
	"time"
)

type CookieConfig struct {
	AccessName  string
	RefreshName string
	Domain      string
	Secure      bool
	SameSite    http.SameSite
}

func SetAccessCookie(write http.ResponseWriter, cfg CookieConfig, token string, exp time.Time) {
	http.SetCookie(write, &http.Cookie{
		Name:     cfg.AccessName,
		Value:    token,
		Path:     "/",
		Expires:  exp,
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
	})
}

func SetRefreshCookie(w http.ResponseWriter, cfg CookieConfig, token string, exp time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     cfg.RefreshName,
		Value:    token,
		Path:     "/auth/refresh",
		Expires:  exp,
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
	})
}

func ClearAuthCookies(w http.ResponseWriter, cfg CookieConfig) {
	http.SetCookie(w, &http.Cookie{
		Name:     cfg.AccessName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     cfg.RefreshName,
		Value:    "",
		Path:     "/auth/refresh",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
	})
}