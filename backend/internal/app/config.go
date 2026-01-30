package app

import (
	"net/http"
	"time"
)

type Config struct {
	HTTP   HTTPConfig
	JWT    JWTConfig
	Cookie CookieConfig
	ImageStore ImageStoreConfig
}

type JWTConfig struct {
	Issuer   string
	AccessSecret  string
	RefreshSecret string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type CookieConfig struct {
	AccessName  string
	RefreshName string
	Domain      string
	Secure      bool
	SameSite    http.SameSite
}

type HTTPConfig struct {
	Addr string
}

type ImageStoreConfig struct {
	Dir      string
	BaseURL  string
	Route    string
	MaxBytes int64
}