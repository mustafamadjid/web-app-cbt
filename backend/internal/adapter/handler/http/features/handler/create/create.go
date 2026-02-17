package httpx

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	sesi_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/sesi/create"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"

)

type CreateSesiHandler struct {
	svc *sesi_service.CreateSesiService
}

func NewCreateSesiHandler(svc *sesi_service.CreateSesiService) *CreateSesiHandler {
	return &CreateSesiHandler{
		svc: svc,
	}
}

func(s *CreateSesiHandler)CreateSesiHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodPost {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request : content type must be application/json")
		return
	}	

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	var dataRequest CreateSesiRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&dataRequest); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		return
	}
}