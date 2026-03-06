package user_service

import "strings"

func sanitizeCreateGuruCmd(cmd CreateGuruCmd) CreateGuruCmd {
	cmd.Username = strings.TrimSpace(cmd.Username)
	cmd.Password = strings.TrimSpace(cmd.Password)
	cmd.NamaLengkap = strings.TrimSpace(cmd.NamaLengkap)
	cmd.JenisKelamin = strings.TrimSpace(cmd.JenisKelamin)
	cmd.Nip = strings.TrimSpace(cmd.Nip)
	cmd.Foto = strings.TrimSpace(cmd.Foto)
	cmd.Jabatan = strings.TrimSpace(cmd.Jabatan)
	cmd.BidangStudi = strings.TrimSpace(cmd.BidangStudi)

	if cmd.Email != nil {
		email := strings.TrimSpace(*cmd.Email)
		cmd.Email = &email
	}

	if cmd.NoHp != nil {
		noHp := strings.TrimSpace(*cmd.NoHp)
		cmd.NoHp = &noHp
	}

	return cmd
}
func sanitizeCreateSiswaCmd(cmd CreateSiswaCmd) CreateSiswaCmd {
	cmd.Username = strings.TrimSpace(cmd.Username)
	cmd.Password = strings.TrimSpace(cmd.Password)
	cmd.NamaLengkap = strings.TrimSpace(cmd.NamaLengkap)
	cmd.JenisKelamin = strings.TrimSpace(cmd.JenisKelamin)
	cmd.Foto = strings.TrimSpace(cmd.Foto)
	cmd.Nisn = strings.TrimSpace(cmd.Nisn)
	cmd.TempatLahir = strings.TrimSpace(cmd.TempatLahir)

	if cmd.Email != nil {
		email := strings.TrimSpace(*cmd.Email)
		cmd.Email = &email
	}

	if cmd.NoHp != nil {
		noHp := strings.TrimSpace(*cmd.NoHp)
		cmd.NoHp = &noHp
	}
	return cmd
}
