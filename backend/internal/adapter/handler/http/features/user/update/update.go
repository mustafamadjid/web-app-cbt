package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	httpx "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/update"
)

type UpdateHandler struct {
	svc        user_service.UpdateTx
	storeImage httpx.ImageStore
}

type requestError struct {
	status  int
	code    string
	message string
}

func (e requestError) Error() string {
	return e.message
}

func NewUpdateUserHandler(svc user_service.UpdateTx) *UpdateHandler {
	return &UpdateHandler{svc: svc}
}

func (h *UpdateHandler) UpdateGuru(write http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		httpResponse.WriteErr(write, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	cmd, err := h.parseUpdateGuruRequest(req)
	if err != nil {
		h.writeRequestError(write, err)
		return
	}

	actor, err := httpx.ActorFromContext(req.Context())
	if err != nil {
		httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	if err := h.svc.UpdateGuru(req.Context(), cmd, actor); err != nil {
		h.writeUpdateError(write, err, "guru")
		return
	}

	httpResponse.WriteOKNoData(write, http.StatusOK, "Success")
}

func (h *UpdateHandler) UpdateSiswa(write http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		httpResponse.WriteErr(write, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	cmd, err := h.parseUpdateSiswaRequest(req)
	if err != nil {
		h.writeRequestError(write, err)
		return
	}

	actor, err := httpx.ActorFromContext(req.Context())
	if err != nil {
		httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	if err := h.svc.UpdateSiswa(req.Context(), cmd, actor); err != nil {
		h.writeUpdateError(write, err, "siswa")
		return
	}

	httpResponse.WriteOKNoData(write, http.StatusOK, "Success")
}

func (h *UpdateHandler) parseUpdateGuruRequest(req *http.Request) (user_service.UpdateGuruCmd, error) {
	ct := req.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "multipart/form-data") {
		return h.parseUpdateGuruForm(req)
	}
	if strings.HasPrefix(ct, "application/json") {
		return h.parseUpdateGuruJSON(req)
	}

	return user_service.UpdateGuruCmd{}, requestError{
		status:  http.StatusBadRequest,
		code:    "BAD_REQUEST",
		message: "Bad request : invalid content type",
	}
}

func (h *UpdateHandler) parseUpdateSiswaRequest(req *http.Request) (user_service.UpdateSiswaCmd, error) {
	ct := req.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "multipart/form-data") {
		return h.parseUpdateSiswaForm(req)
	}
	if strings.HasPrefix(ct, "application/json") {
		return h.parseUpdateSiswaJSON(req)
	}

	return user_service.UpdateSiswaCmd{}, requestError{
		status:  http.StatusBadRequest,
		code:    "BAD_REQUEST",
		message: "Bad request : invalid content type",
	}
}

func (h *UpdateHandler) parseUpdateGuruForm(req *http.Request) (user_service.UpdateGuruCmd, error) {
	if err := req.ParseMultipartForm(10 << 20); err != nil {
		return user_service.UpdateGuruCmd{}, requestError{
			status:  http.StatusBadRequest,
			code:    "BAD_REQUEST",
			message: "Bad request : invalid content type",
		}
	}

	idPengguna, err := parseRequiredID(req.MultipartForm.Value, "id_pengguna")
	if err != nil {
		return user_service.UpdateGuruCmd{}, err
	}

	cmd := user_service.UpdateGuruCmd{IdPengguna: idPengguna}
	applyOptionalStrings(req.MultipartForm.Value, map[string]**string{
		"username":      &cmd.Username,
		"email":         &cmd.Email,
		"nama_lengkap":  &cmd.NamaLengkap,
		"jenis_kelamin": &cmd.JenisKelamin,
		"no_hp":         &cmd.NoHp,
		"nip":           &cmd.Nip,
		"jabatan":       &cmd.Jabatan,
		"bidang_studi":  &cmd.BidangStudi,
	})

	statusAkun, err := parseOptionalStatus(req.MultipartForm.Value, "status_akun")
	if err != nil {
		return user_service.UpdateGuruCmd{}, err
	}
	cmd.StatusAkun = statusAkun

	role, err := parseOptionalRole(req.MultipartForm.Value, "role")
	if err != nil {
		return user_service.UpdateGuruCmd{}, err
	}
	cmd.Role = role

	if relPath, err := h.parseOptionalFoto(req); err != nil {
		return user_service.UpdateGuruCmd{}, err
	} else if relPath != nil {
		cmd.Foto = relPath
	}

	return cmd, nil
}

func (h *UpdateHandler) parseUpdateSiswaForm(req *http.Request) (user_service.UpdateSiswaCmd, error) {
	if err := req.ParseMultipartForm(10 << 20); err != nil {
		return user_service.UpdateSiswaCmd{}, requestError{
			status:  http.StatusBadRequest,
			code:    "BAD_REQUEST",
			message: "Bad request : invalid content type",
		}
	}

	idPengguna, err := parseRequiredID(req.MultipartForm.Value, "id_pengguna")
	if err != nil {
		return user_service.UpdateSiswaCmd{}, err
	}

	cmd := user_service.UpdateSiswaCmd{IdPengguna: idPengguna}
	applyOptionalStrings(req.MultipartForm.Value, map[string]**string{
		"username":      &cmd.Username,
		"email":         &cmd.Email,
		"nama_lengkap":  &cmd.NamaLengkap,
		"jenis_kelamin": &cmd.JenisKelamin,
		"no_hp":         &cmd.NoHp,
		"nisn":          &cmd.Nisn,
		"tempat_lahir":  &cmd.TempatLahir,
	})

	statusAkun, err := parseOptionalStatus(req.MultipartForm.Value, "status_akun")
	if err != nil {
		return user_service.UpdateSiswaCmd{}, err
	}
	cmd.StatusAkun = statusAkun

	role, err := parseOptionalRole(req.MultipartForm.Value, "role")
	if err != nil {
		return user_service.UpdateSiswaCmd{}, err
	}
	cmd.Role = role

	if idTingkat, err := parseOptionalID(req.MultipartForm.Value, "id_tingkat_kelas"); err != nil {
		return user_service.UpdateSiswaCmd{}, err
	} else if idTingkat != nil {
		cmd.IdTingkatKelas = idTingkat
	}

	if idNama, err := parseOptionalID(req.MultipartForm.Value, "id_nama_kelas"); err != nil {
		return user_service.UpdateSiswaCmd{}, err
	} else if idNama != nil {
		cmd.IdNamaKelas = idNama
	}

	if noAbsen, err := parseOptionalInt(req.MultipartForm.Value, "no_absen"); err != nil {
		return user_service.UpdateSiswaCmd{}, err
	} else if noAbsen != nil {
		cmd.NoAbsen = noAbsen
	}

	if angkatan, err := parseOptionalInt(req.MultipartForm.Value, "angkatan"); err != nil {
		return user_service.UpdateSiswaCmd{}, err
	} else if angkatan != nil {
		cmd.Angkatan = angkatan
	}

	if tanggal, err := parseOptionalDate(req.MultipartForm.Value, "tanggal_lahir"); err != nil {
		return user_service.UpdateSiswaCmd{}, err
	} else if tanggal != nil {
		cmd.TanggalLahir = tanggal
	}

	if relPath, err := h.parseOptionalFoto(req); err != nil {
		return user_service.UpdateSiswaCmd{}, err
	} else if relPath != nil {
		cmd.Foto = relPath
	}

	return cmd, nil
}

func (h *UpdateHandler) parseUpdateGuruJSON(req *http.Request) (user_service.UpdateGuruCmd, error) {
	var payload UpdateGuruRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		return user_service.UpdateGuruCmd{}, requestError{
			status:  http.StatusBadRequest,
			code:    "BAD_REQUEST",
			message: "Bad request : invalid request body",
		}
	}
	if payload.IdPengguna == 0 {
		return user_service.UpdateGuruCmd{}, requestError{
			status:  http.StatusBadRequest,
			code:    "BAD_REQUEST",
			message: "Bad request : id_pengguna is required",
		}
	}

	statusAkun, err := parseOptionalStatusFromString(payload.StatusAkun)
	if err != nil {
		return user_service.UpdateGuruCmd{}, err
	}

	role, err := parseOptionalRoleFromString(payload.Role)
	if err != nil {
		return user_service.UpdateGuruCmd{}, err
	}

	return user_service.UpdateGuruCmd{
		IdPengguna:   payload.IdPengguna,
		Username:     payload.Username,
		Email:        payload.Email,
		NamaLengkap:  payload.NamaLengkap,
		JenisKelamin: payload.JenisKelamin,
		NoHp:         payload.NoHp,
		Foto:         payload.Foto,
		StatusAkun:   statusAkun,
		Role:         role,
		Nip:          payload.Nip,
		Jabatan:      payload.Jabatan,
		BidangStudi:  payload.BidangStudi,
	}, nil
}

func (h *UpdateHandler) parseUpdateSiswaJSON(req *http.Request) (user_service.UpdateSiswaCmd, error) {
	var payload UpdateSiswaRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		return user_service.UpdateSiswaCmd{}, requestError{
			status:  http.StatusBadRequest,
			code:    "BAD_REQUEST",
			message: "Bad request : invalid request body",
		}
	}
	if payload.IdPengguna == 0 {
		return user_service.UpdateSiswaCmd{}, requestError{
			status:  http.StatusBadRequest,
			code:    "BAD_REQUEST",
			message: "Bad request : id_pengguna is required",
		}
	}

	statusAkun, err := parseOptionalStatusFromString(payload.StatusAkun)
	if err != nil {
		return user_service.UpdateSiswaCmd{}, err
	}

	role, err := parseOptionalRoleFromString(payload.Role)
	if err != nil {
		return user_service.UpdateSiswaCmd{}, err
	}

	return user_service.UpdateSiswaCmd{
		IdPengguna:     payload.IdPengguna,
		Username:       payload.Username,
		Email:          payload.Email,
		NamaLengkap:    payload.NamaLengkap,
		JenisKelamin:   payload.JenisKelamin,
		NoHp:           payload.NoHp,
		Foto:           payload.Foto,
		StatusAkun:     statusAkun,
		Role:           role,
		IdTingkatKelas: payload.IdTingkatKelas,
		IdNamaKelas:    payload.IdNamaKelas,
		Nisn:           payload.Nisn,
		NoAbsen:        payload.NoAbsen,
		Angkatan:       payload.Angkatan,
		TempatLahir:    payload.TempatLahir,
		TanggalLahir:   payload.TanggalLahir,
	}, nil
}

