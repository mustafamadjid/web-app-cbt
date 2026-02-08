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
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
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
	logger := corelog.FromContext(r.Context())
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
		logger.Error(r.Context(), "failed parsing multipart form", "layer", "adapter.http.handler", "op", "profil_sekolah.update", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid multipart form")
		return
	}

	mf := r.MultipartForm
	getOptional := func(key string) (string, bool) {
		if mf == nil || mf.Value == nil {
			return "", false
		}
		vals, ok := mf.Value[key]
		if !ok || len(vals) == 0 {
			return "", false
		}
		return strings.TrimSpace(vals[0]), true
	}

	var updateRequest updateProfilSekolahRequest

	if email, ok := getOptional("email_sekolah"); ok {
		if email == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: email_sekolah is required")
			return
		}
		if err := validator.ValidateInputSafe(email, "email_sekolah"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
			return
		}
		updateRequest.EmailSekolah = &email
	}
	if noTelp, ok := getOptional("no_telp_sekolah"); ok {
		if noTelp == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: no_telp_sekolah is required")
			return
		}
		if err := validator.ValidateInputSafe(noTelp, "no_telp_sekolah"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
			return
		}
		updateRequest.NoTelpSekolah = &noTelp
	}
	if kepala, ok := getOptional("kepala_sekolah"); ok {
		if kepala == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: kepala_sekolah is required")
			return
		}
		if err := validator.ValidateInputSafe(kepala, "kepala_sekolah"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
			return
		}
		updateRequest.KepalaSekolah = &kepala
	}
	if waka, ok := getOptional("waka_sekolah"); ok {
		if waka == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: waka_sekolah is required")
			return
		}
		if err := validator.ValidateInputSafe(waka, "waka_sekolah"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
			return
		}
		updateRequest.WakaSekolah = &waka
	}
	if nama, ok := getOptional("nama_sekolah"); ok {
		if nama == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: nama_sekolah is required")
			return
		}
		if err := validator.ValidateInputSafe(nama, "nama_sekolah"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
			return
		}
		updateRequest.NamaSekolah = &nama
	}
	if alamat, ok := getOptional("alamat_sekolah"); ok {
		if alamat == "" {
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: alamat_sekolah is required")
			return
		}
		if err := validator.ValidateInputSafe(alamat, "alamat_sekolah"); err != nil {
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
			return
		}
		updateRequest.AlamatSekolah = &alamat
	}

	var logoPtr *string
	file, fh, err := r.FormFile("logo_sekolah")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		logger.Error(r.Context(), "failed reading logo", "layer", "adapter.http.handler", "op", "profil_sekolah.update", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: failed reading logo_sekolah")
		return
	}
	if err == nil {
		defer file.Close()
		relPath, err := h.storeImage.SavePhotoRelative(file, fh)
		if err != nil {
			if errors.Is(err, coreerror.ErrFileTooLarge) {
				logger.Info(r.Context(), "logo too large", "layer", "adapter.http.handler", "op", "profil_sekolah.update", "err", err)
				httpResponse.WriteErr(w, http.StatusBadRequest, "FILE_TOO_LARGE", "file too large")
				return
			}
			logger.Error(r.Context(), "failed saving logo", "layer", "adapter.http.handler", "op", "profil_sekolah.update", "err", err)
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid logo_sekolah")
			return
		}
		logoPtr = &relPath
	}

	if updateRequest.IsEmpty() && logoPtr == nil {
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "no fields to update")
		return
	}

	cmd := profil_sekolah_service.UpdateProfilSekolahCmd{
		IDProfil:      profil_sekolah.IDProfil(1),
		EmailSekolah:  updateRequest.EmailSekolah,
		NoTelpSekolah: updateRequest.NoTelpSekolah,
		KepalaSekolah: updateRequest.KepalaSekolah,
		WakaSekolah:   updateRequest.WakaSekolah,
		NamaSekolah:   updateRequest.NamaSekolah,
		AlamatSekolah: updateRequest.AlamatSekolah,
		LogoSekolah:   logoPtr,
	}

	if err := h.svc.UpdateProfilSekolah(r.Context(), cmd); err != nil {
		logger.Error(r.Context(), "failed updating profil sekolah", "layer", "adapter.http.handler", "op", "profil_sekolah.update", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNoFieldToUpdate):
			httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "no fields to update")
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed update profil sekolah")
		}
		return
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "Success")
}
