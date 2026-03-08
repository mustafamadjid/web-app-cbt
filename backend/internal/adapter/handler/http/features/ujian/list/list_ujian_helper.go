package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

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

func toListUjianResponses(items []ujian.ListUjian) []ListUjianResponse {
	response := make([]ListUjianResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toListUjianResponse(item))
	}

	return response
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
		TglUjian:         httphelper.FormatTanggalIndonesia(item.TanggalUjian),
		TanggalUjian:     httphelper.FormatDateOnly(item.TanggalUjian),
		WaktuMulai:       httphelper.FormatTimeOnly(item.WaktuMulai),
		WaktuSelesai:     httphelper.FormatTimeOnly(item.WaktuSelesai),
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
