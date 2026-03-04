package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	httpResponse "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper/response_envelope"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/get"
)

type GetUjianHandler struct {
	svc *ujian_service.GetUjianService
}

func NewGetUjianHandler(svc *ujian_service.GetUjianService) *GetUjianHandler {
	return &GetUjianHandler{svc: svc}
}

func (h *GetUjianHandler) GetUjianByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	rawID := strings.TrimSpace(params.ByName("idUjian"))
	idUjian, err := strconv.Atoi(rawID)
	if err != nil || idUjian <= 0 {
		logger.Info(r.Context(), "invalid id ujian", "layer", "adapter.http.handler", "op", "ujian.get_by_id", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id ujian")
		return
	}

	item, err := h.svc.GetUjianByIdService(r.Context(), ujian.ID(idUjian))
	if err != nil {
		logger.Error(r.Context(), "failed get ujian by id", "layer", "adapter.http.handler", "op", "ujian.get_by_id", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "MISSING_ID", "missing id")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get ujian")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toUjianResponse(item), "Success")
}

func (h *GetUjianHandler) GetJadwalUjianByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	rawID := strings.TrimSpace(params.ByName("idJadwalUjian"))
	idJadwalUjian, err := strconv.Atoi(rawID)
	if err != nil || idJadwalUjian <= 0 {
		logger.Info(r.Context(), "invalid id jadwal ujian", "layer", "adapter.http.handler", "op", "ujian.get_jadwal_by_id", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id jadwal ujian")
		return
	}

	item, err := h.svc.GetJadwalUjianByIdService(r.Context(), ujian.ID(idJadwalUjian))
	if err != nil {
		logger.Error(r.Context(), "failed get jadwal ujian by id", "layer", "adapter.http.handler", "op", "ujian.get_jadwal_by_id", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "MISSING_ID", "missing id")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get jadwal ujian")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toJadwalUjianResponse(item), "Success")
}

func (h *GetUjianHandler) GetAllPesertaUjian(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	req, err := parseGetPesertaUjianRequest(r)
	if err != nil {
		logger.Info(r.Context(), "invalid peserta ujian filters", "layer", "adapter.http.handler", "op", "ujian.get_all_peserta.filter", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	items, err := h.svc.GetAllPesertaUjianService(r.Context(), toPesertaUjianFilter(req))
	if err != nil {
		logger.Error(r.Context(), "failed get peserta ujian", "layer", "adapter.http.handler", "op", "ujian.get_all_peserta", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get peserta ujian")
		}
		return
	}

	response := make([]PesertaUjianResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toPesertaUjianResponse(item))
	}

	httpResponse.WriteOK(w, http.StatusOK, response, "Success")
}

func (h *GetUjianHandler) GetPesertaUjianBySiswa(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	rawIDSiswa := strings.TrimSpace(params.ByName("idSiswa"))
	idSiswa, err := strconv.Atoi(rawIDSiswa)
	if err != nil || idSiswa <= 0 {
		logger.Info(r.Context(), "invalid id siswa", "layer", "adapter.http.handler", "op", "ujian.get_peserta_by_siswa", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id siswa")
		return
	}

	req, err := parseGetPesertaUjianRequest(r)
	if err != nil {
		logger.Info(r.Context(), "invalid peserta ujian filters", "layer", "adapter.http.handler", "op", "ujian.get_peserta_by_siswa.filter", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	item, err := h.svc.GetPesertaUjianBySiswaService(r.Context(), ujian.ID(idSiswa), toPesertaUjianFilter(req))
	if err != nil {
		logger.Error(r.Context(), "failed get peserta ujian by siswa", "layer", "adapter.http.handler", "op", "ujian.get_peserta_by_siswa", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "MISSING_ID", "missing id")
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get peserta ujian")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toPesertaUjianResponse(item), "Success")
}

func (h *GetUjianHandler) GetAllJawabanUjianSiswa(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	req, err := parseGetJawabanUjianRequest(r)
	if err != nil {
		logger.Info(r.Context(), "invalid jawaban ujian filters", "layer", "adapter.http.handler", "op", "ujian.get_all_jawaban.filter", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	items, err := h.svc.GetAllJawabanUjianSiswaService(r.Context(), toJawabanUjianFilter(req))
	if err != nil {
		logger.Error(r.Context(), "failed get jawaban ujian siswa", "layer", "adapter.http.handler", "op", "ujian.get_all_jawaban", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get jawaban ujian siswa")
		}
		return
	}

	response := make([]JawabanUjianSiswaResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toJawabanUjianSiswaResponse(item))
	}

	httpResponse.WriteOK(w, http.StatusOK, response, "Success")
}

func (h *GetUjianHandler) GetJawabanBySiswa(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	rawIDSiswa := strings.TrimSpace(params.ByName("idSiswa"))
	idSiswa, err := strconv.Atoi(rawIDSiswa)
	if err != nil || idSiswa <= 0 {
		logger.Info(r.Context(), "invalid id siswa", "layer", "adapter.http.handler", "op", "ujian.get_jawaban_by_siswa", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "BAD_REQUEST", "bad request: invalid id siswa")
		return
	}

	req, err := parseGetJawabanUjianRequest(r)
	if err != nil {
		logger.Info(r.Context(), "invalid jawaban ujian filters", "layer", "adapter.http.handler", "op", "ujian.get_jawaban_by_siswa.filter", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	item, err := h.svc.GetJawabanBySiswaService(r.Context(), ujian.ID(idSiswa), toJawabanUjianFilter(req))
	if err != nil {
		logger.Error(r.Context(), "failed get jawaban by siswa", "layer", "adapter.http.handler", "op", "ujian.get_jawaban_by_siswa", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrMissingId):
			httpResponse.WriteErr(w, http.StatusBadRequest, "MISSING_ID", "missing id")
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		case errors.Is(err, coreerror.ErrNotFound):
			httpResponse.WriteErr(w, http.StatusNotFound, "NOT_FOUND", "data not found")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get jawaban ujian siswa")
		}
		return
	}

	httpResponse.WriteOK(w, http.StatusOK, toJawabanUjianSiswaResponse(item), "Success")
}

func parseGetPesertaUjianRequest(r *http.Request) (GetPesertaUjianRequest, error) {
	values := r.URL.Query()
	req := GetPesertaUjianRequest{}

	if raw := strings.TrimSpace(values.Get("id_peserta_ujian")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return GetPesertaUjianRequest{}, errors.New("id_peserta_ujian must be a number")
		}
		req.IDPesertaUjian = &parsed
	}

	if raw := strings.TrimSpace(values.Get("id_jadwal_ujian")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return GetPesertaUjianRequest{}, errors.New("id_jadwal_ujian must be a number")
		}
		req.IDJadwalUjian = &parsed
	}

	if raw := strings.TrimSpace(values.Get("id_siswa")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return GetPesertaUjianRequest{}, errors.New("id_siswa must be a number")
		}
		req.IDSiswa = &parsed
	}

	if raw := strings.TrimSpace(values.Get("waktu_mulai")); raw != "" {
		parsed, err := parseDateTimeValue(raw)
		if err != nil {
			return GetPesertaUjianRequest{}, errors.New("waktu_mulai must be a valid datetime")
		}
		req.WaktuMulai = &parsed
	}

	if raw := strings.TrimSpace(values.Get("waktu_submit")); raw != "" {
		parsed, err := parseDateTimeValue(raw)
		if err != nil {
			return GetPesertaUjianRequest{}, errors.New("waktu_submit must be a valid datetime")
		}
		req.WaktuSubmit = &parsed
	}

	if raw := strings.TrimSpace(values.Get("nilai_ujian")); raw != "" {
		parsed, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return GetPesertaUjianRequest{}, errors.New("nilai_ujian must be a number")
		}
		req.NilaiUjian = &parsed
	}

	return req, nil
}

func parseGetJawabanUjianRequest(r *http.Request) (GetJawabanUjianRequest, error) {
	values := r.URL.Query()
	req := GetJawabanUjianRequest{}

	if raw := strings.TrimSpace(values.Get("id_jawaban")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return GetJawabanUjianRequest{}, errors.New("id_jawaban must be a number")
		}
		req.IDJawaban = &parsed
	}

	if raw := strings.TrimSpace(values.Get("id_peserta_ujian")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return GetJawabanUjianRequest{}, errors.New("id_peserta_ujian must be a number")
		}
		req.IDPesertaUjian = &parsed
	}

	if raw := strings.TrimSpace(values.Get("id_soal")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return GetJawabanUjianRequest{}, errors.New("id_soal must be a number")
		}
		req.IDSoal = &parsed
	}

	if raw := strings.TrimSpace(values.Get("id_pilihan")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return GetJawabanUjianRequest{}, errors.New("id_pilihan must be a number")
		}
		req.IDPilihan = &parsed
	}

	if raw := strings.TrimSpace(values.Get("jawaban_essay")); raw != "" {
		if err := validator.ValidateInputSafe(raw, "jawaban_essay"); err != nil {
			return GetJawabanUjianRequest{}, err
		}
		req.JawabanEssay = &raw
	}

	if raw := strings.TrimSpace(values.Get("is_benar")); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			return GetJawabanUjianRequest{}, errors.New("is_benar must be a boolean")
		}
		req.IsBenar = &parsed
	}

	if raw := strings.TrimSpace(values.Get("waktu_jawab")); raw != "" {
		parsed, err := parseDateTimeValue(raw)
		if err != nil {
			return GetJawabanUjianRequest{}, errors.New("waktu_jawab must be a valid datetime")
		}
		req.WaktuJawab = &parsed
	}

	return req, nil
}

