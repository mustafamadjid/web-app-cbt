package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/middleware"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	userportout "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
	createSvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/create"
	updateSvc "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/update"
)

type UserHandler struct {
	createSvc userportout.CreateUserService
	updateSvc userportout.UpdateUserService
	deleteSvc userportout.DeleteUserService
}

func NewUserHandler(createSvc userportout.CreateUserService, updateSvc userportout.UpdateUserService, deleteSvc userportout.DeleteUserService) *UserHandler {
	return &UserHandler{
		createSvc: createSvc,
		updateSvc: updateSvc,
		deleteSvc: deleteSvc,
	}
}

func (h *UserHandler) CreateGuru(write http.ResponseWriter, req *http.Request) {
	actor, ok := middleware.ActorFromContext(req.Context())
	if !ok {
		httpResponse.WriteErr(write, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	var reqBody CreateGuruRequest
	if err := decodeJSON(req, &reqBody); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "bad request")
		return
	}

	cmd := createSvc.CreateGuruCmd{
		Username:     reqBody.Username,
		Email:        reqBody.Email,
		Password:     reqBody.Password,
		NamaLengkap:  reqBody.NamaLengkap,
		JenisKelamin: reqBody.JenisKelamin,
		NoHp:         reqBody.NoHp,
		Foto:         reqBody.Foto,
		Nip:          reqBody.Nip,
		Jabatan:      reqBody.Jabatan,
		BidangStudi:  reqBody.BidangStudi,
	}

	res, err := h.createSvc.CreateGuru(req.Context(), cmd, actor)
	if err != nil {
		writeUserError(write, err)
		return
	}

	responseData := CreateGuruResponse{
		IdPengguna:   res.IdPengguna,
		IdProfilGuru: res.IdProfilGuru,
	}

	httpResponse.WriteOK(write, http.StatusCreated, responseData, "success")
}

func (h *UserHandler) CreateSiswa(write http.ResponseWriter, req *http.Request) {
	actor, ok := middleware.ActorFromContext(req.Context())
	if !ok {
		httpResponse.WriteErr(write, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	var reqBody CreateSiswaRequest
	if err := decodeJSON(req, &reqBody); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "bad request")
		return
	}

	cmd := createSvc.CreateSiswaCmd{
		Username:       reqBody.Username,
		Email:          reqBody.Email,
		Password:       reqBody.Password,
		NamaLengkap:    reqBody.NamaLengkap,
		JenisKelamin:   reqBody.JenisKelamin,
		NoHp:           reqBody.NoHp,
		Foto:           reqBody.Foto,
		IdTingkatKelas: reqBody.IdTingkatKelas,
		IdNamaKelas:    reqBody.IdNamaKelas,
		Nisn:           reqBody.Nisn,
		NoAbsen:        reqBody.NoAbsen,
		Angkatan:       reqBody.Angkatan,
		TempatLahir:    reqBody.TempatLahir,
		TanggalLahir:   reqBody.TanggalLahir,
	}

	res, err := h.createSvc.CreateSiswa(req.Context(), cmd, actor)
	if err != nil {
		writeUserError(write, err)
		return
	}

	responseData := CreateSiswaResponse{
		IdPengguna:    res.IdPengguna,
		IdProfilSiswa: res.IdProfilSiswa,
	}

	httpResponse.WriteOK(write, http.StatusCreated, responseData, "success")
}

func (h *UserHandler) UpdateGuru(write http.ResponseWriter, req *http.Request) {
	actor, ok := middleware.ActorFromContext(req.Context())
	if !ok {
		httpResponse.WriteErr(write, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	var reqBody UpdateGuruRequest
	if err := decodeJSON(req, &reqBody); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "bad request")
		return
	}

	statusAkun, err := parseStatusAkun(reqBody.StatusAkun)
	if err != nil {
		writeUserError(write, err)
		return
	}

	role, err := parseRole(reqBody.Role)
	if err != nil {
		writeUserError(write, err)
		return
	}

	cmd := updateSvc.UpdateGuruCmd{
		IdPengguna:   reqBody.IdPengguna,
		Username:     reqBody.Username,
		Email:        reqBody.Email,
		NamaLengkap:  reqBody.NamaLengkap,
		JenisKelamin: reqBody.JenisKelamin,
		NoHp:         reqBody.NoHp,
		Foto:         reqBody.Foto,
		StatusAkun:   statusAkun,
		Role:         role,
		Nip:          reqBody.Nip,
		Jabatan:      reqBody.Jabatan,
		BidangStudi:  reqBody.BidangStudi,
	}

	if err := h.updateSvc.UpdateGuru(req.Context(), cmd, actor); err != nil {
		writeUserError(write, err)
		return
	}

	httpResponse.WriteOKNoData(write, http.StatusOK, "success")
}

