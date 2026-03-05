package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

	httpx "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	aktivitas_user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/aktivitas_user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/create"
)

type UserHandler struct {
	svc           *user_service.CreateTx
	storeImage    httpx.ImageStore
	aktivitasUser *aktivitas_user_service.AktivitasUserService
}

func NewCreateUserHandler(svc *user_service.CreateTx, storeImage httpx.ImageStore, aktivitasUser *aktivitas_user_service.AktivitasUserService) *UserHandler {
	return &UserHandler{svc: svc, storeImage: storeImage, aktivitasUser: aktivitasUser}
}

func (h *UserHandler) CreateGuru(write http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(req.Context())
	if req.Method != http.MethodPost {
		httpResponse.WriteErr(write, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	if err := httpx.MultipartHeaderBodyValidator(write, req, 10<<20); err != nil {
		logger.Error(req.Context(), "failed parsing multipart form", "layer", "adapter.http.handler", "op", "user.create_guru", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrContentTypeMustMultipart):
			httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "bad request: content type must be multipart/form-data")
		default:
			httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid multipart form")
		}
		return
	}

	relPathPtr, err := httpx.StoreFileToDisk(req, "foto_profil", true, h.storeImage.SavePhotoRelative)
	if err != nil {
		switch {
		case errors.Is(err, coreerror.ErrMissingField):
			logger.Info(req.Context(), "missing foto file", "layer", "adapter.http.handler", "op", "user.create_guru", "err", err)
			httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : foto is required")
			return
		case errors.Is(err, coreerror.ErrFileTooLarge):
			logger.Info(req.Context(), "file too large", "layer", "adapter.http.handler", "op", "user.create_guru", "err", err)
			httpResponse.WriteErr(write, http.StatusBadRequest, "FILE_TOO_LARGE", "file too large")
			return
		default:
			logger.Error(req.Context(), "failed saving photo", "layer", "adapter.http.handler", "op", "user.create_guru", "err", err)
			httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed save photo")
			return
		}
	}
	relPath := ""
	if relPathPtr != nil {
		relPath = *relPathPtr
	}

	cmd := user_service.CreateGuruCmd{
		Username:     req.FormValue("username"),
		Email:        req.FormValue("email"),
		Password:     req.FormValue("password"),
		NamaLengkap:  req.FormValue("nama_lengkap"),
		JenisKelamin: req.FormValue("jenis_kelamin"),
		NoHp:         req.FormValue("no_hp"),
		Nip:          req.FormValue("nip"),
		Jabatan:      req.FormValue("jabatan"),
		BidangStudi:  req.FormValue("bidang_studi"),
		Foto:         relPath,
	}

	if cmd.Username == "" || cmd.Email == "" || cmd.Password == "" || cmd.NamaLengkap == "" || cmd.JenisKelamin == "" || cmd.NoHp == "" || cmd.Nip == "" || cmd.Jabatan == "" || cmd.BidangStudi == "" {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : invalid request body")
		return
	}
	if err := validator.ValidateInputSafe(cmd.Username, "username"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(cmd.Email, "email"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateEmail(cmd.Email); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidatePassword(cmd.Password); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(cmd.NamaLengkap, "nama_lengkap"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(cmd.JenisKelamin, "jenis_kelamin"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(cmd.NoHp, "no_hp"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(cmd.Nip, "nip"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(cmd.Jabatan, "jabatan"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(cmd.BidangStudi, "bidang_studi"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	actor, ok := middleware.ActorFromContext(req.Context())
	if !ok {
		logger.Error(req.Context(), "missing actor in context", "layer", "adapter.http.handler", "op", "user.create_guru", "err", "actor_not_found")
		httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	res, err := h.svc.CreateGuru(req.Context(), cmd, actor)
	if err != nil {
		logger.Error(req.Context(), "failed creating guru", "layer", "adapter.http.handler", "op", "user.create_guru", "err", err)
		switch {

		case errors.Is(err, coreerror.ErrForbidden):
			httpResponse.WriteErr(write, http.StatusForbidden, "FORBIDDEN", "forbidden")
			return

		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
			return

		case errors.Is(err, coreerror.ErrUsernameTaken):
			httpResponse.WriteErr(write, http.StatusConflict, "USERNAME_TAKEN", "username sudah terdaftar. data yang diinputkan harus unik")
			return
		case errors.Is(err, coreerror.ErrEmailTaken):
			httpResponse.WriteErr(write, http.StatusConflict, "EMAIL_TAKEN", "email sudah terdaftar. data yang diinputkan harus unik")
			return
		case errors.Is(err, coreerror.ErrNoHpTaken):
			httpResponse.WriteErr(write, http.StatusConflict, "NO_HP_TAKEN", "nomor HP sudah terdaftar. data yang diinputkan harus unik")
			return
		case errors.Is(err, coreerror.ErrNipTaken):
			httpResponse.WriteErr(write, http.StatusConflict, "NIP_TAKEN", "NIP sudah terdaftar. data yang diinputkan harus unik")
			return
		case errors.Is(err, coreerror.ErrConflict):
			httpResponse.WriteErr(write, http.StatusConflict, "CONFLICT", "data yang diinputkan sudah ada sebelumnya, harus unik")
			return

		default:
			httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed create guru")
			return
		}
	}

	if h.aktivitasUser != nil {
		aktivitasCmd := aktivitas_user_service.AktivitasUserCmd{
			IdPengguna:  actor.IdPengguna,
			Action:      aktivitas_user.CREATE,
			Description: "Membuat akun guru",
			IpAddress:   httpx.GetClientIP(req),
		}
		if err := h.aktivitasUser.CreateAktivitasUserService(req.Context(), aktivitasCmd); err != nil {
			logger.Error(req.Context(), "failed creating aktivitas user", "layer", "adapter.http.handler", "op", "user.create_guru.activity", "err", err)
		}
	}

	httpResponse.WriteOK(write, http.StatusOK, res, "Success")
}

func (h *UserHandler) CreateSiswa(write http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(req.Context())
	if req.Method != http.MethodPost {
		httpResponse.WriteErr(write, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	if err := httpx.MultipartHeaderBodyValidator(write, req, 10<<20); err != nil {
		logger.Error(req.Context(), "failed parsing multipart form", "layer", "adapter.http.handler", "op", "user.create_siswa", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrContentTypeMustMultipart):
			httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "bad request: content type must be multipart/form-data")
		default:
			httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid multipart form")
		}
		return
	}

	relPathPtr, err := httpx.StoreFileToDisk(req, "foto_profil", true, h.storeImage.SavePhotoRelative)
	if err != nil {
		switch {
		case errors.Is(err, coreerror.ErrMissingField):
			logger.Info(req.Context(), "missing foto file", "layer", "adapter.http.handler", "op", "user.create_siswa", "err", err)
			httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : foto is required")
			return
		case errors.Is(err, coreerror.ErrFileTooLarge):
			logger.Info(req.Context(), "file too large", "layer", "adapter.http.handler", "op", "user.create_siswa", "err", err)
			httpResponse.WriteErr(write, http.StatusBadRequest, "FILE_TOO_LARGE", "file too large")
			return
		default:
			logger.Error(req.Context(), "failed saving photo", "layer", "adapter.http.handler", "op", "user.create_siswa", "err", err)
			httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed save photo")
			return
		}
	}
	relPath := ""
	if relPathPtr != nil {
		relPath = *relPathPtr
	}

	// parse int fields
	noAbsen, err := strconv.Atoi(strings.TrimSpace(req.FormValue("no_absen")))
	if err != nil {
		logger.Info(req.Context(), "invalid no_absen", "layer", "adapter.http.handler", "op", "user.create_siswa", "err", err)
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", "no_absen must be a number")
		return
	}

	angkatan, err := strconv.Atoi(strings.TrimSpace(req.FormValue("angkatan")))
	if err != nil {
		logger.Info(req.Context(), "invalid angkatan", "layer", "adapter.http.handler", "op", "user.create_siswa", "err", err)
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", "angkatan must be a number")
		return
	}

	idNamaKelasRaw := strings.TrimSpace(req.FormValue("id_nama_kelas"))

	idNamaKelasInt, err := strconv.ParseInt(idNamaKelasRaw, 10, 64)
	if err != nil {
		logger.Info(req.Context(), "invalid id_nama_kelas", "layer", "adapter.http.handler", "op", "user.create_siswa", "err", err)
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", "id_nama_kelas must be a number")
		return
	}

	tanggalLahirStr := strings.TrimSpace(req.FormValue("tanggal_lahir"))
	tanggalLahir, err := time.Parse("2006-01-02", tanggalLahirStr)
	if err != nil {
		logger.Info(req.Context(), "invalid tanggal_lahir", "layer", "adapter.http.handler", "op", "user.create_siswa", "err", err)
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", "tanggal_lahir must be in YYYY-MM-DD format")
		return
	}

	cmd := user_service.CreateSiswaCmd{
		Username:     req.FormValue("username"),
		Email:        req.FormValue("email"),
		Password:     req.FormValue("password"),
		NamaLengkap:  req.FormValue("nama_lengkap"),
		JenisKelamin: req.FormValue("jenis_kelamin"),
		NoHp:         req.FormValue("no_hp"),
		Foto:         relPath,

		IdNamaKelas:  user.ID(idNamaKelasInt),
		Nisn:         req.FormValue("nisn"),
		NoAbsen:      noAbsen,
		Angkatan:     angkatan,
		TempatLahir:  req.FormValue("tempat_lahir"),
		TanggalLahir: tanggalLahir,
	}

	actor, ok := middleware.ActorFromContext(req.Context())
	if !ok {
		logger.Error(req.Context(), "missing actor in context", "layer", "adapter.http.handler", "op", "user.create_siswa", "err", "actor_not_found")
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
		cmd.IdNamaKelas == 0 ||
		cmd.Nisn == "" ||
		cmd.NoAbsen <= 0 ||
		cmd.Angkatan <= 0 ||
		cmd.TempatLahir == "" ||
		cmd.TanggalLahir.IsZero() {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request: invalid request body")
		return
	}
	if err := validator.ValidateInputSafe(cmd.Username, "username"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(cmd.Email, "email"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateEmail(cmd.Email); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidatePassword(cmd.Password); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(cmd.NamaLengkap, "nama_lengkap"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(cmd.JenisKelamin, "jenis_kelamin"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(cmd.NoHp, "no_hp"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(cmd.Nisn, "nisn"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}
	if err := validator.ValidateInputSafe(cmd.TempatLahir, "tempat_lahir"); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	res, err := h.svc.CreateSiswa(req.Context(), cmd, actor)
	if err != nil {
		logger.Error(req.Context(), "failed creating siswa", "layer", "adapter.http.handler", "op", "user.create_siswa", "err", err)
		switch {

		case errors.Is(err, coreerror.ErrForbidden):
			httpResponse.WriteErr(write, http.StatusForbidden, "FORBIDDEN", "forbidden")
			return

		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
			return

		case errors.Is(err, coreerror.ErrUsernameTaken):
			httpResponse.WriteErr(write, http.StatusConflict, "USERNAME_TAKEN", "username sudah terdaftar. data yang diinputkan harus unik")
			return
		case errors.Is(err, coreerror.ErrEmailTaken):
			httpResponse.WriteErr(write, http.StatusConflict, "EMAIL_TAKEN", "email sudah terdaftar. data yang diinputkan harus unik")
			return
		case errors.Is(err, coreerror.ErrNoHpTaken):
			httpResponse.WriteErr(write, http.StatusConflict, "NO_HP_TAKEN", "nomor HP sudah terdaftar. data yang diinputkan harus unik")
			return
		case errors.Is(err, coreerror.ErrNisnTaken):
			httpResponse.WriteErr(write, http.StatusConflict, "NISN_TAKEN", "NISN sudah terdaftar. data yang diinputkan harus unik")
			return
		case errors.Is(err, coreerror.ErrConflict):
			httpResponse.WriteErr(write, http.StatusConflict, "CONFLICT", "data yang diinputkan sudah ada sebelumnya, harus unik")
			return

		default:
			httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed create siswa")
			return
		}
	}

	if h.aktivitasUser != nil {
		aktivitasCmd := aktivitas_user_service.AktivitasUserCmd{
			IdPengguna:  actor.IdPengguna,
			Action:      aktivitas_user.CREATE,
			Description: "Membuat akun siswa",
			IpAddress:   httpx.GetClientIP(req),
		}
		if err := h.aktivitasUser.CreateAktivitasUserService(req.Context(), aktivitasCmd); err != nil {
			logger.Error(req.Context(), "failed creating aktivitas user", "layer", "adapter.http.handler", "op", "user.create_siswa.activity", "err", err)
		}
	}

	httpResponse.WriteOK(write, http.StatusOK, res, "Success")
}
