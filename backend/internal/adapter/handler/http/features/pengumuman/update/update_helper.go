package httpx

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"

func ptrPengumumanID(id pengumuman.ID) *pengumuman.ID {
	return &id
}
