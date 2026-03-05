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
	query "github.com/mustafamadjid/web-app-cbt/internal/core/query/ujian"
	ujian_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/ujian/get"
)

type ListUjianHandler struct {
	svc *ujian_service.GetUjianService
}

func NewListUjianHandler(svc *ujian_service.GetUjianService) *ListUjianHandler {
	return &ListUjianHandler{svc: svc}
}

func (h *ListUjianHandler) ListUjian(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	logger := corelog.FromContext(r.Context())
	if r.Method != http.MethodGet {
		httpResponse.WriteErr(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}

	req, err := parseListUjianRequest(r)
	if err != nil {
		logger.Info(r.Context(), "invalid ujian filters", "layer", "adapter.http.handler", "op", "ujian.list", "err", err)
		httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	filter := query.ListUjianFilter{
		Search:         req.Search,
		Limit:          req.Limit,
		Offset:         req.Offset,
		TanggalUjian:   req.Tanggal,
		Tahun:          req.Tahun,
		TingkatKelasID: req.TingkatKelasID,
		TingkatKelas:   req.TingkatKelas,
		RuangUjian:     req.RuangUjianID,
	}

	items, err := h.svc.GetAllUjianService(r.Context(), filter)
	if err != nil {
		logger.Error(r.Context(), "failed listing ujian", "layer", "adapter.http.handler", "op", "ujian.list", "err", err)
		switch {
		case errors.Is(err, coreerror.ErrInvalidInput):
			httpResponse.WriteErr(w, http.StatusBadRequest, "INVALID_INPUT", "invalid input")
		default:
			httpResponse.WriteErr(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "internal server error: failed get jadwal ujian")
		}
		return
	}

	response := make([]ListUjianResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toListUjianResponse(item))
	}

	httpResponse.WriteOK(w, http.StatusOK, response, "Success")
}

func parseListUjianRequest(r *http.Request) (ListUjianRequest, error) {
	values := r.URL.Query()
	req := ListUjianRequest{
		Search: strings.TrimSpace(values.Get("q")),
	}

	if req.Search == "" {
		req.Search = strings.TrimSpace(values.Get("search"))
	}
	if err := validator.ValidateInputSafe(req.Search, "search"); err != nil {
		return ListUjianRequest{}, err
	}

	if tanggal := strings.TrimSpace(values.Get("tanggal")); tanggal != "" {
		req.Tanggal = &tanggal
	} else if tanggalUjian := strings.TrimSpace(values.Get("tanggal_ujian")); tanggalUjian != "" {
		req.Tanggal = &tanggalUjian
	}

	if tahun := strings.TrimSpace(values.Get("tahun")); tahun != "" {
		req.Tahun = &tahun
	}

	if raw := strings.TrimSpace(values.Get("limit")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return ListUjianRequest{}, errors.New("limit must be a number")
		}
		req.Limit = parsed
	}

	if raw := strings.TrimSpace(values.Get("offset")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return ListUjianRequest{}, errors.New("offset must be a number")
		}
		req.Offset = parsed
	}

	tingkatKelasIDRaw := strings.TrimSpace(values.Get("tingkat_kelas_id"))
	if tingkatKelasIDRaw != "" {
		parsed, err := strconv.Atoi(tingkatKelasIDRaw)
		if err != nil {
			return ListUjianRequest{}, errors.New("tingkat_kelas_id must be a number")
		}
		req.TingkatKelasID = &parsed
	}
	tingkatKelasRaw := strings.TrimSpace(values.Get("tingkat_kelas"))
	if tingkatKelasRaw != "" {
		parsed, err := strconv.Atoi(tingkatKelasRaw)
		if err != nil {
			return ListUjianRequest{}, errors.New("tingkat_kelas must be a number")
		}
		req.TingkatKelas = &parsed
	}

	ruangUjianRaw := strings.TrimSpace(values.Get("ruang_ujian_id"))
	if ruangUjianRaw == "" {
		ruangUjianRaw = strings.TrimSpace(values.Get("ruang_ujian"))
	}
	if ruangUjianRaw != "" {
		parsed, err := strconv.Atoi(ruangUjianRaw)
		if err != nil {
			return ListUjianRequest{}, errors.New("ruang_ujian_id must be a number")
		}
		req.RuangUjianID = &parsed
	}

	return req, nil
}

func toListUjianResponse(item ujian.ListUjian) ListUjianResponse {
	namaKelas := ""
	if item.NamaKelas != nil {
		namaKelas = *item.NamaKelas
	}

	status, started := mapStatusUjian(item.StatusUjian)

	return ListUjianResponse{
		ID:               int(item.IdJadwalUjian),
		IDUjian:          int(item.IdUjian),
		IDGuru:           int(item.IdGuru),
		IDPengawas:       int(item.IdPengawas),
		NamaUjian:        item.NamaUjian,
		PengawasUjian:    item.NamaPengawas,
		TglUjian:         formatTanggalIndonesia(item.TanggalUjian),
		TanggalUjian:     item.TanggalUjian.Format("2006-01-02"),
		WaktuMulai:       item.WaktuMulai.Format("15:04"),
		WaktuSelesai:     item.WaktuSelesai.Format("15:04"),
		SesiUjian:        int(item.IdSesi),
		RuangUjian:       item.NamaRuangan,
		IDRuang:          int(item.IdRuangan),
		StatusUjian:      status,
		Started:          started,
		TingkatKelas:     item.TingkatKelas,
		TingkatKelasID:   int(item.IdKelas),
		NamaKelas:        namaKelas,
		PembuatUsername:  item.PembuatUsername,
		PengawasUsername: item.PengawasUsername,
	}
}

func mapStatusUjian(status ujian.StatusUjian) (string, int) {
	switch status {
	case ujian.BELUM_MULAI:
		return "belum_dimulai", 0
	case ujian.MULAI:
		return "berlangsung", 1
	case ujian.SELESAI:
		return "selesai", 1
	case ujian.DIBATALKAN:
		return "dibatalkan", 0
	default:
		return "belum_dimulai", 0
	}
}

func formatTanggalIndonesia(t time.Time) string {
	hari := map[time.Weekday]string{
		time.Sunday:    "Minggu",
		time.Monday:    "Senin",
		time.Tuesday:   "Selasa",
		time.Wednesday: "Rabu",
		time.Thursday:  "Kamis",
		time.Friday:    "Jumat",
		time.Saturday:  "Sabtu",
	}

	bulan := map[time.Month]string{
		time.January:   "Januari",
		time.February:  "Februari",
		time.March:     "Maret",
		time.April:     "April",
		time.May:       "Mei",
		time.June:      "Juni",
		time.July:      "Juli",
		time.August:    "Agustus",
		time.September: "September",
		time.October:   "Oktober",
		time.November:  "November",
		time.December:  "Desember",
	}

	hariNama := hari[t.Weekday()]
	bulanNama := bulan[t.Month()]
	if hariNama == "" || bulanNama == "" {
		return t.Format("2006-01-02")
	}

	return hariNama + ", " + t.Format("02") + " " + bulanNama + " " + t.Format("2006")
}
