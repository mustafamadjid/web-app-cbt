package httpx

import (
	"errors"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	httpx "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/features/profil_sekolah"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	profil_sekolah_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/profil_sekolah/get"
)

type GetProfilSekolahHandler struct {
	svc *profil_sekolah_service.GetProfilSekolahService
}

func NewGetProfilSekolahHandler(svc *profil_sekolah_service.GetProfilSekolahService) *GetProfilSekolahHandler {
	return &GetProfilSekolahHandler{svc: svc}
}



func (h *GetProfilSekolahHandler) GetProfilSekolah(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	profil, err := h.svc.GetProfilSekolah(r.Context())
	if err != nil {
		logger.Error(r.Context(), "failed getting profil sekolah", "layer", "adapter.http.handler", "op", "profil_sekolah.get", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "profil sekolah tidak ditemukan")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get profil sekolah")
		}
		return
	}

	response := httpx.ProfilSekolahResponse{
		IDProfil:      int(profil.IDProfil),
		EmailSekolah:  profil.EmailSekolah,
		NoTelpSekolah: profil.NoTelpSekolah,
		KepalaSekolah: profil.KepalaSekolah,
		WakaSekolah:   profil.WakaSekolah,
		NamaSekolah:   profil.NamaSekolah,
		AlamatSekolah: profil.AlamatSekolah,
		LogoSekolah:   profil.LogoSekolah,
		CreatedAt:     profil.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     profil.UpdatedAt.Format(time.RFC3339),
	}

	httpResponse.WriteOK(w, http.StatusOK, response, "Success")
}