func (h *UserHandler) UpdateSiswa(write http.ResponseWriter, req *http.Request) {
	actor, ok := middleware.ActorFromContext(req.Context())
	if !ok {
		httpResponse.WriteErr(write, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	var reqBody UpdateSiswaRequest
	if err := decodeJSON(req, &reqBody); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "bad request")
		return
	}

	statusAkun, err := parseStatusAkun(reqBody.StatusAkun)
	if err != nil {
		writeUserError(write, err)
		return
	}

	role, err := parseRole(reqBody.Role)
	if err != nil {
		writeUserError(write, err)
		return
	}

	cmd := updateSvc.UpdateSiswaCmd{
		IdPengguna:     reqBody.IdPengguna,
		Username:       reqBody.Username,
		Email:          reqBody.Email,
		NamaLengkap:    reqBody.NamaLengkap,
		JenisKelamin:   reqBody.JenisKelamin,
		NoHp:           reqBody.NoHp,
		Foto:           reqBody.Foto,
		StatusAkun:     statusAkun,
		Role:           role,
		IdTingkatKelas: reqBody.IdTingkatKelas,
		IdNamaKelas:    reqBody.IdNamaKelas,
		Nisn:           reqBody.Nisn,
		NoAbsen:        reqBody.NoAbsen,
		Angkatan:       reqBody.Angkatan,
		TempatLahir:    reqBody.TempatLahir,
		TanggalLahir:   reqBody.TanggalLahir,
	}

	if err := h.updateSvc.UpdateSiswa(req.Context(), cmd, actor); err != nil {
		writeUserError(write, err)
		return
	}

	httpResponse.WriteOKNoData(write, http.StatusOK, "success")
}

func (h *UserHandler) Delete(write http.ResponseWriter, req *http.Request) {
	actor, ok := middleware.ActorFromContext(req.Context())
	if !ok {
		httpResponse.WriteErr(write, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}
	if actor.Role != user.ADMIN {
		httpResponse.WriteErr(write, http.StatusForbidden, "FORBIDDEN", "forbidden")
		return
	}

	var reqBody DeleteUserRequest
	if err := decodeJSON(req, &reqBody); err != nil {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "bad request")
		return
	}

	if reqBody.IdPengguna == 0 {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "id_pengguna required")
		return
	}

	if err := h.deleteSvc.Delete(req.Context(), reqBody.IdPengguna); err != nil {
		writeUserError(write, err)
		return
	}

	httpResponse.WriteOKNoData(write, http.StatusOK, "success")
}

func decodeJSON(req *http.Request, dest any) error {
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dest)
}

func parseStatusAkun(status *string) (*user.StatusAkun, error) {
	if status == nil {
		return nil, nil
	}
	normalized := strings.ToUpper(strings.TrimSpace(*status))
	switch normalized {
	case string(user.AKTIF), string(user.NONAKTIF):
		value := user.StatusAkun(normalized)
		return &value, nil
	default:
		return nil, coreerror.ErrInvalidStatusAkun
	}
}

func parseRole(role *string) (*user.Role, error) {
	if role == nil {
		return nil, nil
	}
	normalized := user.Role(strings.ToUpper(strings.TrimSpace(*role)))
	if !normalized.ValidRole() {
		return nil, coreerror.ErrInvalidInput
	}
	return &normalized, nil
}

func writeUserError(write http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, coreerror.ErrForbidden):
		httpResponse.WriteErr(write, http.StatusForbidden, "FORBIDDEN", "forbidden")
	case errors.Is(err, coreerror.ErrUnauthorized):
		httpResponse.WriteErr(write, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
	case errors.Is(err, coreerror.ErrNotFound):
		httpResponse.WriteErr(write, http.StatusNotFound, "NOT_FOUND", "not found")
	case errors.Is(err, coreerror.ErrUsernameTaken):
		httpResponse.WriteErr(write, http.StatusConflict, "USERNAME_TAKEN", "username already taken")
	case errors.Is(err, coreerror.ErrNipTaken):
		httpResponse.WriteErr(write, http.StatusConflict, "NIP_TAKEN", "nip already taken")
	case errors.Is(err, coreerror.ErrNisnTaken):
		httpResponse.WriteErr(write, http.StatusConflict, "NISN_TAKEN", "nisn already taken")
	case errors.Is(err, coreerror.ErrNoFieldToUpdate):
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "no field to update")
	case errors.Is(err, coreerror.ErrInvalidStatusAkun):
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "invalid status_akun")
	case errors.Is(err, user.ErrInvalidEmail):
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "invalid email")
	case errors.Is(err, user.ErrInvalidNIP):
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "invalid nip")
	case errors.Is(err, user.ErrInvalidNISN):
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "invalid nisn")
	case errors.Is(err, user.ErrInvalidAbsen):
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "invalid no_absen")
	case errors.Is(err, user.ErrInvalidAngkatan):
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "invalid angkatan")
	case isBadRequestMessage(err):
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", err.Error())
	default:
		httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error")
	}
}

func isBadRequestMessage(err error) bool {
	switch err.Error() {
	case "Id pengguna required", "username cannot be empty", "nama_lengkap cannot be empty":
		return true
	default:
		return false
	}
}
