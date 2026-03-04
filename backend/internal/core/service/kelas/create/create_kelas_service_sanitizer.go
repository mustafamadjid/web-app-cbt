package kelas_service

import "strings"

func sanitizeCreateNamaKelasCmd(cmd CreateNamaKelasCmd) CreateNamaKelasCmd {
	cmd.NamaKelas = strings.TrimSpace(cmd.NamaKelas)
	return cmd
}
