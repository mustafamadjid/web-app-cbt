package httpx

import (
	"errors"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	profil_sekolah_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/profil_sekolah/get"
)

type GetProfilSekolahHandler struct {
	svc *profil_sekolah_service.GetProfilSekolahService
}

func NewGetProfilSekolahHandler(svc *profil_sekolah_service.GetProfilSekolahService) *GetProfilSekolahHandler {
	return &GetProfilSekolahHandler{svc: svc}
}

type profilSekolahResponse struct {
	IDProfil      int     `json:"id_profil"`
	EmailSekolah  string  `json:"email_sekolah"`
	NoTelpSekolah string  `json:"no_telp_sekolah"`
	KepalaSekolah string  `json:"kepala_sekolah"`
	WakaSekolah   string  `json:"waka_sekolah"`
	NamaSekolah   string  `json:"nama_sekolah"`
	AlamatSekolah string  `json:"alamat_sekolah"`
	LogoSekolah   *string `json:"logo_sekolah"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

func (h *GetProfilSekolahHandler) GetProfilSekolah(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	profil, err := h.svc.GetProfilSekolah(r.Context())
	if err != nil {
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "profil sekolah tidak ditemukan")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get profil sekolah")
		}
		return
	}

	response := profilSekolahResponse{
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
