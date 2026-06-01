package user_service

import (
	"context"

	coreerror "github.com/mustafamadjid/web-app-cbt/internal/core/core_error"
	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
)

func (u *UpdateTx) UpdateGuru(ctx context.Context, cmd UpdateGuruCmd, actor user.Actor) error {
	if err := validateUpdateUserActor(actor); err != nil {
		return err
	}

	if err := validateUpdateUserID(cmd.IdPengguna); err != nil {
		return err
	}

	cmd, emailVO, err := sanitizeAndValidateUpdateGuruCmd(cmd)
	if err != nil {
		return err
	}

	if cmd.Foto != nil {
		if err := u.deleteOldUserPhoto(ctx, cmd.IdPengguna, "user.update", "user.update_user.delete_file_foto", "failed deleting old phot user"); err != nil {
			return err
		}
	}

	if hasNoFieldToUpdateGuru(cmd) {
		return coreerror.ErrNoFieldToUpdate
	}

	return u.updateGuruTx(ctx, cmd, emailVO)
}

func (u *UpdateTx) UpdateSiswa(ctx context.Context, cmd UpdateSiswaCmd, actor user.Actor) error {
	if err := validateUpdateUserActor(actor); err != nil {
		return err
	}

	if err := validateUpdateUserID(cmd.IdPengguna); err != nil {
		return err
	}

	cmd, emailVO, err := sanitizeAndValidateUpdateSiswaCmd(cmd)
	if err != nil {
		return err
	}

	if cmd.Foto != nil {
		if err := u.deleteOldUserPhoto(ctx, cmd.IdPengguna, "user.update_siswa", "user.update_file_foto", "failed deleting user"); err != nil {
			return err
		}
	}

	if hasNoFieldToUpdateSiswa(cmd) {
		return coreerror.ErrNoFieldToUpdate
	}

	return u.updateSiswaTx(ctx, cmd, emailVO)
}
