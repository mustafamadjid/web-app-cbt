package middleware

import (
	"net/http"
	"os"
	"strings"
)

func PreventCSRF(next http.Handler) http.Handler {
	cop := http.NewCrossOriginProtection()

	raw := os.Getenv("TRUSTED_ORIGINS")
	origins := strings.Split(raw, ",")

	for _, o := range origins {
		o = strings.TrimSpace(o)
		if o == "" {
			continue
		}
		if err := cop.AddTrustedOrigin(o); err != nil {
			// TODO : Tambahkan log 
		}
	}

	cop.SetDenyHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "CSRF check failed", http.StatusForbidden)
	}))

	return cop.Handler(next)
}
