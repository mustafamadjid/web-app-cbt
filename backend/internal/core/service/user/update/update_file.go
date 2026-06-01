package user_service

import (
	"context"

	"github.com/mustafamadjid/web-app-cbt/internal/core/domain/user"
	corelog "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/log"
)

func (u *UpdateTx) deleteOldUserPhoto(ctx context.Context, idPengguna user.ID, findOp string, deleteOp string, deleteErrMsg string) error {
	logger := corelog.FromContext(ctx)

	userData, err := u.users.FindUserByID(ctx, idPengguna)
	if err != nil {
		logger.Error(ctx, "failed finding user", "layer", "core.service", "op", findOp, "user_id", idPengguna, "err", err)
		return err
	}

	if err := u.deleteFile.DeleteFile(ctx, userData.Foto); err != nil {
		logger.Error(ctx, deleteErrMsg, "layer", "core.service", "op", deleteOp, "user_id", idPengguna, "err", err)
		return err
	}

	return nil
}
