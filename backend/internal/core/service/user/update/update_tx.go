package user_service

import (
	"context"
	"errors"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

func (u *UpdateTx) updateGuruTx(ctx context.Context, cmd UpdateGuruCmd, emailVO *user.Email) error {
	logger := corelog.FromContext(ctx)

	tx, err := u.txm.Begin(ctx)
	if err != nil {
		logger.Error(ctx, "failed starting transaction", "layer", "core.service", "op", "user.update_guru", "err", err)
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	penggunaPatch := buildGuruPenggunaPatch(cmd, emailVO)
	if hasPenggunaPatch(penggunaPatch) {
		if err := tx.Pengguna().UpdateUser(ctx, cmd.IdPengguna, penggunaPatch); err != nil {
			if isUpdateUserConstraintErr(err) {
				return err
			}

			logger.Error(ctx, "failed updating pengguna", "layer", "core.service", "op", "user.update_guru", "user_id", cmd.IdPengguna, "err", err)
			return err
		}
	}

	profilPatch := buildProfilGuruPatch(cmd)
	if hasProfilPatch(profilPatch) {
		if err := tx.ProfilGuru().UpdateProfilGuru(ctx, cmd.IdPengguna, profilPatch); err != nil {
			logger.Error(ctx, "failed updating profil guru", "layer", "core.service", "op", "user.update_guru", "user_id", cmd.IdPengguna, "err", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Error(ctx, "failed committing transaction", "layer", "core.service", "op", "user.update_guru", "user_id", cmd.IdPengguna, "err", err)
		return err
	}

	return nil
}

func (u *UpdateTx) updateSiswaTx(ctx context.Context, cmd UpdateSiswaCmd, emailVO *user.Email) error {
	logger := corelog.FromContext(ctx)

	tx, err := u.txm.Begin(ctx)
	if err != nil {
		logger.Error(ctx, "failed starting transaction", "layer", "core.service", "op", "user.update_siswa", "err", err)
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	penggunaPatch := buildSiswaPenggunaPatch(cmd, emailVO)
	if hasPenggunaPatch(penggunaPatch) {
		if err := tx.Pengguna().UpdateUser(ctx, cmd.IdPengguna, penggunaPatch); err != nil {
			if isUpdateUserConstraintErr(err) {
				return err
			}

			logger.Error(ctx, "failed updating pengguna", "layer", "core.service", "op", "user.update_siswa", "user_id", cmd.IdPengguna, "err", err)
			return err
		}
	}

	profilPatch := buildProfilSiswaPatch(cmd)
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

	return nil
}

func isUpdateUserConstraintErr(err error) bool {
	return errors.Is(err, coreerror.ErrUsernameTaken) ||
		errors.Is(err, coreerror.ErrEmailTaken) ||
		errors.Is(err, coreerror.ErrNoHpTaken)
}
