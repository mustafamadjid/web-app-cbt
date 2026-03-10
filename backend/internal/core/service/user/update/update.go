package user_service

import (
	"context"
	"errors"
	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	delete_file_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/delete_file_system"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
	txout "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
	updatepatch "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/update_patch"
	outuser "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/user"
)

type UpdateTx struct {
	txm        txout.TxManager
	sessions   out.SessionRepository
	deleteFile delete_file_repo.DeleteFileRepo
	users      outuser.UserRepository
}

func NewUpdateUserService(txm txout.TxManager, session out.SessionRepository, deleteFile delete_file_repo.DeleteFileRepo, users outuser.UserRepository) *UpdateTx {
	return &UpdateTx{txm: txm, sessions: session, deleteFile: deleteFile, users: users}
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

func (u *UpdateTx) UpdateGuru(ctx context.Context, cmd UpdateGuruCmd, actor user.Actor) error {
	logger := corelog.FromContext(ctx)
	if err := validateUpdateUserActor(actor); err != nil {
		return err
	}

	if err := validateUpdateUserID(cmd.IdPengguna); err != nil {
		return err
	}

	var (
		emailVO *user.Email
		err     error
	)
	cmd, emailVO, err = sanitizeAndValidateUpdateGuruCmd(cmd)
	if err != nil {
		return err
	}

	if cmd.Foto != nil {
		userData, err := u.users.FindUserByID(ctx, cmd.IdPengguna)
		if err != nil {
			logger.Error(ctx, "failed finding user", "layer", "core.service", "op", "user.update", "user_id", cmd.IdPengguna, "err", err)
			return err
		}

		if err := u.deleteFile.DeleteFile(ctx, userData.Foto); err != nil {
			logger.Error(ctx, "failed deleting old phot user", "layer", "core.service", "op", "user.update_user.delete_file_foto", "user_id", cmd.IdPengguna, "err", err)
			return err
		}
	}

	if hasNoFieldToUpdateGuru(cmd) {
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
		Username:     cmd.Username,
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
			if errors.Is(error, coreerror.ErrUsernameTaken) ||
				errors.Is(error, coreerror.ErrEmailTaken) ||
				errors.Is(error, coreerror.ErrNoHpTaken) {
				return error
			}

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

	if err := u.sessions.RevokeSessionAllbyUser(ctx, cmd.IdPengguna); err != nil {
		logger.Error(ctx, "failed revoking sessions", "layer", "core.service", "op", "user.update_guru", "user_id", cmd.IdPengguna, "err", err)
		return err
	}
	return nil
}

func (u *UpdateTx) UpdateSiswa(ctx context.Context, cmd UpdateSiswaCmd, actor user.Actor) error {
	logger := corelog.FromContext(ctx)
	if err := validateUpdateUserActor(actor); err != nil {
		return err
	}

	if err := validateUpdateUserID(cmd.IdPengguna); err != nil {
		return err
	}

	var (
		emailVO *user.Email
		err     error
	)
	cmd, emailVO, err = sanitizeAndValidateUpdateSiswaCmd(cmd)
	if err != nil {
		return err
	}

	if cmd.Foto != nil {
		userData, err := u.users.FindUserByID(ctx, cmd.IdPengguna)
		if err != nil {
			logger.Error(ctx, "failed finding user", "layer", "core.service", "op", "user.update_siswa", "user_id", cmd.IdPengguna, "err", err)
			return err
		}

		if err := u.deleteFile.DeleteFile(ctx, userData.Foto); err != nil {
			logger.Error(ctx, "failed deleting user", "layer", "core.service", "op", "user.update_file_foto", "user_id", cmd.IdPengguna, "err", err)
			return err
		}
	}

	if hasNoFieldToUpdateSiswa(cmd) {
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
		Username:     cmd.Username,
		NamaLengkap:  cmd.NamaLengkap,
		Email:        emailVO,
		NoHp:         cmd.NoHp,
		Foto:         cmd.Foto,
		StatusAkun:   cmd.StatusAkun,
		Role:         cmd.Role,
		JenisKelamin: cmd.JenisKelamin,
	}

	if hasPenggunaPatch(penggunaPatch) {
		if err := tx.Pengguna().UpdateUser(ctx, cmd.IdPengguna, penggunaPatch); err != nil {
			if errors.Is(err, coreerror.ErrUsernameTaken) ||
				errors.Is(err, coreerror.ErrEmailTaken) ||
				errors.Is(err, coreerror.ErrNoHpTaken) {
				return err
			}

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

	if err := u.sessions.RevokeSessionAllbyUser(ctx, cmd.IdPengguna); err != nil {
		logger.Error(ctx, "failed revoking sessions", "layer", "core.service", "op", "user.update_guru", "user_id", cmd.IdPengguna, "err", err)
		return err
	}
	return nil
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
