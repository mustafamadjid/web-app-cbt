package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/rs/cors"
)

func CORSPolicy(next http.Handler) http.Handler {
	raw := os.Getenv("TRUSTED_ORIGINS")
	parts := strings.Split(raw, ",")

	allowed := make([]string, 0, len(parts))
	for _, o := range parts {
		o = strings.TrimSpace(o)
		if o == "" {
			continue
		}
		allowed = append(allowed, o)
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   allowed,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return c.Handler(next)
}
