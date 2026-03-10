package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/update"
)

type requestError struct {
	status  int
	code    string
	message string
}

func (e requestError) Error() string { return e.message }

func (h *UpdateHandler) parseUpdateGuruMultipart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (user_service.UpdateGuruCmd, error) {
	idPengguna, err := parseIDFromURLParam(ps, "id")
	if err != nil {
		return user_service.UpdateGuruCmd{}, err
	}

	if err := httphelper.MultipartHeaderBodyValidator(w, r, 10<<20); err != nil {
		if errors.Is(err, coreerror.ErrContentTypeMustMultipart) {
			return user_service.UpdateGuruCmd{}, requestError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "bad request: content type must be multipart/form-data",
			}
		}
		return user_service.UpdateGuruCmd{}, requestError{
			status:  http.StatusBadRequest,
			code:    "BAD_REQUEST",
			message: "bad request: invalid multipart form",
		}
	}

	cmd := user_service.UpdateGuruCmd{IdPengguna: idPengguna}
	if err := applyOptionalStrings(r.MultipartForm.Value, map[string]**string{
		"username":      &cmd.Username,
		"email":         &cmd.Email,
		"nama_lengkap":  &cmd.NamaLengkap,
		"jenis_kelamin": &cmd.JenisKelamin,
		"no_hp":         &cmd.NoHp,
		"nip":           &cmd.Nip,
		"jabatan":       &cmd.Jabatan,
		"bidang_studi":  &cmd.BidangStudi,
	}); err != nil {
		return user_service.UpdateGuruCmd{}, err
	}

	statusAkun, err := parseOptionalStatus(r.MultipartForm.Value, "status_akun")
	if err != nil {
		return user_service.UpdateGuruCmd{}, err
	}
	cmd.StatusAkun = statusAkun

	role, err := parseOptionalRole(r.MultipartForm.Value, "role")
	if err != nil {
		return user_service.UpdateGuruCmd{}, err
	}
	cmd.Role = role

	if relPath, err := h.parseOptionalFoto(r); err != nil {
		return user_service.UpdateGuruCmd{}, err
	} else {
		cmd.Foto = relPath
	}

	return cmd, nil
}

func (h *UpdateHandler) parseUpdateSiswaMultipart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (user_service.UpdateSiswaCmd, error) {
	idPengguna, err := parseIDFromURLParam(ps, "id")
	if err != nil {
		return user_service.UpdateSiswaCmd{}, err
	}

	if err := httphelper.MultipartHeaderBodyValidator(w, r, 10<<20); err != nil {
		if errors.Is(err, coreerror.ErrContentTypeMustMultipart) {
			return user_service.UpdateSiswaCmd{}, requestError{
				status:  http.StatusBadRequest,
				code:    "BAD_REQUEST",
				message: "bad request: content type must be multipart/form-data",
			}
		}
		return user_service.UpdateSiswaCmd{}, requestError{
			status:  http.StatusBadRequest,
			code:    "BAD_REQUEST",
			message: "bad request: invalid multipart form",
		}
	}

	cmd := user_service.UpdateSiswaCmd{IdPengguna: idPengguna}
	if err := applyOptionalStrings(r.MultipartForm.Value, map[string]**string{
		"username":      &cmd.Username,
		"email":         &cmd.Email,
		"nama_lengkap":  &cmd.NamaLengkap,
		"jenis_kelamin": &cmd.JenisKelamin,
		"no_hp":         &cmd.NoHp,
		"nisn":          &cmd.Nisn,
		"tempat_lahir":  &cmd.TempatLahir,
	}); err != nil {
		return user_service.UpdateSiswaCmd{}, err
	}

	statusAkun, err := parseOptionalStatus(r.MultipartForm.Value, "status_akun")
	if err != nil {
		return user_service.UpdateSiswaCmd{}, err
	}
	cmd.StatusAkun = statusAkun

	role, err := parseOptionalRole(r.MultipartForm.Value, "role")
	if err != nil {
		return user_service.UpdateSiswaCmd{}, err
	}
	cmd.Role = role

	if idTingkat, err := parseOptionalID(r.MultipartForm.Value, "id_tingkat_kelas"); err != nil {
		return user_service.UpdateSiswaCmd{}, err
	} else {
		cmd.IdTingkatKelas = idTingkat
	}

	if idNama, err := parseOptionalID(r.MultipartForm.Value, "id_nama_kelas"); err != nil {
		return user_service.UpdateSiswaCmd{}, err
	} else {
		cmd.IdNamaKelas = idNama
	}

	if noAbsen, err := parseOptionalInt(r.MultipartForm.Value, "no_absen"); err != nil {
		return user_service.UpdateSiswaCmd{}, err
	} else {
		cmd.NoAbsen = noAbsen
	}

	if angkatan, err := parseOptionalInt(r.MultipartForm.Value, "angkatan"); err != nil {
		return user_service.UpdateSiswaCmd{}, err
	} else {
		cmd.Angkatan = angkatan
	}

	if tanggal, err := parseOptionalDate(r.MultipartForm.Value, "tanggal_lahir"); err != nil {
		return user_service.UpdateSiswaCmd{}, err
	} else {
		cmd.TanggalLahir = tanggal
	}

	if relPath, err := h.parseOptionalFoto(r); err != nil {
		return user_service.UpdateSiswaCmd{}, err
	} else {
		cmd.Foto = relPath
	}

	return cmd, nil
}

