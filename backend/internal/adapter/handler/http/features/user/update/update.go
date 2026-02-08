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

	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	aktivitas_user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/aktivitas_user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/update"
)

type UpdateHandler struct {
	svc           *user_service.UpdateTx
	storeImage    httphelper.ImageStore
	aktivitasUser *aktivitas_user_service.AktivitasUserService
}

type requestError struct {
	status  int
	code    string
	message string
}

func (e requestError) Error() string { return e.message }

func NewUpdateUserHandler(svc *user_service.UpdateTx, storeImage httphelper.ImageStore, aktivitasUser *aktivitas_user_service.AktivitasUserService) *UpdateHandler {
	return &UpdateHandler{svc: svc, storeImage: storeImage, aktivitasUser: aktivitasUser}
}

func (h *UpdateHandler) UpdateGuru(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodPatch {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	cmd, err := h.parseUpdateGuruMultipart(r, ps)
	if err != nil {
		logger.Error(r.Context(), "failed parsing update guru request", "layer", "adapter.http.handler", "op", "user.update_guru", "err", err)
		h.writeRequestError(w, err)
		return
	}

	actor, ok := middleware.ActorFromContext(r.Context())
	if !ok {
		logger.Error(r.Context(), "missing actor in context", "layer", "adapter.http.handler", "op", "user.update_guru", "err", "actor_not_found")
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	if err := h.svc.UpdateGuru(r.Context(), cmd, actor); err != nil {
		logger.Error(r.Context(), "failed updating guru", "layer", "adapter.http.handler", "op", "user.update_guru", "user_id", cmd.IdPengguna, "err", err)
		h.writeUpdateError(w, err, "guru")
		return
	}

	if h.aktivitasUser != nil {
		aktivitasCmd := aktivitas_user_service.AktivitasUserCmd{
			IdPengguna:  actor.IdPengguna,
			Action:      aktivitas_user.UPDATE,
			Description: "Memperbarui data guru",
			IpAddress:   httphelper.GetClientIP(r),
		}
		if err := h.aktivitasUser.CreateAktivitasUserService(r.Context(), aktivitasCmd); err != nil {
			logger.Error(r.Context(), "failed creating aktivitas user", "layer", "adapter.http.handler", "op", "user.update_guru.activity", "err", err)
		}
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "Success")
}

func (h *UpdateHandler) UpdateSiswa(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodPatch {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	cmd, err := h.parseUpdateSiswaMultipart(r, ps)
	if err != nil {
		logger.Error(r.Context(), "failed parsing update siswa request", "layer", "adapter.http.handler", "op", "user.update_siswa", "err", err)
		h.writeRequestError(w, err)
		return
	}

	actor, ok := middleware.ActorFromContext(r.Context())
	if !ok {
		logger.Error(r.Context(), "missing actor in context", "layer", "adapter.http.handler", "op", "user.update_siswa", "err", "actor_not_found")
		httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error : failed get actor from context")
		return
	}

	if err := h.svc.UpdateSiswa(r.Context(), cmd, actor); err != nil {
		logger.Error(r.Context(), "failed updating siswa", "layer", "adapter.http.handler", "op", "user.update_siswa", "user_id", cmd.IdPengguna, "err", err)
		h.writeUpdateError(w, err, "siswa")
		return
	}

	if h.aktivitasUser != nil {
		aktivitasCmd := aktivitas_user_service.AktivitasUserCmd{
			IdPengguna:  actor.IdPengguna,
			Action:      aktivitas_user.UPDATE,
			Description: "Memperbarui data siswa",
			IpAddress:   httphelper.GetClientIP(r),
		}
		if err := h.aktivitasUser.CreateAktivitasUserService(r.Context(), aktivitasCmd); err != nil {
			logger.Error(r.Context(), "failed creating aktivitas user", "layer", "adapter.http.handler", "op", "user.update_siswa.activity", "err", err)
		}
	}

	httpResponse.WriteOKNoData(w, http.StatusOK, "Success")
}

func (h *UpdateHandler) parseUpdateGuruMultipart(r *http.Request, ps httprouter.Params) (user_service.UpdateGuruCmd, error) {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		return user_service.UpdateGuruCmd{}, requestError{
			status:  http.StatusBadRequest,
			code:    "BAD_REQUEST",
			message: "Bad request : content type must be multipart/form-data",
		}
	}

	idPengguna, err := parseIDFromURLParam(ps, "id")
	if err != nil {
		return user_service.UpdateGuruCmd{}, err
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return user_service.UpdateGuruCmd{}, requestError{
			status:  http.StatusBadRequest,
			code:    "BAD_REQUEST",
			message: "Bad request : invalid multipart form",
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
	} else if relPath != nil {
		cmd.Foto = relPath
	}

	return cmd, nil
}

func (h *UpdateHandler) parseUpdateSiswaMultipart(r *http.Request, ps httprouter.Params) (user_service.UpdateSiswaCmd, error) {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		return user_service.UpdateSiswaCmd{}, requestError{
			status:  http.StatusBadRequest,
			code:    "BAD_REQUEST",
			message: "Bad request : content type must be multipart/form-data",
		}
	}

	idPengguna, err := parseIDFromURLParam(ps, "id")
	if err != nil {
		return user_service.UpdateSiswaCmd{}, err
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return user_service.UpdateSiswaCmd{}, requestError{
			status:  http.StatusBadRequest,
			code:    "BAD_REQUEST",
			message: "Bad request : invalid multipart form",
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
	} else if idTingkat != nil {
		cmd.IdTingkatKelas = idTingkat
	}

	if idNama, err := parseOptionalID(r.MultipartForm.Value, "id_nama_kelas"); err != nil {
		return user_service.UpdateSiswaCmd{}, err
	} else if idNama != nil {
		cmd.IdNamaKelas = idNama
	}

	if noAbsen, err := parseOptionalInt(r.MultipartForm.Value, "no_absen"); err != nil {
		return user_service.UpdateSiswaCmd{}, err
	} else if noAbsen != nil {
		cmd.NoAbsen = noAbsen
	}

	if angkatan, err := parseOptionalInt(r.MultipartForm.Value, "angkatan"); err != nil {
		return user_service.UpdateSiswaCmd{}, err
	} else if angkatan != nil {
		cmd.Angkatan = angkatan
	}

	if tanggal, err := parseOptionalDate(r.MultipartForm.Value, "tanggal_lahir"); err != nil {
		return user_service.UpdateSiswaCmd{}, err
	} else if tanggal != nil {
		cmd.TanggalLahir = tanggal
	}

	if relPath, err := h.parseOptionalFoto(r); err != nil {
		return user_service.UpdateSiswaCmd{}, err
	} else if relPath != nil {
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
	file, fh, err := r.FormFile("foto")
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
	case errors.Is(err, coreerror.ErrUsernameTaken):
		httpResponse.WriteErr(w, http.StatusConflict, "USERNAME_TAKEN", "username already taken")
	case errors.Is(err, coreerror.ErrNipTaken):
		httpResponse.WriteErr(w, http.StatusConflict, "NIP_TAKEN", "nip already taken")
	case errors.Is(err, coreerror.ErrNisnTaken):
		httpResponse.WriteErr(w, http.StatusConflict, "NISN_TAKEN", "nisn already taken")
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
		val := strings.TrimSpace(raw[0])
		if err := validator.ValidateInputSafe(val, key); err != nil {
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
