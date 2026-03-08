package httpx

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	validator "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/validation"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

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
		parsed, err := ParseDateTimeValue(raw)
		if err != nil {
			return GetPesertaUjianRequest{}, errors.New("waktu_mulai must be a valid datetime")
		}
		req.WaktuMulai = &parsed
	}

	if raw := strings.TrimSpace(values.Get("waktu_submit")); raw != "" {
		parsed, err := ParseDateTimeValue(raw)
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
		parsed, err := ParseDateTimeValue(raw)
		if err != nil {
			return GetJawabanUjianRequest{}, errors.New("waktu_jawab must be a valid datetime")
		}
		req.WaktuJawab = &parsed
	}

	return req, nil
}

func ToPesertaUjianFilter(req GetPesertaUjianRequest) ujian.PesertaUjian {
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

func ToJawabanUjianFilter(req GetJawabanUjianRequest) ujian.JawabanUjianSiswa {
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

func ParseDateTimeValue(raw string) (time.Time, error) {
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

func UjianIDToIntPtr(value *ujian.ID) *int {
	if value == nil {
		return nil
	}

	parsed := int(*value)
	return &parsed
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

func toPesertaUjianResponses(items []ujian.PesertaUjian) []PesertaUjianResponse {
	response := make([]PesertaUjianResponse, 0, len(items))
	for _, item := range items {
		response = append(response, ToPesertaUjianResponse(item))
	}

	return response
}

func ToPesertaUjianResponse(item ujian.PesertaUjian) PesertaUjianResponse {
	return PesertaUjianResponse{
		IDPesertaUjian: int(item.IdPesertaUjian),
		IDJadwalUjian:  int(item.IdJadwalUjian),
		IDSiswa:        int(item.IdSiswa),
		WaktuMulai:     httphelper.FormatRFC3339Ptr(item.WaktuMulai),
		WaktuSubmit:    httphelper.FormatRFC3339Ptr(item.WaktuSubmit),
		NilaiUjian:     item.NilaiUjian,
		CreatedAt:      httphelper.FormatRFC3339(item.CreatedAt),
		UpdatedAt:      httphelper.FormatRFC3339Ptr(item.UpdatedAt),
	}
}

func toJawabanUjianSiswaResponses(items []ujian.JawabanUjianSiswa) []JawabanUjianSiswaResponse {
	response := make([]JawabanUjianSiswaResponse, 0, len(items))
	for _, item := range items {
		response = append(response, ToJawabanUjianSiswaResponse(item))
	}

	return response
}

func ToJawabanUjianSiswaResponse(item ujian.JawabanUjianSiswa) JawabanUjianSiswaResponse {
	return JawabanUjianSiswaResponse{
		IDJawaban:      int(item.IdJawaban),
		IDPesertaUjian: int(item.IdPesertaUjian),
		IDSoal:         int(item.IdSoal),
		IDPilihan:      UjianIDToIntPtr(item.IdPilihan),
		JawabanEssay:   item.JawabanEssay,
		IsBenar:        item.IsBenar,
		WaktuJawab:     httphelper.FormatRFC3339Ptr(item.WaktuJawab),
	}
}

func ToUjianResponse(item ujian.ListUjian) ListUjianByIdResponse {
	status, started := mapStatusUjian(item.StatusUjian)

	namaKelas := ""
	if item.NamaKelas != nil {
		namaKelas = *item.NamaKelas
	}

	deskripsiUjian := ""
	if item.DeskripsiUjian != nil {
		deskripsiUjian = *item.DeskripsiUjian
	}

	return ListUjianByIdResponse{
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
		DeskripsiUjian:   deskripsiUjian,
		Token:            item.Token,
	}
}
