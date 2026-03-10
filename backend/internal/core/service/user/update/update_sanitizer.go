package user_service

import (
	"errors"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"strings"
)

func sanitizeAndValidateUpdateGuruCmd(cmd UpdateGuruCmd) (UpdateGuruCmd, *user.Email, error) {
	if cmd.Username != nil {
		trimmedUsername := strings.TrimSpace(*cmd.Username)
		if err := user.CheckUsernameLength(trimmedUsername); err != nil {
			return cmd, nil, err
		}
		cmd.Username = &trimmedUsername
	}
	if cmd.NamaLengkap != nil {
		trimmedNamaLengkap := strings.TrimSpace(*cmd.NamaLengkap)
		if trimmedNamaLengkap == "" {
			return cmd, nil, errors.New("nama_lengkap cannot be empty")
		}
		cmd.NamaLengkap = &trimmedNamaLengkap
	}
	var emailVO *user.Email
	if cmd.Email != nil {
		email, err := user.CheckNewEmail(cmd.Email)
		if err != nil {
			return cmd, nil, err
		}
		emailVO = &email
	}
	if cmd.NoHp != nil {
		trimmedNoHp := strings.TrimSpace(*cmd.NoHp)
		cmd.NoHp = &trimmedNoHp
	}
	if cmd.Foto != nil {
		trimmedFoto := strings.TrimSpace(*cmd.Foto)
		if trimmedFoto == "" {
			return cmd, nil, coreerror.ErrMissingField
		}
		cmd.Foto = &trimmedFoto
	}
	if cmd.Nip != nil {
		trimmedNip := strings.TrimSpace(*cmd.Nip)
		cmd.Nip = &trimmedNip
	}
	if cmd.Jabatan != nil {
		trimmedJabatan := strings.TrimSpace(*cmd.Jabatan)
		cmd.Jabatan = &trimmedJabatan
	}
	if cmd.BidangStudi != nil {
		trimmedBidangStudi := strings.TrimSpace(*cmd.BidangStudi)
		cmd.BidangStudi = &trimmedBidangStudi
	}
	if cmd.JenisKelamin != nil {
		trimmedJenisKelamin := strings.TrimSpace(*cmd.JenisKelamin)
		cmd.JenisKelamin = &trimmedJenisKelamin
	}
	return cmd, emailVO, nil
}
func sanitizeAndValidateUpdateSiswaCmd(cmd UpdateSiswaCmd) (UpdateSiswaCmd, *user.Email, error) {
	if cmd.Username != nil {
		trimmedUsername := strings.TrimSpace(*cmd.Username)
		if err := user.CheckUsernameLength(trimmedUsername); err != nil {
			return cmd, nil, err
		}
		cmd.Username = &trimmedUsername
	}
	if cmd.NamaLengkap != nil {
		trimmedNamaLengkap := strings.TrimSpace(*cmd.NamaLengkap)
		if trimmedNamaLengkap == "" {
			return cmd, nil, errors.New("nama_lengkap cannot be empty")
		}
		cmd.NamaLengkap = &trimmedNamaLengkap
	}
	var emailVO *user.Email
	if cmd.Email != nil {
		email, err := user.CheckNewEmail(cmd.Email)
		if err != nil {
			return cmd, nil, err
		}
		emailVO = &email
	}
	if cmd.NoHp != nil {
		trimmedNoHp := strings.TrimSpace(*cmd.NoHp)
		cmd.NoHp = &trimmedNoHp
	}
	if cmd.Foto != nil {
		trimmedFoto := strings.TrimSpace(*cmd.Foto)
		if trimmedFoto == "" {
			return cmd, nil, coreerror.ErrMissingField
		}
		cmd.Foto = &trimmedFoto
	}
	if cmd.Nisn != nil {
		nisn, err := user.CheckNewNISN(*cmd.Nisn)
		if err != nil {
			return cmd, nil, err
		}
		validatedNisn := string(nisn)
		cmd.Nisn = &validatedNisn
	}
	if cmd.NoAbsen != nil {
		if err := user.CheckAbsen(*cmd.NoAbsen); err != nil {
			return cmd, nil, err
		}
	}
	if cmd.Angkatan != nil {
		if err := user.CheckAngkatan(*cmd.Angkatan); err != nil {
			return cmd, nil, err
		}
	}
	if cmd.TempatLahir != nil {
		trimmedTempatLahir := strings.TrimSpace(*cmd.TempatLahir)
		cmd.TempatLahir = &trimmedTempatLahir
	}
	return cmd, emailVO, nil
}
