package user_service

import (
	"github.com/mustafamadjid/web-app-cbt/internal/core/port/out"
	delete_file_repo "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/delete_file_system"
	txout "github.com/mustafamadjid/web-app-cbt/internal/core/port/out/tx"
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
