package user_service

import "github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"

func buildCreateGuruUser(cmd CreateGuruCmd, email *user.Email, hashedPassword string) user.Pengguna {
	return user.Pengguna{
		Username:       cmd.Username,
		Email:          email,
		PasswordHashed: hashedPassword,
		NamaLengkap:    cmd.NamaLengkap,
		JenisKelamin:   cmd.JenisKelamin,
		NoHp:           cmd.NoHp,
		Foto:           cmd.Foto,
		Role:           user.GURU,
		StatusAkun:     user.AKTIF,
	}
}

func buildCreateSiswaUser(cmd CreateSiswaCmd, email *user.Email, hashedPassword string) user.Pengguna {
	return user.Pengguna{
		Username:       cmd.Username,
		Email:          email,
		PasswordHashed: hashedPassword,
		NamaLengkap:    cmd.NamaLengkap,
		JenisKelamin:   cmd.JenisKelamin,
		NoHp:           cmd.NoHp,
		Foto:           cmd.Foto,
		Role:           user.SISWA,
		StatusAkun:     user.AKTIF,
	}
}

func buildProfilGuru(cmd CreateGuruCmd, idPengguna user.ID, nip user.NIP) user.ProfilGuru {
	return user.ProfilGuru{
		IdPengguna:  idPengguna,
		Nip:         nip,
		Jabatan:     cmd.Jabatan,
		BidangStudi: cmd.BidangStudi,
	}
}

func buildProfilSiswa(cmd CreateSiswaCmd, idPengguna user.ID, nisn user.NISN) user.ProfilSiswa {
	return user.ProfilSiswa{
		IdPengguna:   idPengguna,
		IdNamaKelas:  cmd.IdNamaKelas,
		Nisn:         nisn,
		NoAbsen:      cmd.NoAbsen,
		Angkatan:     cmd.Angkatan,
		TempatLahir:  cmd.TempatLahir,
		TanggalLahir: cmd.TanggalLahir,
	}
}
