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

func parseListUjianSiswaRequest(r *http.Request) (ListUjianSiswaRequest, error) {
	values := r.URL.Query()
	req := ListUjianSiswaRequest{
		Search: strings.TrimSpace(values.Get("q")),
	}

	if req.Search == "" {
		req.Search = strings.TrimSpace(values.Get("search"))
	}
	search, err := validator.ValidateOptionalPrintableText(req.Search, "search")
	if err != nil {
		return ListUjianSiswaRequest{}, err
	}
	req.Search = search

	if tanggal := strings.TrimSpace(values.Get("tanggal")); tanggal != "" {
		req.Tanggal = &tanggal
	} else if tanggalUjian := strings.TrimSpace(values.Get("tanggal_ujian")); tanggalUjian != "" {
		req.Tanggal = &tanggalUjian
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
			return ListUjianSiswaRequest{}, errors.New("limit must be a number")
		}
		req.Limit = parsed
	}

	if raw := strings.TrimSpace(values.Get("offset")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return ListUjianSiswaRequest{}, errors.New("offset must be a number")
		}
		req.Offset = parsed
	}

	tingkatKelasIDRaw := strings.TrimSpace(values.Get("tingkat_kelas_id"))
	if tingkatKelasIDRaw != "" {
		parsed, err := strconv.Atoi(tingkatKelasIDRaw)
		if err != nil {
			return ListUjianSiswaRequest{}, errors.New("tingkat_kelas_id must be a number")
		}
		req.TingkatKelasID = &parsed
	}

	tingkatKelasRaw := strings.TrimSpace(values.Get("tingkat_kelas"))
	if tingkatKelasRaw != "" {
		parsed, err := strconv.Atoi(tingkatKelasRaw)
		if err != nil {
			return ListUjianSiswaRequest{}, errors.New("tingkat_kelas must be a number")
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
			return ListUjianSiswaRequest{}, errors.New("ruang_ujian_id must be a number")
		}
		req.RuangUjianID = &parsed
	}

	idMapelRaw := strings.TrimSpace(values.Get("id_mapel"))
	if idMapelRaw == "" {
		idMapelRaw = strings.TrimSpace(values.Get("mapel"))
	}
	if idMapelRaw != "" {
		parsed, err := strconv.Atoi(idMapelRaw)
		if err != nil {
			return ListUjianSiswaRequest{}, errors.New("id_mapel must be a number")
		}
		req.IDMapel = &parsed
	}

	req.KategoriUjian = strings.TrimSpace(values.Get("kategori_ujian"))
	if req.KategoriUjian == "" {
		req.KategoriUjian = strings.TrimSpace(values.Get("kategori"))
	}

	return req, nil
}

func toListUjianSiswaResponses(items []ujian.ListUjian) []ListUjianSiswaResponse {
	response := make([]ListUjianSiswaResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toListUjianSiswaResponse(item))
	}

	return response
}

func toListUjianSiswaResponse(item ujian.ListUjian) ListUjianSiswaResponse {
	namaKelas := ""
	if item.NamaKelas != nil {
		namaKelas = *item.NamaKelas
	}

	deskripsiUjian := ""
	if item.DeskripsiUjian != nil {
		deskripsiUjian = *item.DeskripsiUjian
	}

	status, started := mapStatusUjian(item.StatusUjian)

	return ListUjianSiswaResponse{
		ID:                  int(item.IdJadwalUjian),
		IDUjian:             int(item.IdUjian),
		IDBankSoal:          int(item.IdBankSoal),
		IDGuru:              int(item.IdGuru),
		IDPengawas:          int(item.IdPengawas),
		NamaUjian:           item.NamaUjian,
		PengawasUjian:       item.NamaPengawas,
		TglUjian:            httphelper.FormatTanggalIndonesia(item.TanggalUjian),
		TanggalUjian:        httphelper.FormatDateOnly(item.TanggalUjian),
		WaktuMulai:          httphelper.FormatTimeOnly(item.WaktuMulai),
		WaktuSelesai:        httphelper.FormatTimeOnly(item.WaktuSelesai),
		SesiUjian:           int(item.IdSesi),
		NamaSesi:            item.NamaSesi,
		RuangUjian:          item.NamaRuangan,
		IDRuang:             int(item.IdRuangan),
		StatusUjian:         status,
		Started:             started,
		TingkatKelas:        item.TingkatKelas,
		TingkatKelasID:      int(item.IdKelas),
		NamaKelas:           namaKelas,
		PengawasNamaLengkap: item.PengawasNamaLengkap,
		DeskripsiUjian:      deskripsiUjian,
		AcakSoal:            item.AcakSoal,
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
