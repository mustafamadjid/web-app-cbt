package httpx

import (
	"net/http"
	"strings"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
)

func JsonHeaderBodyValidator(w http.ResponseWriter, r *http.Request) error {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		return coreerror.ErrContentTypeMustJson
	}

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	return nil
}
