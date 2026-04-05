package httpx

import (
	httphelper "github.com/mustafamadjid/web-app-cbt/internal/adapter/handler/http/helper"
	authsession "github.com/mustafamadjid/web-app-cbt/internal/core/domain/auth/session"
)

func toListActiveLoginSessionResponse(items []authsession.SessionWithUser) ListActiveLoginSessionResponse {
	response := ListActiveLoginSessionResponse{
		Items: make([]ActiveLoginSessionResponse, 0, len(items)),
	}

	for _, item := range items {
		response.Items = append(response.Items, toActiveLoginSessionResponse(item))
	}

	return response
}

func toActiveLoginSessionResponse(item authsession.SessionWithUser) ActiveLoginSessionResponse {
	return ActiveLoginSessionResponse{
		Session: ActiveSessionResponse{
			SessionID:  item.Session.SessionID,
			IdPengguna: item.Session.UserID,
			Role:       item.Session.Role,
			Revoked:    item.Session.Revoked,
			ExpiresAt:  httphelper.FormatRFC3339(item.Session.ExpiresAt),
		},
		Pengguna: ActiveSessionUserResponse{
			IdPengguna:   item.Pengguna.ID,
			Username:     item.Pengguna.Username,
			Email:        item.Pengguna.Email,
			NamaLengkap:  item.Pengguna.NamaLengkap,
			JenisKelamin: item.Pengguna.JenisKelamin,
			NoHp:         item.Pengguna.NoHp,
			Role:         item.Pengguna.Role,
			StatusAkun:   item.Pengguna.StatusAkun,
			Foto:         item.Pengguna.Foto,
		},
	}
}
