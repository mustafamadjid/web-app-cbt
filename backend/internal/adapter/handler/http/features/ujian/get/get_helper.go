package httpx

import (
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	ujian "github.com/mustafamadjid/web-app-cbt/internal/core/domain/ujian_siswa"
)

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
		SesiUjian:        item.NamaSesi,
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