func parseIDFromURLParam(ps httprouter.Params, key string) (user.ID, error) {
	raw := strings.TrimSpace(ps.ByName(key))
	if raw == "" {
		return 0, requestError{
			status:  http.StatusBadRequest,
			code:    "BAD_REQUEST",
			message: "Bad request : " + key + " is required",
		}
	}
	idInt, err := strconv.Atoi(raw)
	if err != nil {
		return 0, requestError{
			status:  http.StatusBadRequest,
			code:    "INVALID_INPUT",
			message: key + " must be a number",
		}
	}
	if idInt <= 0 {
		return 0, requestError{
			status:  http.StatusBadRequest,
			code:    "INVALID_INPUT",
			message: key + " must be a positive number",
		}
	}
	return user.ID(idInt), nil
}

func (h *UpdateHandler) parseOptionalFoto(r *http.Request) (*string, error) {
	relPath, err := httphelper.StoreFileToDisk(r, "foto", false, h.storeImage.SavePhotoRelative)
	if err != nil {
		if errors.Is(err, coreerror.ErrFileTooLarge) {
			return nil, requestError{
				status:  http.StatusBadRequest,
				code:    "FILE_TOO_LARGE",
				message: "file too large",
			}
		}

		return nil, requestError{
			status:  http.StatusBadRequest,
			code:    "BAD_REQUEST",
			message: "Bad request : failed reading foto",
		}
	}

	return relPath, nil
}

func (h *UpdateHandler) writeRequestError(w http.ResponseWriter, err error) {
	if reqErr, ok := err.(requestError); ok {
		httpResponse.WriteErr(w, reqErr.status, reqErr.code, reqErr.message)
		return
	}
	httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "Bad request : invalid request body")
}

func (h *UpdateHandler) writeUpdateError(w http.ResponseWriter, err error, role string) {
	switch {
	case errors.Is(err, coreerror.ErrForbidden):
		httpResponse.WriteErr(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
	case errors.Is(err, coreerror.ErrNoFieldToUpdate):
		httpResponse.WriteErr(w, http.StatusBadRequest, "NO_FIELD_TO_UPDATE", "no field to update")
	case errors.Is(err, coreerror.ErrUsernameLengthInvalid):
		httpResponse.WriteErr(w, http.StatusBadRequest, "USERNAME_LENGTH_INVALID", "username length invalid")
	case errors.Is(err, coreerror.ErrUsernameTaken):
		httpResponse.WriteErr(w, http.StatusConflict, "USERNAME_TAKEN", "username sudah terdaftar. data yang diinputkan harus unik")
	case errors.Is(err, coreerror.ErrEmailTaken):
		httpResponse.WriteErr(w, http.StatusConflict, "EMAIL_TAKEN", "email sudah terdaftar. data yang diinputkan harus unik")
	case errors.Is(err, coreerror.ErrNoHpTaken):
		httpResponse.WriteErr(w, http.StatusConflict, "NO_HP_TAKEN", "nomor HP sudah terdaftar. data yang diinputkan harus unik")
	case errors.Is(err, coreerror.ErrNipTaken):
		httpResponse.WriteErr(w, http.StatusConflict, "NIP_TAKEN", "NIP sudah terdaftar. data yang diinputkan harus unik")
	case errors.Is(err, coreerror.ErrNisnTaken):
		httpResponse.WriteErr(w, http.StatusConflict, "NISN_TAKEN", "NISN sudah terdaftar. data yang diinputkan harus unik")
	case errors.Is(err, coreerror.ErrConflict):
		httpResponse.WriteErr(w, http.StatusConflict, "CONFLICT", "data yang diinputkan sudah ada sebelumnya, harus unik")
	case errors.Is(err, coreerror.ErrInvalidInput),
		errors.Is(err, coreerror.ErrInvalidStatusAkun),
		errors.Is(err, user.ErrInvalidEmail),
		errors.Is(err, user.ErrInvalidNISN),
		errors.Is(err, user.ErrInvalidAbsen),
		errors.Is(err, user.ErrInvalidAngkatan):
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
	case errors.Is(err, coreerror.ErrNotFound):
		httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
	default:
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed update "+role)
	}
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
	if err != nil || idInt <= 0 {
		return nil, requestError{
			status:  http.StatusBadRequest,
			code:    "INVALID_INPUT",
			message: key + " must be a positive number",
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

func applyOptionalStrings(values map[string][]string, fields map[string]**string) error {
	for key, target := range fields {
		raw, ok := values[key]
		if !ok || len(raw) == 0 {
			continue
		}
		val := raw[0]
		var err error
		switch key {
		case "username":
			val, err = validator.ValidateUsername(val)
		case "email":
			val, err = validator.ValidateEmailAddress(val, key)
		default:
			val, err = validator.ValidateRequiredPrintableText(val, key)
		}
		if err != nil {
			return requestError{
				status:  http.StatusBadRequest,
				code:    "INVALID_INPUT",
				message: err.Error(),
			}
		}

		*target = &val
	}
	return nil
}

func parseOptionalStatus(values map[string][]string, key string) (*user.StatusAkun, error) {
	raw, ok := values[key]
	if !ok || len(raw) == 0 {
		return nil, nil
	}
	val := raw[0]
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
	val := raw[0]
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
