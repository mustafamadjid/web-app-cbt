package aktivitas_user_service

import "strings"

func sanitizeAktivitasUserCmd(cmd AktivitasUserCmd) AktivitasUserCmd {
	cmd.Description = strings.TrimSpace(cmd.Description)
	cmd.IpAddress = strings.TrimSpace(cmd.IpAddress)
	return cmd
}