func (h *UpdateHandler) parseOptionalFoto(req *http.Request) (*string, error) {
	file, fh, err := req.FormFile("foto")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return nil, nil
		}
		return nil, requestError{
			status:  http.StatusBadRequest,
			code:    "BAD_REQUEST",
			message: "Bad request : failed reading foto",
		}
	}
	defer file.Close()

	relPath, err := h.storeImage.SavePhotoRelative(file, fh)
	if err != nil {
		if errors.Is(err, coreerror.ErrFileTooLarge) {
			return nil, requestError{
				status:  http.StatusBadRequest,
				code:    "FILE_TOO_LARGE",
				message: "file too large",
			}
		}
		return nil, requestError{
			status:  http.StatusInternalServerError,
			code:    "INTERNAL_SERVER_ERROR",
			message: "internal server error : failed save photo",
		}
	}

	return &relPath, nil
}

func (h *UpdateHandler) writeRequestError(write http.ResponseWriter, err error) {
	if reqErr, ok := err.(requestError); ok {
		httpResponse.WriteErr(write, reqErr.status, reqErr.code, reqErr.message)
		return
	}
	httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : invalid request body")
}

func (h *UpdateHandler) writeUpdateError(write http.ResponseWriter, err error, role string) {
	switch {
	case errors.Is(err, coreerror.ErrForbidden):
		httpResponse.WriteErr(write, http.StatusForbidden, "FORBIDDEN", "forbidden")
	case errors.Is(err, coreerror.ErrNoFieldToUpdate):
		httpResponse.WriteErr(write, http.StatusBadRequest, "NO_FIELD_TO_UPDATE", "no field to update")
	case errors.Is(err, coreerror.ErrUsernameTaken):
		httpResponse.WriteErr(write, http.StatusConflict, "USERNAME_TAKEN", "username already taken")
	case errors.Is(err, coreerror.ErrNipTaken):
		httpResponse.WriteErr(write, http.StatusConflict, "NIP_TAKEN", "nip already taken")
	case errors.Is(err, coreerror.ErrNisnTaken):
		httpResponse.WriteErr(write, http.StatusConflict, "NISN_TAKEN", "nisn already taken")
	case errors.Is(err, coreerror.ErrInvalidInput),
		errors.Is(err, coreerror.ErrInvalidStatusAkun),
		errors.Is(err, user.ErrInvalidEmail),
		errors.Is(err, user.ErrInvalidNISN),
		errors.Is(err, user.ErrInvalidAbsen),
		errors.Is(err, user.ErrInvalidAngkatan):
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
	case errors.Is(err, coreerror.ErrNotFound):
		httpResponse.WriteErr(write, http.StatusNotFound, "NOT_FOUND", "data not found")
	default:
		httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed update "+role)
	}
}

