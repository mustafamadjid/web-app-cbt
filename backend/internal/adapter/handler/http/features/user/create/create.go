package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	httpx "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

	// out "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/create"
)

type UserHandler struct {
	svc        *user_service.CreateTx
	storeImage httpx.ImageStore
}

func NewCreateUserHandler(svc *user_service.CreateTx, storeImage httpx.ImageStore) *UserHandler {
	return &UserHandler{svc: svc, storeImage: storeImage}
}

func (h *UserHandler) CreateGuru(write http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Method != http.MethodPost {
		httpResponse.WriteErr(write, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	ct := req.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : invalid content type")
		return
	}

	if err := req.ParseMultipartForm(10 << 20); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : invalid content type")
		return
	}

	file, fh, err := req.FormFile("foto")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : foto is required")
			return
		}

		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : failed reading foto")
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
		Foto:         relPath,
	}

	if cmd.Username == "" || cmd.Email == "" || cmd.Password == "" || cmd.NamaLengkap == "" || cmd.JenisKelamin == "" || cmd.NoHp == "" || cmd.Nip == "" || cmd.Jabatan == "" || cmd.BidangStudi == "" || cmd.Foto == "" {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : invalid request body")
		return
	}
	if err := validateInputSafe(cmd.Username, "username"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validateInputSafe(cmd.Email, "email"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validateInputSafe(cmd.NamaLengkap, "nama_lengkap"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validateInputSafe(cmd.JenisKelamin, "jenis_kelamin"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validateInputSafe(cmd.NoHp, "no_hp"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validateInputSafe(cmd.Nip, "nip"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validateInputSafe(cmd.Jabatan, "jabatan"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validateInputSafe(cmd.BidangStudi, "bidang_studi"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	actor, ok := middleware.ActorFromContext(req.Context())
	if !ok {
		httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	res, err := h.svc.CreateGuru(req.Context(), cmd, actor)
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

	httpResponse.WriteOK(write, http.StatusOK, res, "Success")
}

func (h *UserHandler) CreateSiswa(write http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Method != http.MethodPost {
		httpResponse.WriteErr(write, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	ct := req.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : invalid content type")
		return
	}

	if err := req.ParseMultipartForm(10 << 20); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : invalid content type")
		return
	}

	file, fh, err := req.FormFile("foto")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : foto is required")
			return
		}

		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : failed reading foto")
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

	// parse int fields
	noAbsen, err := strconv.Atoi(strings.TrimSpace(req.FormValue("no_absen")))
	if err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", "no_absen must be a number")
		return
	}

	angkatan, err := strconv.Atoi(strings.TrimSpace(req.FormValue("angkatan")))
	if err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", "angkatan must be a number")
		return
	}

	idTingkatKelasRaw := strings.TrimSpace(req.FormValue("id_tingkat_kelas"))
	idNamaKelasRaw := strings.TrimSpace(req.FormValue("id_nama_kelas"))

	idTingkatKelasInt, err := strconv.ParseInt(idTingkatKelasRaw, 10, 64)
	if err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", "id_tingkat_kelas must be a number")
		return
	}

	idNamaKelasInt, err := strconv.ParseInt(idNamaKelasRaw, 10, 64)
	if err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", "id_nama_kelas must be a number")
		return
	}

	tanggalLahirStr := strings.TrimSpace(req.FormValue("tanggal_lahir"))
	tanggalLahir, err := time.Parse("2006-01-02", tanggalLahirStr)
	if err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", "tanggal_lahir must be in YYYY-MM-DD format")
		return
	}

	cmd := user_service.CreateSiswaCmd{
		Username:     strings.TrimSpace(req.FormValue("username")),
		Email:        strings.TrimSpace(req.FormValue("email")),
		Password:     req.FormValue("password"), // jangan Trim password
		NamaLengkap:  strings.TrimSpace(req.FormValue("nama_lengkap")),
		JenisKelamin: strings.TrimSpace(req.FormValue("jenis_kelamin")),
		NoHp:         strings.TrimSpace(req.FormValue("no_hp")),
		Foto:         relPath,

		IdTingkatKelas: user.ID(idTingkatKelasInt),
		IdNamaKelas:    user.ID(idNamaKelasInt),
		Nisn:           strings.TrimSpace(req.FormValue("nisn")),
		NoAbsen:        noAbsen,
		Angkatan:       angkatan,
		TempatLahir:    strings.TrimSpace(req.FormValue("tempat_lahir")),
		TanggalLahir:   tanggalLahir,
	}

	actor, ok := middleware.ActorFromContext(req.Context())
	if !ok {
		httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}
	if cmd.Username == "" ||
		cmd.Email == "" ||
		cmd.Password == "" ||
		cmd.NamaLengkap == "" ||
		cmd.JenisKelamin == "" ||
		cmd.NoHp == "" ||
		cmd.Foto == "" ||
		cmd.IdTingkatKelas == 0 ||
		cmd.IdNamaKelas == 0 ||
		cmd.Nisn == "" ||
		cmd.NoAbsen <= 0 ||
		cmd.Angkatan <= 0 ||
		cmd.TempatLahir == "" ||
		cmd.TanggalLahir.IsZero() {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request: invalid request body")
		return
	}
	if err := validateInputSafe(cmd.Username, "username"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validateInputSafe(cmd.Email, "email"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validateInputSafe(cmd.NamaLengkap, "nama_lengkap"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validateInputSafe(cmd.JenisKelamin, "jenis_kelamin"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validateInputSafe(cmd.NoHp, "no_hp"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validateInputSafe(cmd.Nisn, "nisn"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validateInputSafe(cmd.TempatLahir, "tempat_lahir"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	res, err := h.svc.CreateSiswa(req.Context(), cmd, actor)
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
		case errors.Is(err, coreerror.ErrNisnTaken):
			httpResponse.WriteErr(write, http.StatusConflict, "NISN_TAKEN", "NISN already taken")
			return

		default:
			httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed create guru")
			return
		}
	}

	httpResponse.WriteOK(write, http.StatusOK, res, "Success")
}
