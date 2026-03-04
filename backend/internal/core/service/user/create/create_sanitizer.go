package user_service

import "strings"

func sanitizeCreateGuruCmd(cmd CreateGuruCmd) CreateGuruCmd {
	cmd.Username = strings.TrimSpace(cmd.Username)
	cmd.Email = strings.TrimSpace(cmd.Email)
	cmd.Password = strings.TrimSpace(cmd.Password)
	cmd.NamaLengkap = strings.TrimSpace(cmd.NamaLengkap)
	cmd.JenisKelamin = strings.TrimSpace(cmd.JenisKelamin)
	cmd.NoHp = strings.TrimSpace(cmd.NoHp)
	cmd.Nip = strings.TrimSpace(cmd.Nip)
	cmd.Foto = strings.TrimSpace(cmd.Foto)
	cmd.Jabatan = strings.TrimSpace(cmd.Jabatan)
	cmd.BidangStudi = strings.TrimSpace(cmd.BidangStudi)
	return cmd
}
func sanitizeCreateSiswaCmd(cmd CreateSiswaCmd) CreateSiswaCmd {
	cmd.Username = strings.TrimSpace(cmd.Username)
	cmd.Email = strings.TrimSpace(cmd.Email)
	cmd.Password = strings.TrimSpace(cmd.Password)
	cmd.NamaLengkap = strings.TrimSpace(cmd.NamaLengkap)
	cmd.JenisKelamin = strings.TrimSpace(cmd.JenisKelamin)
	cmd.NoHp = strings.TrimSpace(cmd.NoHp)
	cmd.Foto = strings.TrimSpace(cmd.Foto)
	cmd.Nisn = strings.TrimSpace(cmd.Nisn)
	cmd.TempatLahir = strings.TrimSpace(cmd.TempatLahir)
	return cmd
}