func toPesertaUjianFilter(req GetPesertaUjianRequest) ujian.PesertaUjian {
	filter := ujian.PesertaUjian{
		WaktuMulai:  req.WaktuMulai,
		WaktuSubmit: req.WaktuSubmit,
		NilaiUjian:  req.NilaiUjian,
	}

	if req.IDPesertaUjian != nil {
		filter.IdPesertaUjian = ujian.ID(*req.IDPesertaUjian)
	}
	if req.IDJadwalUjian != nil {
		filter.IdJadwalUjian = ujian.ID(*req.IDJadwalUjian)
	}
	if req.IDSiswa != nil {
		filter.IdSiswa = ujian.ID(*req.IDSiswa)
	}

	return filter
}

func toJawabanUjianFilter(req GetJawabanUjianRequest) ujian.JawabanUjianSiswa {
	filter := ujian.JawabanUjianSiswa{
		JawabanEssay: req.JawabanEssay,
		IsBenar:      req.IsBenar,
		WaktuJawab:   req.WaktuJawab,
	}

	if req.IDJawaban != nil {
		filter.IdJawaban = ujian.ID(*req.IDJawaban)
	}
	if req.IDPesertaUjian != nil {
		filter.IdPesertaUjian = ujian.ID(*req.IDPesertaUjian)
	}
	if req.IDSoal != nil {
		filter.IdSoal = ujian.ID(*req.IDSoal)
	}
	if req.IDPilihan != nil {
		idPilihan := ujian.ID(*req.IDPilihan)
		filter.IdPilihan = &idPilihan
	}

	return filter
}

