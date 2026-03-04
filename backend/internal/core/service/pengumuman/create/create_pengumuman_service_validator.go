package pengumuman_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/pengumuman"
	pengumuman_service "github.com/mustafamadjid/web-app-cbt/internal/core/service/pengumuman/date_validation"
)

func validateCreatePengumumanID(data pengumuman.Pengumuman) error {
	if data.IdPengguna <= 0 {
		return coreerror.ErrMissingId
	}
	return nil
}
func validateTanggalRilisPengumuman(data pengumuman.Pengumuman) error {
	return pengumuman_service.ValidateDate(data.TanggalRilisPengumuman)
}
func validateTanggalSelesaiPengumuman(data pengumuman.Pengumuman) error {
	return pengumuman_service.ValidateDate(data.TanggalSelesaiPengumuman)
}
