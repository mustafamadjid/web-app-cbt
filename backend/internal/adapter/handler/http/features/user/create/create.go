package httpx

import (
	"errors"
	"net/http"
	"strings"

	httpx "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	// out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/create"
)

type UserHandler struct {
	svc user_service.CreateTx
	storeImage httpx.ImageStore
}


func NewCreateUserHandler(svc user_service.CreateTx) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) CreateGuru(write http.ResponseWriter, req *http.Request){
	if req.Method != http.MethodPost {
		httpResponse.WriteErr(write,http.StatusMethodNotAllowed,"METHOD_NOT_ALLOWED","method not allowed")
		return
	}

	ct := req.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		httpResponse.WriteErr(write,http.StatusBadRequest,"BAD_REQUEST","Bad request : invalid content type")
		return
	}

	if err := req.ParseMultipartForm(10 << 20); err != nil {
		httpResponse.WriteErr(write,http.StatusBadRequest,"BAD_REQUEST","Bad request : invalid content type")
		return
	}

	file, fh,err := req.FormFile("foto")
		if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			httpResponse.WriteErr(write,http.StatusBadRequest,"BAD_REQUEST","Bad request : foto is required")
			return
		}

		httpResponse.WriteErr(write,http.StatusBadRequest,"BAD_REQUEST","Bad request : failed reading foto")
		return
	}
	defer file.Close()

	relPath, err := h.storeImage.SavePhotoRelative(file, fh)
	if err != nil {
		if errors.Is(err, coreerror.ErrFileTooLarge) {
			httpResponse.WriteErr(write, http.StatusBadRequest, "FILE_TOO_LARGE", "file too large")
			return
		}

		httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed save photo")
		return
	}
	

	cmd := user_service.CreateGuruCmd{
		Username:     strings.TrimSpace(req.FormValue("username")),
		Email:        strings.TrimSpace(req.FormValue("email")),
		Password:     req.FormValue("password"),
		NamaLengkap:  strings.TrimSpace(req.FormValue("nama_lengkap")),
		JenisKelamin: strings.TrimSpace(req.FormValue("jenis_kelamin")),
		NoHp:         strings.TrimSpace(req.FormValue("no_hp")),
		Nip:          strings.TrimSpace(req.FormValue("nip")),
		Jabatan:      strings.TrimSpace(req.FormValue("jabatan")),
		BidangStudi:  strings.TrimSpace(req.FormValue("bidang_studi")),
		Foto: relPath,
	}

	if cmd.Username == "" || cmd.Email == "" || cmd.Password == "" || cmd.NamaLengkap == "" || cmd.JenisKelamin == "" || cmd.NoHp == "" || cmd.Nip == "" || cmd.Jabatan == "" || cmd.BidangStudi == "" || cmd.Foto == "" {
		httpResponse.WriteErr(write,http.StatusBadRequest,"BAD_REQUEST","Bad request : invalid request body")
		return
	}
	actor,err := httpx.ActorFromContext(req.Context())
	if err != nil {
		httpResponse.WriteErr(write,http.StatusInternalServerError,"INTERNAL_SERVER_ERROR","internal server error : failed get actor from context")
		return
	}

	res, err := h.svc.CreateGuru(req.Context(),cmd,actor)
	if err != nil {
	switch {
	
	case errors.Is(err, coreerror.ErrForbidden):
		httpResponse.WriteErr(write, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return

	
	case errors.Is(err, coreerror.ErrInvalidInput):
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		return

	case errors.Is(err, coreerror.ErrUsernameTaken):
		httpResponse.WriteErr(write, http.StatusConflict, "USERNAME_TAKEN", "username already taken")
		return
	case errors.Is(err, coreerror.ErrNipTaken):
		httpResponse.WriteErr(write, http.StatusConflict, "NIP_TAKEN", "nip already taken")
		return

	default:
		httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed create guru")
		return
	}
}

httpResponse.WriteOK(write,http.StatusOK,res,"Success")
}