func toUjianResponse(item ujian.Ujian) UjianResponse {
	return UjianResponse{
		IDUjian:        int(item.IdUjian),
		IDBankSoal:     int(item.IdBankSoal),
		IDKelas:        int(item.IdKelas),
		IDNamaKelas:    ujianIDToIntPtr(item.IdNamaKelas),
		IDGuru:         int(item.IdGuru),
		NamaUjian:      item.NamaUjian,
		DeskripsiUjian: item.DeskripsiUjian,
		AcakSoal:       item.AcakSoal,
		CreatedAt:      item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      timeToRFC3339Ptr(item.UpdatedAt),
	}
}

func toJadwalUjianResponse(item ujian.JadwalUjian) JadwalUjianResponse {
	updatedAt := timeToRFC3339Ptr(item.UpdatedAt)

	return JadwalUjianResponse{
		IDJadwalUjian: int(item.IdJadwalUjian),
		IDUjian:       int(item.IdUjian),
		IDSesi:        int(item.IdSesi),
		IDRuangan:     int(item.IdRuangan),
		IDPengawas:    int(item.IdPengawas),
		Token:         item.Token,
		TanggalUjian:  item.TanggalUjian.Format("2006-01-02"),
		WaktuMulai:    item.WaktuMulai.Format("15:04"),
		WaktuSelesai:  item.WaktuSelesai.Format("15:04"),
		StatusUjian:   mapStatusUjianToClient(item.StatusUjian),
		CreatedAt:     item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     updatedAt,
	}
}

func toPesertaUjianResponse(item ujian.PesertaUjian) PesertaUjianResponse {
	return PesertaUjianResponse{
		IDPesertaUjian: int(item.IdPesertaUjian),
		IDJadwalUjian:  int(item.IdJadwalUjian),
		IDSiswa:        int(item.IdSiswa),
		WaktuMulai:     timeToRFC3339Ptr(item.WaktuMulai),
		WaktuSubmit:    timeToRFC3339Ptr(item.WaktuSubmit),
		NilaiUjian:     item.NilaiUjian,
		CreatedAt:      item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      timeToRFC3339Ptr(item.UpdatedAt),
	}
}

func toJawabanUjianSiswaResponse(item ujian.JawabanUjianSiswa) JawabanUjianSiswaResponse {
	return JawabanUjianSiswaResponse{
		IDJawaban:      int(item.IdJawaban),
		IDPesertaUjian: int(item.IdPesertaUjian),
		IDSoal:         int(item.IdSoal),
		IDPilihan:      ujianIDToIntPtr(item.IdPilihan),
		JawabanEssay:   item.JawabanEssay,
		IsBenar:        item.IsBenar,
		WaktuJawab:     timeToRFC3339Ptr(item.WaktuJawab),
	}
}

func parseDateTimeValue(raw string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
	}

	for _, layout := range layouts {
		parsed, err := time.Parse(layout, raw)
		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, errors.New("invalid datetime")
}

func mapStatusUjianToClient(status ujian.StatusUjian) string {
	switch status {
	case ujian.BELUM_MULAI:
		return "belum_dimulai"
	case ujian.MULAI:
		return "berlangsung"
	case ujian.SELESAI:
		return "selesai"
	case ujian.DIBATALKAN:
		return "dibatalkan"
	default:
		return "belum_dimulai"
	}
}

func timeToRFC3339Ptr(value *time.Time) *string {
	if value == nil {
		return nil
	}

	formatted := value.Format(time.RFC3339)
	return &formatted
}

func ujianIDToIntPtr(value *ujian.ID) *int {
	if value == nil {
		return nil
	}

	parsed := int(*value)
	return &parsed
}
