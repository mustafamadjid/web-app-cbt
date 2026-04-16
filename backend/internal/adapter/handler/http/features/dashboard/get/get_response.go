package httpx

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/dashboard"

type GetDashboardStatistikResponse struct {
	TotalSiswa           int `json:"total_siswa"`
	TotalGuru            int `json:"total_guru"`
	TotalUjianTerlaksana int `json:"total_ujian_terlaksana"`
	TotalBankSoal        int `json:"total_bank_soal"`
	TotalMapelAktif      int `json:"total_mapel_aktif"`
}

func toGetDashboardStatistikResponse(item dashboard.DashboardStatistik) GetDashboardStatistikResponse {
	return GetDashboardStatistikResponse{
		TotalSiswa:           item.TotalSiswa,
		TotalGuru:            item.TotalGuru,
		TotalUjianTerlaksana: item.TotalUjianTerlaksana,
		TotalBankSoal:        item.TotalBankSoal,
		TotalMapelAktif:      item.TotalMapelAktif,
	}
}
