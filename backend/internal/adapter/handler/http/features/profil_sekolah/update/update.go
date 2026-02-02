package httpx

import (
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/profil_sekolah"
	profil_sekolah_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/profil_sekolah/update"
)

type UpdateProfilSekolahHandler struct {
	svc        *profil_sekolah_service.UpdateProfilSekolahService
	storeImage httphelper.ImageStore
}

func NewUpdateProfilSekolahHandler(svc *profil_sekolah_service.UpdateProfilSekolahService, storeImage httphelper.ImageStore) *UpdateProfilSekolahHandler {
	return &UpdateProfilSekolahHandler{svc: svc, storeImage: storeImage}
}

func (h *UpdateProfilSekolahHandler) UpdateProfilSekolah(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != http.MethodPatch {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: content type must be multipart/form-data")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid multipart form")
		return
	}

	email := strings.TrimSpace(r.FormValue("email_sekolah"))
	noTelp := strings.TrimSpace(r.FormValue("no_telp_sekolah"))
	kepala := strings.TrimSpace(r.FormValue("kepala_sekolah"))
	waka := strings.TrimSpace(r.FormValue("waka_sekolah"))
	nama := strings.TrimSpace(r.FormValue("nama_sekolah"))
	alamat := strings.TrimSpace(r.FormValue("alamat_sekolah"))

	if email == "" || noTelp == "" || kepala == "" || waka == "" || nama == "" || alamat == "" {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid request body")
		return
	}

	if err := validator.ValidateInputSafe(email, "email_sekolah"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(noTelp, "no_telp_sekolah"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(kepala, "kepala_sekolah"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(waka, "waka_sekolah"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(nama, "nama_sekolah"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(alamat, "alamat_sekolah"); err != nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	var logo *string
	file, fh, err := r.FormFile("logo_sekolah")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: failed reading logo_sekolah")
		return
	}
	if err == nil {
		defer file.Close()
		relPath, err := h.storeImage.SavePhotoRelative(file, fh)
		if err != nil {
			if errors.Is(err, coreerror.ErrFileTooLarge) {
				httpResponse.WriteErr(w, http.StatusBadRequest, "FILE_TOO_LARGE", "file too large")
				return
			}
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid logo_sekolah")
			return
		}
		logo = &relPath
	}

	cmd := profil_sekolah_service.UpdateProfilSekolahCmd{
		IDProfil:      profil_sekolah.IDProfil(1),
		EmailSekolah:  email,
		NoTelpSekolah: noTelp,
		KepalaSekolah: kepala,
		WakaSekolah:   waka,
		NamaSekolah:   nama,
		AlamatSekolah: alamat,
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
