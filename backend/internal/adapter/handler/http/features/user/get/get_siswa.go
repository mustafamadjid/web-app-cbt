package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/user"
	user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/user/get"
)

type GetSiswaHandler struct {
	svc *user_service.GetSiswaService
}

func NewGetSiswaHandler(svc *user_service.GetSiswaService) *GetSiswaHandler {
	return &GetSiswaHandler{svc: svc}
}

func (h *GetSiswaHandler) ListSiswa(write http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(req.Context())
	if req.Method != http.MethodGet {
		httpResponse.WriteErr(write, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	filters, err := parseListSiswaFilters(req)
	if err != nil {
		logger.Info(req.Context(), "invalid siswa filters", "layer", "adapter.http.handler", "op", "user.list_siswa", "err", err)
		httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	items, err := h.svc.ListSiswa(req.Context(), filters)
	if err != nil {
		logger.Error(req.Context(), "failed listing siswa", "layer", "adapter.http.handler", "op", "user.list_siswa", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(write, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		default:
			httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get list siswa")
		}
		return
	}

	responseData := make([]SiswaResponseItem, 0, len(items))
	for _, item := range items {
		var tanggalLahir string
		if !item.TanggalLahir.IsZero() {
			tanggalLahir = item.TanggalLahir.Format("2006-01-02")
		}
		responseData = append(responseData, SiswaResponseItem{
			IdPengguna:   item.IdPengguna,
			Role:         user.SISWA,
			Username:     item.Username,
			NoHp:         item.NoHp,
			Email:        item.Email,
			NamaLengkap:  item.NamaLengkap,
			StatusAkun:   item.StatusAkun,
			NoAbsen:      item.NoAbsen,
			Angkatan:     item.Angkatan,
			TempatLahir:  item.TempatLahir,
			TanggalLahir: tanggalLahir,
			NamaKelas:    item.NamaKelas,
			TingkatKelas: item.TingkatKelas,
			Foto:         item.Foto,
			JenisKelamin: item.JenisKelamin,
			Nisn:         item.Nisn,
		})
	}

	httpResponse.WriteOK(write, http.StatusOK, responseData, "Success")
}

func (h *GetSiswaHandler) GetSiswaByID(write http.ResponseWriter, req *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(req.Context())
	if req.Method != http.MethodGet {
		httpResponse.WriteErr(write, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	rawID := strings.TrimSpace(params.ByName("id"))
	if rawID == "" {
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : id is required")
		return
	}

	parsedID, err := strconv.Atoi(rawID)
	if err != nil || parsedID <= 0 {
		logger.Info(req.Context(), "invalid siswa id", "layer", "adapter.http.handler", "op", "user.get_siswa", "err", err)
		httpResponse.WriteErr(write, http.StatusBadRequest, "BAD_REQUEST", "Bad request : invalid id")
		return
	}

	result, err := h.svc.FindProfilSiswaByID(req.Context(), user.ID(parsedID))
	if err != nil {
		logger.Error(req.Context(), "failed getting siswa", "layer", "adapter.http.handler", "op", "user.get_siswa", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(write, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(write, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get siswa")
		}
		return
	}

	var tanggalLahir string
	if !result.TanggalLahir.IsZero() {
		tanggalLahir = result.TanggalLahir.Format("2006-01-02")
	}

	responseData := SiswaResponseItem{
		IdPengguna:   result.IdPengguna,
		Role:         result.Role,
		Username:     result.Username,
		NoHp:         result.NoHp,
		Email:        user.Email(result.Email),
		NamaLengkap:  result.NamaLengkap,
		StatusAkun:   result.StatusAkun,
		Nisn:         result.Nisn,
		NoAbsen:      result.NoAbsen,
		Angkatan:     result.Angkatan,
		TempatLahir:  result.TempatLahir,
		TanggalLahir: tanggalLahir,
		NamaKelas:    result.NamaKelas,
		TingkatKelas: result.TingkatKelas,
		Foto:         result.Foto,
		JenisKelamin: result.JenisKelamin,
	}

	httpResponse.WriteOK(write, http.StatusOK, responseData, "Success")
}

func parseListSiswaFilters(req *http.Request) (query.ListSiswaFilter, error) {
	values := req.URL.Query()
	filters := query.ListSiswaFilter{}

	filters.Search = strings.TrimSpace(values.Get("q"))
	if filters.Search == "" {
		filters.Search = strings.TrimSpace(values.Get("search"))
	}
	if err := validator.ValidateInputSafe(filters.Search, "search"); err != nil {
		return query.ListSiswaFilter{}, err
	}

	if statusRaw := strings.TrimSpace(values.Get("status")); statusRaw != "" {
		if err := validator.ValidateInputSafe(statusRaw, "status"); err != nil {
			return query.ListSiswaFilter{}, err
		}
		status := user.StatusAkun(strings.ToUpper(statusRaw))
		filters.Status = &status
	}

	if limitRaw := strings.TrimSpace(values.Get("limit")); limitRaw != "" {
		limit, err := strconv.Atoi(limitRaw)
		if err != nil {
			return query.ListSiswaFilter{}, errors.New("limit must be a number")
		}
		filters.Limit = limit
	}

	if offsetRaw := strings.TrimSpace(values.Get("offset")); offsetRaw != "" {
		offset, err := strconv.Atoi(offsetRaw)
		if err != nil {
			return query.ListSiswaFilter{}, errors.New("offset must be a number")
		}
		filters.Offset = offset
	}

	filters.SortBy = strings.TrimSpace(values.Get("sort_by"))
	if err := validator.ValidateInputSafe(filters.SortBy, "sort_by"); err != nil {
		return query.ListSiswaFilter{}, err
	}

	if sortDescRaw := strings.TrimSpace(values.Get("sort_desc")); sortDescRaw != "" {
		sortDesc, err := strconv.ParseBool(sortDescRaw)
		if err != nil {
			return query.ListSiswaFilter{}, errors.New("sort_desc must be a boolean")
		}
		filters.SortDesc = sortDesc
	}

	if angkatanRaw := strings.TrimSpace(values.Get("angkatan")); angkatanRaw != "" {
		angkatan, err := strconv.Atoi(angkatanRaw)
		if err != nil {
			return query.ListSiswaFilter{}, errors.New("angkatan must be a number")
		}
		filters.Angkatan = &angkatan
	}

	if tingkatKelasRaw := strings.TrimSpace(values.Get("tingkat_kelas")); tingkatKelasRaw != "" {
		tingkatKelas, err := strconv.Atoi(tingkatKelasRaw)
		if err != nil {
			return query.ListSiswaFilter{}, errors.New("tingkat_kelas must be a number")
		}
		filters.TingkatKelas = &tingkatKelas
	}

	if jenisKelaminRaw := strings.TrimSpace(values.Get("jenis_kelamin")); jenisKelaminRaw != "" {
		jenisKelamin, err := strconv.Atoi(jenisKelaminRaw)
		if err != nil {
			return query.ListSiswaFilter{}, errors.New("jenis_kelamin must be a number")
		}
		filters.JenisKelamin = &jenisKelamin
	}

	if idTingkatKelasRaw := strings.TrimSpace(values.Get("id_tingkat_kelas")); idTingkatKelasRaw != "" {
		idTingkatKelas, err := strconv.Atoi(idTingkatKelasRaw)
		if err != nil {
			return query.ListSiswaFilter{}, errors.New("id_tingkat_kelas must be a number")
		}
		filters.IdTingkatKelas = &idTingkatKelas
	}

	if idNamaKelasRaw := strings.TrimSpace(values.Get("id_nama_kelas")); idNamaKelasRaw != "" {
		idNamaKelas, err := strconv.Atoi(idNamaKelasRaw)
		if err != nil {
			return query.ListSiswaFilter{}, errors.New("id_nama_kelas must be a number")
		}
		filters.IdNamaKelas = &idNamaKelas
	}

	return filters, nil
}
