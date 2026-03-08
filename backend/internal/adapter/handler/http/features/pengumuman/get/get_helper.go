package httpx

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"

func toPengumumanResponses(items []pengumuman.Pengumuman) []PengumumanGetResponse {
	response := make([]PengumumanGetResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toPengumumanResponse(item))
	}

	return response
}

func toPengumumanResponse(item pengumuman.Pengumuman) PengumumanGetResponse {
	return PengumumanGetResponse{
		IdPengumuman:             int(item.IdPengumuman),
		IdPengguna:               int(item.IdPengguna),
		JudulPengumuman:          item.JudulPengumuman,
		IsiPengumuman:            item.IsiPengumuman,
		TanggalRilisPengumuman:   item.TanggalRilisPengumuman,
		TanggalSelesaiPengumuman: item.TanggalSelesaiPengumuman,
		DokumenPengumuman:        item.DokumenPengumuman,
	}
}
