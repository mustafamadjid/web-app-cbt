package httpx

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	siswaujianlist_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/siswa_ujian/list_soal_ujian_selesai"
)

type ListUjianSelesaiSiswaHandler struct {
	svc *siswaujianlist_service.ListUjianSelesaiSiswaService
}

func NewListUjianSelesaiSiswaHandler(svc *siswaujianlist_service.ListUjianSelesaiSiswaService) *ListUjianSelesaiSiswaHandler {
	return &ListUjianSelesaiSiswaHandler{svc: svc}
}

func (h *ListUjianSelesaiSiswaHandler) ListUjianSelesaiSiswa(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	actor, ok := middleware.ActorFromContext(r.Context())
	if !ok {
		logger.Error(r.Context(), "missing actor in context", "layer", "adapter.http.handler", "op", "siswa_ujian.list_ujian_selesai", "err", "actor_not_found")
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get actor from context")
		return
	}

	items, err := h.svc.ListUjianSelesaiSiswa(r.Context(), int(actor.IdPengguna))
	if err != nil {
		logger.Error(r.Context(), "failed listing ujian selesai siswa", "layer", "adapter.http.handler", "op", "siswa_ujian.list_ujian_selesai", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get list ujian selesai siswa")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toListUjianSelesaiSiswaResponses(items), "Success")
}
