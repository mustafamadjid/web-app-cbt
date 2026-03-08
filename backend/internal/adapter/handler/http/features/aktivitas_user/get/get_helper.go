package httpx

import (
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
)

func toAktivitasUserResponses(items []aktivitas_user.AktivitasUser) []AktivitasUserResponse {
	response := make([]AktivitasUserResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toAktivitasUserResponse(item))
	}

	return response
}

func toAktivitasUserResponse(item aktivitas_user.AktivitasUser) AktivitasUserResponse {
	return AktivitasUserResponse{
		IdAktivitas: item.IdAktivitas,
		IdPengguna:  item.IdPengguna,
		Username:    item.Username,
		Role:        item.Role,
		Action:      item.Action,
		Description: item.Description,
		IpAddress:   item.IpAddress,
		CreatedAt:   httphelper.FormatTanggalWaktuIndonesia(item.CreatedAt),
	}
}
