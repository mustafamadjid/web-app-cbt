package httpmodule

import (
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func ServeStaticFrom(dir string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		rel := ps.ByName("filepath")
		if rel == "/" || rel == "" || rel == "." {
			http.NotFound(w, r)
			return
		}

		clean := path.Clean(rel)
		if strings.Contains(clean, "..") {
			http.Error(w, "bad path", http.StatusBadRequest)
			return
		}

		clean = strings.TrimPrefix(clean, "/")
		if clean == "." || clean == "" {
			http.NotFound(w, r)
			return
		}

		full := filepath.Join(dir, clean)
		http.ServeFile(w, r, full)
	}
}
