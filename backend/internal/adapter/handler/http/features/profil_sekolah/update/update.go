package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/profil_sekolah"
	profil_sekolah_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/profil_sekolah/update"
)

type UpdateProfilSekolahHandler struct {
	svc *profil_sekolah_service.UpdateProfilSekolahService
}

func NewUpdateProfilSekolahHandler(svc *profil_sekolah_service.UpdateProfilSekolahService) *UpdateProfilSekolahHandler {
	return &UpdateProfilSekolahHandler{svc: svc}
}

type updateProfilSekolahRequest struct {
	EmailSekolah  string  `json:"email_sekolah"`
	NoTelpSekolah string  `json:"no_telp_sekolah"`
	KepalaSekolah string  `json:"kepala_sekolah"`
	WakaSekolah   string  `json:"waka_sekolah"`
	NamaSekolah   string  `json:"nama_sekolah"`
	AlamatSekolah string  `json:"alamat_sekolah"`
	LogoSekolah   *string `json:"logo_sekolah"`
}

func (h *UpdateProfilSekolahHandler) UpdateProfilSekolah(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != http.MethodPatch {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	if ct := r.Header.Get("Content-Type"); ct != "" && !strings.HasPrefix(ct, "application/json") {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: content type must be application/json")
		return
	}

	var req updateProfilSekolahRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid JSON body")
		return
	}

	req.EmailSekolah = strings.TrimSpace(req.EmailSekolah)
	req.NoTelpSekolah = strings.TrimSpace(req.NoTelpSekolah)
	req.KepalaSekolah = strings.TrimSpace(req.KepalaSekolah)
	req.WakaSekolah = strings.TrimSpace(req.WakaSekolah)
	req.NamaSekolah = strings.TrimSpace(req.NamaSekolah)
	req.AlamatSekolah = strings.TrimSpace(req.AlamatSekolah)

	if req.EmailSekolah == "" || req.NoTelpSekolah == "" || req.KepalaSekolah == "" || req.WakaSekolah == "" || req.NamaSekolah == "" || req.AlamatSekolah == "" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		return
	}

	if err := validator.ValidateInputSafe(req.EmailSekolah, "email_sekolah"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(req.NoTelpSekolah, "no_telp_sekolah"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(req.KepalaSekolah, "kepala_sekolah"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(req.WakaSekolah, "waka_sekolah"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(req.NamaSekolah, "nama_sekolah"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(req.AlamatSekolah, "alamat_sekolah"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	var logo *string
	if req.LogoSekolah != nil {
		trimmed := strings.TrimSpace(*req.LogoSekolah)
		if trimmed == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input: logo_sekolah cannot be blank")
			return
		}
		if err := validator.ValidateInputSafe(trimmed, "logo_sekolah"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
			return
		}
		logo = &trimmed
	}

	cmd := profil_sekolah_service.UpdateProfilSekolahCmd{
		IDProfil:      profil_sekolah.IDProfil(1),
		EmailSekolah:  req.EmailSekolah,
		NoTelpSekolah: req.NoTelpSekolah,
		KepalaSekolah: req.KepalaSekolah,
		WakaSekolah:   req.WakaSekolah,
		NamaSekolah:   req.NamaSekolah,
		AlamatSekolah: req.AlamatSekolah,
		LogoSekolah:   logo,
	}

	if err := h.svc.UpdateProfilSekolah(r.Context(), cmd); err != nil {
		switch {
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed update profil sekolah")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "Success")
}
