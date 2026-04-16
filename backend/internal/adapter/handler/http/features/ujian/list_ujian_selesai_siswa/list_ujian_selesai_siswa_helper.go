package httpx

import (
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

func toListUjianSelesaiSiswaResponses(items []ujian.ListUjian) []ListUjianSelesaiSiswaResponse {
	response := make([]ListUjianSelesaiSiswaResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toListUjianSelesaiSiswaResponse(item))
	}

	return response
}

func toListUjianSelesaiSiswaResponse(item ujian.ListUjian) ListUjianSelesaiSiswaResponse {
	namaKelas := ""
	if item.NamaKelas != nil {
		namaKelas = *item.NamaKelas
	}

	deskripsiUjian := ""
	if item.DeskripsiUjian != nil {
		deskripsiUjian = *item.DeskripsiUjian
	}

	status, started := mapStatusUjian(item.StatusUjian)

	return ListUjianSelesaiSiswaResponse{
		ID:                  int(item.IdJadwalUjian),
		IDAttempt:           int(item.IdAttempt),
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
