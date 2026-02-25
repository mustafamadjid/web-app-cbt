package user_service

import (
	"context"
	"errors"
	"strings"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	delete_file_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/delete_file_system"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	txout "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

type UpdateTx struct {
	txm txout.TxManager
	sessions out.SessionRepository
	deleteFile delete_file_repo.DeleteFileRepo
	users outuser.UserRepository
}

func NewUpdateUserService(txm txout.TxManager,session out.SessionRepository, deleteFile delete_file_repo.DeleteFileRepo, users outuser.UserRepository) *UpdateTx {
	return &UpdateTx{txm: txm,sessions: session, deleteFile: deleteFile, users: users}
}


func hasPenggunaPatch(p updatepatch.Pengguna) bool {
	return p.NamaLengkap != nil || p.Email != nil || p.NoHp != nil || p.Foto != nil || p.StatusAkun != nil || p.Role != nil
}

func hasProfilPatch(p updatepatch.ProfilGuru) bool {
	return p.Nip != nil || p.Jabatan != nil || p.BidangStudi != nil
}

func hasProfilSiswaPatch(p updatepatch.ProfilSiswa) bool {
	return p.IdTingkatKelas != nil || p.IdNamaKelas != nil || p.Nisn != nil || p.NoAbsen != nil || p.Angkatan != nil || p.TempatLahir != nil || p.TanggalLahir != nil
}

func (u *UpdateTx) UpdateGuru(ctx context.Context, cmd UpdateGuruCmd, actor user.Actor) error {
	logger := corelog.FromContext(ctx)
	if actor.Role != user.ADMIN {
		return coreerror.ErrForbidden
	}

	if cmd.IdPengguna == 0 {
		return errors.New("Id pengguna required")
	}

	

	// Normalisasi dan Validasi
	if cmd.Username != nil {
		s := strings.TrimSpace(*cmd.Username)
		if s == "" {
			return errors.New("username cannot be empty")
		}
		cmd.Username = &s
	}

	if cmd.NamaLengkap != nil {
		s := strings.TrimSpace(*cmd.NamaLengkap)
		if s == "" {
			return errors.New("nama_lengkap cannot be empty")
		}
		cmd.NamaLengkap = &s
	}
	var emailVO *user.Email
	if cmd.Email != nil {
		e, err := user.CheckNewEmail(*cmd.Email)
		if err != nil {
			return err
		}
		emailVO = &e
	}
	if cmd.NoHp != nil {
		s := strings.TrimSpace(*cmd.NoHp)
		cmd.NoHp = &s
	}
	if cmd.Foto != nil {
		s := strings.TrimSpace(*cmd.Foto)

		if s == "" {
			return coreerror.ErrMissingField
		}// user
		userData,err := u.users.FindUserByID(ctx,cmd.IdPengguna)
		if err != nil {
			logger.Error(ctx, "failed finding user", "layer", "core.service", "op", "user.update", "user_id", cmd.IdPengguna, "err", err)
			return err
		}
		
		if err := u.deleteFile.DeleteFile(ctx,userData.Foto); err != nil {
			logger.Error(ctx, "failed deleting old phot user", "layer", "core.service", "op", "user.update_user.delete_file_foto", "user_id", cmd.IdPengguna, "err", err)
			return err
		}

		cmd.Foto = &s
	}
	if cmd.Nip != nil {
		s := strings.TrimSpace(*cmd.Nip)
		cmd.Nip = &s
	}
	if cmd.Jabatan != nil {
		s := strings.TrimSpace(*cmd.Jabatan)
		cmd.Jabatan = &s
	}
	if cmd.BidangStudi != nil {
		s := strings.TrimSpace(*cmd.BidangStudi)
		cmd.BidangStudi = &s
	}

	if cmd.JenisKelamin != nil {
		s := strings.TrimSpace(*cmd.JenisKelamin)
		cmd.JenisKelamin = &s
	}

	if cmd.NamaLengkap == nil &&
		cmd.Email == nil &&
		cmd.NoHp == nil &&
		cmd.Foto == nil &&
		cmd.StatusAkun == nil &&
		cmd.Nip == nil &&
		cmd.Jabatan == nil &&
		cmd.BidangStudi == nil &&
		cmd.JenisKelamin == nil {
		return coreerror.ErrNoFieldToUpdate
	}

	// Transaksi
	tx, error := u.txm.Begin(ctx)
	if error != nil {
		logger.Error(ctx, "failed starting transaction", "layer", "core.service", "op", "user.update_guru", "err", error)
		return error
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Update patch pengguna
	penggunaPatch := updatepatch.Pengguna{
		NamaLengkap:  cmd.NamaLengkap,
		Email:        emailVO,
		NoHp:         cmd.NoHp,
		Foto:         cmd.Foto,
		StatusAkun:   cmd.StatusAkun,
		Role:         cmd.Role,
		JenisKelamin: cmd.JenisKelamin,
	}

	if hasPenggunaPatch(penggunaPatch) {
		if error := tx.Pengguna().UpdateUser(ctx, cmd.IdPengguna, penggunaPatch); error != nil {
			logger.Error(ctx, "failed updating pengguna", "layer", "core.service", "op", "user.update_guru", "user_id", cmd.IdPengguna, "err", error)
			return error
		}
	}

	// Update patch profil guru
	profilPatch := updatepatch.ProfilGuru{
		Nip:         cmd.Nip,
		Jabatan:     cmd.Jabatan,
		BidangStudi: cmd.BidangStudi,
	}

	if hasProfilPatch(profilPatch) {
		if error := tx.ProfilGuru().UpdateProfilGuru(ctx, cmd.IdPengguna, profilPatch); error != nil {
			logger.Error(ctx, "failed updating profil guru", "layer", "core.service", "op", "user.update_guru", "user_id", cmd.IdPengguna, "err", error)
			return error
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Error(ctx, "failed committing transaction", "layer", "core.service", "op", "user.update_guru", "user_id", cmd.IdPengguna, "err", err)
		return err
	}

	if err := u.sessions.RevokeSessionAllbyUser(ctx,cmd.IdPengguna); err != nil {
		logger.Error(ctx, "failed revoking sessions", "layer", "core.service", "op", "user.update_guru", "user_id", cmd.IdPengguna, "err", err)
		return err
	}
	return nil
}

func (u *UpdateTx) UpdateSiswa(ctx context.Context, cmd UpdateSiswaCmd, actor user.Actor) error {
	logger := corelog.FromContext(ctx)
	if actor.Role != user.ADMIN {
		return coreerror.ErrForbidden
	}

	if cmd.IdPengguna == 0 {
		return errors.New("Id pengguna required")
	}

	

	if cmd.Username != nil {
		s := strings.TrimSpace(*cmd.Username)
		if s == "" {
			return errors.New("username cannot be empty")
		}
		cmd.Username = &s
	}

	if cmd.NamaLengkap != nil {
		s := strings.TrimSpace(*cmd.NamaLengkap)
		if s == "" {
			return errors.New("nama_lengkap cannot be empty")
		}
		cmd.NamaLengkap = &s
	}

	var emailVO *user.Email
	if cmd.Email != nil {
		e, err := user.CheckNewEmail(*cmd.Email)
		if err != nil {
			return err
		}
		emailVO = &e
	}
	if cmd.NoHp != nil {
		s := strings.TrimSpace(*cmd.NoHp)
		cmd.NoHp = &s
	}
	if cmd.Foto != nil {
		s := strings.TrimSpace(*cmd.Foto)
		if s == "" {
			return coreerror.ErrMissingField
		}
		userData,err := u.users.FindUserByID(ctx,cmd.IdPengguna)
		if err != nil {
			logger.Error(ctx, "failed finding user", "layer", "core.service", "op", "user.update_siswa", "user_id", cmd.IdPengguna, "err", err)
			return err
		}

		if err := u.deleteFile.DeleteFile(ctx,userData.Foto); err != nil {
			logger.Error(ctx, "failed deleting user", "layer", "core.service", "op", "user.update_file_foto", "user_id", cmd.IdPengguna, "err", err)
			return err
		}

		cmd.Foto = &s
	}
	if cmd.Nisn != nil {
		nisn, err := user.CheckNewNISN(*cmd.Nisn)
		if err != nil {
			return err
		}
		s := string(nisn)
		cmd.Nisn = &s
	}
	if cmd.NoAbsen != nil {
		if err := user.CheckAbsen(*cmd.NoAbsen); err != nil {
			return err
		}
	}
	if cmd.Angkatan != nil {
		if err := user.CheckAngkatan(*cmd.Angkatan); err != nil {
			return err
		}
	}
	if cmd.TempatLahir != nil {
		s := strings.TrimSpace(*cmd.TempatLahir)
		cmd.TempatLahir = &s
	}

	if cmd.NamaLengkap == nil &&
		cmd.Email == nil &&
		cmd.NoHp == nil &&
		cmd.Foto == nil &&
		cmd.StatusAkun == nil &&
		cmd.Role == nil &&
		cmd.IdTingkatKelas == nil &&
		cmd.IdNamaKelas == nil &&
		cmd.Nisn == nil &&
		cmd.NoAbsen == nil &&
		cmd.Angkatan == nil &&
		cmd.TempatLahir == nil &&
		cmd.TanggalLahir == nil {
		return coreerror.ErrNoFieldToUpdate
	}

	tx, err := u.txm.Begin(ctx)
	if err != nil {
		logger.Error(ctx, "failed starting transaction", "layer", "core.service", "op", "user.update_siswa", "err", err)
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	penggunaPatch := updatepatch.Pengguna{
		NamaLengkap: cmd.NamaLengkap,
		Email:       emailVO,
		NoHp:        cmd.NoHp,
		Foto:        cmd.Foto,
		StatusAkun:  cmd.StatusAkun,
		Role:        cmd.Role,
	}

	if hasPenggunaPatch(penggunaPatch) {
		if err := tx.Pengguna().UpdateUser(ctx, cmd.IdPengguna, penggunaPatch); err != nil {
			logger.Error(ctx, "failed updating pengguna", "layer", "core.service", "op", "user.update_siswa", "user_id", cmd.IdPengguna, "err", err)
			return err
		}
	}

	profilPatch := updatepatch.ProfilSiswa{
		IdTingkatKelas: cmd.IdTingkatKelas,
		IdNamaKelas:    cmd.IdNamaKelas,
		Nisn:           cmd.Nisn,
		NoAbsen:        cmd.NoAbsen,
		Angkatan:       cmd.Angkatan,
		TempatLahir:    cmd.TempatLahir,
		TanggalLahir:   cmd.TanggalLahir,
	}

	if hasProfilSiswaPatch(profilPatch) {
		if err := tx.ProfilSiswa().UpdateProfilSiswa(ctx, cmd.IdPengguna, profilPatch); err != nil {
			logger.Error(ctx, "failed updating profil siswa", "layer", "core.service", "op", "user.update_siswa", "user_id", cmd.IdPengguna, "err", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Error(ctx, "failed committing transaction", "layer", "core.service", "op", "user.update_siswa", "user_id", cmd.IdPengguna, "err", err)
		return err
	}

	if err := u.sessions.RevokeSessionAllbyUser(ctx,cmd.IdPengguna); err != nil {
		logger.Error(ctx, "failed revoking sessions", "layer", "core.service", "op", "user.update_guru", "user_id", cmd.IdPengguna, "err", err)
		return err
	}
	return nil
}
