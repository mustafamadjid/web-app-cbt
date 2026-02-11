package httpx

import (
	aktivitas_user_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/aktivitas_user"
	kelas_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/kelas/update"
)

type UpdateKelasHandler struct {
	svc *kelas_service.UpdateKelasService
	aktivitasUser *aktivitas_user_service.AktivitasUserService
}