func parseRequiredID(values map[string][]string, key string) (user.ID, error) {
	raw, ok := values[key]
	if !ok || len(raw) == 0 {
		return 0, requestError{
			status:  http.StatusBadRequest,
			code:    "BAD_REQUEST",
			message: "Bad request : " + key + " is required",
		}
	}
	val := strings.TrimSpace(raw[0])
	if val == "" {
		return 0, requestError{
			status:  http.StatusBadRequest,
			code:    "BAD_REQUEST",
			message: "Bad request : " + key + " is required",
		}
	}
	idInt, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, requestError{
			status:  http.StatusBadRequest,
			code:    "INVALID_INPUT",
			message: key + " must be a number",
		}
	}
	return user.ID(idInt), nil
}

func parseOptionalID(values map[string][]string, key string) (*user.ID, error) {
	raw, ok := values[key]
	if !ok || len(raw) == 0 {
		return nil, nil
	}
	val := strings.TrimSpace(raw[0])
	if val == "" {
		return nil, requestError{
			status:  http.StatusBadRequest,
			code:    "INVALID_INPUT",
			message: key + " must be a number",
		}
	}
	idInt, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return nil, requestError{
			status:  http.StatusBadRequest,
			code:    "INVALID_INPUT",
			message: key + " must be a number",
		}
	}
	id := user.ID(idInt)
	return &id, nil
}

