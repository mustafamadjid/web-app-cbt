package aktivitas_user_service

import (
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/aktivitas_user"
)

func validateAktivitasAction(cmd AktivitasUserCmd) error {
	if !cmd.Action.ValidAction() {
		return coreerror.ErrInvalidActionActivity
	}
	return nil
}
func validateAktivitasIPAddress(cmd AktivitasUserCmd) error {
	if !aktivitas_user.ValidIpAddress(cmd.IpAddress) {
		return coreerror.ErrInvalidIpAddress
	}
	return nil
}
