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

func parseListUjianEssayUngradedRequest(r *http.Request) (ListUjianEssayUngradedRequest, error) {
	values := r.URL.Query()
	req := ListUjianEssayUngradedRequest{
		Search: strings.TrimSpace(values.Get("q")),
	}

	if req.Search == "" {
		req.Search = strings.TrimSpace(values.Get("search"))
	}

	search, err := validator.ValidateOptionalPrintableText(req.Search, "search")
	if err != nil {
		return ListUjianEssayUngradedRequest{}, err
	}
	req.Search = search

	if tanggal := strings.TrimSpace(values.Get("tanggal")); tanggal != "" {
		req.TanggalUjian = &tanggal
	} else if tanggalUjian := strings.TrimSpace(values.Get("tanggal_ujian")); tanggalUjian != "" {
		req.TanggalUjian = &tanggalUjian
	}

	if tahun := strings.TrimSpace(values.Get("tahun")); tahun != "" {
		req.Tahun = &tahun
	}

	if bulan := strings.TrimSpace(values.Get("bulan")); bulan != "" {
		req.Bulan = &bulan
	}

	if raw := strings.TrimSpace(values.Get("limit")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return ListUjianEssayUngradedRequest{}, errors.New("limit must be a number")
		}
		req.Limit = parsed
	}

	if raw := strings.TrimSpace(values.Get("offset")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return ListUjianEssayUngradedRequest{}, errors.New("offset must be a number")
		}
		req.Offset = parsed
	}

	if raw := strings.TrimSpace(values.Get("tingkat_kelas_id")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return ListUjianEssayUngradedRequest{}, errors.New("tingkat_kelas_id must be a number")
		}
		req.TingkatKelasID = &parsed
	}

	if raw := strings.TrimSpace(values.Get("nama_kelas_id")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return ListUjianEssayUngradedRequest{}, errors.New("nama_kelas_id must be a number")
		}
		req.NamaKelasID = &parsed
	}

	mapelIDRaw := strings.TrimSpace(values.Get("id_mapel"))
	if mapelIDRaw == "" {
		mapelIDRaw = strings.TrimSpace(values.Get("mapel"))
	}
	if mapelIDRaw != "" {
		parsed, err := strconv.Atoi(mapelIDRaw)
		if err != nil {
			return ListUjianEssayUngradedRequest{}, errors.New("id_mapel must be a number")
		}
		req.MapelID = &parsed
	}

	sesiIDRaw := strings.TrimSpace(values.Get("sesi_id"))
	if sesiIDRaw == "" {
		sesiIDRaw = strings.TrimSpace(values.Get("sesi"))
	}
	if sesiIDRaw != "" {
		parsed, err := strconv.Atoi(sesiIDRaw)
		if err != nil {
			return ListUjianEssayUngradedRequest{}, errors.New("sesi_id must be a number")
		}
		req.SesiID = &parsed
	}

	return req, nil
}

func toListUjianEssayUngradedResponses(items []ujian.ListUjian) []ListUjianEssayUngradedResponse {
	response := make([]ListUjianEssayUngradedResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toListUjianEssayUngradedResponse(item))
	}

	return response
}

func toListUjianEssayUngradedResponse(item ujian.ListUjian) ListUjianEssayUngradedResponse {
	namaKelas := ""
	if item.NamaKelas != nil {
		namaKelas = *item.NamaKelas
	}

	status, started := mapStatusUjian(item.StatusUjian)

	return ListUjianEssayUngradedResponse{
		ID:               int(item.IdJadwalUjian),
		IDUjian:          int(item.IdUjian),
		IDBankSoal:       int(item.IdBankSoal),
		IDGuru:           int(item.IdGuru),
		IDPengawas:       int(item.IdPengawas),
		NamaUjian:        item.NamaUjian,
		PengawasUjian:    item.NamaPengawas,
		TglUjian:         httphelper.FormatTanggalIndonesia(item.TanggalUjian),
		TanggalUjian:     httphelper.FormatDateOnly(item.TanggalUjian),
		WaktuMulai:       httphelper.FormatTimeOnly(item.WaktuMulai),
		WaktuSelesai:     httphelper.FormatTimeOnly(item.WaktuSelesai),
		IDSesi:           int(item.IdSesi),
		NamaSesi:         item.NamaSesi,
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

func mapStatusUjian(status *ujian.StatusUjian) (string, int) {
	if status == nil {
		return "belum_dimulai", 0
	}

	switch *status {
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