func parseOptionalInt(values map[string][]string, key string) (*int, error) {
	raw, ok := values[key]
	if !ok || len(raw) == 0 {
		return nil, nil
	}
	val := strings.TrimSpace(raw[0])
	if val == "" {
		return nil, requestError{
			status:  http.StatusBadRequest,
			code:    "INVALID_INPUT",
			message: key + " must be a number",
		}
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		return nil, requestError{
			status:  http.StatusBadRequest,
			code:    "INVALID_INPUT",
			message: key + " must be a number",
		}
	}
	return &parsed, nil
}

func parseOptionalDate(values map[string][]string, key string) (*time.Time, error) {
	raw, ok := values[key]
	if !ok || len(raw) == 0 {
		return nil, nil
	}
	val := strings.TrimSpace(raw[0])
	if val == "" {
		return nil, requestError{
			status:  http.StatusBadRequest,
			code:    "INVALID_INPUT",
			message: key + " must be in YYYY-MM-DD format",
		}
	}
	parsed, err := time.Parse("2006-01-02", val)
	if err != nil {
		return nil, requestError{
			status:  http.StatusBadRequest,
			code:    "INVALID_INPUT",
			message: key + " must be in YYYY-MM-DD format",
		}
	}
	return &parsed, nil
}

func applyOptionalStrings(values map[string][]string, fields map[string]**string) {
	for key, target := range fields {
		raw, ok := values[key]
		if !ok || len(raw) == 0 {
			continue
		}
		val := strings.TrimSpace(raw[0])
		*target = &val
	}
}

func parseOptionalStatus(values map[string][]string, key string) (*user.StatusAkun, error) {
	raw, ok := values[key]
	if !ok || len(raw) == 0 {
		return nil, nil
	}
	val := strings.TrimSpace(raw[0])
	if val == "" {
		return nil, requestError{
			status:  http.StatusBadRequest,
			code:    "INVALID_INPUT",
			message: "status_akun cannot be empty",
		}
	}
	status := user.StatusAkun(val)
	switch status {
	case user.AKTIF, user.NONAKTIF:
		return &status, nil
	default:
		return nil, requestError{
			status:  http.StatusBadRequest,
			code:    "INVALID_INPUT",
			message: "invalid status_akun",
		}
	}
}

func parseOptionalRole(values map[string][]string, key string) (*user.Role, error) {
	raw, ok := values[key]
	if !ok || len(raw) == 0 {
		return nil, nil
	}
	val := strings.TrimSpace(raw[0])
	if val == "" {
		return nil, requestError{
			status:  http.StatusBadRequest,
			code:    "INVALID_INPUT",
			message: "role cannot be empty",
		}
	}
	role := user.Role(val)
	if !role.ValidRole() {
		return nil, requestError{
			status:  http.StatusBadRequest,
			code:    "INVALID_INPUT",
			message: "invalid role",
		}
	}
	return &role, nil
}

func parseOptionalStatusFromString(raw *string) (*user.StatusAkun, error) {
	if raw == nil {
		return nil, nil
	}
	val := strings.TrimSpace(*raw)
	if val == "" {
		return nil, requestError{
			status:  http.StatusBadRequest,
			code:    "INVALID_INPUT",
			message: "status_akun cannot be empty",
		}
	}
	status := user.StatusAkun(val)
	switch status {
	case user.AKTIF, user.NONAKTIF:
		return &status, nil
	default:
		return nil, requestError{
			status:  http.StatusBadRequest,
			code:    "INVALID_INPUT",
			message: "invalid status_akun",
		}
	}
}

func parseOptionalRoleFromString(raw *string) (*user.Role, error) {
	if raw == nil {
		return nil, nil
	}
	val := strings.TrimSpace(*raw)
	if val == "" {
		return nil, requestError{
			status:  http.StatusBadRequest,
			code:    "INVALID_INPUT",
			message: "role cannot be empty",
		}
	}
	role := user.Role(val)
	if !role.ValidRole() {
		return nil, requestError{
			status:  http.StatusBadRequest,
			code:    "INVALID_INPUT",
			message: "invalid role",
		}
	}
	return &role, nil
}
