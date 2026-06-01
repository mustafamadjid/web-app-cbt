package user_service

import (
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
)

func buildGuruPenggunaPatch(cmd UpdateGuruCmd, emailVO *user.Email) updatepatch.Pengguna {
	return updatepatch.Pengguna{
		Username:     cmd.Username,
		NamaLengkap:  cmd.NamaLengkap,
		Email:        emailVO,
		NoHp:         cmd.NoHp,
		Foto:         cmd.Foto,
		StatusAkun:   cmd.StatusAkun,
		Role:         cmd.Role,
		JenisKelamin: cmd.JenisKelamin,
	}
}

func buildSiswaPenggunaPatch(cmd UpdateSiswaCmd, emailVO *user.Email) updatepatch.Pengguna {
	return updatepatch.Pengguna{
		Username:     cmd.Username,
		NamaLengkap:  cmd.NamaLengkap,
		Email:        emailVO,
		NoHp:         cmd.NoHp,
		Foto:         cmd.Foto,
		StatusAkun:   cmd.StatusAkun,
		Role:         cmd.Role,
		JenisKelamin: cmd.JenisKelamin,
	}
}

func buildProfilGuruPatch(cmd UpdateGuruCmd) updatepatch.ProfilGuru {
	return updatepatch.ProfilGuru{
		Nip:         cmd.Nip,
		Jabatan:     cmd.Jabatan,
		BidangStudi: cmd.BidangStudi,
	}
}

func buildProfilSiswaPatch(cmd UpdateSiswaCmd) updatepatch.ProfilSiswa {
	return updatepatch.ProfilSiswa{
		IdTingkatKelas: cmd.IdTingkatKelas,
		IdNamaKelas:    cmd.IdNamaKelas,
		Nisn:           cmd.Nisn,
		NoAbsen:        cmd.NoAbsen,
		Angkatan:       cmd.Angkatan,
		TempatLahir:    cmd.TempatLahir,
		TanggalLahir:   cmd.TanggalLahir,
	}
}

func hasPenggunaPatch(p updatepatch.Pengguna) bool {
	return p.Username != nil || p.NamaLengkap != nil || p.Email != nil || p.NoHp != nil || p.Foto != nil || p.StatusAkun != nil || p.Role != nil || p.JenisKelamin != nil
}

func hasProfilPatch(p updatepatch.ProfilGuru) bool {
	return p.Nip != nil || p.Jabatan != nil || p.BidangStudi != nil
}

func hasProfilSiswaPatch(p updatepatch.ProfilSiswa) bool {
	return p.IdTingkatKelas != nil || p.IdNamaKelas != nil || p.Nisn != nil || p.NoAbsen != nil || p.Angkatan != nil || p.TempatLahir != nil || p.TanggalLahir != nil
}

func hasNoFieldToUpdateGuru(cmd UpdateGuruCmd) bool {
	return cmd.Username == nil &&
		cmd.NamaLengkap == nil &&
		cmd.Email == nil &&
		cmd.NoHp == nil &&
		cmd.Foto == nil &&
		cmd.StatusAkun == nil &&
		cmd.Role == nil &&
		cmd.Nip == nil &&
		cmd.Jabatan == nil &&
		cmd.BidangStudi == nil &&
		cmd.JenisKelamin == nil
}

func hasNoFieldToUpdateSiswa(cmd UpdateSiswaCmd) bool {
	return cmd.Username == nil &&
		cmd.NamaLengkap == nil &&
		cmd.Email == nil &&
		cmd.NoHp == nil &&
		cmd.Foto == nil &&
		cmd.StatusAkun == nil &&
		cmd.Role == nil &&
		cmd.JenisKelamin == nil &&
		cmd.IdTingkatKelas == nil &&
		cmd.IdNamaKelas == nil &&
		cmd.Nisn == nil &&
		cmd.NoAbsen == nil &&
		cmd.Angkatan == nil &&
		cmd.TempatLahir == nil &&
		cmd.TanggalLahir == nil
}
